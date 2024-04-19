package sipmsg

// HType header type
type HType uint8

// headers' types enum
const (
	HGeneric HType = iota
	HRequest
	HResponse
	HAccept
	HAcceptEncoding
	HAcceptLanguage
	HAlertInfo
	HAllow
	HAuthenticationInfo
	HAuthorization
	HCallID
	HCallInfo
	HContact
	HContentDisposition
	HContentEncoding
	HContentLanguage
	HContentLength
	HContentType
	HCSeq
	HDate
	HErrorInfo
	HExpires
	HFrom
	HInReplyTo
	HMaxForwards
	HMIMEVersion
	HMinExpires
	HOrganization
	HPriority
	HProxyAuthenticate
	HProxyAuthorization
	HProxyRequire
	HRecordRoute
	HReplyTo
	HRequire
	HRetryAfter
	HRoute
	HServer
	HSubject
	HSupported
	HTimestamp
	HTo
	HUnsupported
	HUserAgent
	HVia
	HWarning
	HWWWAuthenticate
)

// AnyHeader interface for any SIP header in Message
type AnyHeader interface {
	Copy() AnyHeader
	Len() int
	Name() string
	String() string
	Stringify(*Stringer)
	Type() HType
}

// Headers list
type Headers []AnyHeader

// Len headers list length
func (hdrs Headers) Len() int { return len(hdrs) }

// HeaderGeneric generic SIP message header
type HeaderGeneric struct {
	T          HType
	HeaderName string
	Value      string
}

// Copy HeaderGeneric and returns new header pointer
func (hg *HeaderGeneric) Copy() AnyHeader {
	return &HeaderGeneric{
		T:          hg.T,
		HeaderName: hg.HeaderName,
		Value:      hg.Value,
	}
}

// Type returns HVia type
// @impl AnyHeader interface
func (hg *HeaderGeneric) Type() HType { return hg.T }

// Name returns header name as string
// @impl AnyHeader interface
func (hg *HeaderGeneric) Name() string { return hg.HeaderName }

// String method to build string representation
// @impl AnyHeader interface
func (hg *HeaderGeneric) String() string {
	buf := NewStringer(hg.Len())
	hg.Stringify(buf)
	return buf.String()
}

// Stringify puts HeaderGeneric as a string into Stringer buffer
// @impl AnyHeader interface
func (hg *HeaderGeneric) Stringify(buf *Stringer) {
	buf.Print(hg.HeaderName, ": ", hg.Value)
}

// Len returns size of the HeaderGeneric length as a string
// @impl AnyHeader interface
func (hg *HeaderGeneric) Len() int {
	return len(hg.HeaderName) + len(hg.Value) + 2 // +2 is for : and space after name
}

// HeaderVia SIP Via header with a pointer to linked
// comma separated list
type HeaderVia struct {
	HeaderName string
	Proto      string
	Transp     string
	Host       string
	Port       string
	Branch     string
	Recvd      string
	Params     Params
	Next       *HeaderVia // linked list for comma separated list of vias in the same header
}

// NewHeaderVia create new header with name
func (msg *Message) NewHeaderVia(name string) *HeaderVia {
	via := &HeaderVia{
		HeaderName: name,
	}
	msg.Headers = append(msg.Headers, via)
	return via
}

// Copy create copy of the Via header and return its pointer
func (via *HeaderVia) Copy() AnyHeader {
	return via.copy()
}

func (via *HeaderVia) copy() *HeaderVia {
	v := &HeaderVia{
		HeaderName: via.HeaderName,
		Proto:      via.Proto,
		Transp:     via.Transp,
		Host:       via.Host,
		Port:       via.Port,
		Branch:     via.Branch,
		Recvd:      via.Recvd,
		Params:     via.Params,
	}
	if via.Next == nil {
		return v
	}
	v.Next = via.Next.copy() // recursively copy all linked headers

	return v
}

// LinkNext create *HeaderVia and link it to the caller
func (via *HeaderVia) LinkNext() *HeaderVia {
	linkVia := &HeaderVia{}
	via.Next = linkVia
	return linkVia
}

// String method to build string representation
func (via *HeaderVia) String() string {
	buf := NewStringer(via.Len())
	via.Stringify(buf)
	return buf.String()
}

// Stringify push Via header as a string into Stringer buffer
func (via *HeaderVia) Stringify(buf *Stringer) {
	if len(via.HeaderName) > 0 {
		buf.Print(via.HeaderName, ": ")
	}
	buf.Print(via.Proto, via.Transp, " ", via.Host)

	if len(via.Port) > 0 {
		buf.Print(":", via.Port)
	}

	if via.Params.Len() > 0 {
		buf.Print(via.Params.String())
	}

	if via.Next != nil {
		// call linked via build
		buf.Print(",")
		via.Next.Stringify(buf)
	}
}

// Len returns size of the Via header as a string
func (via *HeaderVia) Len() int {
	l := 0
	if len(via.HeaderName) > 0 {
		l += len(via.HeaderName) + 2 // name and colon with space
	}

	l += len(via.Proto) + len(via.Transp) + len(via.Host) + 1

	if len(via.Port) > 0 {
		l += len(via.Port) + 1 // port and colon
	}
	if via.Params.Len() > 0 {
		l += via.Params.Len() + 1 // semi and params
	}
	if via.Next != nil {
		return l + via.Next.Len() + 1 // coma and extra via header
	}
	return l
}

// Type returns HVia type
// @impl anyHeader interface
func (via *HeaderVia) Type() HType { return HVia }

// Name returns header name as string
// @impl anyHeader interface
func (via *HeaderVia) Name() string { return via.HeaderName }
