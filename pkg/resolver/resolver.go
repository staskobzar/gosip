// Package resolver provides DNS lookup for NAPTR, SRV and A/AAA records
package resolver

import (
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/miekg/dns"
)

// module errors
var (
	Error     = errors.New("resolver")
	ErrConfig = fmt.Errorf("%w: config read", Error)
)

// Resolver serves to locate SIP servers
// via DNS queries and SIP addresses as
// specified in rfc3263
type Resolver struct {
	conf *dns.ClientConfig
}

// NAPTR record
type NAPTR struct {
	Flags   string
	Service string
	Replace string
	Order   int
	Pref    int
}

// SRV record
type SRV struct {
	Target   string
	Port     int
	Priority int
	Weight   int
}

// NewResolver creates new Resolver
func NewResolver(path string) (*Resolver, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConfig, err)
	}
	return NewResolverReader(file)
}

// NewResolverReader reads nameserevers from Reader and creates new Resolver
func NewResolverReader(r io.Reader) (*Resolver, error) {
	conf, err := dns.ClientConfigFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConfig, err)
	}

	if len(conf.Servers) == 0 {
		return nil, fmt.Errorf("%w: no servers found in the config file", ErrConfig)
	}
	resolver := &Resolver{conf: conf}

	return resolver, nil
}

// LookupNAPTR makes DNS NAPTR lookup and returns list of naptre records
func (r *Resolver) LookupNAPTR(domain string) []*NAPTR {
	logger.Log("naptr request for %q", domain)
	m := new(dns.Msg).SetQuestion(dns.Fqdn(domain), dns.TypeNAPTR)

	procAnsw := func(rr []dns.RR) []*NAPTR {
		answ := make([]*NAPTR, 0, len(rr))
		for _, r := range rr {
			if data, ok := r.(*dns.NAPTR); ok {
				naptr := &NAPTR{
					Order:   int(data.Order),
					Pref:    int(data.Preference),
					Flags:   data.Flags,
					Service: data.Service,
					Replace: data.Replacement,
				}
				answ = append(answ, naptr)
			}
		}
		return answ
	}

	for _, s := range r.conf.Servers {
		srv := net.JoinHostPort(s, r.conf.Port)
		resp, err := dns.Exchange(m, srv)
		if err != nil {
			logger.Err("failed lookup naptr at %q: %s", srv, err)
			continue
		}
		return procAnsw(resp.Answer)
	}
	return nil
}

// LookupSRV makes DNS SRV request and returns list of SRV targets
func (r *Resolver) LookupSRV(target string) []*SRV {
	logger.Log("srv request for %q", target)

	srv := make([]*SRV, 0)
	m := new(dns.Msg).SetQuestion(dns.Fqdn(target), dns.TypeSRV)
	for _, s := range r.conf.Servers {
		namesrv := net.JoinHostPort(s, r.conf.Port)
		resp, err := dns.Exchange(m, namesrv)
		if err != nil {
			logger.Err("failed lookup naptr at %q: %s", srv, err)
			continue
		}
		fmt.Printf("%#v\n", resp)
		for _, rr := range resp.Answer {
			fmt.Printf("%#v\n", rr)
		}
		return srv
	}

	return nil
}
