package txnlayer

import (
	"errors"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createMock() (*mockEndPoint, *mockTransp, *mockMsg, netip.AddrPort) {
	ep := &mockEndPoint{msg: make([]Message, 0)}
	tr := &mockTransp{msg: make([]Message, 0)}
	msg := &mockMsg{}
	addr, _ := netip.ParseAddrPort("127.0.0.1:5670")
	return ep, tr, msg, addr
}

func TestCreateClientInvTnx(t *testing.T) {
	endpoint, transp, msg, _ := createMock()

	txn := createClientInvTxn(transp, endpoint, msg)
	assert.Equal(t, Unknown, txn.state.Load())
}

func TestTxnClientInviteCalling(t *testing.T) {
	t.Run("retransactions for not reliable transport", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()
		txn := createClientInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond

		assert.False(t, txn.transp.IsReliable())

		txn.calling(msg)

		<-time.After(140 * time.Millisecond)

		assert.Equal(t, Calling, txn.state.Load())
		assert.Equal(t, 8, transp.msgLen())
	})

	t.Run("no retransactions for reliable transport", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()
		txn := createClientInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond
		transp.isReliable = true

		txn.calling(msg)
		<-time.After(40 * time.Millisecond)
		assert.Equal(t, 1, transp.msgLen())
	})
}

func TestTxnClientInviteInit(t *testing.T) {
	t.Run("init with reliable transport", func(t *testing.T) {
		endpoint, transp, msg, addr := createMock()
		transp.isReliable = true
		txn := createClientInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond

		txn.Init(msg, addr)
		assert.Equal(t, Calling, txn.state.Load())
		<-time.After(65 * time.Millisecond)
		assert.Equal(t, Terminated, txn.state.Load())
		assert.Equal(t, 1, transp.msgLen())
		assert.Equal(t, "127.0.0.1:5670", transp.addr.String())
	})

	t.Run("init with non-reliable transport", func(t *testing.T) {
		endpoint, transp, msg, addr := createMock()
		txn := createClientInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond

		txn.Init(msg, addr)
		assert.Equal(t, Calling, txn.state.Load())
		<-time.After(100 * time.Millisecond)
		assert.Equal(t, Terminated, txn.state.Load())
		assert.True(t, transp.msgLen() > 5)
		assert.Equal(t, "127.0.0.1:5670", transp.addr.String())
	})
}

