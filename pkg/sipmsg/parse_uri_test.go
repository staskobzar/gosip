package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURIValidAddresses(t *testing.T) {
	tests := []struct {
		input                                       string
		scheme, userinfo, hostport, params, headers string
	}{
		{
			"sip:alice@atlanta.com",
			"sip", "alice", "atlanta.com", "", "",
		}, {
			"sip:alice:secretword@atlanta.com;transport=tcp",
			"sip", "alice:secretword", "atlanta.com", "transport=tcp", "",
		}, {
			"sip:unres-d_.d!!w*'(city)'@atlanta.com",
			"sip", "unres-d_.d!!w*'(city)'", "atlanta.com", "", "",
		}, {
			"sips:usr+=in&d$1,2;5?with/here@atlanta.com",
			"sips", "usr+=in&d$1,2;5?with/here", "atlanta.com", "", "",
		}, {
			"sip:alice%20Doe?NY:p-s_d.F!r~y*'j(U)%20%25&Y+=11$p,V@atlanta.com",
			"sip", "alice%20Doe?NY:p-s_d.F!r~y*'j(U)%20%25&Y+=11$p,V", "atlanta.com", "", "",
		}, {
			"sips:alice@atlanta.com?subject=project%20x&priority=urgent",
			"sips", "alice", "atlanta.com", "", "subject=project%20x&priority=urgent",
		}, {
			"sip:+1-212-555-1212:1234@gateway.com;user=phone",
			"sip", "+1-212-555-1212:1234", "gateway.com", "user=phone", "",
		}, {
			"sips:gateway.com",
			"sips", "", "gateway.com", "", "",
		}, {
			"sip:alice@192.0.2.4:8899",
			"sip", "alice", "192.0.2.4:8899", "", "",
		}, {
			"sip:bob@88.123.44.56",
			"sip", "bob", "88.123.44.56", "", "",
		}, {
			"sip:bob+carl@[fe80::9aef:5bad:992a:f54e]",
			"sip", "bob+carl", "[fe80::9aef:5bad:992a:f54e]", "", "",
		}, {
			"sip:[2001:db8::1:0:0:1]:8899;user=phone;time=2000%2004",
			"sip", "", "[2001:db8::1:0:0:1]:8899", "user=phone;time=2000%2004", "",
		}, {
			"sip:john@[2001:db8::1:0:0:1:1:208.8.8.101]:5061?time=now",
			"sip", "john", "[2001:db8::1:0:0:1:1:208.8.8.101]:5061", "", "time=now",
		}, {
			"sip:atlanta.com;method=REGISTER?to=alice%40atlanta.com",
			"sip", "", "atlanta.com", "method=REGISTER", "to=alice%40atlanta.com",
		}, {
			"sips:gateway-s1.com.?param=&foo=bar",
			"sips", "", "gateway-s1.com.", "", "param=&foo=bar",
			// TODO: fix username with ";"
			// }, {
			// 	"sips:alice;day=tuesday@atlanta.com",
			// 	"sips", "alice;day=tuesday", "atlanta.com", "", "",
		}, {
			"sip:vivekg@chair-dnrc.example.com;unknownparam",
			"sip", "vivekg", "chair-dnrc.example.com", "unknownparam", "",
		}, {
			"sip:unres-d_.d!!w*'(city)':(pa0)_.-!~*'*'@[::1:199.9.9.1]:9060;foo[1]=+id",
			"sip", "unres-d_.d!!w*'(city)':(pa0)_.-!~*'*'", "[::1:199.9.9.1]:9060",
			"foo[1]=+id", "",
		}, {
			"sips:atlanta.com;a[1-3].f&$2=z:!*%20;b[2~9]=t('l*8')&f;~!~",
			"sips", "", "atlanta.com", "a[1-3].f&$2=z:!*%20;b[2~9]=t('l*8')&f;~!~", "",
		}, {
			"sip:alice@10.9.8.1:90;a[1-3]=z_%20u!?foo[1-4!3]=&[*bar](3'?r'/0)=$f:0+9",
			"sip", "alice", "10.9.8.1:90", "a[1-3]=z_%20u!",
			"foo[1-4!3]=&[*bar](3'?r'/0)=$f:0+9",
		}, {
			"sips:188.23.100.22:9009",
			"sips", "", "188.23.100.22:9009", "", "",
		},
	}
	for _, tc := range tests {
		uri, err := ParseURI(tc.input)
		assert.Nil(t, err)
		assert.Equal(t, tc.scheme, uri.Scheme)
		assert.Equal(t, tc.userinfo, uri.Userinfo)
		assert.Equal(t, tc.hostport, uri.Hostport)
		assert.Equal(t, tc.params, uri.Params)
		assert.Equal(t, tc.headers, uri.Headers)
	}
}

func TestParseURIFail(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"foo"},
		{"sisp:atlanta.com"},
		{"sip:1.1.1.1:a22"},
		{"sip:atlanta.com;foo\""},
		{"sip:atlanta.com;foo?bar"},
		{"sip:;foo?bar"},
		{"sip:?foo"},
		{"sip:alice%t@example.com"},
		{"sip:alice%1i@example.com"},
		{"sip:alice-[a]@example.com"},
		{"sip:alice:foo//@example.com"},
		{"sip:alice@123.123..33"},
		{"sip:alice@123.123.9999.33"},
		{"sip:alice@12.123.88.33:"},
		{"sip:alice@atlanta.com:123456"},
		{"sip:alice@[2001:db8::1:0:0:1:192.11..33]:8899"},
		{"sip:alice@[2001:db8::1:0:0:1:192.11.0987.33]:8899"},
		{"sip:john@[2001:db8::1:0:0:1:1:208.8.8.101]:506198"},
		{"sip:alice@atlanta.com:123456"},
		{"sip:alice@atlanta.com:1234;"},
		{"sip:alice@atlanta.com:1234;foo="},
		{"sip:alice@atlanta.com:1234;foo=\"bar"},
		{"sip:alice@atlanta.com:1234;foo|=bar"},
		{"sip:alice@atlanta.com?"},
		{"sip:alice@atlanta.com?contact=&here"},
		{"sip:10.0.09:6600;time=now?contact=&here"},
		{"sip:[2001:db8::1:0:0:1:192.11.0987.33]:6600?contact=&here"},
		{"sip::pawd@[::1]"},
		{"sip:alice%20Doe?NY:p-s_d.F!r~y*'j(U)%20%25&Y+=11$p,V@[atlanta.com"},
		{"sip:alice%20Doe?NY:p-s_d.F!r~y*'j(U)%20%25&Y+=11$p,V@[::1]:990088"},
		{"sip:unres-d_.d!!w*'(city)':(pa0)_.-!~*'*' @atlanta.com"},
		{"sip:unres-d_.d!!w*'(city)':(pa0)_.-!~*'*'@[::1:199.9.9.1]:9060;foo[1]=+id;"},
		{"sip:alice@10.9.8.1:90;a[1-3]=z_%20u!?foo[1-4!3]=&[*bar](3'?r'/0)=$f:0+9="},
	}

	for _, tc := range tests {
		uri, err := ParseURI(tc.input)
		assert.NotNil(t, err)
		assert.Nil(t, uri)
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

func BenchmarkParseURI(b *testing.B) {
	b.ResetTimer()
	input := "sip:alice:_pa55w0Rd@biloxi.com:5062;method=REGISTER;transport=tcp?to=sip:bob%40biloxi.com&subject=renew"
	// input := "sips:bob:pa55w0rd@example.com:8080;user=phone?X-t=foo"
	for i := 0; i < b.N; i++ {
		_, _ = ParseURI(input)
	}
}
