package transport

import (
	"context"
	"fmt"
	"net"

	"gosip/pkg/logger"
)

// TCP transport connection.
type TCP struct {
	conn net.Conn
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

func (tcp *TCP) write(name string, msg []byte) error {
	n, err := tcp.conn.Write(msg)
	if err != nil {
		return fmt.Errorf("%w: failed to write to conn %q: %w",
			ErrSend, name, err)
	}

	dbg("sent %d bytes to %q", n, name)

	return nil
}
