// Package sipmsg provides SIP Message parse and generate
package sipmsg

import (
	"errors"
	"net"
)

// errors
var (
	ErrURIParse = errors.New("SIP URI parse")
	ErrMsgParse = errors.New("SIP Message parse")
)

type Decoder struct {
	rcv chan Packet
	err chan error
}

type Packet struct {
	Laddr   net.Addr
	Raddr   net.Addr
	Message *Message
}

func NewDecoder() *Decoder {
	return &Decoder{
		rcv: make(chan Packet),
		err: make(chan error),
	}
}

func (d *Decoder) Recv() <-chan Packet {
	return d.rcv
}

func (d *Decoder) Err() <-chan error {
	return d.err
}

// Decode parse received SIP message into Message
func (d *Decoder) Decode(payload []byte, laddr, raddr net.Addr) {
	go func(payload []byte, laddr, raddr net.Addr) {
		msg, err := Parse(string(payload))
		if err != nil {
			d.err <- err
			return
		}
		d.rcv <- Packet{
			Laddr:   laddr,
			Raddr:   raddr,
			Message: msg,
		}
	}(payload, laddr, raddr)
}
