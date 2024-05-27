// Package sipmsg provides SIP Message parse and generate
package sipmsg

import (
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"net"
	"time"
)

// errors
var (
	ErrURIParse = errors.New("SIP URI parse")
	ErrMsgParse = errors.New("SIP Message parse")
	ErrDecoder  = errors.New("SIP Decoder")
)

// Decoder provides channels to return parsed SIP message or error
type Decoder struct {
	rcv chan Packet
	err chan error
}

// Packet represents decoded SIP Message and network information
type Packet struct {
	Laddr   net.Addr
	Raddr   net.Addr
	Message *Message
}

// NewDecoder creates and returns new Decoder
func NewDecoder() *Decoder {
	return &Decoder{
		rcv: make(chan Packet, 1),
		err: make(chan error, 1),
	}
}

// Recv returns receiver channel
func (d *Decoder) Recv() <-chan Packet {
	return d.rcv
}

// Err returns errors channel
func (d *Decoder) Err() <-chan error {
	return d.err
}

func (d *Decoder) toErr(err error) {
	select {
	case d.err <- err:
	case <-time.After(500 * time.Millisecond):
		logger.Err("txn:layer: failed to send error. Channel blocked")
	}
}

func (d *Decoder) toRecv(msg *Message, laddr, raddr net.Addr) {
	pack := Packet{
		Laddr:   laddr,
		Raddr:   raddr,
		Message: msg,
	}
	select {
	case d.rcv <- pack:
	case <-time.After(500 * time.Millisecond):
		logger.Err("txn:layer: failed to send decoded packet. Channel blocked")
	}

}

// Decode parse received SIP message into Message
func (d *Decoder) Decode(payload []byte, laddr, raddr net.Addr) {
	go func(payload []byte, laddr, raddr net.Addr) {
		msg, err := Parse(string(payload))
		if err != nil {
			d.toErr(err)
			return
		}

		via := msg.TopVia()
		if via == nil {
			d.toErr(fmt.Errorf("%w: invalid top Via header <nil>", ErrDecoder))
			return
		}

		if raddr == nil {
			d.toErr(fmt.Errorf("%w: invalid remote address <nil>", ErrDecoder))
			return
		}

		remoteHost, _, _ := net.SplitHostPort(raddr.String())

		if net.ParseIP(via.Host) == nil || via.Host != remoteHost {
			via.Recvd = remoteHost
		}

		d.toRecv(msg, laddr, raddr)
	}(payload, laddr, raddr)
}
