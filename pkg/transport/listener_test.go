package transport

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestManagerListenProto(t *testing.T) {
	t.Parallel()

	mgr := Init()

	assert.ErrorContains(t,
		mgr.ListenTCP(context.Background(), "foobar"),
		"failed to resolve TCP address")

	assert.ErrorContains(t,
		mgr.ListenUDP(context.Background(), "foobar"),
		"failed to resolve UDP address")

	assert.ErrorContains(t,
		mgr.ListenTCP(context.Background(), "0.0.0.0:5060"),
		"can not be empty or unspecified")

	assert.ErrorContains(t,
		mgr.ListenUDP(context.Background(), "0.0.0.0:5060"),
		"can not be empty or unspecified")
}

//nolint:paralleltest
func TestManagerListenTransportAndReceive(t *testing.T) {
	send := func(address string) {
		network, addr, _ := strings.Cut(address, ":")

		conn, err := net.Dial(network, addr)
		if err != nil {
			panic(err)
		}

		_, err = fmt.Fprint(conn, strings.Repeat("x", 100))
		if err != nil {
			panic(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("ListenTCP success", func(t *testing.T) {
		mgr := Init()
		err := mgr.ListenTCP(ctx, "127.0.0.1:0")
		assert.Nil(t, err)

		<-time.After(time.Millisecond) // let store listen address in goroutine

		addr, ok := mgr.sock.Get("")
		assert.True(t, ok)

		go send(addr.key())

		assert.Equal(t, bytes.Repeat([]byte{0x78}, 100), (<-mgr.Recv()).Payload)
	})

	t.Run("ListenUDP success", func(t *testing.T) {
		mgr := Init()
		err := mgr.ListenUDP(ctx, "127.0.0.1:0")
		assert.Nil(t, err)

		<-time.After(time.Millisecond) // let store listen address in goroutine

		addr, ok := mgr.sock.Get("")
		assert.True(t, ok)

		go send(addr.key())

		assert.Equal(t, bytes.Repeat([]byte{0x78}, 100), (<-mgr.Recv()).Payload)
	})
}

func TestListenerListen(t *testing.T) {
	t.Parallel()

	success := func(ln Listener) {
		t.Helper()
		assert.Nil(t, ln.listen(context.Background()))
	}

	failOnAddr := func(ln Listener) {
		t.Helper()

		err := ln.listen(context.Background())
		assert.ErrorContains(t, err, "failed start")
	}

	failOnCtx := func(ln Listener) {
		t.Helper()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := ln.listen(ctx)
		assert.ErrorContains(t, err, "context canceled")
	}

	success(&UDP{laddr: &net.UDPAddr{IP: net.IP{0, 0, 0, 0}, Port: 0}})
	success(&TCPListener{laddr: &net.TCPAddr{IP: net.IP{0, 0, 0, 0}, Port: 0}})

	failOnAddr(&UDP{laddr: &net.UDPAddr{IP: net.IP{0, 0, 0, 0}, Port: -5}})
	failOnAddr(&TCPListener{laddr: &net.TCPAddr{IP: net.IP{0, 0, 0, 0}, Port: -5}})

	failOnCtx(&UDP{laddr: &net.UDPAddr{IP: net.IP{0, 0, 0, 0}, Port: 0}})
	failOnCtx(&TCPListener{laddr: &net.TCPAddr{IP: net.IP{0, 0, 0, 0}, Port: 0}})
}

func TestManagerListenersCancelOnContext(t *testing.T) {
	t.Parallel()

	mgr := Init()
	ctx, cancel := context.WithCancel(context.Background())

	assert.Nil(t, mgr.ListenTCP(ctx, "127.0.0.1:0"))
	assert.Nil(t, mgr.ListenTCP(ctx, "127.0.0.1:0"))

	assert.Nil(t, mgr.ListenUDP(ctx, "127.0.0.1:0"))
	assert.Nil(t, mgr.ListenUDP(ctx, "127.0.0.1:0"))

	<-time.After(time.Millisecond)
	assert.Equal(t, 4, mgr.sock.Len())

	cancel()

	for range len(mgr.sock.pool) {
		ln, found := mgr.sock.Get("")
		assert.True(t, found)
		ln.close()
		<-time.After(time.Millisecond)
	}

	assert.Equal(t, 0, mgr.sock.Len())
}
