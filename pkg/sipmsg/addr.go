package sipmsg

import "strings"

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

// NameAddrSpec is a base structure for other Name-Addr structures
type NameAddrSpec struct {
	HeaderName  string
	DisplayName string
	Addr        *URI
	Params      string
}

// NewNameAddrSpec creates new NameAddrSpec
func NewNameAddrSpec(name string) NameAddrSpec {
	return NameAddrSpec{
		HeaderName: name,
		Addr:       &URI{},
	}
}

// Name returns header name as string and implements AnyHeader
func (naddr *NameAddrSpec) Name() string              { return naddr.HeaderName }
func (naddr *NameAddrSpec) setDisplayName(val string) { naddr.DisplayName = strings.TrimSpace(val) }
func (naddr *NameAddrSpec) setURIScheme(val string)   { naddr.Addr.Scheme = val }
func (naddr *NameAddrSpec) setURIUserinfo(val string) { naddr.Addr.Userinfo = val }
func (naddr *NameAddrSpec) setURIHostport(val string) { naddr.Addr.Hostport = val }
func (naddr *NameAddrSpec) setURIParams(val string)   { naddr.Addr.Params = val }
func (naddr *NameAddrSpec) setURIHeaders(val string)  { naddr.Addr.Headers = val }
func (naddr *NameAddrSpec) setParam(_, val string)    { naddr.Params = val }

// NameAddr headers From and To
type NameAddr struct {
	NameAddrSpec
	T   HType
	Tag string
}

// NewNameAddr creates new NameAddr header
// t must be HFrom or HTo
func NewNameAddr(t HType, name string) *NameAddr {
	return &NameAddr{
		T:            t,
		NameAddrSpec: NewNameAddrSpec(name),
	}
}

// String represents NameAddr header as string
// @impl anyHeader interface
func (naddr *NameAddr) String() string {
	hdr := naddr.HeaderName + ": "
	if len(naddr.DisplayName) > 0 {
		hdr += naddr.DisplayName + " "
	}

	hdr += "<" + naddr.Addr.String() + ">" + naddr.Params

	return hdr
}

// Type returns NameAddr type
// @impl anyHeader interface
func (naddr *NameAddr) Type() HType { return naddr.T }

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

// NewHeaderContact creates HeaderContact
func NewHeaderContact(name string) *HeaderContact {
	return &HeaderContact{
		NameAddrSpec: NewNameAddrSpec(name),
	}
}

// String represents HeaderContact header as string
// @impl anyHeader interface
func (cnt *HeaderContact) String() string {
	var hdr string
	if len(cnt.HeaderName) > 0 {
		hdr = cnt.HeaderName + ": "
	}

	if cnt.Params == "*" {
		return hdr + "*"
	}

	if len(cnt.DisplayName) > 0 {
		hdr += cnt.DisplayName + " "
	}

	hdr += "<" + cnt.Addr.String() + ">" + cnt.Params

	if cnt.Next != nil {
		hdr += "," + cnt.Next.String()
	}
	return hdr
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

// Type returns HContact type
// @impl anyHeader interface
func (cnt *HeaderContact) Type() HType { return HContact }

// Route structure that represence Route or Record-Route headers
type Route struct {
	T HType
	NameAddrSpec
	Next *Route
}

// NewRoute creates new Route header
// parameter "t" must be HRoute or HRecordRoute
func NewRoute(t HType, name string) *Route {
	return &Route{
		T:            t,
		NameAddrSpec: NewNameAddrSpec(name),
	}
}

// String represents Route header as string
// @impl anyHeader interface
func (r *Route) String() string {
	var hdr string
	if len(r.HeaderName) > 0 {
		hdr = r.HeaderName + ": "
	}

	// NOTE: ANBF from RFC3261 provides display name
	// and parameters but I have not seen it match in real life
	if len(r.DisplayName) > 0 {
		hdr += r.DisplayName
	}
	hdr += "<" + r.Addr.String() + ">" + r.Params

	if r.Next != nil {
		hdr += "," + r.Next.String()
	}
	return hdr
}

// Type returns Route type
// @impl anyHeader interface
func (r *Route) Type() HType { return r.T }
