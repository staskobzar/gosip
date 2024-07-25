package transport

import (
	"net"
	"testing"

	"gosip/pkg/dns"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"

	"github.com/stretchr/testify/assert"
)

type mockDNS struct {
	naptr []*dns.NAPTR
	srv   []*dns.SRV
	ips   []net.IP
}

func (m *mockDNS) LookupNAPTR(string) []*dns.NAPTR { return m.naptr }
func (m *mockDNS) LookupSRV(string) []*dns.SRV     { return m.srv }
func (m *mockDNS) LookupAddr(string) []net.IP      { return m.ips }

func TestManagerResolve(t *testing.T) {
	t.Parallel()
	// helper closures
	toURI := func(sipuri string) *sipmsg.URI {
		uri, _ := sipmsg.ParseURI(sipuri)
		return uri
	}
	toUDP := func(ip string, port int) net.Addr {
		return &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	}
	toTCP := func(ip string, port int) net.Addr {
		return &net.TCPAddr{IP: net.ParseIP(ip), Port: port}
	}

	tests := map[string]struct {
		support  tTransp
		dns      sip.DNS
		uri      *sipmsg.URI
		wantAddr []net.Addr
		wantErr  error
	}{
		`empty dns member`: {
			uri: toURI("sip:alice@pbx.com"), wantErr: ErrResolv,
		},
		`empty input uri`: {
			dns: &mockDNS{}, uri: nil, wantErr: ErrResolv,
		},
		`uri with host IP and port and transport parameter`: {
			uri:      toURI("sip:alice@10.2.0.1:5457;transport=UDP"),
			dns:      &mockDNS{},
			wantAddr: []net.Addr{toUDP("10.2.0.1", 5457)},
		},
		`uri with host IP default udp port and transport parameter`: {
			uri:      toURI("sip:alice@10.2.0.100;transport=UDP"),
			dns:      &mockDNS{},
			wantAddr: []net.Addr{toUDP("10.2.0.100", 5060)},
		},
		`uri with host domain and port and tcp transport parameter`: {
			uri:      toURI("sip:alice@pbx.office.com:15099;transport=tcp"),
			dns:      &mockDNS{ips: []net.IP{{101, 25, 26, 10}, {101, 25, 26, 11}}},
			wantAddr: []net.Addr{toTCP("101.25.26.10", 15099), toTCP("101.25.26.11", 15099)},
		},
		`uri with host domain default tls port and transport parameter and no srv`: {
			uri:      toURI("sip:alice@pbx.com;transport=tls"),
			dns:      &mockDNS{ips: []net.IP{{10, 20, 30, 1}}},
			wantAddr: []net.Addr{toTCP("10.20.30.1", 5061)},
		},
		`uri with transport param and srv resolve rfc2782`: {
			uri: toURI("sip:alice@pbx.com;transport=tls"),
			dns: &mockDNS{
				srv: []*dns.SRV{{Target: "ssl.pbx.com", Port: 5062, Priority: 0, Weight: 0}},
				ips: []net.IP{{101, 0, 0, 100}},
			},
			wantAddr: []net.Addr{toTCP("101.0.0.100", 5062)},
		},
		`uri host is IP and default port and udp for sip scheme`: {
			uri:      toURI("sip:alice@192.168.0.1"),
			dns:      &mockDNS{},
			wantAddr: []net.Addr{toUDP("192.168.0.1", 5060)},
		},
		`uri host is IP and default port and tcp for sips scheme`: {
			uri:      toURI("sips:bob@10.200.202.10"),
			dns:      &mockDNS{},
			wantAddr: []net.Addr{toTCP("10.200.202.10", 5060)},
		},
		`uri host is IP and port and sip scheme`: {
			uri:      toURI("sip:alice@192.168.0.1:5588"),
			dns:      &mockDNS{},
			wantAddr: []net.Addr{toUDP("192.168.0.1", 5588)},
		},
		`uri host is domain resolve with port and sips scheme`: {
			uri:      toURI("sips:alice@h1.pbx.com:5589"),
			dns:      &mockDNS{ips: []net.IP{{110, 201, 30, 1}}},
			wantAddr: []net.Addr{toTCP("110.201.30.1", 5589)},
		},
		`naptr returns transport not supported`: {
			uri:     toURI("sip:alice@h1.pbx.com"),
			support: tTCP,
			dns:     &mockDNS{naptr: []*dns.NAPTR{{Flags: "s", Service: "SIP+D2U"}}},
			wantErr: ErrResolv,
		},
		`naptr resolves but srv is nil`: {
			uri:     toURI("sip:alice@h1.pbx.com"),
			support: tUDP,
			dns: &mockDNS{
				naptr: []*dns.NAPTR{{Flags: "s", Service: "SIP+D2U", Replace: "_sip._udp.pbx.com"}},
			},
			wantErr: ErrResolv,
		},
		`naptr resolves srv record`: {
			uri:     toURI("sip:alice@h1.pbx.com"),
			support: tUDP,
			dns: &mockDNS{
				naptr: []*dns.NAPTR{{Flags: "s", Service: "SIP+D2U", Replace: "_sip._udp.pbx.com"}},
				srv:   []*dns.SRV{{Target: "h1.pbx.com", Port: 8060}},
				ips:   []net.IP{{110, 2, 30, 100}},
			},
			wantAddr: []net.Addr{toUDP("110.2.30.100", 8060)},
		},
		`no naptr try srv match udp`: {
			uri:     toURI("sip:alice@h1.pbx.com"),
			support: tUDP,
			dns: &mockDNS{
				srv: []*dns.SRV{{Target: "h1.pbx.com", Port: 9060}},
				ips: []net.IP{{110, 2, 30, 1}},
			},
			wantAddr: []net.Addr{toUDP("110.2.30.1", 9060)},
		},
		`no naptr try srv match tls`: {
			uri:     toURI("sips:alice@h1.pbx.com"),
			support: tUDP | tSCTP | tTLS | tTCP,
			dns: &mockDNS{
				srv: []*dns.SRV{{Target: "h1.pbx.com", Port: 7061}},
				ips: []net.IP{{110, 2, 30, 1}},
			},
			wantAddr: []net.Addr{toTCP("110.2.30.1", 7061)},
		},
		`no naptr and no srv do tcp addr lookup for sips`: {
			uri:     toURI("sips:alice@h1.pbx.com"),
			support: tUDP | tTCP,
			dns: &mockDNS{
				ips: []net.IP{{110, 2, 30, 2}},
			},
			wantAddr: []net.Addr{toTCP("110.2.30.2", 5060)},
		},
		`no naptr and no srv do udp addr lookup for sip`: {
			uri:     toURI("sip:alice@h1.pbx.com"),
			support: tUDP | tTCP,
			dns: &mockDNS{
				ips: []net.IP{{110, 2, 30, 3}},
			},
			wantAddr: []net.Addr{toUDP("110.2.30.3", 5060)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mgr := &Manager{dns: tc.dns, support: tc.support}

			addrs, err := mgr.Resolve(tc.uri)
			assert.ErrorIs(t, err, tc.wantErr)

			assert.Equal(t, len(tc.wantAddr), len(addrs))

			for _, addr := range addrs {
				assert.Contains(t, tc.wantAddr, addr)
			}
		})
	}
}

