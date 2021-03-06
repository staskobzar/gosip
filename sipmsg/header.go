package sipmsg

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// HdrType type header ID
type HdrType int

// SIP Header identifiers
const (
	MsgEOF HdrType = iota
	SIPHdrGeneric
	SIPHdrRequestLine
	SIPHdrStatusLine
	SIPHdrAccept
	SIPHdrAcceptEncoding
	SIPHdrAcceptLanguage
	SIPHdrAlertInfo
	SIPHdrAllow
	SIPHdrAuthenticationInfo
	SIPHdrAuthorization
	SIPHdrCallID
	SIPHdrCallInfo
	SIPHdrContact
	SIPHdrContentDisposition
	SIPHdrContentEncoding
	SIPHdrContentLanguage
	SIPHdrContentLength
	SIPHdrContentType
	SIPHdrCSeq
	SIPHdrDate
	SIPHdrErrorInfo
	SIPHdrExpires
	SIPHdrFrom
	SIPHdrInReplyTo
	SIPHdrMaxForwards
	SIPHdrMIMEVersion
	SIPHdrMinExpires
	SIPHdrOrganization
	SIPHdrPriority
	SIPHdrProxyAuthenticate
	SIPHdrProxyAuthorization
	SIPHdrProxyRequire
	SIPHdrRecordRoute
	SIPHdrReplyTo
	SIPHdrRequire
	SIPHdrRetryAfter
	SIPHdrRoute
	SIPHdrServer
	SIPHdrSubject
	SIPHdrSupported
	SIPHdrTimestamp
	SIPHdrTo
	SIPHdrUnsupported
	SIPHdrUserAgent
	SIPHdrVia
	SIPHdrWarning
	SIPHdrWWWAuthenticate
)

// HeadersList SIP headers list
type HeadersList struct {
	*list.List
}

func initHeadersList(msg *Message) {
	msg.Headers = HeadersList{list.New()}
}

// Count number of headers
func (l HeadersList) Count() int {
	return l.Len()
}

// FindByName find header by name
func (l HeadersList) FindByName(name string) *Header {
	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value.(*Header)
		if strings.EqualFold(name, h.Name()) {
			return h
		}
	}
	return nil
}

// Find find header by ID
func (l HeadersList) Find(id HdrType) *Header {
	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value.(*Header)
		if h.ID() == id {
			return h
		}
	}
	return nil
}

// FindAll find all headers by ID and returns array of headers
func (l HeadersList) FindAll(id HdrType) []*Header {
	headers := make([]*Header, 0)
	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value.(*Header)
		if h.ID() == id {
			headers = append(headers, h)
		}
	}
	return headers
}

// ForEach call first class function on each header
func (l HeadersList) ForEach(callback func(h *Header)) {
	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value.(*Header)
		callback(h)
	}
}

func (l HeadersList) exists(buf []byte) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value.(*Header)
		if bytes.Equal(buf, h.buf) {
			return true
		}
	}
	return false
}

func (l HeadersList) push(h *Header) {
	l.PushBack(h)
}

func (l HeadersList) remove(name string) bool {
	found := false
	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value.(*Header)
		if strings.EqualFold(h.Name(), name) {
			l.Remove(e)
			found = true
		}
	}
	return found
}

// Header SIP header
type Header struct {
	buf   []byte
	id    HdrType
	name  pl
	value pl
}

// ID SIP header ID
func (h *Header) ID() HdrType {
	return h.id
}

// Name SIP header name
func (h *Header) Name() string {
	return string(h.buf[h.name.p:h.name.l])
}

// Value SIP header value
func (h *Header) Value() string {
	return string(h.buf[h.value.p:h.value.l])
}

// CSeq SIP sequence number
type CSeq struct {
	Num    uint
	Method string
}

func searchParam(name string, buf []byte, params []pl) (string, bool) {
	for _, p := range params {
		prm := bytes.SplitN(buf[p.p:p.l], []byte("="), 2)
		if bytes.EqualFold([]byte(name), prm[0]) {
			if len(prm) < 2 {
				return "", true
			}
			return string(prm[1]), true
		}
	}
	return "", false
}

// local helper functions and structures
func randomString() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Uint32())
}

func randomStringPrefix(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, randomString())
}

func hashString() string {
	b := make([]byte, 18)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func headerValue(vals ...string) ([]byte, pl, pl) {
	pn, pv := pl{}, pl{}
	b := buffer{}
	b.WriteString(vals[0])
	pn.l = b.plen()
	b.WriteByte(':')

	pv.p = b.plen() + 1
	for _, v := range vals[1:] {
		b.WriteByte(' ')
		b.WriteString(v)
	}
	pv.l = b.plen()
	b.crlf()
	return b.Bytes(), pn, pv
}

// local buffer extended structrue
type buffer struct {
	bytes.Buffer
}

func (b *buffer) init(data []byte) {
	b.Write(data)
}

func (b *buffer) plen() ptr {
	return ptr(b.Len())
}

// field name with colon and space prepended and pl set.
func (b *buffer) name(name string, p *pl) {
	b.WriteString(name)
	p.l = b.plen()
	b.WriteString(": ")
}

func (b *buffer) write(val string, p *pl) {
	if p != nil {
		p.p = b.plen()
	}
	b.WriteString(val)
	if p != nil {
		p.l = b.plen()
	}
}

func (b *buffer) writeBytePrefix(prefix byte, value string, p *pl) {
	b.WriteByte(prefix)
	b.write(value, p)
}

// write parameter (name=value) to buffer and prepend ";"
// pl pointer for parameter is set only for value.
// if name == value then single word parameter is written: ;param
func (b *buffer) paramVal(name, value string, p *pl) {
	b.WriteByte(';')
	b.WriteString(name)
	if name != value {
		b.WriteByte('=')
		b.write(value, p)
	}
}

func (b *buffer) param(name, value string) pl {
	c := pl{}
	b.WriteByte(';')
	c.p = b.plen()
	b.WriteString(name)
	if name != value {
		b.WriteByte('=')
		b.write(value, nil)
	}
	c.l = b.plen()
	return c
}

// write and wrap
// if plInside is true then set pl only around value, otherwise all with wrapper
func (b *buffer) wwrap(wrapper, value string, p *pl, plInside bool) {
	p.p = b.plen()
	b.WriteByte(wrapper[0])
	b.WriteString(value)
	b.WriteByte(wrapper[1])
	p.l = b.plen()

	if plInside {
		p.p++
		p.l--
	}
}

func (b *buffer) appendPort(port int, p *pl) error {
	if port == 0 {
		return nil
	}

	if port < 0 || port > 65535 {
		return ErrorURI.msg("Invalid port %d", port)
	}

	b.writeBytePrefix(':', strconv.Itoa(port), p)

	return nil
}

func (b *buffer) byt(p pl) []byte {
	buf := b.Bytes()
	return buf[p.p:p.l]
}

func (b *buffer) str(p pl) string {
	return string(b.byt(p))
}

func (b *buffer) uncrlf() {
	b.Truncate(b.Len() - 2)
}

func (b *buffer) crlf() []byte {
	b.WriteString("\r\n")
	return b.Bytes()
}
