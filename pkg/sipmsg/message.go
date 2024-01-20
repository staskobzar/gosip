package sipmsg

import "strings"

type Message struct {
	t      HType
	Method string
	RURI   *URI
	Code   string
	Reason string

	CallID string
	CSeq   int
	MaxFwd int
	From   *NameAddr
	To     *NameAddr

	Headers Headers
	Body    string
}

func NewMessage() *Message {
	return &Message{
		Headers: make(Headers, 0, 32),
	}
}

func (msg *Message) IsResponse() bool {
	return msg.t == HResponse
}

// HLen returns number of headers in the SIP message
func (msg *Message) HLen() int {
	return msg.Headers.Len()
}

func (msg *Message) Find(t HType) anyHeader {
	return msg.find(func(h anyHeader) bool { return h.t() == t })
}

func (msg *Message) FindAll(t HType) Headers {
	return msg.findAll(func(h anyHeader) bool { return h.t() == t })
}

func (msg *Message) FindByName(name string) anyHeader {
	return msg.find(func(h anyHeader) bool { return h.name() == name })
}

func (msg *Message) FindByNameAll(name string) Headers {
	return msg.findAll(func(h anyHeader) bool { return h.name() == name })
}

func (msg *Message) String() string {
	buf := make([]string, 0, msg.HLen()+1)
	buf = append(buf, msg.firstLine())

	for _, hdr := range msg.Headers {
		h := hdr.String()
		if hdr.t() == HCSeq {
			h += " " + msg.Method
		}
		buf = append(buf, h)
	}

	return strings.Join(buf, "\r\n") + "\r\n\r\n" + msg.Body
}

func (msg *Message) firstLine() string {
	if msg.t == HRequest {
		return msg.Method + " " + msg.RURI.String() + " SIP/2.0"
	}
	return "SIP/2.0" + msg.Code + " " + msg.Reason
}

func (msg *Message) find(match func(h anyHeader) bool) anyHeader {
	for _, h := range msg.Headers {
		if match(h) {
			return h
		}
	}
	return nil
}

func (msg *Message) findAll(match func(h anyHeader) bool) Headers {
	list := make(Headers, 0)
	for _, h := range msg.Headers {
		if match(h) {
			list = append(list, h)
		}
	}
	return list
}

func (msg *Message) pushHeader(t HType, name, value string) {
	generic := &HeaderGeneric{
		Type:  t,
		Name:  name,
		Value: value,
	}
	msg.Headers = append(msg.Headers, generic)
}
