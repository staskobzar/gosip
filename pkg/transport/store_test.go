package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Parallel()
	t.Run("listeners store", func(t *testing.T) {
		t.Parallel()

		store := NewStore[Listener]()

		udpln := &UDP{}
		tcpln := &TCPListener{}

		store.Put("udp:10.1.1.1:5858", udpln)
		store.Put("tcp:10.1.1.1:5858", tcpln)
		assert.Len(t, store.pool, 2)

		item, ok := store.Get("udp:10.1.1.1:5858")
		assert.True(t, ok)
		assert.Same(t, item, udpln)

		item, ok = store.Get("udp:10.2.2.100:5060")
		assert.False(t, ok)
		assert.Nil(t, item)

		store.Del("udp:10.1.1.1:5858")
		assert.Len(t, store.pool, 1)
		store.Del("tcp:10.1.1.1:5858")
		assert.Len(t, store.pool, 0)

		item, ok = store.Get("udp:10.1.1.1:5858")
		assert.False(t, ok)
		assert.Nil(t, item)
	})

	t.Run("connections store", func(t *testing.T) {
		t.Parallel()

		store := NewStore[Conn]()

		udpconn := &UDP{}
		tcpconn := &TCP{}

		store.Put("udp:10.1.1.1:5858", udpconn)
		store.Put("tcp:10.1.1.1:5858", tcpconn)
		assert.Len(t, store.pool, 2)

		item, ok := store.Get("tcp:10.1.1.1:5858")
		assert.True(t, ok)
		assert.Same(t, item, tcpconn)
	})

	t.Run("get first available listener", func(t *testing.T) {
		store := NewStore[Listener]()

		udpln := &UDP{}

		item, ok := store.Get("")
		assert.False(t, ok)
		assert.Nil(t, item)

		store.Put("udp:10.1.1.1:5858", udpln)
		item, ok = store.Get("")
		assert.True(t, ok)
		assert.Same(t, udpln, item)
	})
}
