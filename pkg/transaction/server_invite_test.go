package transaction

import (
	"gosip/pkg/sip"
	"gosip/pkg/transaction/timer"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitServerInvite(t *testing.T) {
	setup := func() *ServerInvite {
		layer, pack := mockInvite()
		return initServerInvite(pack, layer)
	}

	t.Run("set state to proceeding and pass to TU", func(t *testing.T) {
		txn := setup()
		assert.True(t, txn.state.IsProceeding())
		assert.Same(t, txn.req, <-txn.layer.sndTU)
	})

	t.Run("send 100 if TU will not in 200ms and resend on invite retransmission", func(t *testing.T) {
		txn := setup()
		assert.Len(t, txn.layer.sndTransp, 0)
		assert.Nil(t, txn.response)
		<-time.After(100 * time.Millisecond)
		assert.Len(t, txn.layer.sndTransp, 0)
		assert.Nil(t, txn.response)
		<-time.After(110 * time.Millisecond)
		assert.Nil(t, txn.response)
		assert.Equal(t, "100", (<-txn.layer.sndTransp).Message.Code)
	})
}

func TestServerInviteConsume(t *testing.T) {
	setup := func() *ServerInvite {
		layer, pack := mockInvite()
		return initServerInvite(pack, layer)
	}
	toResp := func(pack *sip.Packet, code int, reason string) *sip.Packet {
		return &sip.Packet{
			Message: pack.Message.Response(code, reason),
		}
	}

	t.Run("ignore is missing sip message", func(t *testing.T) {
		txn := setup()
		assert.True(t, txn.state.IsProceeding())
		assert.NotPanics(t, func() {
			txn.Consume(&sip.Packet{})
		})
		assert.Len(t, txn.layer.sndTU, 1)
		assert.Len(t, txn.layer.sndTransp, 0)
	})

	t.Run("Proceeding", func(t *testing.T) {
		t.Run("do not send 100 if TU sends 1xx within 200ms", func(t *testing.T) {
			txn := setup()
			assert.Len(t, txn.layer.sndTransp, 0)
			resp := &sip.Packet{Message: txn.req.Message.Response(180, "Ringing")}
			txn.Consume(resp)
			assert.Same(t, resp, txn.response)
			assert.Len(t, txn.layer.sndTransp, 1)
			<-time.After(210 * time.Millisecond)
			assert.Same(t, resp, txn.response)
			assert.Len(t, txn.layer.sndTransp, 1)
		})

		t.Run("on request retransmission resend default 100 if no response from TU", func(t *testing.T) {
			txn := setup()
			// invite retransaction
			txn.Consume(txn.req)
			assert.Equal(t, "100", (<-txn.layer.sndTransp).Message.Code)
			txn.Consume(txn.req)
			assert.Equal(t, "100", (<-txn.layer.sndTransp).Message.Code)
			<-time.After(210 * time.Millisecond)
			assert.Len(t, txn.layer.sndTransp, 0)
		})

		t.Run("on request retransmission resend last response from TU", func(t *testing.T) {
			txn := setup()
			resp := toResp(txn.req, 183, "Session in progress")
			txn.Consume(resp)
			assert.Equal(t, "183", (<-txn.layer.sndTransp).Message.Code)
			assert.Same(t, resp, txn.response)
			txn.Consume(txn.req)
			assert.Equal(t, "183", (<-txn.layer.sndTransp).Message.Code)
		})

		t.Run("on 2xx response pass to transport and terminates", func(t *testing.T) {
			txn := setup()
			resp := toResp(txn.req, 200, "OK")
			txn.Consume(resp)
			assert.Same(t, resp, txn.response)
			assert.Equal(t, "200", (<-txn.layer.sndTransp).Message.Code)
			assert.True(t, txn.state.IsTerminated())
		})

		t.Run("300 to 699 response transmit to completed", func(t *testing.T) {
			txn := setup()
			resp := toResp(txn.req, 404, "Not Found")
			txn.Consume(resp)
			assert.Same(t, resp, txn.response)
			assert.Equal(t, "404", (<-txn.layer.sndTransp).Message.Code)
			assert.True(t, txn.state.IsCompleted())
		})
	})

	t.Run("Completed", func(t *testing.T) {
		t.Run("start H timer for any transport", func(t *testing.T) {
			layer, pack := mockInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.H = 10 * time.Millisecond
				return t
			}
			txn := initServerInvite(pack, layer)
			resp := toResp(pack, 404, "Not Found")

			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())

			<-time.After(50 * time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
			assert.Equal(t, ErrTxnFail, <-txn.layer.Err())
		})

		t.Run("on request retransmission resend response", func(t *testing.T) {
			txn := setup()
			req := txn.req
			resp := toResp(req, 404, "Not Found")
			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Equal(t, "404", (<-txn.layer.sndTransp).Message.Code)

			txn.Consume(req)
			assert.Equal(t, "404", (<-txn.layer.sndTransp).Message.Code)

			txn.Consume(req)
			assert.Equal(t, "404", (<-txn.layer.sndTransp).Message.Code)
		})

		t.Run("resend response on G timer", func(t *testing.T) {
			layer, pack := mockInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.T1 = 1 * time.Millisecond
				t.T2 = 16 * time.Millisecond
				return t
			}
			txn := initServerInvite(pack, layer)
			resp := toResp(pack, 404, "Not Found")
			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Equal(t, "404", (<-txn.layer.sndTransp).Message.Code)

			<-time.After(34 * time.Millisecond)
			assert.Len(t, txn.layer.sndTransp, 5)
		})

		t.Run("on ACK transition to the confirmed state", func(t *testing.T) {
			txn := setup()
			req := txn.req
			resp := toResp(req, 404, "Not Found")
			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Equal(t, "404", (<-txn.layer.sndTransp).Message.Code)

			ack := &sip.Packet{
				Message: req.Message.ACK(resp.Message),
			}
			txn.Consume(ack)
			assert.True(t, txn.state.IsConfirmed())
		})
	})

	t.Run("Confirmed", func(t *testing.T) {
		t.Run("absorb ACK and ignore INVITE retransmissions", func(t *testing.T) {
			txn := setup()
			pack := txn.req
			resp := toResp(pack, 404, "Not Found")
			ack := &sip.Packet{Message: pack.Message.ACK(resp.Message)}
			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			txn.Consume(ack)
			assert.True(t, txn.state.IsConfirmed())
			assert.Len(t, txn.layer.sndTU, 1)
			assert.Len(t, txn.layer.sndTransp, 1)
			txn.Consume(pack)
			txn.Consume(pack)
			txn.Consume(pack)
			txn.Consume(ack)
			txn.Consume(ack)
			txn.Consume(ack)
			assert.Len(t, txn.layer.sndTU, 1)
			assert.Len(t, txn.layer.sndTransp, 1)
		})

		t.Run("when timer I fired then terminated for unreliable transport", func(t *testing.T) {
			layer, pack := mockInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.T4 = 10 * time.Millisecond
				return t
			}
			txn := initServerInvite(pack, layer)
			resp := toResp(pack, 404, "Not Found")
			ack := &sip.Packet{Message: pack.Message.ACK(resp.Message)}
			txn.Consume(resp)
			txn.Consume(ack)
			assert.True(t, txn.state.IsConfirmed())
			<-time.After(20 * time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
		})

		t.Run("do not fire timer I for reliable transport", func(t *testing.T) {
			layer, pack := mockInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.T1 = time.Millisecond
				t.B = t.T1 * 64
				return t
			}
			pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
			txn := initServerInvite(pack, layer)
			resp := toResp(pack, 404, "Not Found")
			ack := &sip.Packet{Message: pack.Message.ACK(resp.Message)}
			txn.Consume(resp)
			txn.Consume(ack)
			<-time.After(time.Millisecond) // to ensure goroutin worked
			assert.True(t, txn.state.IsTerminated())
		})
	})
}
