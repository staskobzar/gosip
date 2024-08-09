package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

// TLSListener tls protocol listener.
type TLSListener struct {
	*TCPListener
	conf *tls.Config
}

func newTLSListener(addr *net.TCPAddr, conf *tls.Config) *TLSListener {
	return &TLSListener{
		TCPListener: &TCPListener{laddr: &TLSAddr{addr: addr}},
		conf:        conf,
	}
}

func (tlsln *TLSListener) listen(ctx context.Context) error {
	ln, err := tls.Listen("tcp", tlsln.laddr.String(), tlsln.conf)
	if err != nil {
		return fmt.Errorf("failed start tls listener: %w", err)
	}

	return tlsln.lnSetup(ctx, ln)
}

func (tlsln *TLSListener) accept(ctx context.Context) (<-chan Conn, <-chan error) {
	return tlsln.acceptLoop(ctx, func(conn net.Conn) Conn {
		return &TLS{TCP: &TCP{conn: conn}, conf: tlsln.conf}
	})
}
