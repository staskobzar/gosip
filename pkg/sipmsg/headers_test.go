package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseViaHeaders(t *testing.T) {
	t.Run("single header", func(t *testing.T) {
		tests := []struct {
			hdr  string
			want HeaderVia
		}{
			{
				"Via: SIP/2.0/UDP pbx.com ;branch=z9hG4bKnashds7",
				HeaderVia{
					HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "pbx.com",
					Port: "", Branch: "z9hG4bKnashds7", Params: "branch=z9hG4bKnashds7",
				},
			},
			{
				"VIA: SIP/ 2.0/ TCP 10.0.0.1: 15060; branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120",
				HeaderVia{
					HeaderName: "VIA", Proto: "SIP/ 2.0/ ", Transp: "TCP", Host: "10.0.0.1",
					Port: "15060", Branch: "z9hG4bKna", Recvd: "10.0.0.100",
					Params: "branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120",
				},
			},
			{
				"v: SIP/2.0/TLS [fe80::2e8d:b1ff:fef3:8a40] :6060;branch=z9hG4bKff;rl",
				HeaderVia{
					HeaderName: "v", Proto: "SIP/2.0/", Transp: "TLS", Host: "[fe80::2e8d:b1ff:fef3:8a40]",
					Port: "6060", Branch: "z9hG4bKff", Params: "branch=z9hG4bKff;rl",
				},
			},
		}
		for _, tc := range tests {
			msg, err := Parse(toMsg([]string{tc.hdr}))
			assert.Nil(t, err)

			via := msg.Find(HVia).(*HeaderVia)

			assert.Equal(t, tc.want.HeaderName, via.HeaderName)
			assert.Equal(t, tc.want.Proto, via.Proto)
			assert.Equal(t, tc.want.Transp, via.Transp)
			assert.Equal(t, tc.want.Host, via.Host)
			assert.Equal(t, tc.want.Port, via.Port)
			assert.Equal(t, tc.want.Branch, via.Branch)
			assert.Equal(t, tc.want.Recvd, via.Recvd)
			assert.Equal(t, tc.want.Params, via.Params)
			assert.Nil(t, via.Next)
		}
	})

	t.Run("linked header", func(t *testing.T) {
		hdr := "Via: SIP/2.0/UDP h1.pbx.com;branch=z9hG4bKnashd" +
			", SIP/2.0/UDP h2.pbx.com;branch=z9hG4bKnasff;received=10.0.0.1" +
			", SIP/2.0/TCP h3.pbx.com:11609;branch=z9hG4bKnasfe;maddr=10.0.0.1"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)

		via := msg.Find(HVia).(*HeaderVia)
		assert.Equal(t, "Via", via.HeaderName)
		assert.Equal(t, "SIP/2.0/", via.Proto)
		assert.Equal(t, "UDP", via.Transp)
		assert.Equal(t, "h1.pbx.com", via.Host)
		assert.Equal(t, "z9hG4bKnashd", via.Branch)

		via = via.Next
		assert.Equal(t, "", via.HeaderName)
		assert.Equal(t, "SIP/2.0/", via.Proto)
		assert.Equal(t, "UDP", via.Transp)
		assert.Equal(t, "h2.pbx.com", via.Host)
		assert.Equal(t, "z9hG4bKnasff", via.Branch)
		assert.Equal(t, "10.0.0.1", via.Recvd)

		via = via.Next
		assert.Equal(t, "TCP", via.Transp)
		assert.Equal(t, "h3.pbx.com", via.Host)
		assert.Equal(t, "z9hG4bKnasfe", via.Branch)
		assert.Equal(t, ";branch=z9hG4bKnasfe;maddr=10.0.0.1", via.Params.String())
		assert.Nil(t, via.Next)
	})

	t.Run("multiple vias", func(t *testing.T) {
		//nolint:goconst
		input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
			"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
			"V: SIP / 2.0 / UDP first.example.com: 4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1\r\n" +
			"Max-Forwards: 70\r\n" +
			"To: Bob <sip:bob@biloxi.com>\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
			"Call-ID: 843817637684230@998sdasdh09\r\n" +
			"CSeq: 1826 REGISTER\r\n" +
			"Accept \t: */*\r\n" +
			"VIA: SIP/2.0/TCP h1.biloxi.com;branch=z9hG4bKna.7\r\n" +
			"Contact: <sip:bob@192.0.2.4>\r\n" +
			"Expires: 7200\r\n" +
			"Content-Type: application/simple-message-summary ; foo/bar;xy/123d\r\n" +
			"v: SIP/2.0/TLS h2.biloxi.com;branch=z9hG4bKna.9\r\n" +
			"Content-Length: 1234\r\n\r\n"
		msg, err := Parse(input)
		assert.Nil(t, err)
		hdrs := msg.FindAll(HVia)
		assert.Equal(t, 4, hdrs.Len())
	})
}

