package sipmsg

import (
	"bytes"
)

// URI represents SIP URI structure
type URI struct {
	Scheme   string
	Userinfo string
	Hostport string
	Params   Params
	Headers  string
}

// String representation of URI
func (uri *URI) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 255))
	buf.WriteString(uri.Scheme)
	buf.WriteByte(':')

	if len(uri.Userinfo) > 0 {
		buf.WriteString(uri.Userinfo)
		buf.WriteByte('@')
	}

	buf.WriteString(uri.Hostport)

	if len(uri.Params) > 0 {
		buf.WriteString(uri.Params.String())
	}

	if len(uri.Headers) > 0 {
		buf.WriteByte('?')
		buf.WriteString(uri.Headers)
	}

	return buf.String()
}
