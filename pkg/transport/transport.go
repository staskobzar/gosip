package transport

import (
	"context"
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"net"
	"time"
)

// module errors
var (
	Error     = errors.New("transport")
	ErrListen = fmt.Errorf("%w: listner", Error)
	ErrSend   = fmt.Errorf("%w: send", Error)
	ErrResolv = fmt.Errorf("%w: dns resovl", Error)
)

const reconnTout = time.Second * 3

type tTransp uint8

const (
	tUnknown tTransp = 0
	// https://github.com/ishidawataru/sctp
	tSCTP tTransp = 1 << iota
	tTCP
	tTLS
	tUDP
)

func (transp tTransp) String() string {
	switch transp {
	case tSCTP:
		return "sctp"
	case tTCP:
		return "tcp"
	case tTLS:
		return "tls"
	case tUDP:
		return "udp"
	default:
		return "unknown"
	}
}

type Manager struct {
	sock    *Store[Listener]
	conn    *Store[Conn]
	rcv     chan Packet
	resolv  chan sip.Packet
	err     chan error
	support tTransp
	dns     sip.DNS
}

type Listener interface {
	listen(ctx context.Context) error
	accept(ctx context.Context) (<-chan Conn, <-chan error)
	key() string
	close()
}

type Conn interface {
	consume(ctx context.Context, rcv chan<- Packet, store *Store[Conn])
}

type Packet struct {
	Payload []byte
	Laddr   net.Addr
	Raddr   net.Addr
}

func dbg(pattern string, args ...any) {
	logger.Log("transp: "+pattern, args...)
}
func wrn(pattern string, args ...any) {
	logger.Wrn("transp: "+pattern, args...)
}
func perr(pattern string, args ...any) {
	logger.Err("transp: "+pattern, args...)
}
func dbgr(pattern string, args ...any) {
	logger.Log("transp:dns: "+pattern, args...)
}
func wrnr(pattern string, args ...any) {
	logger.Wrn("transp:dns: "+pattern, args...)
}

func Init() *Manager {
	dbg("init manager")
	return &Manager{
		sock:   NewStore[Listener](),
		conn:   NewStore[Conn](),
		rcv:    make(chan Packet, 32),
		resolv: make(chan sip.Packet, 32),
		err:    make(chan error, 32),
	}
}

func (mgr *Manager) WithResolver(dns sip.DNS) {
	mgr.dns = dns
}

func (mgr *Manager) ListenTCP(ctx context.Context, addrport string) error {
	dbg("start TCP listener on %q", addrport)
	addr, err := net.ResolveTCPAddr("tcp", addrport)
	if err != nil {
		return fmt.Errorf("%w: failed to resolve TCP address %q: %s",
			ErrListen, addrport, err)
	}

	ln := &TCPListener{laddr: addr}
	go mgr.listen(ctx, ln)

	mgr.support |= tTCP
	return nil
}

func (mgr *Manager) ListenUDP(ctx context.Context, addrport string) error {
	dbg("start UDP listener on %q", addrport)
	addr, err := net.ResolveUDPAddr("udp", addrport)
	if err != nil {
		return fmt.Errorf("%w: failed to resolve UDP address %q: %s",
			ErrListen, addrport, err)
	}
	ln := &UDP{laddr: addr}
	go mgr.listen(ctx, ln)

	mgr.support |= tUDP
	return nil
}

func (mgr *Manager) Send(src, dst net.Addr, msg *sipmsg.Message) error {
	switch src.Network() {
	case "udp":
		return mgr.SendUDP(src, dst, msg)
	case "tcp":
		return mgr.SendTCP(src, dst, msg)
	default:
		return fmt.Errorf("%w: invalid or unsupported network %q", ErrSend, src.Network())
	}
}

func (mgr *Manager) SendUDP(src, dst net.Addr, msg *sipmsg.Message) error {
	name := sockName(src)
	ln, found := mgr.sock.Get(name)
	if !found {
		return fmt.Errorf("%w: connection %q was not found", ErrSend, name)
	}

	udp, ok := ln.(*UDP)
	if !ok {
		return fmt.Errorf("%w: found %q socket but type is not UDPConn",
			ErrSend, name)
	}

	n, err := udp.conn.WriteTo(msg.Byte(), dst)
	if err != nil {
		return fmt.Errorf("%w: failed to write to udp:%s from %q: %s",
			ErrSend, dst, name, err)
	}
	dbg("successfully sent to udp %s from %q %d bytes", dst, src, n)
	return nil
}

func (mgr *Manager) SendTCP(src, dst net.Addr, msg *sipmsg.Message) error {
	name := connName(src, dst)
	cn, found := mgr.conn.Get(name)
	if !found {
		return fmt.Errorf("%w: connection %q not found", ErrSend, name)
	}

	tcp, ok := cn.(*TCP)
	if !ok {
		return fmt.Errorf("%w: found %q socket but type is not UDPConn", ErrSend, name)
	}

	n, err := tcp.conn.Write(msg.Byte())
	if err != nil {
		return fmt.Errorf("%w: failed to write to conn %q: %s", ErrSend, name, err)
	}
	dbg("sent %d bytes to %q", n, name)

	return nil
}

func (mgr *Manager) Err() <-chan error {
	return mgr.err
}

func (mgr *Manager) Recv() <-chan Packet {
	return mgr.rcv
}

func (mgr *Manager) Resolved() <-chan sip.Packet {
	return mgr.resolv
}

func (mgr *Manager) listen(ctx context.Context, ln Listener) {
	for {
		if ctx.Err() != nil {
			wrn("stop listener on context: %s", ctx.Err())
			return
		}

		if err := ln.listen(ctx); err != nil {
			perr("failed to start listener: %s", err)
			wrn("restart in %v", reconnTout)
			<-time.After(reconnTout)
			continue
		}
		mgr.sock.Put(ln.key(), ln)

		connCh, chErr := ln.accept(ctx)
	connLoop:
		for {
			select {
			case conn := <-connCh:
				go conn.consume(ctx, mgr.rcv, mgr.conn)
			case err := <-chErr:
				perr("accept connection err: %s", err)
				break connLoop
			}
		}

		wrn("close and restart listener")
		mgr.sock.Del(ln.key())
		ln.close()
	}
}

func rcvPacket(rcv chan<- Packet, buf []byte, laddr, raddr net.Addr) {
	size := len(buf)
	payload := make([]byte, size)
	copy(payload, buf)

	pack := Packet{
		Payload: payload,
		Laddr:   laddr,
		Raddr:   raddr,
	}

	select {
	case rcv <- pack:
		dbg("sent pack with payload of %d bytes sock %q remote address %q", size, laddr, raddr)
	case <-time.After(100 * time.Millisecond):
		perr("failed to send pack on blocked chan for sock %q remote address %q", laddr, raddr)
	}
}

func sockName(addr net.Addr) string {
	return addr.Network() + ":" + addr.String()
}

func connName(laddr, raddr net.Addr) string {
	lsock := sockName(laddr)
	rsock := sockName(raddr)

	return lsock + "<@>" + rsock

}
