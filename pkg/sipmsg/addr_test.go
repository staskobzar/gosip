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
			{"Via: SIP/2.0/UDP pbx.com ;branch=z9hG4bKnashds7",
				HeaderVia{Name: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "pbx.com",
					Port: "", Branch: "z9hG4bKnashds7", Params: " ;branch=z9hG4bKnashds7"},
			},
			{"VIA: SIP/ 2.0/ TCP 10.0.0.1: 15060; branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120",
				HeaderVia{Name: "VIA", Proto: "SIP/ 2.0/ ", Transp: "TCP", Host: "10.0.0.1",
					Port: "15060", Branch: "z9hG4bKna", Recvd: "10.0.0.100",
					Params: "; branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120"},
			},
			{"v: SIP/2.0/TLS [fe80::2e8d:b1ff:fef3:8a40] :6060;branch=z9hG4bKff;rl",
				HeaderVia{Name: "v", Proto: "SIP/2.0/", Transp: "TLS", Host: "[fe80::2e8d:b1ff:fef3:8a40]",
					Port: "6060", Branch: "z9hG4bKff", Params: ";branch=z9hG4bKff;rl"},
			},
		}
		for _, tc := range tests {
			msg, err := Parse(toMsg([]string{tc.hdr}))
			assert.Nil(t, err)

			via := msg.Find(HVia).(*HeaderVia)

			assert.Equal(t, tc.want.Name, via.Name)
			assert.Equal(t, tc.want.Proto, via.Proto)
			assert.Equal(t, tc.want.Transp, via.Transp)
			assert.Equal(t, tc.want.Host, via.Host)
			assert.Equal(t, tc.want.Port, via.Port)
			assert.Equal(t, tc.want.Branch, via.Branch)
			assert.Equal(t, tc.want.Recvd, via.Recvd)
			assert.Equal(t, tc.want.Params, via.Params)
			assert.Nil(t, via.Via)
		}
	})

	t.Run("linked header", func(t *testing.T) {
		hdr := "Via: SIP/2.0/UDP h1.pbx.com;branch=z9hG4bKnashd" +
			", SIP/2.0/UDP h2.pbx.com;branch=z9hG4bKnasff;received=10.0.0.1" +
			", SIP/2.0/TCP h3.pbx.com:11609;branch=z9hG4bKnasfe;maddr=10.0.0.1"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)

		via := msg.Find(HVia).(*HeaderVia)
		assert.Equal(t, "Via", via.Name)
		assert.Equal(t, "SIP/2.0/", via.Proto)
		assert.Equal(t, "UDP", via.Transp)
		assert.Equal(t, "h1.pbx.com", via.Host)
		assert.Equal(t, "z9hG4bKnashd", via.Branch)

		via = via.Via
		assert.Equal(t, "", via.Name)
		assert.Equal(t, "SIP/2.0/", via.Proto)
		assert.Equal(t, "UDP", via.Transp)
		assert.Equal(t, "h2.pbx.com", via.Host)
		assert.Equal(t, "z9hG4bKnasff", via.Branch)
		assert.Equal(t, "10.0.0.1", via.Recvd)

		via = via.Via
		assert.Equal(t, "TCP", via.Transp)
		assert.Equal(t, "h3.pbx.com", via.Host)
		assert.Equal(t, "z9hG4bKnasfe", via.Branch)
		assert.Equal(t, ";branch=z9hG4bKnasfe;maddr=10.0.0.1", via.Params)
		assert.Nil(t, via.Via)
	})

	t.Run("multiple vias", func(t *testing.T) {
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

func TestParseContactHeaders(t *testing.T) {
	t.Run("single header", func(t *testing.T) {
		hdr := "Contact: \"Bob C\" <sip:bob@192.0.2.4>;q=0.5;foo=bar;expires=1800"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)

		cnt := msg.Find(HContact).(*HeaderContact)
		assert.Equal(t, "Contact", cnt.HeaderName)
		assert.Equal(t, "\"Bob C\"", cnt.DisplayName)
		assert.Equal(t, "sip:bob@192.0.2.4", cnt.Addr.String())
		assert.Equal(t, ";q=0.5;foo=bar;expires=1800", cnt.Params)
		assert.Equal(t, "0.5", cnt.Q)
		assert.Equal(t, "1800", cnt.Expires)
		assert.Nil(t, cnt.Next)
	})

	t.Run("linked header", func(t *testing.T) {
		hdr := "m: <sip:100@192.0.2.4>,<sip:100@10.0.0.4:4555>;q=0.5,<sip:100@[2041:0:140F::875B:131B]>;q=0.8"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)

		cnt := msg.Find(HContact).(*HeaderContact)
		assert.Equal(t, "m", cnt.HeaderName)
		assert.Equal(t, "sip:100@192.0.2.4", cnt.Addr.String())
		assert.Equal(t, "", cnt.Params)

		cnt = cnt.Next
		assert.Equal(t, "sip:100@10.0.0.4:4555", cnt.Addr.String())
		assert.Equal(t, "0.5", cnt.Q)

		cnt = cnt.Next
		assert.Equal(t, "sip:100@[2041:0:140F::875B:131B]", cnt.Addr.String())
		assert.Equal(t, "0.8", cnt.Q)
		assert.Nil(t, cnt.Next)
	})

	t.Run("multiple headers", func(t *testing.T) {
		hdr := "m: <sip:100@192.0.2.4>\r\n" +
			"Contact: <sip:caller@u1.space.com>;q=0.1\r\n" +
			"Contact: <sips:caller@u2.space.com>;q=0.3"

		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)
		list := msg.FindAll(HContact)
		assert.Equal(t, 3, list.Len())
	})
}

