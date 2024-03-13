package sip

import (
	"gosip/pkg/dns"
	"gosip/pkg/sipmsg"
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

type Packet struct {
	ReqAddrs   []net.Addr
	LocalSock  net.Addr
	RemoteSock net.Addr
	Message    *sipmsg.Message
}

// Transport SIP
type Transport interface {
	Send(addr netip.AddrPort, msg *sipmsg.Message) error
	IsReliable() bool
}

// Transaction SIP
type Transaction interface{}
