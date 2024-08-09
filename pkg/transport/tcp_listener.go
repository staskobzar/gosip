package transport

import (
	"context"
	"fmt"
	"net"
	"time"

	"gosip/pkg/logger"
)

// TCPListener tcp protocol listener.
type TCPListener struct {
	laddr net.Addr
	ln    net.Listener
}

func (tcpln *TCPListener) listen(ctx context.Context) error {
	ln, err := net.Listen("tcp", tcpln.laddr.String())
	if err != nil {
		return fmt.Errorf("failed start tcp listener: %w", err)
	}

	return tcpln.lnSetup(ctx, ln)
}

func (tcpln *TCPListener) lnSetup(ctx context.Context, ln net.Listener) error {
	tcpln.laddr = ln.Addr()

	logger.Log("start TCP listener on %q", tcpln.laddr)

	tcpln.ln = ln

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: listener is already stopped by context: %w", ErrListen, err)
	}

	return nil
}

func (tcpln *TCPListener) accept(ctx context.Context) (<-chan Conn, <-chan error) {
	return tcpln.acceptLoop(ctx, func(conn net.Conn) Conn {
		return &TCP{conn: conn}
	})
}

func (tcpln *TCPListener) acceptLoop(ctx context.Context, toConn func(conn net.Conn) Conn) (<-chan Conn, <-chan error) {
	connCh := make(chan Conn, 32)
	errCh := make(chan error)

	go func() {
		for {
			tcpconn, err := tcpln.ln.Accept()
			if err != nil {
				errCh <- fmt.Errorf("failed to accept connection for %q: %w", tcpln.laddr, err)

				break
			}

			if ctx.Err() != nil {
				errCh <- fmt.Errorf("accept routine %q is terminated by context: %w",
					tcpln.laddr, ctx.Err())

				break
			}

			select {
			case connCh <- toConn(tcpconn):
				logger.Log("new tcp connection accepted from %q", tcpconn.RemoteAddr())
			case <-time.After(time.Millisecond * 100):
				logger.Err("failed to send connection for %q on blocked channel", tcpln.laddr)
			}
		}
		close(connCh)
		close(errCh)
	}()

	return connCh, errCh
}

func (tcpln *TCPListener) key() string {
	return sockName(tcpln.laddr)
}

func (tcpln *TCPListener) close() {
	logger.Wrn("closing listener %q", tcpln.laddr)

	if err := tcpln.ln.Close(); err != nil {
		logger.Err("failed to close TCP listener %q: %s", tcpln.laddr, err)
	}
}
