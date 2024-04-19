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

// Transport SIP
type Transport interface {
	Send(addr netip.AddrPort, msg *sipmsg.Message) error
	IsReliable() bool
}

type Packet struct {
	SendTo     []net.Addr
	ReqAddrs   []net.Addr
	LocalSock  net.Addr
	RemoteSock net.Addr
	Message    *sipmsg.Message
}
