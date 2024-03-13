package transport

import (
	"fmt"
	"gosip/pkg/dns"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"net"
	"strconv"
	"strings"
)

func (mgr *Manager) ResolveRURI(pack sipmsg.Packet) {
	go func(pack sipmsg.Packet) {
		if pack.Message == nil {
			mgr.err <- fmt.Errorf("%w: cannot resolve RURI: packet message is nil", ErrResolv)
			return
		}
		addrs, err := mgr.Resolve(pack.Message.RURI)
		if err != nil {
			mgr.err <- err
			return
		}
		sippack := sip.Packet{
			ReqAddrs:   addrs,
			LocalSock:  pack.Laddr,
			RemoteSock: pack.Raddr,
			Message:    pack.Message,
		}
		mgr.resolv <- sippack
	}(pack)
}

func (mgr *Manager) Resolve(uri *sipmsg.URI) ([]net.Addr, error) {
	dbgr("resolving uri %q", uri)
	if mgr.dns == nil {
		return nil, fmt.Errorf("%w: manager dns object is nil", ErrResolv)
	}
	if uri == nil {
		return nil, fmt.Errorf("%w: uri is nil", ErrResolv)
	}

	if len(uri.Transport) > 0 {
		dbgr("uri has transport param %q", uri.Transport)
		transp := strToTransport(uri.Transport)
		if transp == tUnknown {
			return nil, fmt.Errorf("%w: invalid or unsupported transport %q", ErrResolv, uri.Transport)
		}
		return mgr.getPortIP(transp, uri, nil)
	}

	// if no transport protocol is specified, but the TARGET is a
	// numeric IP address, the client SHOULD use UDP for a SIP URI, and TCP
	// for a SIPS URI
	if ipaddr := net.ParseIP(uri.Hostport); ipaddr != nil {
		dbgr("uri host is an IP address %q", ipaddr)
		return mgr.getPortIP(schemeToTransp(uri.Scheme), uri, nil)
	}

	// Similarly, if no transport protocol is specified,
	// and the TARGET is not numeric, but an explicit port is provided, the
	// client SHOULD use UDP for a SIP URI, and TCP for a SIPS URI.
	if host, port, err := net.SplitHostPort(uri.Hostport); err == nil {
		dbgr("uri host part has port %q", uri.Hostport)
		if ipaddr := net.ParseIP(host); ipaddr != nil {
			dbgr("hostpart is IP %q", ipaddr)
			return netAddr(schemeToTransp(uri.Scheme), []string{host}, port)
		}
		dbgr("hostpart is domain %q", host)
		return mgr.lookupAddr(schemeToTransp(uri.Scheme), host, port)
	}

	// Otherwise, if no transport protocol or port is specified, and the
	// target is not a numeric IP address, the client SHOULD perform a NAPTR
	// query for the domain in the URI.
	dbgr("trying to perform NAPTR lookup")
	if naptr := mgr.dns.LookupNAPTR(uri.Hostport); len(naptr) > 0 {
		dbgr("found %d NAPTR records", len(naptr))
		transp, srvtarget, err := mgr.naptrSrvRec(naptr)
		if err != nil {
			return nil, err
		}
		srvrr := mgr.dns.LookupSRV(srvtarget)
		if len(srvrr) == 0 {
			return nil, fmt.Errorf("%w: failed resolve SRV rec %q", ErrResolv, srvtarget)
		}
		return mgr.srvToAddr(transp, srvrr)
	}

	// If no NAPTR records are found, the client constructs SRV queries for
	// those transport protocols it supports, and does a query for each.
	dbgr("no NAPTR records found. trying SRV per supported transport")
	for _, transp := range []tTransp{tUDP, tTCP, tTLS, tSCTP} {
		if transp != tTLS && uri.Scheme == "sips" {
			continue
		}
		if mgr.support&transp != transp {
			continue
		}
		srv := mgr.lookupSRV(transp, uri.Scheme, uri.Hostport)
		if len(srv) == 0 {
			wrnr("no SRV records found")
			continue
		}
		return mgr.srvToAddr(transp, srv)
	}

	// If no SRV records are found, the client SHOULD use TCP for a SIPS
	// URI, and UDP for a SIP URI.
	dbgr("no SRV records found")
	return mgr.lookupAddr(schemeToTransp(uri.Scheme), uri.Hostport, "")
}

// TODO: RFC3263#5 Server Usage
func (mgr *Manager) ResolveVia(via *sipmsg.HeaderVia) ([]net.Addr, error) {
	via = via
	return nil, nil
}