func TestManagerNaptrSrvRec(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		support    tTransp
		input      []*dns.NAPTR
		wantTransp tTransp
		wantTarget string
		wantErr    error
	}{
		`nil input list`: {
			0, nil, 0, "", ErrResolv,
		},
		`transport not supported`: {
			tTCP | tTLS,
			[]*dns.NAPTR{
				{Flags: "s", Service: "SIP+D2U", Replace: "_sip._udp.h1.pbx.com"},
				{Flags: "s", Service: "SIP+D2S", Replace: "_sip._sctp.h2.pbx.com"},
			},
			0, "", ErrResolv,
		},
		`single record match udp`: {
			tUDP | tTCP | tSCTP | tTLS,
			[]*dns.NAPTR{{Flags: "s", Service: "SIP+D2U", Replace: "_sip._udp.h1.pbx.com"}},
			tUDP, "_sip._udp.h1.pbx.com", nil,
		},
		`invalid service`: {
			tUDP,
			[]*dns.NAPTR{{Flags: "s", Service: "PDU", Replace: "_sip._udp.h1.pbx.com"}},
			0, "", ErrResolv,
		},
		`unknown service`: {
			tUDP,
			[]*dns.NAPTR{{Flags: "s", Service: "SIP+D2Y", Replace: "_sip._udp.h1.pbx.com"}},
			0, "", ErrResolv,
		},
		`invalid naptr flag`: {
			tUDP | tTLS,
			[]*dns.NAPTR{{Flags: "e", Service: "SIPS+D2T", Replace: "_sips._tcp.h1.pbx.com"}},
			0, "", ErrResolv,
		},
		`record match tls`: {
			tUDP | tTLS,
			[]*dns.NAPTR{{Flags: "s", Service: "SIPS+D2T", Replace: "_sips._tcp.h1.pbx.com"}},
			tTLS, "_sips._tcp.h1.pbx.com", nil,
		},
		`record match tcp`: {
			tTLS | tTCP,
			[]*dns.NAPTR{{Flags: "s", Service: "SIP+D2T", Replace: "_sip._tcp.h1.pbx.com"}},
			tTCP, "_sip._tcp.h1.pbx.com", nil,
		},
		`record match sctp`: {
			tSCTP | tTCP,
			[]*dns.NAPTR{{Flags: "eS", Service: "SIP+D2S", Replace: "_sip._sctp.h1.pbx.com"}},
			tSCTP, "_sip._sctp.h1.pbx.com", nil,
		},
		`select in multi records`: {
			tTCP,
			[]*dns.NAPTR{
				{Flags: "s", Service: "SIP+D2U", Replace: "_sip._udp.h1.pbx.com"},
				{Flags: "s", Service: "SIP+D2S", Replace: "_sip._sctp.h2.pbx.com"},
				{Flags: "se", Service: "SIP+D2T", Replace: "_sip._tcp.h3.pbx.com"},
				{Flags: "s", Service: "SIPS+D2T", Replace: "_sips._tcp.h3.pbx.com"},
			},
			tTCP, "_sip._tcp.h3.pbx.com", nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mgr := &Manager{support: tc.support}

			transp, target, err := mgr.naptrSrvRec(tc.input)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantTransp, transp)
			assert.Equal(t, tc.wantTarget, target)
		})
	}
}
