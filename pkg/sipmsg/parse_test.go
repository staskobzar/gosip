package sipmsg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequest(t *testing.T) {
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Via: SIP / 2.0 / UDP first.example.com: 4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1\r\n" +
		"Record-Route: <sip:h1.domain.com;lr>;host=one\r\n" +
		"Record-Route: <sip:h2.domain.com;lr>;host=two\r\n" +
		"Route: <sip:p1.example.com;lr>,<sip:p2.domain.com;lr>\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: \"Hello Dolly\" <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Accept \t: */*\r\n" +
		"Contact: <sip:bob@192.0.2.4>;q=0.123;expires=30;foo=bar, Alice <sip:alice@192.0.1.1:5066>\r\n" +
		"m : sips:0001@atlanta.com\r\n" +
		"Expires: 7200\r\n" +
		"Content-Type: application/simple-message-summary ; foo/bar;xy/123d\r\n" +
		"Content-Length: 1234\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)
	fmt.Printf("%#v\n", msg)
	fmt.Printf("%#v\n", msg.RURI)
	for _, h := range msg.Headers {
		fmt.Printf("%#v\n", h)
	}
	for _, v := range msg.Via {
		fmt.Printf("[v] %#v\n", v)
	}
	fmt.Printf("[f] uri: %s from %#v\n", msg.From.Addr, msg.From)
	fmt.Printf("[t] uri: %s   to %#v\n", msg.To.Addr, msg.To)
	for _, c := range msg.Contact {
		fmt.Printf("[m] uri: %s   to %#v\n", c.Addr, c)
	}

	for _, rr := range msg.RecRoute {
		fmt.Printf("[rr] uri: %s rr %#v\n", rr.Addr, rr)
	}

	for r := msg.Route[0]; r != nil; r = r.Next {
		fmt.Printf("[r] uri: %s r %#v\n", r.Addr, r)
	}
}

func TestParseResponse(t *testing.T) {
	input := "SIP/2.0 200 OK\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7;received=192.0.2.4\r\n" +
		"To: Bob <sip:bob@biloxi.com>;tag=2493k59kd\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n" +
		"Content-Length: 0\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)
	fmt.Printf("%#v\n", msg)
}
func TestParseHeaders(t *testing.T) {
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Record-Route: <sip:h1.domain.com;lr>;host=one\r\n" +
		"Route: <sip:p1.example.com;lr>,<sip:p2.domain.com;lr>\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: \"Hello Dolly\" <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Accept \t: */*\r\n" +
		"Contact: <sip:bob@192.0.2.4>;q=0.123;expires=30;foo=bar, Alice <sip:alice@192.0.1.1:5066>\r\n" +
		"Expires: 7200\r\n" +
		"Content-Type: application/simple-message-summary ; foo/bar;xy/123d\r\n" +
		"X-Foo: foo\r\nX-Bar: bar\r\n" +
		"Accept-Encoding: gzip\r\n" +
		"Accept-Language: da, en-gb;q=0.8, en;q=0.7\r\n" +
		"Alert-Info: <http://www.example.com/sounds/moo.wav>\r\n" +
		"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
		"Authentication-Info: nextnonce=\"47364c23432d2e131a5fb210812c\"\r\n" +
		`Authorization: Digest username="Alice", realm="atlanta.com", nonce="84a4cc6f3082121f32b42a2187831a9e",` +
		"response=\"7587245234b3434cc3412213e5f113a5432\"\r\n" +
		"Call-Info: <http://wwww.example.com/alice/photo.jpg> ;purpose=icon," +
		"<http://www.example.com/alice/> ;purpose=info\r\n" +
		"Content-Disposition: session\r\n" +
		"Content-Encoding: gzip\r\n" +
		"Content-Language: fr\r\n" +
		"Date: Sat, 13 Nov 2010 23:29:00 GMT\r\n" +
		"Error-Info: <sip:not-in-service-recording@atlanta.com>\r\n" +
		"In-Reply-To: 70710@saturn.bell-tel.com, 17320@saturn.bell-tel.com\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Min-Expires: 60\r\n" +
		"Organization: Boxes by Bob\r\n" +
		"Priority: emergency\r\n" +
		`Proxy-Authenticate: Digest realm="atlanta.com", domain="sip:ss1.carrier.com", qop="auth",` +
		`nonce="f84f1cec41e6cbe5aea9c8e88d359", opaque="", stale=FALSE, algorithm=MD5` + "\r\n" +
		`Proxy-Authorization: Digest username="Alice", realm="atlanta.com",` +
		`nonce="c60f3082ee1212b402a21831ae", response="245f23415f11432b3434341c022"` + "\r\n" +
		"Proxy-Require: foo, bar\r\n" +
		"Reply-To: Bob <sip:bob@biloxi.com>\r\n" +
		"Require: 100rel\r\n" +
		"Retry-After: 18000;duration=3600\r\n" +
		"Server: HomeServer v2\r\n" +
		"Subject: A tornado is heading our way!\r\n" +
		"Supported: 100rel\r\n" +
		"Timestamp: 54\r\n" +
		"Unsupported: foo\r\n" +
		"User-Agent: Softphone Beta1.5\r\n" +
		"Warning: 307 isi.edu \"Session parameter 'foo' not understood\"\r\n" +
		`WWW-Authenticate: Digest realm="atlanta.com", domain="sip:ss1.carrier.com", qop="auth",` +
		`nonce="f84f1cec41e6cbe5aea9c8e88d359", opaque="", stale=FALSE, algorithm=MD5` + "\r\n" +
		"Content-Length: 1234\r\n\r\n"
	msg, err := Parse(input)
	assert.Nil(t, err)
	for _, h := range msg.Headers {
		fmt.Printf("%#v\n", h)
	}
}

