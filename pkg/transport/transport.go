package transport

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"net"
	"time"
)

const reconnTout = time.Second * 3

type tTransp uint8

const (
	tNone tTransp = 0
	tUDP  tTransp = 1 << iota
	tTCP
)

type Manager struct {
	sock    *Store[Listener]
	conn    *Store[Conn]
	rcv     chan Packet
	support tTransp
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

func Init() *Manager {
	return &Manager{
		sock: NewStore[Listener](),
		conn: NewStore[Conn](),
		rcv:  make(chan Packet, 32),
	}
}

func (mgr *Manager) ListenTCP(ctx context.Context, addrport string) error {
	addr, err := net.ResolveTCPAddr("tcp", addrport)
	if err != nil {
		return fmt.Errorf("failed to resolve TCP address %q: %s", addrport, err)
	}

	ln := &TCPListener{laddr: addr}
	go mgr.listen(ctx, ln)

	mgr.support |= tTCP
	return nil
}

func (mgr *Manager) ListenUDP(ctx context.Context, addrport string) error {
	addr, err := net.ResolveUDPAddr("udp", addrport)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address %q: %s", addrport, err)
	}
	ln := &UDP{laddr: addr}
	go mgr.listen(ctx, ln)

	mgr.support |= tUDP
	return nil
}

func (mgr *Manager) Send(src, dst net.Addr, msg sip.Message) error {
	switch src.Network() {
	case "udp":
		return mgr.SendUDP(src, dst, msg)
	case "tcp":
		return mgr.SendTCP(src, dst, msg)
	default:
		return fmt.Errorf("invalid or unsupported network %q", src.Network())
	}
}

func (mgr *Manager) SendUDP(src, dst net.Addr, msg sip.Message) error {
	name := sockName(src)
	ln, found := mgr.sock.Get(name)
	if !found {
		return fmt.Errorf("connection %q was not found", name)
	}

	udp, ok := ln.(*UDP)
	if !ok {
		return fmt.Errorf("found %q socket but type is not UDPConn", name)
	}

	n, err := udp.conn.WriteTo(msg.Byte(), dst)
	if err != nil {
		return fmt.Errorf("failed to write to udp:%s from %q: %s",
			dst, name, err)
	}
	logger.Log("successfully sent to udp %s from %q %d bytes", dst, src, n)
	return nil
}

func (mgr *Manager) SendTCP(src, dst net.Addr, msg sip.Message) error {
	name := connName(src, dst)
	cn, found := mgr.conn.Get(name)
	if !found {
		return fmt.Errorf("connection %q not found", name)
	}

	tcp, ok := cn.(*TCP)
	if !ok {
		return fmt.Errorf("found %q socket but type is not UDPConn", name)
	}

	n, err := tcp.conn.Write(msg.Byte())
	if err != nil {
		return fmt.Errorf("failed to write to conn %q: %s", name, err)
	}
	logger.Log("sent %d bytes to %q", n, name)

	return nil
}

func (mgr *Manager) Recv() <-chan Packet {
	return mgr.rcv
}

func (mgr *Manager) listen(ctx context.Context, ln Listener) {
	for {
		if ctx.Err() != nil {
			logger.Wrn("stop listener on context: %s", ctx.Err())
			return
		}

		if err := ln.listen(ctx); err != nil {
			logger.Err("failed to start listener: %s", err)
			logger.Wrn("restart in %v", reconnTout)
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
				logger.Err("accept connection err: %s", err)
				break connLoop
			}
		}

		logger.Wrn("close and restart listener")
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
		logger.Log("sent pack with payload of %d bytes sock %q remote address %q", size, laddr, raddr)
	case <-time.After(100 * time.Millisecond):
		logger.Err("failed to send pack on blocked chan for sock %q remote address %q", laddr, raddr)
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
