package transaction

import (
	"gosip/pkg/sipmsg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxnLayerClient(t *testing.T) {
	tests := map[string]struct {
		method   string
		wantType any
	}{
		`invite txn`:     {"INVITE", &ClientInvite{}},
		`non-invite txn`: {"PUBLISH", &ClientNonInvite{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			endpoint, transp, msg, addr := createMock(stubInvite)
			msg.Method = tc.method

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
	endpoint, transp, msg, addr := createMock(stubInvite)
	transp.isReliable = true

	txl := New(endpoint)
	txl.Client(msg, transp, addr)
	assert.Equal(t, 1, len(txl.pool))
	assert.Equal(t, 1, transp.msgLen())

	resp := mockResponse("200", "OK")
	via := resp.Find(sipmsg.HVia).(*sipmsg.HeaderVia)
	via.Branch = msg.TopViaBranch()
	via.Params = via.Params.Set("branch", via.Branch)

	txl.Consume(resp, transp, addr)
	txl.TxnDestroy(endpoint.destroyID)

	assert.Zero(t, len(txl.pool))
	assert.Equal(t, 1, transp.msgLen())
}
