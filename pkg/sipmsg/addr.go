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
	Params      Params
}

// NewNameAddrSpec creates new NameAddrSpec
func NewNameAddrSpec(name string) NameAddrSpec {
	return NameAddrSpec{
		HeaderName: name,
		Addr:       &URI{},
	}
}

// Len returns size of the NameAddrSpec length as a string
// @impl AnyHeader interface
func (naddr *NameAddrSpec) Len() int {
	l := 0
	if len(naddr.HeaderName) > 0 {
		l += len(naddr.HeaderName) + 2
	}
	if len(naddr.DisplayName) > 0 {
		l += len(naddr.DisplayName) + 1
	}

	l += naddr.Addr.Len() + 2 // wrapped in <>
	if naddr.Params.Len() > 0 {
		l += naddr.Params.Len() + 1
	}
	return l
}

// Name returns header name as string and implements AnyHeader
func (naddr *NameAddrSpec) Name() string { return naddr.HeaderName }

// PrintAddr print to buffer name adders as <address>;params
func (naddr *NameAddrSpec) PrintAddr(buf *Stringer) {
	if len(naddr.DisplayName) > 0 {
		buf.Print(naddr.DisplayName, " ")
	}
	buf.Print("<")
	naddr.Addr.Stringify(buf)
	buf.Print(">", naddr.Params.String())
}

func (naddr *NameAddrSpec) setDisplayName(val string) { naddr.DisplayName = strings.TrimSpace(val) }
func (naddr *NameAddrSpec) setURIScheme(val string)   { naddr.Addr.Scheme = val }
func (naddr *NameAddrSpec) setURIUserinfo(val string) { naddr.Addr.Userinfo = val }
func (naddr *NameAddrSpec) setURIHostport(val string) { naddr.Addr.Hostport = val }
func (naddr *NameAddrSpec) setURIParams(val string)   { naddr.Addr.Params = Params(val).setup() }
func (naddr *NameAddrSpec) setURIHeaders(val string)  { naddr.Addr.Headers = val }
func (naddr *NameAddrSpec) setParam(_, val string)    { naddr.Params = Params(val).setup() }

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

// SetTag sets tag parameter to NameAddr header
func (naddr *NameAddr) SetTag(tag string) {
	naddr.Tag = tag
	if _, ok := naddr.Params.Get("tag"); ok {
		naddr.Params = naddr.Params.Set("tag", tag)
		return
	}
	naddr.Params = naddr.Params.Add("tag", tag)
}

// String represents NameAddr header as string
// @impl anyHeader interface
func (naddr *NameAddr) String() string {
	buf := NewStringer(naddr.Len())
	naddr.Stringify(buf)
	return buf.String()
}

// Stringify puts NameAddr as a string into Stringer buffer
// @impl anyHeader interface
func (naddr *NameAddr) Stringify(buf *Stringer) {
	buf.Print(naddr.HeaderName, ": ")

	naddr.PrintAddr(buf)
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
		naddr.Params = Params(val).setup()
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

// Len returns size of the HeaderContact length as a string
// @impl AnyHeader interface
func (cnt *HeaderContact) Len() int {
	if cnt.Params == "*" {
		return len(cnt.HeaderName) + 3
	}
	l := cnt.NameAddrSpec.Len()

	if cnt.Next != nil {
		return l + 1 + cnt.Next.Len() // +1 for ,
	}
	return l
}

// String represents HeaderContact header as string
// @impl AnyHeader interface
func (cnt *HeaderContact) String() string {
	buf := NewStringer(cnt.Len())
	cnt.Stringify(buf)
	return buf.String()
}

// Stringify puts HeaderContact as a string into Stringer buffer
// @impl AnyHeader interface
func (cnt *HeaderContact) Stringify(buf *Stringer) {
	if len(cnt.HeaderName) > 0 {
		buf.Print(cnt.HeaderName, ": ")
	}

	if cnt.Params == "*" {
		buf.Print("*")
		return
	}

	cnt.PrintAddr(buf)

	if cnt.Next != nil {
		buf.Print(",")
		cnt.Next.Stringify(buf)
	}
}

// Type returns HContact type
// @impl AnyHeader interface
func (cnt *HeaderContact) Type() HType { return HContact }

// override method from NameAddrSpec
func (cnt *HeaderContact) setParam(name, val string) {
	switch name {
	case "q":
		cnt.Q = val
	case "expires":
		cnt.Expires = val
	case "params":
		cnt.Params = Params(val).setup()
	}
}

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

// Len returns size of the HeaderContact length as a string
// @impl AnyHeader interface
func (r *Route) Len() int {
	l := r.NameAddrSpec.Len()
	if r.Next != nil {
		return l + 1 + r.Next.Len()
	}
	return l
}

// String represents Route header as string
// @impl anyHeader interface
func (r *Route) String() string {
	buf := NewStringer(r.Len())
	r.Stringify(buf)
	return buf.String()
}

// Stringify puts HeaderContact as a string into Stringer buffer
// @impl AnyHeader interface
func (r *Route) Stringify(buf *Stringer) {
	if len(r.HeaderName) > 0 {
		buf.Print(r.HeaderName, ": ")
	}

	r.PrintAddr(buf)

	if r.Next != nil {
		buf.Print(",")
		r.Next.Stringify(buf)
	}
}

// Type returns Route type
// @impl anyHeader interface
func (r *Route) Type() HType { return r.T }