func TestTxnClientInviteConsume(t *testing.T) {
	endPointMsg := func(txn *TxnClientInvite, index int) *mockMsg {
		return txn.endpoint.(*mockEndPoint).msg[index].(*mockMsg)
	}
	transpMsg := func(txn *TxnClientInvite, index int) *mockMsg {
		return txn.transp.(*mockTransp).msg[index].(*mockMsg)
	}

	t.Run("ignore if Message is not response", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()

		txn := createClientInvTxn(transp, endpoint, msg)
		txn.Consume(&mockMsg{})

		assert.Equal(t, Unknown, txn.state.Load())
	})

	t.Run("on calling", func(t *testing.T) {
		tests := map[string]struct {
			respCode  int
			wantState uint32
			wantSentN int
			lastAck   bool
		}{
			`got early 1XX respons`:         {180, Proceeding, 1, false},
			`got confirm 2XX response`:      {200, Terminated, 1, false},
			`got client error 4XX response`: {404, Completed, 2, true},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				endpoint, transp, msg, addr := createMock()

				txn := createClientInvTxn(transp, endpoint, msg)
				txn.Init(msg, addr)
				assert.Equal(t, Calling, txn.state.Load())

				// make sure first msg is send from A timer loop
				<-time.After(10 * time.Millisecond)

				resp := &mockMsg{code: tc.respCode, isResp: true}
				txn.Consume(resp)

				assert.Equal(t, tc.wantState, txn.state.Load())
				assert.Equal(t, resp, endPointMsg(txn, 0))

				assert.Equal(t, tc.wantSentN, transp.msgLen())
				assert.Equal(t, tc.lastAck,
					transpMsg(txn, tc.wantSentN-1).isAck)
			})
		}

		t.Run("re-transmission transport error to TU", func(t *testing.T) {
			endpoint, transp, msg, _ := createMock()
			transp.senderr = errors.New("failed to send")

			txn := createClientInvTxn(transp, endpoint, msg)
			txn.state.Store(Calling)

			txn.Consume(&mockMsg{code: 300, isResp: true})
			assert.Equal(t, Terminated, txn.state.Load())
			assert.Error(t, endpoint.err, "failed to send")
		})

		t.Run("terminates on reliable transport", func(t *testing.T) {
			endpoint, transp, msg, _ := createMock()
			transp.isReliable = true
			txn := createClientInvTxn(transp, endpoint, msg)
			txn.state.Store(Calling)
			txn.Consume(&mockMsg{code: 404, isResp: true})
			assert.Equal(t, Terminated, txn.state.Load())
		})
	})

	t.Run("on proceeding", func(t *testing.T) {
		tests := map[string]struct {
			respCode  int
			wantState uint32
			wantSentN int
			lastAck   bool
		}{
			`got early 1XX respons`:         {100, Proceeding, 0, false},
			`got confirm 2XX response`:      {200, Terminated, 0, false},
			`got client error 5XX response`: {504, Completed, 1, true},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				endpoint, transp, msg, _ := createMock()

				txn := createClientInvTxn(transp, endpoint, msg)
				txn.state.Store(Proceeding)

				resp := &mockMsg{code: tc.respCode, isResp: true}
				txn.Consume(resp)
				assert.Equal(t, tc.wantState, txn.state.Load())
				assert.Equal(t, tc.wantSentN, transp.msgLen())
				if tc.wantSentN > 0 {
					assert.Equal(t, tc.lastAck,
						transpMsg(txn, tc.wantSentN-1).isAck)
				}
			})
		}

		t.Run("terminates on reliable transport", func(t *testing.T) {
			endpoint, transp, msg, _ := createMock()
			transp.isReliable = true
			txn := createClientInvTxn(transp, endpoint, msg)
			txn.state.Store(Proceeding)
			txn.Consume(&mockMsg{code: 404, isResp: true})
			assert.Equal(t, Terminated, txn.state.Load())
		})
	})

	t.Run("on completed", func(t *testing.T) {
		t.Run("absorb 300-699 responses", func(t *testing.T) {
			endpoint, transp, msg, _ := createMock()

			txn := createClientInvTxn(transp, endpoint, msg)
			txn.state.Store(Completed)

			txn.Consume(&mockMsg{code: 100, isResp: true})

			assert.Nil(t, endpoint.err)
			assert.Equal(t, Completed, txn.state.Load())
			assert.Equal(t, 0, transp.msgLen())

			txn.Consume(&mockMsg{code: 600, isResp: true})
			assert.Equal(t, Completed, txn.state.Load())
			assert.Equal(t, 1, transp.msgLen())
			assert.True(t, transpMsg(txn, 0).isAck)
			assert.Nil(t, endpoint.err)

			txn.Consume(&mockMsg{code: 600, isResp: true})
			assert.Equal(t, Completed, txn.state.Load())
			assert.Equal(t, 2, transp.msgLen())
			assert.Nil(t, endpoint.err)
		})

		t.Run("send ACK transport error to TU", func(t *testing.T) {
			endpoint, transp, msg, _ := createMock()
			transp.senderr = errors.New("failed to send")

			txn := createClientInvTxn(transp, endpoint, msg)
			txn.state.Store(Completed)

			txn.Consume(&mockMsg{code: 500, isResp: true})
			assert.Equal(t, Terminated, txn.state.Load())
			assert.Error(t, endpoint.err, "failed to send")
		})

		t.Run("terminates on timer D fired", func(t *testing.T) {
			endpoint, transp, msg, _ := createMock()
			txn := createClientInvTxn(transp, endpoint, msg)
			txn.timer.D = 1 * time.Millisecond

			txn.state.Store(Proceeding)
			txn.Consume(&mockMsg{code: 404, isResp: true})
			txn.state.Store(Completed)
			<-time.After(3 * time.Millisecond)
			assert.Equal(t, Terminated, txn.state.Load())
		})
	})
}
