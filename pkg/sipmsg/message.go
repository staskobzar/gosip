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

type Message struct {
	t      HType
	Method string
	RURI   *URI
	Code   string
	Reason string

	CallID      string
	ContentLen  string
	ContentType string
	CSeq        string
	MaxFwd      string
	Expires     string
	Via         []*HeaderVia
	From        *NameAddr
	To          *NameAddr
	Contact     []*HeaderContact
	RecRoute    []*Route
	Route       []*Route

	Headers Headers
	Body    string
}

func NewMessage() *Message {
	return &Message{
		Headers:  make(Headers, 0, 32),
		Via:      make([]*HeaderVia, 0, 1),     // at least one Via header expected
		Contact:  make([]*HeaderContact, 0, 1), // at least one Contact header expected
		RecRoute: make([]*Route, 0),
		Route:    make([]*Route, 0),
	}
}

func (msg *Message) pushHeader(t HType, name, value string) {
	generic := &HeaderGeneric{
		Type:  t,
		Name:  name,
		Value: value,
	}
	msg.Headers = append(msg.Headers, generic)
}

type Headers []any

type HeaderGeneric struct {
	Type  HType
	Name  string
	Value string
}

type HeaderVia struct {
	Name   string
	Proto  string
	Transp string
	Host   string
	Port   string
	Branch string
	Recvd  string
	Params string
	Via    *HeaderVia // linked list for comma separated list of vias in the same header
}

func NewHeaderVia(name string) *HeaderVia {
	return &HeaderVia{
		Name: name,
	}
}
