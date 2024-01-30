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
	Type() HType
	Name() string
	String() string
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

// Type returns HVia type
// @impl anyHeader interface
func (hg *HeaderGeneric) Type() HType { return hg.T }

// Name returns header name as string
// @impl anyHeader interface
func (hg *HeaderGeneric) Name() string { return hg.HeaderName }

// String method to build string representation
// @impl anyHeader interface
func (hg *HeaderGeneric) String() string {
	return hg.HeaderName + ": " + hg.Value
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
	Params     string
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

// LinkNext create *HeaderVia and link it to the caller
func (via *HeaderVia) LinkNext() *HeaderVia {
	linkVia := &HeaderVia{}
	via.Next = linkVia
	return linkVia
}

// String method to build string representation
func (via *HeaderVia) String() string {
	var hdr string
	if len(via.HeaderName) > 0 {
		hdr = via.HeaderName + ": "
	}
	hdr += via.Proto + via.Transp + " " + via.Host

	if len(via.Port) > 0 {
		hdr += ":" + via.Port
	}

	if len(via.Params) > 0 {
		hdr += via.Params
	}

	if via.Next != nil {
		// call linked via build
		return hdr + "," + via.Next.String()
	}

	return hdr
}

// Type returns HVia type
// @impl anyHeader interface
func (via *HeaderVia) Type() HType { return HVia }

// Name returns header name as string
// @impl anyHeader interface
func (via *HeaderVia) Name() string { return via.HeaderName }
