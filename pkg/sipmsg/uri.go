package sipmsg

import "errors"

// errors
var (
	ErrURIParse = errors.New("SIP URI parse")
)

// URI represents SIP URI structure
type URI struct {
	Scheme   string
	Userinfo string
	Hostport string
	Params   string
	Headers  string
}
