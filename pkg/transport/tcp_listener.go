package transport

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"net"
	"time"
)

type TCPListener struct {
	laddr *net.TCPAddr
	ln    *net.TCPListener
}

func (tcpln *TCPListener) listen(ctx context.Context) error {
	ln, err := net.ListenTCP("tcp", tcpln.laddr)
	if err != nil {
		return err
	}

	logger.Log("start TCP listener on %q", tcpln.laddr)

	tcpln.ln = ln
	return ctx.Err()
}

func (tcpln *TCPListener) accept(ctx context.Context) (<-chan Conn, <-chan error) {
	connCh := make(chan Conn, 32)
	errCh := make(chan error)

	go func() {
		for {
			tcpconn, err := tcpln.ln.AcceptTCP()
			if err != nil {
				errCh <- fmt.Errorf("failed to accept connection for %q: %s", tcpln.laddr, err)
				break
			}

			if ctx.Err() != nil {
				logger.Wrn("accept routine %q is terminated by context: %s",
					tcpln.laddr, ctx.Err())
				break
			}

			select {
			case connCh <- &TCP{conn: tcpconn}:
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
