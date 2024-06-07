package transaction

import (
	"gosip/pkg/sip"
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

	t.Run("ignore is missing sip message", func(t *testing.T) {
		txn := setup()
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
			resp := txn.gen100()
			resp.Message.Code = "183"
			resp.Message.Reason = "Session in progress"
			txn.Consume(resp)
			assert.Equal(t, "183", (<-txn.layer.sndTransp).Message.Code)
			assert.Same(t, resp, txn.response)
			txn.Consume(txn.req)
			assert.Equal(t, "183", (<-txn.layer.sndTransp).Message.Code)
		})
	})
}
