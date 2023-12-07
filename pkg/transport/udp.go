package transport

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"net"
)

type UDP struct {
	laddr *net.UDPAddr
	conn  *net.UDPConn
}

func (udp *UDP) Write(raddr net.Addr, payload []byte) error {
	n, err := udp.conn.WriteTo(payload, raddr)
	if err != nil {
		return err
	}
	logger.Log("sent %d bytes to %s", n, raddr)
	return nil
}

func (udp *UDP) listen() error {
	conn, err := net.ListenUDP("udp", udp.laddr)
	if err != nil {
		return fmt.Errorf("%w: failed to start on %q: %s",
			ErrUDPListener, udp.laddr, err)
	}

	logger.Log("starting UDP listener on %q", udp.laddr)

	udp.conn = conn
	return nil
}

func (udp *UDP) recv(ctx context.Context, rcv chan<- Packet) error {
	buf := make([]byte, 1<<16)

	for {
		n, raddr, err := udp.conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("%w: failed read from %q: %s",
				ErrUDPListener, udp.conn.LocalAddr(), err)
		}
		payload := make([]byte, n)
		copy(payload, buf[:n])
		logger.Log("udp read %d bytes from %q", n, raddr)
		pack := Packet{
			Payload: buf[:n],
			Laddr:   udp.conn.LocalAddr(),
			Raddr:   raddr,
		}
		consume(rcv, pack)
	}
}

func (udp *UDP) close() {
	udp.conn.Close()
}

func (udp *UDP) id() sID {
	return sIDBuild(udp.laddr)
}
