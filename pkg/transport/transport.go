// Package transport provides management for different
// transport protocol connectoins and listenters
package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
)

// module errors.
var (
	Error     = errors.New("transport") //nolint:errname
	ErrListen = fmt.Errorf("%w: listner", Error)
	ErrSend   = fmt.Errorf("%w: send", Error)
	ErrResolv = fmt.Errorf("%w: dns resovl", Error)
)

// ErrChan error chan structure.
type ErrChan struct {
	Err  error
	Pack *sip.Packet
}

type tTransp uint8

const reconnTout = time.Second * 3

// transport types.
const (
	tUnknown tTransp = 0
	// https://github.com/ishidawataru/sctp
	tSCTP tTransp = 1 << iota
	tTCP
	tTLS
	tUDP
)

// String stringifier method.
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

// Manager transport core object.
type Manager struct {
	sock    *Store[Listener]
	conn    *Store[Conn]
	rcv     chan Packet
	resolv  chan sip.Packet
	err     chan ErrChan
	support tTransp
	dns     sip.DNS
}

// Listener interface for different transport protocols (udp, tcp etc).
type Listener interface {
	listen(ctx context.Context) error
	accept(ctx context.Context) (<-chan Conn, <-chan error)
	key() string
	close()
}

// Conn connection interface for a transport protocol.
type Conn interface {
	consume(ctx context.Context, rcv chan<- Packet, store *Store[Conn])
}

// Packet represents transport packet with addresses and payload.
type Packet struct {
	Payload []byte
	Laddr   net.Addr
	Raddr   net.Addr
}

// Init transport Manager.
func Init() *Manager {
	dbg("transp: init manager")

	return &Manager{
		sock:   NewStore[Listener](),
		conn:   NewStore[Conn](),
		rcv:    make(chan Packet),     // 32),
		resolv: make(chan sip.Packet), // 32),
		err:    make(chan ErrChan),    // 32),
	}
}

// ListenTCP starts TCP listener and push it into the manager's store.
func (mgr *Manager) ListenTCP(ctx context.Context, addrport string) error {
	dbg("start TCP listener on %q", addrport)

	addr, err := net.ResolveTCPAddr("tcp", addrport)
	if err != nil {
		return fmt.Errorf("%w: failed to resolve TCP address %q: %w",
			ErrListen, addrport, err)
	}

	if err := mgr.validIP(addr.IP, tTCP); err != nil {
		return err
	}

	ln := &TCPListener{laddr: addr}
	go mgr.listen(ctx, ln)

	return nil
}

// ListenUDP starts UDP listener and push it into the manager's store.
func (mgr *Manager) ListenUDP(ctx context.Context, addrport string) error {
	dbg("start UDP listener on %q", addrport)
	addr, err := net.ResolveUDPAddr("udp", addrport)

	if err != nil {
		return fmt.Errorf("%w: failed to resolve UDP address %q: %w",
			ErrListen, addrport, err)
	}

	if err := mgr.validIP(addr.IP, tTCP); err != nil {
		return err
	}

	ln := &UDP{laddr: addr}
	go mgr.listen(ctx, ln)

	return nil
}

// Send sip Packet to network.
func (mgr *Manager) Send(pack *sip.Packet) {
	if pack.Message == nil {
		perr("invalid SIP Message <nil> when trying to send packet")

		return
	}

	dbg("sending pack with SIP message %q", pack.Message.FirstLine())

	go mgr.send(pack)
}

// SendUDP sends SIP message to UDP transport.
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

	src = udp.laddr // in case if src is a first ln address

	if err := msg.SetViaTransp("UDP", udp.laddr); err != nil {
		return fmt.Errorf("failed set via transport: %w", err)
	}

	dbg("set top Via transport to UDP and sent-by to %q", udp.laddr)

	n, err := udp.conn.WriteTo(msg.Byte(), dst)

	if err != nil {
		return fmt.Errorf("%w: failed to write to udp:%s from %q: %w",
			ErrSend, dst, name, err)
	}

	dbg("successfully sent to udp %s from %q %d bytes", dst, src, n)

	return nil
}

