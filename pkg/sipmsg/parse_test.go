package sipmsg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toMsg(hdrs []string) string {
	return "NOTIFY sip:bob@biloxi.com SIP/2.0\r\n" +
		strings.Join(hdrs, "\r\n") + "\r\n\r\n"
}

func TestParseRequestToMessage(t *testing.T) {
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: \"Hello World\" <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n" +
		"Content-Length: 0\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, "REGISTER", msg.Method)
	assert.Equal(t, HRequest, msg.t)
	assert.Equal(t, "sip:registrar.biloxi.com", msg.RURI.String())
	assert.Empty(t, msg.Code)
	assert.Empty(t, msg.Reason)
	assert.Equal(t, "843817637684230@998sdasdh09", msg.CallID)
	assert.Equal(t, 1826, msg.CSeq)
	assert.Equal(t, 70, msg.MaxFwd)
	assert.Equal(t, "456248", msg.From.Tag)
	assert.Equal(t, "sip:bob@biloxi.com", msg.To.Addr.String())
	assert.Equal(t, 9, msg.HLen())
	assert.Empty(t, msg.Body)
	assert.False(t, msg.IsResponse())

	t.Run("error", func(t *testing.T) {
		input := "REGISTER sip:registrar.biloxi.com FOO/2.0\r\n" +
			"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
			"Max-Forwards: 70\r\n" +
			"To: \"Hello World\" <sip:bob@biloxi.com>\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
			"Call-ID: 843817637684230@998sdasdh09\r\n" +
			"CSeq: 1826 REGISTER\r\n" +
			"Contact: <sip:bob@192.0.2.4>\r\n" +
			"Expires: 7200\r\n" +
			"Content-Length: 0\r\n\r\n"
		_, err := Parse(input)
		assert.ErrorContains(t, err, `gistrar.biloxi.com ">>>"FOO/2.0`)
	})
}

func TestParseResponseToMessage(t *testing.T) {
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
	assert.Equal(t, "REGISTER", msg.Method)
	assert.Equal(t, HResponse, msg.t)
	assert.Empty(t, msg.RURI)
	assert.Equal(t, "200", msg.Code)
	assert.Equal(t, "OK", msg.Reason)
	assert.True(t, msg.IsResponse())
}

