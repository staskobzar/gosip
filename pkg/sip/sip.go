package sip

import (
	"gosip/pkg/dns"
	"gosip/pkg/sipmsg"
	"net"
)

// DNS interface to a module that implements
// resolving NAPTR/DRV and A records
type DNS interface {
	LookupNAPTR(target string) []*dns.NAPTR
	LookupSRV(target string) []*dns.SRV
	LookupAddr(target string) []net.IP
}

// Packet with SIP Message and transport
type Packet struct {
	SendTo     []net.Addr
	ReqAddrs   []net.Addr
	LocalSock  net.Addr
	RemoteSock net.Addr
	Message    *sipmsg.Message
}

// Transaction that all invite/non-invite
// and client/server transactions must impelemnt
type Transaction interface {
	BranchID() string
	Consume(*Packet)
	Match(msg *sipmsg.Message) (Transaction, bool)
	Terminate()
}
