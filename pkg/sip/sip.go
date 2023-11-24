package sip

import "net/netip"

type Transport interface {
	Send(addr netip.AddrPort, msg Message) error
	IsReliable() bool
}

type Transaction interface{}

type Message interface {
	Ack() Message
	IsResponse() bool
	TopViaBranch() string
	Method() string
	ResponseCode() int
}
