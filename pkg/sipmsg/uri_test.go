package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURILen(t *testing.T) {
	tests := []struct {
		want int
		uri  *URI
	}{
		{13, &URI{Scheme: "sip", Hostport: "localhost"}},
		{26, &URI{Scheme: "sips", Hostport: "10.0.1.100:5050", Userinfo: "alice"}},
		{23, &URI{Scheme: "sip", Hostport: "pbx.com", Userinfo: "bob:p123", Params: "rl"}},
		{34, &URI{Scheme: "sip", Hostport: "pbx.com", Userinfo: "bob:p123", Params: "rl", Headers: "subj=renew"}},
		{25, &URI{Scheme: "sip", Hostport: "pbx.com", Params: "rl", Headers: "subj=renew"}},
		{14, &URI{Scheme: "sip", Hostport: "pbx.com", Params: "rl"}},
		{11, &URI{Scheme: "sip", Hostport: "pbx.com"}},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.want, tc.uri.Len())
	}
}

func TestURIString(t *testing.T) {
	tests := []struct {
		uri string
	}{
		{"sip:atlanta.com"},
		{"sip:alice@atlanta.com"},
		{"sips:atlanta.com;transport=TLS;rl"},
		{"sips:bob@pbx.com;rl?user=Bob"},
		{"sip:pbx.com?phone=polycom&v=vvx"},
		{"sip:100@sip.ca:1223;foo=bar"},
		{"sip:200@10.0.0.1:1223"},
		{"sip:300@[2001:db8::1:0:0:1:1:208.8.8.101]:5061"},
		{"sip:300@[2001:db8::1:0:0:1:1];user=bob?phone=yealink"},
	}

	for _, tc := range tests {
		uri, err := ParseURI(tc.uri)
		assert.Nil(t, err)
		assert.Equal(t, tc.uri, uri.String())
	}
}

func BenchmarkURIString(b *testing.B) {
	uri := &URI{
		Scheme:   "sip",
		Userinfo: "alice:_pa55w0Rd",
		Hostport: "biloxi.com:5062",
		Params:   "method=REGISTER;transport=tcp",
		Headers:  "to=sip:bob%40biloxi.com&subject=renew",
	}

	for i := 0; i < b.N; i++ {
		_ = uri.String()
	}
}

func BenchmarkURIStringify(b *testing.B) {
	uri := &URI{
		Scheme:   "sip",
		Userinfo: "alice:_pa55w0Rd",
		Hostport: "biloxi.com:5062",
		Params:   "method=REGISTER;transport=tcp",
		Headers:  "to=sip:bob%40biloxi.com&subject=renew",
	}
	buf := NewStringer(uri.Len())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uri.Stringify(buf)
		buf.buf.Reset()
		// buf.buf = buf.buf[:0]
	}
}
