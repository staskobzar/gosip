package pool

import (
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTxn struct {
	branch string
	match  bool
}

func (*mockTxn) Consume(*sip.Packet) {}
func (t *mockTxn) Match(*sipmsg.Message) (sip.Transaction, bool) {
	if t.match {
		return t, true
	}
	return nil, false
}
func (t *mockTxn) BranchID() string { return t.branch }
func (t *mockTxn) Terminate()       {}

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
		assert.Equal(t, 2, p.Len())
		p.Delete("bar")
		assert.Equal(t, 1, p.Len())
	})

	t.Run("delete txn that does not exist", func(t *testing.T) {
		assert.Equal(t, 1, p.Len())
		p.Delete("xyz")
		assert.Equal(t, 1, p.Len())
	})
}

func TestPoolMatch(t *testing.T) {
	inputReq := "INVITE sip:bob@biloxi.example.com SIP/2.0\r\n" +
		"Via: SIP/2.0/TCP atlanta.com;branch=z9hG4bK74b43\r\n" +
		"From: Alice <sip:alice@atlanta.example.com>;tag=9fxced76sl\r\n" +
		"To: Bob <sip:bob@biloxi.example.com>\r\n" +
		"Call-ID: 384827@e.com\r\n" +
		"CSeq: 1 INVITE\r\n\r\n"
	msg, _ := sipmsg.Parse(inputReq)

	t.Run("false when txn no found", func(t *testing.T) {
		p := New()
		assert.Equal(t, 0, p.Len())
		assert.Equal(t, 0, p.Len())
		transac, ok := p.Match(msg)
		assert.False(t, ok)
		assert.Nil(t, transac)
	})

	t.Run("true when match txn", func(t *testing.T) {
		p := New()
		assert.Equal(t, 0, p.Len())
		txn := &mockTxn{branch: "z9hG4bK74b43", match: true}
		err := p.Add(txn)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.Len())
		transac, ok := p.Match(msg)
		assert.True(t, ok)
		assert.Same(t, txn, transac)
	})

	t.Run("false when txn not match", func(t *testing.T) {
		p := New()
		assert.Equal(t, 0, p.Len())
		txn := &mockTxn{branch: "foo", match: false}
		err := p.Add(txn)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.Len())
		transac, ok := p.Match(msg)
		assert.False(t, ok)
		assert.Nil(t, transac)
	})
}
