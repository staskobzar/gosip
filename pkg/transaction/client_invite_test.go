package transaction

import (
	"gosip/pkg/sip"
	"gosip/pkg/transaction/state"
	"gosip/pkg/transaction/timer"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientInviteInit(t *testing.T) {
	t.Run("init transaction in Calling state", func(t *testing.T) {
		layer, pack := mockInvite()
		txn := initClientInvite(pack, layer)

		assert.True(t, txn.state.IsCalling())
		assert.Same(t, pack, <-layer.SendTransp())
	})

	t.Run("timer A", func(t *testing.T) {
		setup := func() (*Layer, *sip.Packet) {
			layer, pack := mockInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.T1 = time.Millisecond
				return t
			}
			return layer, pack
		}

		t.Run("retransmit requests for unreliable transport", func(t *testing.T) {
			layer, pack := setup()
			txn := initClientInvite(pack, layer)
			<-time.After(20 * time.Millisecond)
			assert.Len(t, txn.layer.sndTransp, 5)
		})

		t.Run("no requests retransmitions for reliable transport", func(t *testing.T) {
			layer, pack := setup()
			pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
			txn := initClientInvite(pack, layer)
			<-time.After(20 * time.Millisecond)
			assert.Len(t, txn.layer.sndTransp, 1)
		})
	})

	t.Run("timer B", func(t *testing.T) {
		setup := func(sendAddr net.Addr) *ClientInvite {
			layer, pack := mockInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.T1 = time.Millisecond
				t.B = t.T1 * 64
				return t
			}
			pack.SendTo = []net.Addr{sendAddr}
			return initClientInvite(pack, layer)
		}

		t.Run("starts for unreliable transport and send timeout error", func(t *testing.T) {
			txn := setup(&net.UDPAddr{IP: net.IPv4(192, 168, 1, 100), Port: 5060})
			assert.True(t, txn.state.IsCalling())
			<-time.After(70 * time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
			assert.InDelta(t, 7, len(txn.layer.sndTransp), 1)
			assert.ErrorIs(t, <-txn.layer.Err(), ErrTimeout)
		})

		t.Run("starts for reliable transport and send timeout error", func(t *testing.T) {
			txn := setup(&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060})
			assert.True(t, txn.state.IsCalling())
			<-time.After(70 * time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
			assert.Len(t, txn.layer.sndTransp, 1)
			assert.ErrorIs(t, <-txn.layer.Err(), ErrTimeout)
		})

		t.Run("do not send timeout error when state is not calling", func(t *testing.T) {
			txn := setup(&net.UDPAddr{IP: net.IPv4(192, 168, 1, 100), Port: 5060})
			txn.state.Set(state.Proceeding)
			<-time.After(70 * time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
			assert.Len(t, txn.layer.err, 0)
		})
	})
}

func TestClientInviteConsume(t *testing.T) {
	setup := func(sendAddr net.Addr) *ClientInvite {
		layer, pack := mockInvite()
		if sendAddr != nil {
			pack.SendTo = []net.Addr{sendAddr}
		} else {
			pack.SendTo = []net.Addr{&net.UDPAddr{IP: net.IPv4(192, 168, 1, 100), Port: 5060}}
		}
		return initClientInvite(pack, layer)
	}
	toResp := func(pack *sip.Packet, code int, reason string) *sip.Packet {
		return &sip.Packet{
			Message: pack.Message.Response(code, reason),
		}
	}

	t.Run("ignore on SIP Message is nil", func(t *testing.T) {
		txn := setup(nil)
		assert.NotPanics(t, func() {
			txn.Consume(&sip.Packet{})
		})
		assert.Len(t, txn.layer.sndTU, 0)
	})

	t.Run("ignore on message is request", func(t *testing.T) {
		txn := setup(nil)
		assert.NotPanics(t, func() {
			txn.Consume(txn.req)
		})
		assert.Len(t, txn.layer.sndTU, 0)
	})

	t.Run("Calling", func(t *testing.T) {
		t.Run("recv provisional transitions to Proceeding", func(t *testing.T) {
			txn := setup(nil)
			resp := toResp(txn.req, 100, "Trying")
			assert.Equal(t, "INVITE", (<-txn.layer.SendTransp()).Message.Method, "initial invite to transport")

			txn.Consume(resp)
			assert.True(t, txn.state.IsProceeding())
			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Len(t, txn.layer.sndTransp, 0, "no ACK to transport")
		})

		t.Run("recv 2xx terminates transaction", func(t *testing.T) {
			txn := setup(nil)
			resp := toResp(txn.req, 200, "OK")

			assert.Len(t, txn.layer.sndTransp, 1, "send initial invite request")
			txn.Consume(resp)
			assert.True(t, txn.state.IsTerminated())
			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Len(t, txn.layer.sndTransp, 1, "")
		})

		t.Run("recv 3xx-699 transitions to Completed", func(t *testing.T) {
			txn := setup(nil)
			resp := toResp(txn.req, 404, "OK")
			assert.Equal(t, "INVITE", (<-txn.layer.SendTransp()).Message.Method, "initial invite to transport")

			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Len(t, txn.layer.sndTransp, 1, "on 3xx-699 ACK should be sent by transaction")
			assert.Equal(t, "ACK", (<-txn.layer.SendTransp()).Message.Method)
		})

		t.Run("recv 3xx-699 for reliable transport terminates immediately", func(t *testing.T) {
			txn := setup(&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060})
			resp := toResp(txn.req, 404, "OK")
			assert.Equal(t, "INVITE", (<-txn.layer.SendTransp()).Message.Method, "initial invite to transport")

			txn.Consume(resp)
			assert.True(t, txn.state.IsTerminated())
			assert.Len(t, txn.layer.sndTransp, 1)
			assert.Equal(t, "ACK", (<-txn.layer.SendTransp()).Message.Method)
		})
	})

	t.Run("Proceeding", func(t *testing.T) {
		t.Run("any provisional response pass to TU", func(t *testing.T) {
			txn := setup(nil)
			resp := toResp(txn.req, 180, "Ringing")
			txn.state.Set(state.Proceeding)

			txn.Consume(resp)
			txn.Consume(resp)
			txn.Consume(resp)

			assert.True(t, txn.state.IsProceeding())

			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Same(t, resp, <-txn.layer.SendTU())
		})

		t.Run("recv 2xx terminates transaction", func(t *testing.T) {
			txn := setup(nil)
			resp := toResp(txn.req, 200, "OK")
			assert.Len(t, txn.layer.sndTransp, 1, "send initial invite request")
			txn.state.Set(state.Proceeding)

			txn.Consume(resp)
			assert.True(t, txn.state.IsTerminated())
			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Len(t, txn.layer.sndTransp, 1, "")
		})

		t.Run("recv 3xx-699 transitions to Completed", func(t *testing.T) {
			txn := setup(nil)
			resp := toResp(txn.req, 404, "OK")
			assert.Equal(t, "INVITE", (<-txn.layer.SendTransp()).Message.Method, "initial invite to transport")
			txn.state.Set(state.Proceeding)

			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Same(t, resp, <-txn.layer.SendTU())
			assert.Len(t, txn.layer.sndTransp, 1, "on 3xx-699 ACK should be sent by transaction")
			assert.Equal(t, "ACK", (<-txn.layer.SendTransp()).Message.Method)
		})
	})

	t.Run("Completed ACK any response retransmissions", func(t *testing.T) {
		txn := setup(nil)
		resp := toResp(txn.req, 500, "Server Not Ready")
		assert.Equal(t, "INVITE", (<-txn.layer.SendTransp()).Message.Method, "initial invite to transport")

		txn.Consume(resp)
		assert.True(t, txn.state.IsCompleted())
		assert.Same(t, resp, <-txn.layer.SendTU())
		assert.Equal(t, "ACK", (<-txn.layer.SendTransp()).Message.Method)

		txn.Consume(resp)
		assert.Equal(t, "ACK", (<-txn.layer.SendTransp()).Message.Method)

		txn.Consume(resp)
		assert.Equal(t, "ACK", (<-txn.layer.SendTransp()).Message.Method)

		assert.Len(t, txn.layer.sndTU, 0)
	})
}
