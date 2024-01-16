package sipmsg

// parser helper interface
type nameAddr interface {
	setDisplayName(string)
	setURIScheme(string)
	setURIUserinfo(string)
	setURIHostport(string)
	setURIParams(string)
	setURIHeaders(string)
	setParam(string, string)
}

type NameAddrSpec struct {
	HeaderName  string
	DisplayName string
	Addr        *URI
	Params      string
}

func NewNameAddrSpec(hdrName string) NameAddrSpec {
	return NameAddrSpec{
		HeaderName: hdrName,
		Addr:       &URI{},
	}
}

func (naddr *NameAddrSpec) setDisplayName(val string) { naddr.DisplayName = val }
func (naddr *NameAddrSpec) setURIScheme(val string)   { naddr.Addr.Scheme = val }
func (naddr *NameAddrSpec) setURIUserinfo(val string) { naddr.Addr.Userinfo = val }
func (naddr *NameAddrSpec) setURIHostport(val string) { naddr.Addr.Hostport = val }
func (naddr *NameAddrSpec) setURIParams(val string)   { naddr.Addr.Params = val }
func (naddr *NameAddrSpec) setURIHeaders(val string)  { naddr.Addr.Headers = val }
func (naddr *NameAddrSpec) setParam(_, val string)    { naddr.Params = val }

type NameAddr struct {
	NameAddrSpec
	Tag string
}

func NewNameAddr(name string) *NameAddr {
	return &NameAddr{
		NameAddrSpec: NewNameAddrSpec(name),
	}
}

// override method from NameAddrSpec
func (naddr *NameAddr) setParam(name, val string) {
	switch name {
	case "tag":
		naddr.Tag = val
	case "params":
		naddr.Params = val
	}
}

// HeaderContact represents Contact SIP header single value
// Next element provides link to any other contact separated by comma
// If Contact header value is "*" (STAR) then it will have only
// Param element with value "*". Other elements are nil or empty
type HeaderContact struct {
	NameAddrSpec
	Q       string
	Expires string
	Next    *HeaderContact
}

func NewHeaderContact(name string) *HeaderContact {
	return &HeaderContact{
		NameAddrSpec: NewNameAddrSpec(name),
	}
}

// override method from NameAddrSpec
func (cnt *HeaderContact) setParam(name, val string) {
	switch name {
	case "q":
		cnt.Q = val
	case "expires":
		cnt.Expires = val
	case "params":
		cnt.Params = val
	}
}

type Route struct {
	NameAddrSpec
	Next *Route
}

func NewRoute(name string) *Route {
	return &Route{
		NameAddrSpec: NewNameAddrSpec(name),
	}
}
