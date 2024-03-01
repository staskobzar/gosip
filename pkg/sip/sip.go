package sip

import (
	"gosip/pkg/dns"
	"net"
	"net/netip"
)

// DNS interface to a module that implements
// resolving NAPTR/DRV and A records
type DNS interface {
	LookupNAPTR(target string) []*dns.NAPTR
	LookupSRV(target string) []*dns.SRV
	LookupAddr(target string) []net.IP
}

// Message interface for SIP requests or responses
type Message interface {
	// generate ACK for initial SIP request
	// this is part of the transaction ACK
	// for 3xx-6xx responses: rfc3261#17.1.1.3
	Ack(Message) Message
	IsResponse() bool
	TopViaBranch() string
	SIPMethod() string
	ResponseCode() int
	Byte() []byte
}

// Transport SIP
type Transport interface {
	Send(addr netip.AddrPort, msg Message) error
	IsReliable() bool
}

// Transaction SIP
type Transaction interface{}