func TestAnyHeaderCopy(t *testing.T) {
	t.Run("HeaderVia", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			svia := "v: SIP/2.0/UDP example.com:4000;ttl=16;maddr=224.2.0.1;" +
				"branch=z9hG4bKa7c6a8dlze.1;received=10.0.0.100;rport"
			msg, err := Parse(toMsg([]string{svia}))
			assert.Nil(t, err)
			via := msg.Find(HVia).(*HeaderVia)
			v := via.Copy().(*HeaderVia)
			assert.NotSame(t, via, v)
			assert.Equal(t, "v", v.HeaderName)
			assert.Equal(t, "SIP/2.0/", v.Proto)
			assert.Equal(t, "UDP", v.Transp)
			assert.Equal(t, "example.com", v.Host)
			assert.Equal(t, "4000", v.Port)
			assert.Equal(t, "z9hG4bKa7c6a8dlze.1", v.Branch)
			assert.Equal(t, "10.0.0.100", v.Recvd)
			assert.Equal(t, ";ttl=16;maddr=224.2.0.1;branch=z9hG4bKa7c6a8dlze.1;received=10.0.0.100;rport",
				v.Params.String())
			assert.Nil(t, v.Next)
		})

		t.Run("with linked next", func(t *testing.T) {
			svia := "via: SIP/2.0/SCTP atlanta.com;branch=ff00aa;ttl=60" +
				",SIP/2.0/UDP u.pbx.com;branch=z9hG4bKnasff"
			msg, err := Parse(toMsg([]string{svia}))
			assert.Nil(t, err)
			via := msg.Find(HVia).(*HeaderVia)
			v := via.Copy().(*HeaderVia)
			assert.NotSame(t, via, v)
			assert.Equal(t, "via", v.HeaderName)
			assert.Equal(t, "SIP/2.0/", v.Proto)
			assert.Equal(t, "SCTP", v.Transp)
			assert.Equal(t, "atlanta.com", v.Host)
			assert.Equal(t, "", v.Port)
			assert.Equal(t, "ff00aa", v.Branch)
			assert.Equal(t, "", v.Recvd)
			assert.Equal(t, ";branch=ff00aa;ttl=60", v.Params.String())
			assert.NotNil(t, v.Next)

			v = v.Next
			assert.Equal(t, "", v.HeaderName)
			assert.Equal(t, "SIP/2.0/", v.Proto)
			assert.Equal(t, "UDP", v.Transp)
			assert.Equal(t, "u.pbx.com", v.Host)
			assert.Equal(t, "", v.Port)
			assert.Equal(t, "z9hG4bKnasff", v.Branch)
			assert.Equal(t, "", v.Recvd)
			assert.Equal(t, ";branch=z9hG4bKnasff", v.Params.String())
			assert.Nil(t, v.Next)
		})

		t.Run("linked multiple", func(t *testing.T) {
			svia := "Via: SIP/2.0/UDP h1.pbx.com;branch=z9hG4bKnashd" +
				",SIP/2.0/UDP h2.pbx.com;branch=z9hG4bKnasff;received=10.0.0.1" +
				", SIP/2.0/TCP h3.pbx.com:11609;branch=z9hG4bKnasfe;maddr=10.0.0.1"
			msg, err := Parse(toMsg([]string{svia}))
			assert.Nil(t, err)
			via := msg.Find(HVia).(*HeaderVia)
			v := via.Copy().(*HeaderVia)
			assert.NotSame(t, via, v)
			assert.Equal(t, "Via", v.HeaderName)
			assert.Equal(t, "SIP/2.0/", v.Proto)
			assert.Equal(t, "UDP", v.Transp)
			assert.Equal(t, "h1.pbx.com", v.Host)
			assert.Equal(t, "", v.Port)
			assert.Equal(t, "z9hG4bKnashd", v.Branch)
			assert.Equal(t, "", v.Recvd)
			assert.Equal(t, ";branch=z9hG4bKnashd", v.Params.String())
			assert.NotNil(t, v.Next)

			v = v.Next
			assert.Equal(t, "", v.HeaderName)
			assert.Equal(t, "SIP/2.0/", v.Proto)
			assert.Equal(t, "UDP", v.Transp)
			assert.Equal(t, "h2.pbx.com", v.Host)
			assert.Equal(t, "", v.Port)
			assert.Equal(t, "z9hG4bKnasff", v.Branch)
			assert.Equal(t, "10.0.0.1", v.Recvd)
			assert.Equal(t, ";branch=z9hG4bKnasff;received=10.0.0.1", v.Params.String())
			assert.NotNil(t, v.Next)

			v = v.Next
			assert.Equal(t, "", v.HeaderName)
			assert.Equal(t, "SIP/2.0/", v.Proto)
			assert.Equal(t, "TCP", v.Transp)
			assert.Equal(t, "h3.pbx.com", v.Host)
			assert.Equal(t, "11609", v.Port)
			assert.Equal(t, "z9hG4bKnasfe", v.Branch)
			assert.Equal(t, "", v.Recvd)
			assert.Equal(t, ";branch=z9hG4bKnasfe;maddr=10.0.0.1", v.Params.String())
			assert.Nil(t, v.Next)
		})
	})

	t.Run("HeaderGeneric", func(t *testing.T) {
		t.Run("P-Asserted-Identity", func(t *testing.T) {
			hdr := &HeaderGeneric{HeaderName: "P-Asserted-Identity", Value: "Alice <sip:alice@atlanta.com>", T: HGeneric}

			pai, ok := hdr.Copy().(*HeaderGeneric)
			assert.True(t, ok)
			assert.Equal(t, "P-Asserted-Identity", pai.HeaderName)
			assert.Equal(t, "Alice <sip:alice@atlanta.com>", pai.Value)
			assert.Equal(t, HGeneric, pai.T)
		})
		t.Run("Allow", func(t *testing.T) {
			hdr := &HeaderGeneric{HeaderName: "Allow", Value: "INVITE, ACK, OPTIONS, CANCEL, BYE", T: HAllow}
			allow, ok := hdr.Copy().(*HeaderGeneric)
			assert.True(t, ok)
			assert.Equal(t, "Allow", allow.HeaderName)
			assert.Equal(t, "INVITE, ACK, OPTIONS, CANCEL, BYE", allow.Value)
			assert.Equal(t, HAllow, allow.T)
		})
	})

	t.Run("NameAddr", func(t *testing.T) {
		headers := []string{"To: Alice <sip:alice@biloxi.com;transport=udp>",
			"From: Bob <sip:bob@atlanta.com>;tag=456248"}
		msg, err := Parse(toMsg(headers))
		assert.Nil(t, err)

		t.Run("To header", func(t *testing.T) {
			hdr, ok := msg.To.Copy().(*NameAddr)
			assert.True(t, ok)
			assert.NotSame(t, hdr, msg.To)
			assert.Equal(t, "To", hdr.HeaderName)
			assert.Equal(t, "Alice", hdr.DisplayName)
			assert.NotSame(t, hdr.Addr, msg.To.Addr)
			assert.Equal(t, "sip:alice@biloxi.com;transport=udp", hdr.Addr.String())
			assert.Equal(t, "", hdr.Params.str())
			assert.Equal(t, HTo, hdr.T)
			assert.Equal(t, "", hdr.Tag)
		})

		t.Run("From header", func(t *testing.T) {
			hdr, ok := msg.From.Copy().(*NameAddr)
			assert.True(t, ok)
			assert.NotSame(t, hdr, msg.From)
			assert.Equal(t, "From", hdr.HeaderName)
			assert.Equal(t, "Bob", hdr.DisplayName)
			assert.NotSame(t, hdr.Addr, msg.From.Addr)
			assert.Equal(t, "sip:bob@atlanta.com", hdr.Addr.String())
			assert.Equal(t, ";tag=456248", hdr.Params.String())
			assert.Equal(t, HFrom, hdr.T)
			assert.Equal(t, "456248", hdr.Tag)
		})
	})

	t.Run("HeaderRoute", func(t *testing.T) {
		t.Run("single header", func(t *testing.T) {
			hdrs := []string{"Record-Route: <sip:h1.domain.com;lr>;host=one"}
			msg, err := Parse(toMsg(hdrs))
			assert.Nil(t, err)
			hdr, ok := msg.Find(HRecordRoute).(*HeaderRoute)
			assert.True(t, ok)
			rr, ok := hdr.Copy().(*HeaderRoute)
			assert.True(t, ok)
			assert.Equal(t, HRecordRoute, rr.T)
			assert.Equal(t, "Record-Route", rr.HeaderName)
			assert.Equal(t, "", rr.DisplayName)
			assert.NotSame(t, rr.Addr, hdr.Addr)
			assert.Equal(t, "sip:h1.domain.com;lr", rr.Addr.String())
			assert.Equal(t, ";host=one", rr.Params.String())

			assert.Nil(t, rr.Next)
		})

		t.Run("multiple linked headers", func(t *testing.T) {
			hdrs := []string{
				"Route: HQ1 <sip:h2.domain.com;lr>, <sip:dd1.pbx.com>;user=pbx,<sips:dd2.pbx.com>",
			}
			msg, err := Parse(toMsg(hdrs))
			assert.Nil(t, err)
			hdr, ok := msg.Find(HRoute).(*HeaderRoute)
			assert.True(t, ok)
			rr, ok := hdr.Copy().(*HeaderRoute)
			assert.True(t, ok)
			assert.Equal(t, HRoute, rr.T)
			assert.Equal(t, "Route", rr.HeaderName)
			assert.Equal(t, "HQ1", rr.DisplayName)
			assert.NotSame(t, rr.Addr, hdr.Addr)
			assert.Equal(t, "sip:h2.domain.com;lr", rr.Addr.String())
			assert.Equal(t, "", rr.Params.String())
			assert.NotNil(t, rr.Next)

			rr = rr.Next
			hdr = hdr.Next
			assert.Equal(t, HRoute, rr.T)
			assert.Equal(t, "", rr.HeaderName)
			assert.Equal(t, "", rr.DisplayName)
			assert.NotSame(t, rr.Addr, hdr.Addr)
			assert.Equal(t, "sip:dd1.pbx.com", rr.Addr.String())
			assert.Equal(t, ";user=pbx", rr.Params.String())
			assert.NotNil(t, rr.Next)

			rr = rr.Next
			hdr = hdr.Next
			assert.Equal(t, HRoute, rr.T)
			assert.Equal(t, "", rr.HeaderName)
			assert.Equal(t, "", rr.DisplayName)
			assert.NotSame(t, rr.Addr, hdr.Addr)
			assert.Equal(t, "sips:dd2.pbx.com", rr.Addr.String())
			assert.Equal(t, "", rr.Params.String())
			assert.Nil(t, rr.Next)
		})
	})

	t.Run("HeaderContact", func(t *testing.T) {
		t.Run("single header", func(t *testing.T) {
			hdr := "Contact: \"Bob C\" <sip:bob@192.0.2.4>;q=0.5;foo=bar;expires=1800"
			msg, err := Parse(toMsg([]string{hdr}))
			assert.Nil(t, err)
			hct, ok := msg.Find(HContact).(*HeaderContact)
			assert.True(t, ok)
			ct, ok := hct.Copy().(*HeaderContact)

			assert.True(t, ok)
			assert.Equal(t, "Contact", ct.HeaderName)
			assert.Equal(t, "\"Bob C\"", ct.DisplayName)
			assert.NotSame(t, ct.Addr, hct.Addr)
			assert.Equal(t, "sip:bob@192.0.2.4", ct.Addr.String())
			assert.Equal(t, ";q=0.5;foo=bar;expires=1800", ct.Params.String())
			assert.Equal(t, "0.5", ct.Q)
			assert.Equal(t, "1800", ct.Expires)
			assert.Nil(t, ct.Next)
		})

		t.Run("multiple linked headers", func(t *testing.T) {
			hdr := "m: <sip:100@192.0.2.4>;expires=800,<sip:100@10.0.0.4:4555>;q=0.5, <sip:100@[2041:0:140F::875B:131B]>"
			msg, err := Parse(toMsg([]string{hdr}))
			assert.Nil(t, err)

			hct, ok := msg.Find(HContact).(*HeaderContact)
			assert.True(t, ok)
			ct, ok := hct.Copy().(*HeaderContact)
			assert.True(t, ok)
			assert.Equal(t, "m", ct.HeaderName)
			assert.Equal(t, "", ct.DisplayName)
			assert.NotSame(t, ct.Addr, hct.Addr)
			assert.Equal(t, "sip:100@192.0.2.4", ct.Addr.String())
			assert.Equal(t, ";expires=800", ct.Params.String())
			assert.Equal(t, "", ct.Q)
			assert.Equal(t, "800", ct.Expires)
			assert.NotNil(t, ct.Next)

			ct = ct.Next
			hct = hct.Next
			assert.Equal(t, "", ct.HeaderName)
			assert.Equal(t, "", ct.DisplayName)
			assert.NotSame(t, ct.Addr, hct.Addr)
			assert.Equal(t, "sip:100@10.0.0.4:4555", ct.Addr.String())
			assert.Equal(t, ";q=0.5", ct.Params.String())
			assert.Equal(t, "0.5", ct.Q)
			assert.Equal(t, "", ct.Expires)
			assert.NotNil(t, ct.Next)

			ct = ct.Next
			hct = hct.Next
			assert.Equal(t, "", ct.HeaderName)
			assert.Equal(t, "", ct.DisplayName)
			assert.NotSame(t, ct.Addr, hct.Addr)
			assert.Equal(t, "sip:100@[2041:0:140F::875B:131B]", ct.Addr.String())
			assert.Equal(t, "", ct.Params.String())
			assert.Equal(t, "", ct.Q)
			assert.Equal(t, "", ct.Expires)
			assert.Nil(t, ct.Next)
		})
	})
}