func (mgr *Manager) getPortIP(transp tTransp, uri *sipmsg.URI, srv []*dns.SRV) ([]net.Addr, error) {
	if host, port, err := net.SplitHostPort(uri.Hostport); err == nil {
		dbgr("uri host part has port %q", port)
		// If TARGET is a numeric IP address, the client uses that address.  If
		// the URI also contains a port, it uses that port.  If no port is
		// specified, it uses the default port for the particular transport
		// protocol.
		if ipaddr := net.ParseIP(host); ipaddr != nil {
			dbgr("uri host is IP %q", ipaddr)
			return netAddr(transp, []string{host}, port)
		}
		// If the TARGET was not a numeric IP address, but a port is present in
		// the URI, the client performs an A or AAAA record lookup of the domain
		// name.
		return mgr.lookupAddr(transp, host, port)
	}

	if ipaddr := net.ParseIP(uri.Hostport); ipaddr != nil {
		dbgr("uri host is IP %q", ipaddr)
		return netAddr(transp, []string{uri.Hostport}, "")
	}

	// If the TARGET was not a numeric IP address, and no port was present
	// in the URI, the client performs an SRV query on the record returned
	// from the NAPTR processing of Section 4.1 if such processing was
	// performed
	if len(srv) == 0 {
		// If it was not, because a transport was specified explicitly, the
		// client performs an SRV query for that specific transport
		dbgr("no srv records for ip:port lookup. trying srv")
		srv = mgr.lookupSRV(transp, uri.Scheme, uri.Hostport) // if we got here then Hostport is host only
	}

	if len(srv) == 0 {
		dbgr("still no srv records. lookup A record for %q", uri.Hostport)
		return mgr.lookupAddr(transp, uri.Hostport, "")
	}

	// TODO srv failover https://datatracker.ietf.org/doc/html/rfc3263#section-4.3
	// For SIP requests, failure occurs if the transaction layer reports a
	// 503 error response or a transport failure of some sort (generally,
	// due to fatal ICMP errors in UDP or connection failures in TCP).

	return mgr.srvToAddr(transp, srv)
}

func (mgr *Manager) lookupSRV(transp tTransp, scheme, domain string) []*dns.SRV {
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

	target := proto() + service() + domain
	dbgr("srv lookup for %q", target)

	// dns package should ensure sort and suffle srv records
	// by priority and weight as rfc2782 says
	return mgr.dns.LookupSRV(target)
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
			wrnr("NAPTR record has incompatible flags %q", rr.Flags)
			continue
		}
		if transp := matchService(mgr.support, rr.Service); transp != tUnknown {
			dbgr("match NAPTR record service %q, target: %q", rr.Service, rr.Replace)
			return transp, rr.Replace, nil
		}
		wrnr("NAPTR record service %q is not supported by manager transports", rr.Service)
	}

	return 0, "", fmt.Errorf("%w: not found supported NAPTR record", ErrResolv)
}

func (mgr *Manager) srvToAddr(transp tTransp, srv []*dns.SRV) ([]net.Addr, error) {
	addrs := make([]net.Addr, 0, len(srv))
	for _, rr := range srv {
		dbgr("using srv record %#v", rr)
		addr, err := mgr.lookupAddr(transp, rr.Target, strconv.Itoa(rr.Port))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr...)
	}
	return addrs, nil
}

func (mgr *Manager) lookupAddr(transp tTransp, target, port string) ([]net.Addr, error) {
	ips := mgr.dns.LookupAddr(target)
	if len(ips) == 0 {
		wrnr("no addresses found for %q. trying lookup host", target)
		return mgr.lookupHost(transp, target, port)
	}
	addrs := make([]string, len(ips))
	for i, host := range ips {
		addrs[i] = host.String()
	}
	return netAddr(transp, addrs, port)
}

func (mgr *Manager) lookupHost(transp tTransp, target, port string) ([]net.Addr, error) {
	addrs, err := net.LookupHost(target)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to resolv host %q: %s", ErrResolv, target, err)
	}
	return netAddr(transp, addrs, port)
}

func netAddr(transp tTransp, hosts []string, port string) ([]net.Addr, error) {
	dbgr("trying to resolve %d addresses", len(hosts))
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
		dbgr("resolved address %q", addr)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse host %q: %s",
				ErrResolv, host, err)
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
		wrnr("unknown transport param %q", transp)
		return tUnknown
	}
}

func transpPort(transp tTransp) int {
	if transp == tTLS {
		return 5061
	}
	return 5060
}

func schemeToTransp(scheme string) tTransp {
	if strings.EqualFold(scheme, "sips") {
		dbgr("using TCP transport for %q uri scheme", scheme)
		return tTCP
	}
	dbgr("using UDP transport for %q uri scheme", scheme)
	return tUDP
}