func TestParseFromToHeader(t *testing.T) {
	t.Run("parse to member and headers list successfully", func(t *testing.T) {
		hdrs := "To: \"Alice Home\" <sip:alice@biloxi.com>;user=phone;tag=ff00aa\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456248;day=monday;free"
		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		assert.Equal(t, 1, msg.FindAll(HFrom).Len())
		assert.Equal(t, 1, msg.FindAll(HTo).Len())

		from := msg.Find(HFrom).(*NameAddr)
		assert.Same(t, msg.From, from)
		to := msg.Find(HTo).(*NameAddr)
		assert.Same(t, msg.To, to)

		assert.Equal(t, "To", to.HeaderName)
		assert.Equal(t, `"Alice Home"`, to.DisplayName)
		assert.Equal(t, "sip:alice@biloxi.com", to.Addr.String())
		assert.Equal(t, "ff00aa", to.Tag)
		assert.Equal(t, ";user=phone;tag=ff00aa", to.Params)

		assert.Equal(t, "From", from.HeaderName)
		assert.Equal(t, "Bob ", from.DisplayName)
		assert.Equal(t, "sip:bob@biloxi.com", from.Addr.String())
		assert.Equal(t, "456248", from.Tag)
		assert.Equal(t, ";tag=456248;day=monday;free", from.Params)
	})

	t.Run("fail when more then one To header", func(t *testing.T) {
		hdrs := "To: \"Alice Home\" <sip:alice@biloxi.com>;tag=ff00aa\r\n" +
			"t: <sip:alice@biloxi.com>\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456248;day=monday;free"

		_, err := Parse(toMsg([]string{hdrs}))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "more then one To headers found")
	})

	t.Run("fail when more then one From header", func(t *testing.T) {
		hdrs := "To: \"Alice Home\" <sip:alice@biloxi.com>;tag=ff00aa\r\n" +
			"f: <sip:alice@biloxi.com>\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456248;day=monday;free"

		_, err := Parse(toMsg([]string{hdrs}))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "more then one From headers found")
	})
}

func TestParseRoutingHeaders(t *testing.T) {
	t.Run("multiple headers", func(t *testing.T) {
		hdrs := "Record-Route: <sip:h1.domain.com;lr>;host=one\r\n" +
			"Record-Route: H2 <sip:h2.domain.com>\r\n" +
			"Record-Route: <sip:h3.domain.com>\r\n" +
			"Route: <sip:s1.pbx.com>\r\n" +
			"Route: <sip:s2.pbx.com>"

		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		assert.Equal(t, 3, msg.FindAll(HRecordRoute).Len())
		assert.Equal(t, 2, msg.FindAll(HRoute).Len())
	})

	t.Run("linked record-route headers", func(t *testing.T) {
		hdrs := "Record-Route: <sip:h1.domain.com;lr>;host=one\r\n" +
			"Record-Route: <sip:h2.domain.com;lr>,<sip:dd1.pbx.com>;user=pbx,<sips:dd2.pbx.com>\r\n" +
			"Record-Route: <sip:h3.domain.com>"
		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		list := msg.FindAll(HRecordRoute)
		assert.Equal(t, 3, list.Len())
		r := list[1].(*Route)

		assert.Equal(t, "Record-Route", r.HeaderName)
		assert.Equal(t, "sip:h2.domain.com;lr", r.Addr.String())
		assert.Equal(t, "", r.Params)
		assert.NotNil(t, r.Next)

		r = r.Next
		assert.Equal(t, "sip:dd1.pbx.com", r.Addr.String())
		assert.Equal(t, ";user=pbx", r.Params)
		assert.NotNil(t, r.Next)

		r = r.Next
		assert.Equal(t, "sips:dd2.pbx.com", r.Addr.String())
		assert.Equal(t, "", r.Params)
		assert.Nil(t, r.Next)
	})

	t.Run("linked record-route headers", func(t *testing.T) {
		hdrs := "Route: <sip:s1.pbx.com;lr>,<sip:h100.sip.com:5060>;now\r\n" +
			"Route: <sip:s2.pbx.com>"

		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		list := msg.FindAll(HRoute)
		assert.Equal(t, 2, list.Len())
		r := list[0].(*Route)

		assert.Equal(t, "Route", r.HeaderName)
		assert.Equal(t, "sip:s1.pbx.com;lr", r.Addr.String())
		assert.Equal(t, "", r.Params)
		assert.NotNil(t, r.Next)

		r = r.Next
		assert.Equal(t, "sip:h100.sip.com:5060", r.Addr.String())
		assert.Equal(t, ";now", r.Params)
		assert.Nil(t, r.Next)
	})
}
