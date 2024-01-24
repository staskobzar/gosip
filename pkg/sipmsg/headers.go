package sipmsg

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

type anyHeader interface {
	t() HType
	name() string

	String() string
}

type Headers []anyHeader

func (hdrs Headers) Len() int { return len(hdrs) }

type HeaderGeneric struct {
	Type  HType
	Name  string
	Value string
}

func (hg *HeaderGeneric) t() HType     { return hg.Type }
func (hg *HeaderGeneric) name() string { return hg.Name }
func (hg *HeaderGeneric) String() string {
	return hg.Name + ": " + hg.Value
}

// HeaderVia SIP Via header with a pointer to linked
// comma separated list
type HeaderVia struct {
	Name   string
	Proto  string
	Transp string
	Host   string
	Port   string
	Branch string
	Recvd  string
	Params string
	Next   *HeaderVia // linked list for comma separated list of vias in the same header
}

// NewHeaderVia create new header with name
func (msg *Message) NewHeaderVia(name string) *HeaderVia {
	via := &HeaderVia{
		Name: name,
	}
	msg.Headers = append(msg.Headers, via)
	return via
}

func (via *HeaderVia) LinkNext() *HeaderVia {
	linkVia := &HeaderVia{}
	via.Next = linkVia
	return linkVia
}

// String method to build string representation
func (via *HeaderVia) String() string {
	var hdr string
	if len(via.Name) > 0 {
		hdr = via.Name + ": "
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

func (via *HeaderVia) t() HType     { return HVia }
func (via *HeaderVia) name() string { return via.Name }
