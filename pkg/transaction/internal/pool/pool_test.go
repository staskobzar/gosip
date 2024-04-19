package pool

import (
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTxn struct {
	branch string
}

func (*mockTxn) Consume(*sip.Packet)                       {}
func (*mockTxn) Match(*sipmsg.Message) (Transaction, bool) { return nil, false }
func (t *mockTxn) BranchID() string                        { return t.branch }

func TestPool(t *testing.T) {
	p := New()
	assert.Equal(t, 0, p.Len())

	txn := &mockTxn{}

	t.Run("successfully add first txn", func(t *testing.T) {
		txn.branch = "foo"
		err := p.Add(txn)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.Len())
	})

	t.Run("failed to add new when key exists", func(t *testing.T) {
		err := p.Add(&mockTxn{branch: "foo"})
		assert.ErrorContains(t, err, "already exists")
		assert.Equal(t, 1, p.Len())
	})

	t.Run("successfully add new", func(t *testing.T) {
		err := p.Add(&mockTxn{branch: "bar"})
		assert.Nil(t, err)
		assert.Equal(t, 2, p.Len())
	})

	t.Run("successfully get txn", func(t *testing.T) {
		val, exists := p.Get("foo")
		assert.True(t, exists)
		assert.Same(t, txn, val)
	})

	t.Run("get not existsing txn", func(t *testing.T) {
		val, exists := p.Get("xyz")
		assert.False(t, exists)
		assert.Nil(t, val)
	})

	t.Run("delete txn", func(t *testing.T) {
		txn := &mockTxn{branch: "bar"}
		assert.Equal(t, 2, p.Len())
		p.Delete(txn)
		assert.Equal(t, 1, p.Len())
	})

	t.Run("delete txn that does not exist", func(t *testing.T) {
		txn := &mockTxn{branch: "xyz"}
		assert.Equal(t, 1, p.Len())
		p.Delete(txn)
		assert.Equal(t, 1, p.Len())
	})
}
