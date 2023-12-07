package transport

import (
	"context"
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"net"
	"time"
)

var (
	ErrManager     = errors.New("ERR Transp Manager")
	ErrUDPListener = errors.New("UDP listener")
	ErrTCPListener = errors.New("TCP listener")
)

type Packet struct {
	Payload []byte
	Laddr   net.Addr // local address
	Raddr   net.Addr // remote address
}

type Manager struct {
	sock *Store[Listener]
	conn *Store[Conn]
	rcv  chan Packet
}

type Listener interface {
	listen() error
	recv(context.Context, chan<- Packet) error
	id() sID
	close()
}

type Conn interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

func InitManager() *Manager {
	logger.Log("init transport manager")
	return &Manager{
		sock: NewStore[Listener](),
		conn: NewStore[Conn](),
		rcv:  make(chan Packet, 12), // buffered channel of messages
	}
}

func (mgm *Manager) ListenUDP(ctx context.Context, addrport string) error {
	addr, err := net.ResolveUDPAddr("udp", addrport)
	if err != nil {
		return fmt.Errorf("%w: address %q: %s", ErrUDPListener, addrport, err)
	}

	ln := &UDP{laddr: addr}

	go mgm.listen(ctx, ln)

	return nil
}

// func (mgm *Manager) ListenTCP(ctx context.Context, addrport string) error {
// 	addr, err := net.ResolveTCPAddr("tcp", addrport)
// 	if err != nil {
// 		return fmt.Errorf("%w: address %q: %s", ErrTCPListener, addrport, err)
// 	}

// 	ln := &TCPListener{laddr: addr}

// 	go mgm.listen(ctx, ln)

// 	return nil
// }

func (mgm *Manager) Send(laddr, raddr net.Addr, msg sip.Message) error {
	logger.Log("sending to %q from %q", raddr, laddr)
	if raddr == nil {
		return fmt.Errorf("remote address can not be nil")
	}

	switch raddr.Network() {
	case "udp":
		return mgm.SendUDP(laddr, raddr, msg)
	default:
		return fmt.Errorf("invalid transport to send %q", raddr.Network())
	}
}

func (mgm *Manager) SendUDP(laddr, raddr net.Addr, msg sip.Message) error {
	logger.Log("send via UDP from %q to %q", laddr, raddr)
	if laddr == nil {
		return mgm.SendUDPFromAny(raddr, msg)
	}
	key := fmt.Sprintf("%s:%s", laddr.Network(), laddr.String())
	ln, found := mgm.sock.Get(sID(key))
	if !found {
		return fmt.Errorf("socket for %q not found", laddr)
	}

	conn := ln.(*UDP)
	logger.Log("found UDP interface %v", conn)
	return conn.Write(raddr, msg.Byte())
}

func (mgm *Manager) SendUDPFromAny(raddr net.Addr, msg sip.Message) error {
	logger.Log("lookup any suitable UDP socket to send to %q", raddr)
	mgm.sock.mu.RLock()
	defer mgm.sock.mu.RUnlock()
	for _, ln := range mgm.sock.pool {
		switch trp := ln.(type) {
		case *UDP:
			logger.Log("found transport on %q", trp.laddr)
			return trp.Write(raddr, msg.Byte())
		}
	}
	return fmt.Errorf("no listener found for %q", raddr)
}

func (mgm *Manager) Recv() <-chan Packet {
	return mgm.rcv
}

func (mgm *Manager) listen(ctx context.Context, ln Listener) {
	for {
		if ctx.Err() != nil {
			logger.Wrn("listen context Err is set. Exit loop")
			return
		}

		if err := ln.listen(); err != nil {
			logger.Err("failed to start listen (restart in 3sec): %s", err)
			<-time.After(3 * time.Second) // TODO use Timer Tick instead of After
			continue
		}

		mgm.sock.Put(ln.id(), ln)

		err := mgm.accept(ctx, ln)
		logger.Err("manager accept loop returns an error: %s", err)
		// err := ln.recv(ctx, mgm.rcv)
		// logger.Err("manager loop return error: %s", err)

		ln.close()
		mgm.sock.Del(ln.id())
	}
}

func (mgm *Manager) accept(ctx context.Context, ln Listener) error {
	for {
		conn, err := ln.accept(ctx)
		if err != nil {
			return err
		}

		mgm.conn.Put(conn.id(), conn)

		go conn.consume()
	}
}

func consume(rcv chan<- Packet, pack Packet) {
	select {
	case rcv <- pack:
		logger.Log("successfully sent pack to TU")
	default:
		logger.Err("mgm failed to send to rcv channel. blocked...")
	}
}

func sIDBuild(addr net.Addr) sID {
	return sID(addr.Network() + ":" + addr.String())
}
