package sipmsg

// URI represents SIP URI structure
type URI struct {
	Scheme   string
	Userinfo string
	Hostport string
	Params   string
	Headers  string
}

func (uri *URI) String() string {
	addr := uri.Scheme + ":"
	if len(uri.Userinfo) > 0 {
		addr += uri.Userinfo + "@"
	}
	addr += uri.Hostport
	if len(uri.Params) > 0 {
		addr += ";" + uri.Params
	}
	if len(uri.Headers) > 0 {
		addr += "?" + uri.Headers
	}
	return addr
}
