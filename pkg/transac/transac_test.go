package txnlayer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxnLayerClient(t *testing.T) {
	tests := map[string]struct {
		method   string
		wantType any
	}{
		`invite txn`:     {"INVITE", &TxnClientInvite{}},
		`non-invite txn`: {"PUBLISH", &TxnClientNonInvite{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			endpoint, transp, msg, addr := createMock()
			msg.method = tc.method

			txl := New(endpoint)
			assert.Zero(t, len(txl.pool))
			assert.Zero(t, transp.msgLen())

			txl.Client(msg, transp, addr)
			assert.Equal(t, 1, len(txl.pool))
			assert.True(t, transp.msgLen() > 0)
			assert.IsType(t, tc.wantType, txl.pool[msg.TopViaBranch()])
		})
	}
}

func TestTxnLayerConsume(t *testing.T) {
	endpoint, transp, msg, addr := createMock()
	transp.isReliable = true

	txl := New(endpoint)
	txl.Client(msg, transp, addr)
	assert.Equal(t, 1, len(txl.pool))
	assert.Equal(t, 1, transp.msgLen())

	resp := &mockMsg{code: 200, branch: msg.TopViaBranch()}
	txl.Consume(resp, transp, addr)
	txl.Destroy(endpoint.destroyID)

	assert.Zero(t, len(txl.pool))
	assert.Equal(t, 1, transp.msgLen())
	assert.Equal(t, 1, endpoint.msgLen())
}