func TestParseHeaderNameAddr(t *testing.T) {
	tests := map[string]struct {
		from                  string
		name, addr, tag, prms string
	}{
		`only address`: {
			"sip:alice@atlanta.com", "", "sip:alice@atlanta.com", "", "",
		},
		`only address with tag`: {
			"sip:100@pbx.com ;tag=0axff34", "", "sip:100@pbx.com", "0axff34", " ;tag=0axff34",
		},
		`with display name`: {
			"Bob <sip:bob@atlanta.com>", "Bob ", "sip:bob@atlanta.com", "", "",
		},
		`with display name as quoted string`: {
			"\"Big Co.\" <sip:big@sip.com>", "\"Big Co.\"", "sip:big@sip.com", "", "",
		},
		`with params`: {
			"<sip:100@sip.com>;user=alice;x-foo", "", "sip:100@sip.com", "", ";user=alice;x-foo",
		},
		`with tag`: {
			"Foo <sip:100@sip.com>;tag=2493k59kd", "Foo ", "sip:100@sip.com", "2493k59kd", ";tag=2493k59kd",
		},
		`with tag and params`: {
			"<sip:100@sip.com>;tag=2493k59kd;user=phone", "", "sip:100@sip.com", "2493k59kd", ";tag=2493k59kd;user=phone",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
				"From: " + tc.from + "\r\n\r\n"
			msg, err := Parse(input)
			assert.Nil(t, err)
			assert.Equal(t, tc.name, msg.From.DisplayName)
			assert.Equal(t, tc.addr, msg.From.Addr.String())
			assert.Equal(t, tc.tag, msg.From.Tag)
			assert.Equal(t, tc.prms, msg.From.Params)
		})
	}
}

// TODO: test Via
// TODO: test Contact
// TODO: test From/To (single etc)
// TODO: test Route and Record-Route
// TODO: test Generic headers

func BenchmarkParse(b *testing.B) {
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Via: SIP / 2.0 / UDP first.example.com: 4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Accept \t: */*\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n" +
		"Content-Type: application/simple-message-summary ; foo/bar;xy/123d\r\n" +
		"Content-Length: 1234\r\n\r\n"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}
