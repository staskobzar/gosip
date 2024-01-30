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
					Port: "", Branch: "z9hG4bKnashds7", Params: " ;branch=z9hG4bKnashds7",
				},
			},
			{
				"VIA: SIP/ 2.0/ TCP 10.0.0.1: 15060; branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120",
				HeaderVia{
					HeaderName: "VIA", Proto: "SIP/ 2.0/ ", Transp: "TCP", Host: "10.0.0.1",
					Port: "15060", Branch: "z9hG4bKna", Recvd: "10.0.0.100",
					Params: "; branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120",
				},
			},
			{
				"v: SIP/2.0/TLS [fe80::2e8d:b1ff:fef3:8a40] :6060;branch=z9hG4bKff;rl",
				HeaderVia{
					HeaderName: "v", Proto: "SIP/2.0/", Transp: "TLS", Host: "[fe80::2e8d:b1ff:fef3:8a40]",
					Port: "6060", Branch: "z9hG4bKff", Params: ";branch=z9hG4bKff;rl",
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
		assert.Equal(t, ";branch=z9hG4bKnasfe;maddr=10.0.0.1", via.Params)
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

func TestHeaderViaString(t *testing.T) {
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
				Params: ";branch=ff00aa",
			},
			"Via: SIP/2.0/UDP atlanta.com;branch=ff00aa",
		},
		`with port and params`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.0.0.100",
				Port: "8060", Params: ";branch=ff00aa",
			},
			"Via: SIP/2.0/UDP 10.0.0.100:8060;branch=ff00aa",
		},
		`via with one linked`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.100", Params: ";branch=z9Hffa",
				Next: &HeaderVia{Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.200", Port: "5067", Params: ";branch=z9Hff0"},
			},
			"Via: SIP/2.0/UDP 10.1.1.100;branch=z9Hffa,SIP/2.0/UDP 10.1.1.200:5067;branch=z9Hff0",
		},
		`via with two linked`: {
			&HeaderVia{
				HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.1", Params: ";branch=z9Hf.2",
				Next: &HeaderVia{
					Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.2", Params: ";branch=z9Hf.1",
					Next: &HeaderVia{Proto: "SIP/2.0/", Transp: "UDP", Host: "10.1.1.3", Params: ";branch=z9Hf.0"},
				},
			},
			"Via: SIP/2.0/UDP 10.1.1.1;branch=z9Hf.2,SIP/2.0/UDP 10.1.1.2;branch=z9Hf.1,SIP/2.0/UDP 10.1.1.3;branch=z9Hf.0",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.via.String())
		})
	}
}
