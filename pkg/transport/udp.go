package transport

import (
	"context"
	"gosip/pkg/logger"
	"net"
)

type UDP struct {
	laddr *net.UDPAddr
	conn  *net.UDPConn
}

func (udp *UDP) listen(ctx context.Context) error {
	conn, err := net.ListenUDP("udp", udp.laddr)
	if err != nil {
		return err
	}

	logger.Log("start UDP listener on %q", udp.laddr)

	udp.conn = conn
	return ctx.Err()
}

func (udp *UDP) accept(_ context.Context) (<-chan Conn, <-chan error) {
	chConn, chErr := make(chan Conn), make(chan error)
	go func() {
		chConn <- udp
	}()
	return chConn, chErr
}

func (udp *UDP) consume(ctx context.Context, rcv chan<- Packet, _ *Store[Conn]) {
	buf := make([]byte, 1<<12)
	name := udp.key()
	for {
		if ctx.Err() != nil {
			logger.Wrn("connection %q is terminated by context", name)
			break
		}

		n, raddr, err := udp.conn.ReadFrom(buf)
		if err != nil {
			logger.Err("failed to read conn %q from %q: %s", name, raddr, err)
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
