package transport

import (
	"context"
	"gosip/pkg/logger"
	"net"
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

		rcvPacket(rcv, buf[:n], tcp.conn.LocalAddr(), tcp.conn.RemoteAddr())
	}
}

func (tcp *TCP) key() string {
	return connName(tcp.conn.LocalAddr(), tcp.conn.RemoteAddr())
}