func TestParseStoreHeadersWithType(t *testing.T) {
	tests := []struct {
		hdr  string
		want HType
	}{
		{"X-Foo: foo\r\n", HGeneric},
		{"Record-Route: <sip:h1.domain.com;lr>;host=one\r\n", HRecordRoute},
		{"Route: <sip:p1.example.com;lr>,<sip:p2.domain.com;lr>\r\n", HRoute},
		{"Max-Forwards: 70\r\n", HMaxForwards},
		{"To: \"Hello Dolly\" <sip:bob@biloxi.com>\r\n", HTo},
		{"From: Bob <sip:bob@biloxi.com?user=phone>;tag=456248\r\n", HFrom},
		{"Call-ID: 843817637684230@998sdasdh09\r\n", HCallID},
		{"CSeq: 1826 REGISTER\r\n", HCSeq},
		{"Accept \t: */*\r\n", HAccept},
		{"Contact: <sip:bob@192.0.2.4>;q=0.123;expires=30;foo=bar, Alice <sip:alice@192.0.1.1:5066>\r\n", HContact},
		{"Expires: 7200\r\n", HExpires},
		{"Content-Type: application/simple-message-summary ; foo/bar;xy/123d\r\n", HContentType},
		{"Accept-Encoding: gzip\r\n", HAcceptEncoding},
		{"Accept-Language: da, en-gb;q=0.8, en;q=0.7\r\n", HAcceptLanguage},
		{"Alert-Info: <http://www.example.com/sounds/moo.wav>\r\n", HAlertInfo},
		{"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n", HAllow},
		{"Authentication-Info: nextnonce=\"47364c23432d2e131a5fb210812c\"\r\n", HAuthenticationInfo},
		{`Authorization: Digest username="Alice", realm="atlanta.com", nonce="84a4cc6f3082121f32b42a2187831a9e",` +
			"response=\"7587245234b3434cc3412213e5f113a5432\"\r\n", HAuthorization},
		{"Call-Info: <http://wwww.example.com/alice/photo.jpg> ;purpose=icon," +
			"<http://www.example.com/alice/> ;purpose=info\r\n", HCallInfo},
		{"Content-Disposition: session\r\n", HContentDisposition},
		{"Content-Encoding: gzip\r\n", HContentEncoding},
		{"Content-Language: fr\r\n", HContentLanguage},
		{"Date: Sat, 13 Nov 2010 23:29:00 GMT\r\n", HDate},
		{"Error-Info: <sip:not-in-service-recording@atlanta.com>\r\n", HErrorInfo},
		{"In-Reply-To: 70710@saturn.bell-tel.com, 17320@saturn.bell-tel.com\r\n", HInReplyTo},
		{"MIME-Version: 1.0\r\n", HMIMEVersion},
		{"Min-Expires: 60\r\n", HMinExpires},
		{"Organization: Boxes by Bob\r\n", HOrganization},
		{"Priority: emergency\r\n", HPriority},
		{`Proxy-Authenticate: Digest realm="atlanta.com", domain="sip:ss1.carrier.com", qop="auth",` +
			`nonce="f84f1cec41e6cbe5aea9c8e88d359", opaque="", stale=FALSE, algorithm=MD5` + "\r\n",
			HProxyAuthenticate},
		{`Proxy-Authorization: Digest username="Alice", realm="atlanta.com",` +
			`nonce="c60f3082ee1212b402a21831ae", response="245f23415f11432b3434341c022"` + "\r\n",
			HProxyAuthorization},
		{"Proxy-Require: foo, bar\r\n", HProxyRequire},
		{"Reply-To: Bob <sip:bob@biloxi.com>\r\n", HReplyTo},
		{"Require: 100rel\r\n", HRequire},
		{"Retry-After: 18000;duration=3600\r\n", HRetryAfter},
		{"Server: HomeServer v2\r\n", HServer},
		{"Subject: A tornado is heading our way!\r\n", HSubject},
		{"Supported: 100rel\r\n", HSupported},
		{"Timestamp: 54\r\n", HTimestamp},
		{"Unsupported: foo\r\n", HUnsupported},
		{"User-Agent: Softphone Beta1.5\r\n", HUserAgent},
		{"Warning: 307 isi.edu \"Session parameter 'foo' not understood\"\r\n", HWarning},
		{`WWW-Authenticate: Digest realm="atlanta.com", domain="sip:ss1.carrier.com", qop="auth",` +
			`nonce="f84f1cec41e6cbe5aea9c8e88d359", opaque="", stale=FALSE, algorithm=MD5` + "\r\n",
			HWWWAuthenticate},
		{"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n", HVia},
		{"Content-Length: 1234\r\n", HContentLength},
	}
	for _, tc := range tests {
		msg, err := Parse(toMsg([]string{tc.hdr}))
		assert.Nil(t, err)
		h := msg.Find(tc.want)
		assert.Equal(t, tc.want, h.t())
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
			"Bob <sip:bob@atlanta.com>", "Bob", "sip:bob@atlanta.com", "", "",
		},
		`with display name as quoted string`: {
			"\"Big Co.\" <sip:big@sip.com>", "\"Big Co.\"", "sip:big@sip.com", "", "",
		},
		`with params`: {
			"<sip:100@sip.com>;user=alice;x-foo", "", "sip:100@sip.com", "", ";user=alice;x-foo",
		},
		`with tag`: {
			"Foo <sip:100@sip.com>;tag=2493k59kd", "Foo", "sip:100@sip.com", "2493k59kd", ";tag=2493k59kd",
		},
		`with tag and params`: {
			"<sip:100@sip.com>;tag=2493k59kd;user=phone", "", "sip:100@sip.com", "2493k59kd", ";tag=2493k59kd;user=phone",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := Parse(toMsg([]string{"From: " + tc.from}))
			assert.Nil(t, err)
			assert.Equal(t, tc.name, msg.From.DisplayName)
			assert.Equal(t, tc.addr, msg.From.Addr.String())
			assert.Equal(t, tc.tag, msg.From.Tag)
			assert.Equal(t, tc.prms, msg.From.Params)
		})
	}
}

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
