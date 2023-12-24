package transport

import (
	"context"
	"gosip/pkg/logger"
	"net"
	"time"
)

type TCP struct {
	conn *net.TCPConn
}

func (tcp *TCP) consume(ctx context.Context, rcv chan<- Packet, store *Store[Conn]) {
	buf := make([]byte, 1<<12)
	connName := tcp.key()

	store.Put(connName, tcp)
	defer store.Del(connName)

	for {
		if ctx.Err() != nil {
			logger.Wrn("connection %q is terminated by contex", connName)
			break
		}

		n, err := tcp.conn.Read(buf)
		if err != nil {
			logger.Err("failed read conn %q: %s", connName, err)
			break
		}
		payload := make([]byte, n)
		copy(payload, buf[:n])

		pack := Packet{
			Payload: payload,
			Laddr:   tcp.conn.LocalAddr(),
			Raddr:   tcp.conn.RemoteAddr(),
		}

		rcvPacket(rcv, pack, connName)
	}
}

func rcvPacket(rcv chan<- Packet, pack Packet, connName string) {
	select {
	case rcv <- pack:
		logger.Log("sent pack with payload of %d bytes on %q", len(pack.Payload), connName)
	case <-time.After(100 * time.Millisecond):
		logger.Err("failed to send pack on blocked chan for %q", connName)
	}
}

func (tcp *TCP) key() string {
	return connName(tcp.conn.LocalAddr(), tcp.conn.RemoteAddr())
}
