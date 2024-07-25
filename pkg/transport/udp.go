package transport

import (
	"context"
	"fmt"
	"net"

	"gosip/pkg/logger"
)

// UDP transport structure.
type UDP struct {
	laddr *net.UDPAddr
	conn  *net.UDPConn
	chErr chan error
}

func (udp *UDP) listen(ctx context.Context) error {
	conn, err := net.ListenUDP("udp", udp.laddr)
	if err != nil {
		return fmt.Errorf("failed start udp listener: %w", err)
	}

	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		udp.laddr = addr
	}

	logger.Log("start UDP listener on %q", udp.laddr)

	udp.conn = conn

	return ctx.Err()
}

// this is just a dummy method as UDP does not use accept
func (udp *UDP) accept(_ context.Context) (<-chan Conn, <-chan error) {
	chConn := make(chan Conn)
	udp.chErr = make(chan error)

	go func() {
		chConn <- udp
	}()

	return chConn, udp.chErr
}

func (udp *UDP) consume(ctx context.Context, rcv chan<- Packet, _ *Store[Conn]) {
	buf := make([]byte, 1<<12)
	name := udp.key()

	for {
		if ctx.Err() != nil {
			udp.chErr <- fmt.Errorf("%w: connection %q is terminated by context",
				ErrListen, name)

			break
		}

		n, raddr, err := udp.conn.ReadFrom(buf)
		if err != nil {
			udp.chErr <- fmt.Errorf("failed to read conn %q from %q: %w",
				name, raddr, err)

			break
		}

		rcvPacket(rcv, buf[:n], udp.laddr, raddr)
	}
}

func (udp *UDP) key() string {
	return sockName(udp.laddr)
}

func (udp *UDP) close() {
	logger.Wrn("closing listener %q", udp.laddr)

	if err := udp.conn.Close(); err != nil {
		logger.Err("failed to close TCP listener %q: %s", udp.laddr, err)
	}
}
