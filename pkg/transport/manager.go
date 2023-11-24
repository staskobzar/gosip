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
)

type Packet struct {
	payload []byte
	laddr   net.Addr // local address
	raddr   net.Addr // remote address
}

type Manager struct {
	sock []Listener
	conn []net.Conn
	rcv  chan Packet
}

type Listener interface {
	listen() error
	recv(ctx context.Context) (net.Addr, net.Addr, []byte, error)
	close()
}

func InitManager() *Manager {
	logger.Log("init transport manager")
	return &Manager{
		conn: make([]net.Conn, 0),
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

func (mgm *Manager) Send(laddr, raddr net.Addr, msg sip.Message) error {
	logger.Log("sending to %q from %q", raddr, laddr)
	return nil
}

func (mgm *Manager) Recv() <-chan Packet {
	return mgm.rcv
}

func (mgm *Manager) listen(ctx context.Context, ln Listener) {
	for {
		logger.Log("starting listen")
		if ctx.Err() != nil {
			logger.Wrn("listen context Err is set. Exit loop")
			return
		}

		if err := ln.listen(); err != nil {
			logger.Err("failed to start listen (restart in 3sec): %s", err)
			<-time.After(3 * time.Second) // TODO use Timer Tick instead of After
			continue
		}

		err := mgm.loop(ctx, ln)
		logger.Err("manager loop return error: %s", err)

		ln.close()
	}
}

func (mgm *Manager) loop(ctx context.Context, ln Listener) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: loop is terminated by context", ErrManager)
		default:
			raddr, laddr, payload, err := ln.recv(ctx)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrManager, err)
			}
			logger.Log("mgm recv %s:%s", laddr.Network(), laddr.String())
			go mgm.consume(payload, raddr, laddr)
		}
	}
}

func (mgm *Manager) consume(payload []byte, rAddr, lAddr net.Addr) {
	pack := Packet{
		payload: payload,
		laddr:   lAddr,
		raddr:   rAddr,
	}

	select {
	case mgm.rcv <- pack:
	default:
		logger.Err("mgm failed to send to rcv channel. blocked...")
	}
}