func TestHeaderViaStringLen(t *testing.T) {
	tests := map[string]struct {
		via  *HeaderVia
		want string
	}{
		`simple via`: {
			&HeaderVia{HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "atlanta.com"},
			"Via: SIP/2.0/UDP atlanta.com",
		},
		`with port`: {
			&HeaderVia{
				HeaderName: "v", Proto: "SIP/2.0/", Transp: "TCP",
				Host: "atlanta.com", Port: "5061",
			},
			"v: SIP/2.0/TCP atlanta.com:5061",
		},
		`with params`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "atlanta.com",
				Params: "branch=ff00aa",
			},
			"Via: SIP/2.0/UDP atlanta.com;branch=ff00aa",
		},
		`with port and params`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.0.0.100",
				Port: "8060", Params: "branch=ff00aa",
			},
			"Via: SIP/2.0/UDP 10.0.0.100:8060;branch=ff00aa",
		},
		`via with one linked`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.100", Params: "branch=z9Hffa",
				Next: &HeaderVia{Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.200", Port: "5067", Params: "branch=z9Hff0"},
			},
			"Via: SIP/2.0/UDP 10.1.1.100;branch=z9Hffa,SIP/2.0/UDP 10.1.1.200:5067;branch=z9Hff0",
		},
		`via with two linked`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.1", Params: "branch=z9Hf.2",
				Next: &HeaderVia{
					Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.2", Params: "branch=z9Hf.1",
					Next: &HeaderVia{Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.3", Params: "branch=z9Hf.0"},
				},
			},
			"Via: SIP/2.0/UDP 10.1.1.1;branch=z9Hf.2,SIP/2.0/UDP 10.1.1.2;branch=z9Hf.1,SIP/2.0/UDP 10.1.1.3;branch=z9Hf.0",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.via.String())
			assert.Equal(t, len(tc.want), tc.via.Len())
		})
	}
}

func BenchmarkHeaderViaString(b *testing.B) {
	via := &HeaderVia{
		HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.1", Params: "branch=z9Hf.2",
		Next: &HeaderVia{
			Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.2", Params: "branch=z9Hf.1",
			Next: &HeaderVia{Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.3", Params: "branch=z9Hf.0"},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = via.String()
	}
}
