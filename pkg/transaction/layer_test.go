package transaction

import (
	"gosip/pkg/sip"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLayerClientInvite(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		layer, pack := mockInvite()
		assert.Equal(t, 0, layer.pool.Len())
		layer.RecvTU(pack)
		assert.Same(t, pack, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		txn, found := layer.pool.Get(pack.Message.TopViaBranch())
		assert.True(t, found)
		assert.IsType(t, &ClientInvite{}, txn)

		resp := &sip.Packet{Message: pack.Message.Response(200, "OK")}
		layer.RecvTransp(resp)
		assert.Same(t, resp, <-layer.SendTU())
		assert.Equal(t, 0, layer.pool.Len())
	})

	t.Run("client error response 404", func(t *testing.T) {
		layer, pack := mockInvite()
		assert.Equal(t, 0, layer.pool.Len())
		pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
		layer.RecvTU(pack)
		assert.Same(t, pack, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		resp := &sip.Packet{Message: pack.Message.Response(404, "Not Found")}
		layer.RecvTransp(resp)
		assert.Equal(t, "ACK", (<-layer.SendTransp()).Message.Method)
		assert.Equal(t, 0, layer.pool.Len())
	})

	t.Run("destroy on from layer", func(t *testing.T) {
		layer, pack := mockInvite()
		layer.RecvTU(pack)
		assert.Same(t, pack, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		layer.Destroy(pack)
		assert.Equal(t, 0, layer.pool.Len())
	})
}

func TestLayerClientNonInvite(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		layer, pack := mockNonInvite()
		pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
		assert.Equal(t, 0, layer.pool.Len())
		layer.RecvTU(pack)
		assert.Same(t, pack, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		txn, found := layer.pool.Get(pack.Message.TopViaBranch())
		assert.True(t, found)
		assert.IsType(t, &ClientNonInvite{}, txn)

		resp := &sip.Packet{Message: pack.Message.Response(200, "OK")}
		layer.RecvTransp(resp)
		assert.Same(t, resp, <-layer.SendTU())
		assert.Equal(t, 0, layer.pool.Len())
	})

	t.Run("early response to proceed state and then complete", func(t *testing.T) {
		layer, pack := mockNonInvite()
		pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
		layer.RecvTU(pack)
		assert.Same(t, pack, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		respEarly := &sip.Packet{Message: pack.Message.Response(100, "Trying")}
		layer.RecvTransp(respEarly)
		assert.Same(t, respEarly, <-layer.SendTU())
		assert.Equal(t, 1, layer.pool.Len())

		respOk := &sip.Packet{Message: pack.Message.Response(200, "OK")}
		layer.RecvTransp(respOk)
		assert.Same(t, respOk, <-layer.SendTU())
		assert.Equal(t, 0, layer.pool.Len())
	})

	t.Run("destroy from layer", func(t *testing.T) {
		layer, pack := mockNonInvite()
		layer.RecvTU(pack)
		assert.Same(t, pack, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		layer.Destroy(pack)
		assert.Equal(t, 0, layer.pool.Len())
	})
}

func TestLayerServerInvite(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		layer, pack := mockInvite()
		assert.Equal(t, 0, layer.pool.Len())
		layer.RecvTransp(pack)
		assert.Same(t, pack, <-layer.SendTU())
		assert.Equal(t, 1, layer.pool.Len())

		txn, found := layer.pool.Get(pack.Message.TopViaBranch())
		assert.True(t, found)
		assert.IsType(t, &ServerInvite{}, txn)

		respEarly := &sip.Packet{Message: pack.Message.Response(100, "Trying")}
		layer.RecvTU(respEarly)
		assert.Same(t, respEarly, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		respOK := &sip.Packet{Message: pack.Message.Response(200, "OK")}
		layer.RecvTU(respOK)
		assert.Same(t, respOK, <-layer.SendTransp())
		assert.Equal(t, 0, layer.pool.Len())
	})

	t.Run("on error response will ACK", func(t *testing.T) {
		layer, pack := mockInvite()
		pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
		assert.Equal(t, 0, layer.pool.Len())
		layer.RecvTransp(pack)
		assert.Equal(t, 1, layer.pool.Len())

		resp := &sip.Packet{Message: pack.Message.Response(404, "Not Found")}
		layer.RecvTU(resp)
		assert.Same(t, resp, <-layer.SendTransp())
		ack := &sip.Packet{Message: pack.Message.ACK(resp.Message)}
		layer.RecvTransp(ack)
		<-time.After(time.Millisecond)
		assert.Equal(t, 0, layer.pool.Len())
	})

	t.Run("destroy from layer", func(t *testing.T) {
		layer, pack := mockInvite()
		layer.RecvTransp(pack)
		assert.Same(t, pack, <-layer.SendTU())
		assert.Equal(t, 1, layer.pool.Len())

		layer.Destroy(pack)
		assert.Equal(t, 0, layer.pool.Len())
	})
}

func TestLayerServerNonInvite(t *testing.T) {
	t.Run("final response", func(t *testing.T) {
		tests := map[string]struct {
			code   int
			reason string
		}{
			`200 ok`:           {200, "OK"},
			`303 redirect`:     {303, "Redirect"},
			`404 not found`:    {404, "Not Found"},
			`500 server error`: {500, "Server Error"},
			`600 global error`: {600, "Error"},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				layer, pack := mockNonInvite()
				pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
				layer.RecvTransp(pack)
				assert.Same(t, pack, <-layer.SendTU())
				assert.Equal(t, 1, layer.pool.Len())

				txn, found := layer.pool.Get(pack.Message.TopViaBranch())
				assert.True(t, found)
				assert.IsType(t, &ServerNonInvite{}, txn)

				resp := &sip.Packet{Message: pack.Message.Response(tc.code, tc.reason)}
				layer.RecvTU(resp)
				assert.Same(t, resp, <-layer.SendTransp())
				assert.Equal(t, 0, layer.pool.Len())
			})
		}
	})

	t.Run("early response", func(t *testing.T) {
		layer, pack := mockNonInvite()
		pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
		layer.RecvTransp(pack)
		assert.Same(t, pack, <-layer.SendTU())
		assert.Equal(t, 1, layer.pool.Len())

		resp := &sip.Packet{Message: pack.Message.Response(100, "Trying")}
		layer.RecvTU(resp)
		assert.Same(t, resp, <-layer.SendTransp())
		assert.Equal(t, 1, layer.pool.Len())

		resp = &sip.Packet{Message: pack.Message.Response(404, "Not Found")}
		layer.RecvTU(resp)
		assert.Same(t, resp, <-layer.SendTransp())
		assert.Equal(t, 0, layer.pool.Len())
	})
}

func TestDestroy(t *testing.T) {
	layer, pack := mockInvite()
	assert.NotPanics(t, func() {
		layer.Destroy(pack)

	}, "not found")
	pack.Message = nil
	assert.NotPanics(t, func() {
		layer.Destroy(pack)
	}, "Message is nil")
}

func TestRecvTU(t *testing.T) {
	layer, pack := mockInvite()
	pack.Message = nil
	assert.NotPanics(t, func() {
		layer.RecvTU(pack)
	}, "Message is nil")
}

func TestRecvTransp(t *testing.T) {
	layer, pack := mockInvite()
	pack.Message = nil
	assert.NotPanics(t, func() {
		layer.RecvTransp(pack)
	}, "Message is nil")
}
