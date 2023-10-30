package txnlayer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTxnClientNonInviteInit(t *testing.T) {
	t.Run("with unreliable transport", func(t *testing.T) {
		endpoint, transp, msg, addr := createMock()
		txn := createClientNonInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond

		assert.Equal(t, Unknown, txn.state.Load())
		assert.Equal(t, 0, transp.msgLen())

		txn.Init(msg, addr)

		assert.Equal(t, Trying, txn.state.Load())
		assert.Equal(t, 1, transp.msgLen())
		<-time.After(5 * time.Millisecond)
		assert.True(t, transp.msgLen() > 1)
	})

	t.Run("with reliable transport", func(t *testing.T) {
		endpoint, transp, msg, addr := createMock()
		txn := createClientNonInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond
		transp.isReliable = true

		txn.Init(msg, addr)

		assert.Equal(t, 1, transp.msgLen())
		<-time.After(5 * time.Millisecond)
		assert.Equal(t, 1, transp.msgLen())
	})

	t.Run("ignore if Message is not response", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()

		txn := createClientNonInvTxn(transp, endpoint, msg)
		txn.Consume(&mockMsg{})

		assert.Equal(t, Unknown, txn.state.Load())
	})

	t.Run("terminates on trying state on timeout timer F", func(t *testing.T) {
		endpoint, transp, msg, addr := createMock()
		txn := createClientNonInvTxn(transp, endpoint, msg)
		txn.timer.T1 = 1 * time.Millisecond
		transp.isReliable = true

		txn.Init(msg, addr)
		assert.Equal(t, Trying, txn.state.Load())
		assert.False(t, endpoint.isTout())
		<-time.After(80 * time.Millisecond)
		assert.Equal(t, Terminated, txn.state.Load())
		assert.True(t, endpoint.isTout())
	})
}

func TestTxnClientNonInviteConsume(t *testing.T) {
	t.Run("ignore if Message is not response", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()

		txn := createClientInvTxn(transp, endpoint, msg)
		txn.Consume(&mockMsg{})

		assert.Equal(t, Unknown, txn.state.Load())
	})

	t.Run("on trying or proceeding state", func(t *testing.T) {
		tests := map[string]struct {
			onState   uint32
			respCode  int
			wantState uint32
		}{
			`trying response 100 to proceeding`:     {Trying, 100, Proceeding},
			`trying response 200 to completed`:      {Trying, 200, Completed},
			`trying response 303 to completed`:      {Trying, 303, Completed},
			`trying response 404 to completed`:      {Trying, 404, Completed},
			`trying response 550 to completed`:      {Trying, 550, Completed},
			`proceeding response 186 to proceeding`: {Proceeding, 186, Proceeding},
			`proceeding response 202 to completed`:  {Proceeding, 202, Completed},
			`proceeding response 303 to completed`:  {Proceeding, 303, Completed},
			`proceeding response 404 to completed`:  {Proceeding, 404, Completed},
			`proceeding response 550 to completed`:  {Proceeding, 550, Completed},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				endpoint, transp, msg, _ := createMock()
				txn := createClientNonInvTxn(transp, endpoint, msg)
				txn.state.Store(tc.onState)
				txn.Consume(&mockMsg{code: tc.respCode})
				assert.Equal(t, tc.wantState, txn.state.Load())
			})
		}
	})

	t.Run("on trying transport error sends to TU and terminates transaction", func(t *testing.T) {
		endpoint, transp, msg, addr := createMock()
		transp.senderr = fmt.Errorf("transport failed to send")
		txn := createClientNonInvTxn(transp, endpoint, msg)

		assert.Nil(t, endpoint.err)
		txn.Init(msg, addr)
		assert.NotNil(t, endpoint.err)
		assert.ErrorContains(t, endpoint.err, "transport failed")
		assert.Equal(t, Terminated, txn.state.Load())
	})
}

func TestTxnClientNonInviteComplete(t *testing.T) {
	t.Run("with unreliable transport starts timer K", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()
		txn := createClientNonInvTxn(transp, endpoint, msg)
		txn.timer.T4 = 2 * time.Millisecond
		txn.state.Store(Proceeding)
		txn.Consume(&mockMsg{code: 200})
		assert.Equal(t, Completed, txn.state.Load())
		<-time.After(5 * time.Millisecond)
		assert.Equal(t, Terminated, txn.state.Load())
	})

	t.Run("with reliable transport switch to terminated", func(t *testing.T) {
		endpoint, transp, msg, _ := createMock()
		txn := createClientNonInvTxn(transp, endpoint, msg)
		transp.isReliable = true
		txn.state.Store(Proceeding)
		txn.Consume(&mockMsg{code: 200})
		assert.Equal(t, Terminated, txn.state.Load())
	})
}

func TestTickTimerE(t *testing.T) {
	tests := map[string]struct {
		state  []uint32
		t1, t2 time.Duration
		want   []time.Duration
	}{
		`for state trying only`: {
			state: []uint32{Trying, Trying, Trying, Trying, Trying, Trying},
			t1:    500 * time.Millisecond,
			t2:    4 * time.Second,
			want: []time.Duration{
				500 * time.Millisecond,
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
		`with short step`: {
			state: []uint32{Trying, Trying, Trying, Trying},
			t1:    1 * time.Second,
			t2:    4 * time.Second,
			want: []time.Duration{
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
		`next state is proceed`: {
			state: []uint32{Trying, Proceeding, Proceeding, Proceeding},
			t1:    500 * time.Millisecond,
			t2:    4 * time.Second,
			want: []time.Duration{
				500 * time.Millisecond,
				4 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tick := tickTimerE(tc.t1, tc.t2)
			for i := 0; i < len(tc.want); i++ {
				assert.Equal(t, tc.want[i], tick(tc.state[i]))
			}
		})
	}
}
