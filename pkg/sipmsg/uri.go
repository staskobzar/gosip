package sipmsg

// URI represents SIP URI structure
type URI struct {
	Scheme    string
	Userinfo  string
	Hostport  string
	Params    Params
	Transport string
	Headers   string
}

// String representation of URI
func (uri *URI) String() string {
	buf := NewStringer(uri.Len())
	uri.Stringify(buf)
	return buf.String()
}

// Stringify puts uri as a string into Stringer buffer
func (uri *URI) Stringify(buf *Stringer) {
	buf.Print(uri.Scheme, ":")

	if len(uri.Userinfo) > 0 {
		buf.Print(uri.Userinfo, "@")
	}

	buf.Print(uri.Hostport)

	if len(uri.Params) > 0 {
		buf.Print(uri.Params.String())
	}

	if len(uri.Headers) > 0 {
		buf.Print("?", uri.Headers)
	}
}

// Len returns length of the URI string
func (uri *URI) Len() int {
	l := len(uri.Scheme) + 1 // scheme semicolon

	if len(uri.Userinfo) > 0 {
		l += len(uri.Userinfo) + 1
	}

	l += len(uri.Hostport)

	if uri.Params.Len() > 0 {
		l += uri.Params.Len() + 1
	}

	if len(uri.Headers) > 0 {
		l += len(uri.Headers) + 1
	}
	return l
}

func (uri *URI) Copy() *URI {
	return &URI{
		Scheme:    uri.Scheme,
		Userinfo:  uri.Userinfo,
		Hostport:  uri.Hostport,
		Params:    uri.Params,
		Transport: uri.Transport,
		Headers:   uri.Headers,
	}
}
