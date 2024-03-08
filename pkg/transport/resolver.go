package transport

import (
	"fmt"
	"gosip/pkg/dns"
	"gosip/pkg/logger"
	"gosip/pkg/sipmsg"
	"net"
	"strconv"
	"strings"
)

func (mgr *Manager) Resolve(uri *sipmsg.URI) ([]net.Addr, error) {
	if mgr.dns == nil {
		return nil, fmt.Errorf("%w: manager dns object is nil", ErrResolv)
	}
	if uri == nil {
		return nil, fmt.Errorf("%w: uri is nil", ErrResolv)
	}

	if len(uri.Transport) > 0 {
		return mgr.getPortIP(strToTransport(uri.Transport), uri, nil)
	}

	// if no transport protocol is specified, but the TARGET is a
	// numeric IP address, the client SHOULD use UDP for a SIP URI, and TCP
	// for a SIPS URI
	if ipaddr := net.ParseIP(uri.Hostport); ipaddr != nil {
		if strings.EqualFold(uri.Scheme, "sip") {
			return mgr.getPortIP(tUDP, uri, nil)
		}
		return mgr.getPortIP(tTCP, uri, nil)
	}

	// Similarly, if no transport protocol is specified,
	// and the TARGET is not numeric, but an explicit port is provided, the
	// client SHOULD use UDP for a SIP URI, and TCP for a SIPS URI.
	if host, port, err := net.SplitHostPort(uri.Hostport); err == nil {
		transp := tUDP
		if strings.EqualFold(uri.Scheme, "sips") {
			transp = tTCP
		}

		if ipaddr := net.ParseIP(host); ipaddr != nil {
			return netAddr(transp, []string{host}, port)
		}
		return netAddr(transp, mgr.lookupAddr(host), port)
	}

	// Otherwise, if no transport protocol or port is specified, and the
	// target is not a numeric IP address, the client SHOULD perform a NAPTR
	// query for the domain in the URI.
	if naptr := mgr.dns.LookupNAPTR(uri.Hostport); len(naptr) > 0 {
		transp, srvtarget, err := mgr.naptrSrvRec(naptr)
		if err != nil {
			return nil, err
		}
		srvrr := mgr.dns.LookupSRV(srvtarget)
		if len(srvrr) == 0 {
			return nil, fmt.Errorf("%w: failed resolve SRV rec %q", ErrResolv, srvtarget)
		}
		// TODO: failover srv records
		return netAddr(transp, mgr.lookupAddr(srvrr[0].Target), strconv.Itoa(srvrr[0].Port))
	}

	// If no NAPTR records are found, the client constructs SRV queries for
	// those transport protocols it supports, and does a query for each.
	for _, transp := range []tTransp{tUDP, tTCP, tTLS, tSCTP} {
		if transp != tTLS && uri.Scheme == "sips" {
			continue
		}
		if mgr.support&transp != transp {
			continue
		}
		srv := mgr.lookupSRV(transp, uri.Scheme, uri.Hostport)
		if len(srv) == 0 {
			continue
		}
		return netAddr(transp, mgr.lookupAddr(srv[0].Target), strconv.Itoa(srv[0].Port))
	}

	// If no SRV records are found, the client SHOULD use TCP for a SIPS
	// URI, and UDP for a SIP URI.
	if uri.Scheme == "sips" {
		return netAddr(tTCP, mgr.lookupAddr(uri.Hostport), "")
	}

	return netAddr(tUDP, mgr.lookupAddr(uri.Hostport), "")
}

// TODO: RFC3263#5 Server Usage
func (mgr *Manager) ResolveVia(via *sipmsg.HeaderVia) ([]net.Addr, error) {
	via = via
	return nil, nil
}

