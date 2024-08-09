package transport

import (
	"crypto/tls"
	"net"
)

// TLS transport connection.
type TLS struct {
	*TCP
	conf *tls.Config
}

func (tlsconn *TLS) key() string {
	localAddr := &TLSAddr{addr: tlsconn.conn.LocalAddr()}
	remoteAddr := &TLSAddr{addr: tlsconn.conn.RemoteAddr()}

	return connName(localAddr, remoteAddr)
}

// TLSAddr represents TLS address. This is
// a simple wrapper to *net.TCPAddr.
type TLSAddr struct {
	addr net.Addr
}

// Network returns tls for *TLSAddr.
func (addr *TLSAddr) Network() string {
	return "tls"
}

// String returns formatted address:port string.
func (addr *TLSAddr) String() string {
	return addr.addr.String()
}
