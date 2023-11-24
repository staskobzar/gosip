package transport

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"net"
)

type UDP struct {
	laddr *net.UDPAddr
	raddr *net.UDPAddr
	conn  *net.UDPConn
}

func (udp *UDP) listen() error {
	conn, err := net.ListenUDP("udp", udp.laddr)
	if err != nil {
		return fmt.Errorf("%w: failed to start on %q: %s",
			ErrUDPListener, udp.laddr, err)
	}

	udp.conn = conn
	return nil
}

func (udp *UDP) recv(ctx context.Context) (net.Addr, net.Addr, []byte, error) {
	buf := make([]byte, 1<<16)
	n, raddr, err := udp.conn.ReadFrom(buf)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: failed read from %q: %s",
			ErrUDPListener, udp.conn.LocalAddr(), err)
	}
	payload := make([]byte, n)
	copy(payload, buf[:n])
	logger.Log("udp read %d bytes from %q", n, raddr)

	return raddr, udp.laddr, payload, nil
}

func (udp *UDP) close() {
	udp.conn.Close()
}