func (mgr *Manager) getPortIP(transp tTransp, uri *sipmsg.URI, srv []*dns.SRV) ([]net.Addr, error) {
	if transp == tUnknown {
		return nil, fmt.Errorf("TODO:test: invalid or unsupported transport")
	}

	if host, port, err := net.SplitHostPort(uri.Hostport); err == nil {
		// If TARGET is a numeric IP address, the client uses that address.  If
		// the URI also contains a port, it uses that port.  If no port is
		// specified, it uses the default port for the particular transport
		// protocol.
		if ipaddr := net.ParseIP(host); ipaddr != nil {
			return netAddr(transp, []string{host}, port)
		}
		// If the TARGET was not a numeric IP address, but a port is present in
		// the URI, the client performs an A or AAAA record lookup of the domain
		// name.
		return netAddr(transp, mgr.lookupAddr(host), port)
	}

	if ipaddr := net.ParseIP(uri.Hostport); ipaddr != nil {
		return netAddr(transp, []string{uri.Hostport}, "")
	}

	// If the TARGET was not a numeric IP address, and no port was present
	// in the URI, the client performs an SRV query on the record returned
	// from the NAPTR processing of Section 4.1 if such processing was
	// performed
	if len(srv) == 0 {
		// If it was not, because a transport was specified explicitly, the
		// client performs an SRV query for that specific transport
		srv = mgr.lookupSRV(transp, uri.Scheme, uri.Hostport) // if we got here then Hostport is host only
	}

	if len(srv) == 0 {
		return netAddr(transp, mgr.lookupAddr(uri.Hostport), "")
	}

	// TODO srv failover https://datatracker.ietf.org/doc/html/rfc3263#section-4.3
	// For SIP requests, failure occurs if the transaction layer reports a
	// 503 error response or a transport failure of some sort (generally,
	// due to fatal ICMP errors in UDP or connection failures in TCP).

	rr := srv[0]
	return netAddr(transp, mgr.lookupAddr(rr.Target), strconv.Itoa(rr.Port))
}

func (mgr *Manager) lookupSRV(transp tTransp, scheme, target string) []*dns.SRV {
	service := func() string {
		switch transp {
		case tTCP, tTLS:
			return "_tcp."
		case tSCTP:
			return "_sctp."
		default:
			return "_udp."
		}
	}

	proto := func() string {
		if strings.EqualFold(scheme, "sips") || transp == tTLS {
			return "_sips."
		}
		return "_sip."
	}

	// dns package should ensure sort and suffle srv records
	// by priority and weight as rfc2782 says
	return mgr.dns.LookupSRV(proto() + service() + target)
}

func (mgr *Manager) lookupAddr(target string) []string {
	ips := mgr.dns.LookupAddr(target)
	addrs := make([]string, len(ips))
	for i, host := range ips {
		addrs[i] = host.String()
	}
	return addrs
}

func (mgr *Manager) naptrSrvRec(naptr []*dns.NAPTR) (tTransp, string, error) {
	matchService := func(support tTransp, service string) tTransp {
		scheme, transp, found := strings.Cut(strings.ToLower(service), "+")
		if !found {
			return tUnknown
		}
		switch transp {
		case "d2u":
			return support & tUDP
		case "d2s":
			return support & tSCTP
		case "d2t":
			if scheme == "sips" {
				return support & tTLS
			}
			return support & tTCP
		default:
			return tUnknown
		}
	}

	for _, rr := range naptr {
		if !strings.Contains(strings.ToLower(rr.Flags), "s") {
			continue
		}
		if transp := matchService(mgr.support, rr.Service); transp != tUnknown {
			return transp, rr.Replace, nil
		}
	}

	return 0, "", fmt.Errorf("%w: not found supported NAPTR record", ErrResolv)
}

func netAddr(transp tTransp, hosts []string, port string) ([]net.Addr, error) {
	defaultPort := func(port, defPort string) string {
		if len(port) > 0 {
			return port
		}
		return defPort
	}
	toAddr := func(host, port string) (net.Addr, error) {
		switch transp {
		case tTCP:
			return net.ResolveTCPAddr("tcp", net.JoinHostPort(host, defaultPort(port, "5060")))
		case tTLS:
			return net.ResolveTCPAddr("tcp", net.JoinHostPort(host, defaultPort(port, "5061")))
		case tSCTP:
			// default port 5060
			return nil, fmt.Errorf("TODO: not implemented yet")
		case tUDP:
			return net.ResolveUDPAddr("udp", net.JoinHostPort(host, defaultPort(port, "5060")))
		default:
			return nil, fmt.Errorf("TODO: invalid transport %d", transp)
		}
	}
	addrs := make([]net.Addr, len(hosts))
	for i, host := range hosts {
		addr, err := toAddr(host, port)
		if err != nil {
			return nil, fmt.Errorf("TODO: error - failed to parse host %q with port %q: %s",
				host, port, err)
		}
		addrs[i] = addr
	}
	return addrs, nil
}

func strToTransport(transp string) tTransp {
	switch strings.ToLower(transp) {
	case "sctp":
		return tSCTP
	case "tcp":
		return tTCP
	case "tls":
		return tTLS
	case "udp":
		return tUDP
	default:
		logger.Err("transport: dns resolv: unknown transport param %q", transp)
		return tUnknown
	}
}

func transpPort(transp tTransp) int {
	if transp == tTLS {
		return 5061
	}
	return 5060
}
