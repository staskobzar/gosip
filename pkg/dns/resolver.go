// Package dns provides DNS lookup for NAPTR, SRV and A/AAA records
package dns

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

	quest := new(dns.Msg).SetQuestion(dns.Fqdn(domain), dns.TypeNAPTR)

	result := func(list []dns.RR) []*NAPTR {
		naptr := make([]*NAPTR, len(list))
		for i, answ := range list {
			rr := answ.(*dns.NAPTR)
			naptr[i] = &NAPTR{
				Order:   int(rr.Order),
				Pref:    int(rr.Preference),
				Flags:   rr.Flags,
				Service: rr.Service,
				Replace: rr.Replacement,
			}
		}
		return naptr
	}

	return lookup[*NAPTR](quest, r.conf.Servers, r.conf.Port, result)
}

// LookupSRV makes DNS SRV request and returns list of SRV targets
func (r *Resolver) LookupSRV(target string) []*SRV {
	logger.Log("srv request for %q", target)

	quest := new(dns.Msg).SetQuestion(dns.Fqdn(target), dns.TypeSRV)

	result := func(list []dns.RR) []*SRV {
		srv := make([]*SRV, len(list))
		for i, answ := range list {
			rr := answ.(*dns.SRV)
			srv[i] = &SRV{
				Target:   rr.Target,
				Port:     int(rr.Port),
				Priority: int(rr.Priority),
				Weight:   int(rr.Weight),
			}
		}
		return srv
	}

	return lookup[*SRV](quest, r.conf.Servers, r.conf.Port, result)
}

// LookupAddr resolves domain IP address(es)
// TODO: IPv6 AAAA support
func (r *Resolver) LookupAddr(target string) []net.IP {
	logger.Log("address request for %q", target)
	quest := new(dns.Msg).SetQuestion(dns.Fqdn(target), dns.TypeA)

	result := func(list []dns.RR) []net.IP {
		ips := make([]net.IP, len(list))
		for i, answ := range list {
			rr := answ.(*dns.A)
			ips[i] = rr.A
		}
		return ips
	}

	return lookup[net.IP](quest, r.conf.Servers, r.conf.Port, result)
}

func lookup[T any](quest *dns.Msg, namesrv []string, port string, result func([]dns.RR) []T) []T {
	for _, ns := range namesrv {
		resp, err := dns.Exchange(quest, net.JoinHostPort(ns, port))
		if err != nil {
			logger.Err("failed lookup naptr at %q: %s", ns, err)
			continue
		}
		return result(resp.Answer)
	}
	return nil
}
