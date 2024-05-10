package transaction

import (
	"gosip/pkg/sip"
	"gosip/pkg/transaction/timer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitServerNonInvite(t *testing.T) {
	layer := Init()
	pack := &sip.Packet{
		Message: mockRegisterMsg(),
	}
	txn := initServerNonInvite(pack, layer)
	assert.Nil(t, txn.response)
	assert.True(t, txn.state.IsTrying())
	assert.Same(t, pack, <-txn.layer.sndTU)
}

func TestServerNonInviteConsume(t *testing.T) {
	createLayerTxn := func(t *testing.T) *ServerNonInvite {
		layer := Init()
		pack := &sip.Packet{Message: mockRegisterMsg()}
		layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
			t.J = time.Millisecond
			return t
		}
		txn := initServerNonInvite(pack, layer)
		assert.True(t, txn.state.IsTrying())
		assert.Same(t, txn.req, <-txn.layer.sndTU)
		return txn
	}

	t.Run("ignore on message input is nil", func(t *testing.T) {
		txn := createLayerTxn(t)
		assert.NotPanics(t, func() {
			txn.Consume(&sip.Packet{})
		})
	})

	t.Run("consume 1xx while in trying state", func(t *testing.T) {
		txn := createLayerTxn(t)
		resp := &sip.Packet{
			Message: txn.req.Message.Response(180, "Ringing"),
		}

		txn.Consume(resp)
		assert.True(t, txn.state.IsProceeding())
		assert.Same(t, resp, txn.response)
		assert.Same(t, resp, <-txn.layer.sndTransp)
	})

	t.Run("consume 200-699 while in trying state", func(t *testing.T) {
		txn := createLayerTxn(t)
		resp := &sip.Packet{
			Message: txn.req.Message.Response(200, "Ok"),
		}

		txn.Consume(resp)
		assert.True(t, txn.state.IsCompleted())
		assert.Same(t, resp, txn.response)
		assert.Same(t, resp, <-txn.layer.sndTransp)
	})

	t.Run("consume 1xx while in proceeding state", func(t *testing.T) {
		txn := createLayerTxn(t)
		resp100 := &sip.Packet{
			Message: txn.req.Message.Response(100, "Trying"),
		}
		txn.Consume(resp100)
		assert.True(t, txn.state.IsProceeding())
		assert.Same(t, resp100, txn.response)
		assert.Same(t, resp100, <-txn.layer.sndTransp)

		resp180 := &sip.Packet{
			Message: txn.req.Message.Response(180, "Ringing"),
		}
		txn.Consume(resp180)
		assert.True(t, txn.state.IsProceeding())
		assert.Same(t, resp180, txn.response)
		assert.Same(t, resp180, <-txn.layer.sndTransp)
	})

	t.Run("request retransmitting while in processing state", func(t *testing.T) {
		txn := createLayerTxn(t)
		resp100 := &sip.Packet{
			Message: txn.req.Message.Response(100, "Trying"),
		}
		txn.Consume(resp100)
		assert.True(t, txn.state.IsProceeding())
		assert.Same(t, resp100, txn.response)
		assert.Same(t, resp100, <-txn.layer.sndTransp)

		// resend last provisioning
		t.Run("resend last provisioning", func(t *testing.T) {
			req := &sip.Packet{
				Message: txn.req.Message,
			}
			txn.Consume(req)
			assert.True(t, txn.state.IsProceeding())
			assert.Same(t, resp100, <-txn.layer.sndTransp)
		})
	})

	t.Run("transport error while in processing state", func(t *testing.T) {
		t.Skip("TODO with transport")
	})

	t.Run("consume 200-699 while in processing state", func(t *testing.T) {
		txn := createLayerTxn(t)
		resp100 := &sip.Packet{
			Message: txn.req.Message.Response(100, "Trying"),
		}
		txn.Consume(resp100)
		assert.True(t, txn.state.IsProceeding())
		assert.Same(t, resp100, txn.response)
		assert.Same(t, resp100, <-txn.layer.sndTransp)

		resp404 := &sip.Packet{
			Message: txn.req.Message.Response(404, "Not Found"),
		}
		t.Run("enter completed state", func(t *testing.T) {
			txn.Consume(resp404)
			assert.True(t, txn.state.IsCompleted())
			assert.Same(t, resp404, txn.response)
			assert.Same(t, resp404, <-txn.layer.sndTransp)
		})

		t.Run("resend last final response on request retransaction", func(t *testing.T) {
			req := &sip.Packet{
				Message: txn.req.Message,
			}
			txn.Consume(req)
			assert.True(t, txn.state.IsCompleted())
			assert.Same(t, resp404, <-txn.layer.sndTransp)
		})

		t.Run("ignore any final response from TU", func(t *testing.T) {
			txn.Consume(resp404)
			assert.True(t, txn.state.IsCompleted())
			select {
			case <-txn.layer.sndTransp:
				t.Error("should not send to transport")
			default:
				// passed
			}
		})
	})

	t.Run("transport error while in completed state", func(t *testing.T) {
		t.Skip("TODO with transport")
	})

	t.Run("timer j in completed state move to terminated", func(t *testing.T) {
		txn := createLayerTxn(t)
		txn.layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
			t.J = 10 * time.Millisecond // reset to low value
			return t
		}
		txn.layer.pool.Add(txn)
		assert.Equal(t, 1, txn.layer.pool.Len())

		resp := &sip.Packet{
			Message: txn.req.Message.Response(404, "Not Found"),
		}
		txn.Consume(resp)
		assert.True(t, txn.state.IsCompleted())
		assert.Equal(t, 1, txn.layer.pool.Len())
		assert.Eventually(t, func() bool {
			return txn.layer.pool.Len() == 0
		}, 12*time.Millisecond, 2*time.Millisecond)
	})
}

func TestServerNonInviteMatch(t *testing.T) {
	t.Skip("TODO test transaction match")
}
