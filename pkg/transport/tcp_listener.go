package transport

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"net"
)

type TCPListener struct {
	laddr *net.TCPAddr
	ln    *net.TCPListener
}

func (tcpln *TCPListener) listen() error {
	ln, err := net.ListenTCP("tcp", tcpln.laddr)
	if err != nil {
		return fmt.Errorf("%w: failed to start on %q: %s",
			ErrTCPListener, tcpln.laddr, err)
	}

	logger.Log("starting TCP listener on %q", tcpln.laddr)
	tcpln.ln = ln

	return nil
}

func (tcpln *TCPListener) recv(ctx context.Context, rcv chan<- Packet) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("listener %q recv interrupted by context",
				tcpln.laddr)
		default:
			conn, err := tcpln.ln.AcceptTCP()
			if err != nil {
				return fmt.Errorf("accept on %q: %s", tcpln.laddr, err)
			}
			fmt.Printf("%#v\n", conn)
		}
	}
}

func (tcpln *TCPListener) close() {
	tcpln.ln.Close()
}

func (tcpln *TCPListener) id() sID {
	return sIDBuild(tcpln.laddr)
}