// SendTCP sends SIP message to TCP transport.
func (mgr *Manager) SendTCP(src, dst net.Addr, msg *sipmsg.Message) error {
	name := connName(src, dst)

	conn, found := mgr.conn.Get(name)

	if !found {
		addr, err := net.ResolveTCPAddr("tcp", dst.String())
		if err != nil {
			return fmt.Errorf("%w: failed to resolve destination addr %q: %w",
				ErrSend, dst, err)
		}

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return fmt.Errorf("%w: failed to connect %q: %w", ErrSend, addr, err)
		}

		// tcp:=&TCP{conn: conn}
		// tcp consume ???

		n, err := conn.Write(msg.Byte())
		if err != nil {
			return fmt.Errorf("%w: failed to write to conn %q: %w", ErrSend, addr, err)
		}

		dbg("sent %d bytes to %q", n, addr)

		return nil
	}

	tcp, ok := conn.(*TCP)
	if !ok {
		return fmt.Errorf("%w: found %q socket but type is not UDPConn", ErrSend, name)
	}

	n, err := tcp.conn.Write(msg.Byte())
	if err != nil {
		return fmt.Errorf("%w: failed to write to conn %q: %w", ErrSend, name, err)
	}

	dbg("sent %d bytes to %q", n, name)

	return nil
}

// Err returns channel to deliver transport errors.
func (mgr *Manager) Err() <-chan ErrChan {
	return mgr.err
}

// Recv returns channel that delivers SIP packets received
// by listeners from network.
func (mgr *Manager) Recv() <-chan Packet {
	return mgr.rcv
}

// Resolved returns channel that delivers SIP packets with
// resolved destination in RURI.
func (mgr *Manager) Resolved() <-chan sip.Packet {
	return mgr.resolv
}

func (mgr *Manager) send(pack *sip.Packet) {
	if len(pack.SendTo) == 0 {
		dbg("no send-to addresses in the packet")

		addrs, err := mgr.Resolve(pack.Message.RURI)
		if err != nil {
			perr("failed to resolve: %s", err)

			return
		}

		pack.SendTo = addrs
	}

	send := func() error {
		for _, dst := range pack.SendTo {
			logger.Log("send to destination addr %q", dst)

			switch dst.Network() {
			case "udp":
				return mgr.SendUDP(pack.LocalSock, dst, pack.Message)
			case "tcp":
				return mgr.SendTCP(pack.LocalSock, dst, pack.Message)
			default:
				return fmt.Errorf("%w: invalid or unsupported network %q", ErrSend, dst.Network())
			}
		}

		return fmt.Errorf("%w: failed to send message %q", ErrSend, pack.Message.FirstLine())
	}

	if err := send(); err != nil {
		mgr.passErr(err, pack)
	}
}

func (mgr *Manager) passErr(err error, pack *sip.Packet) {
	dbg("sending transport error: %s", err)
	select {
	case <-time.After(time.Second):
		perr("timeout to send error: error channel is blocked")
	case mgr.err <- ErrChan{Err: err, Pack: pack}:
	}
}

func (mgr *Manager) passResolv(pack sip.Packet) {
	dbg("sending resolved packet: %v", pack.ReqAddrs)
	select {
	case <-time.After(time.Second):
		perr("timeout to send resolved packet: channel is blocked")
	case mgr.resolv <- pack:
	}
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
	if addr == nil {
		return ""
	}

	return addr.Network() + ":" + addr.String()
}

func connName(laddr, raddr net.Addr) string {
	if laddr == nil || raddr == nil {
		return ""
	}

	lsock := sockName(laddr)
	rsock := sockName(raddr)

	return lsock + "<@>" + rsock
}

func (mgr *Manager) validIP(ip net.IP, trp tTransp) error {
	if ip == nil || ip.IsUnspecified() {
		return fmt.Errorf("%w: invalid %s address: can not be empty or unspecified address",
			ErrListen, trp)
	}

	mgr.support |= tTCP

	return nil
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
