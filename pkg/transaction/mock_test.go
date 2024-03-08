package transaction

import (
	"errors"
	"gosip/pkg/sipmsg"
	"net/netip"
	"sync"
)

const (
	stubInvite = "INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds.1\r\n" +
		"Via: SIP/2.0/UDP pc32.atlanta.com;branch=z9hG4bK776asdhds.2\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
		"Contact: <sip:alice@pc33.atlanta.com>\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 133\r\n\r\n" +
		"v=0\r\no=alice 2890844526 2890844526 IN IP4 10.1.1.100\r\n" +
		"s=\r\nc=IN IP4 10.1.1.100\r\nt=0 0\r\n" +
		"m=audio 49170 RTP/AVP 0\r\na=rtpmap:0 PCMU/8000\r\n"
	stubNonInvite = "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n\r\n"
	stubResp = "SIP/2.0 200 OK\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7;received=192.0.2.4\r\n" +
		"Route: <sip:alice@atlanta.com>,<sip:bob@biloxi.com>\r\n" +
		"To: Bob <sip:bob@biloxi.com>;tag=2493k59kd\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 INVITE\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Content-Length: 0\r\n\r\n"
)

func stubInviteMsg() *sipmsg.Message {
	msg, err := sipmsg.Parse(stubInvite)
	if err != nil {
		panic(err)
	}
	return msg
}

func mockResponse(code, reason string) *sipmsg.Message {
	msg, err := sipmsg.Parse(stubResp)
	if err != nil {
		panic(err)
	}
	msg.Code = code
	msg.Reason = reason
	return msg
}

type mockTransp struct {
	isReliable bool
	addr       netip.AddrPort
	msg        []*sipmsg.Message
	mu         sync.Mutex
	senderr    error
}

func (t *mockTransp) Send(addr netip.AddrPort, msg *sipmsg.Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.addr = addr
	t.msg = append(t.msg, msg)
	return t.senderr
}
func (t *mockTransp) IsReliable() bool { return t.isReliable }
func (t *mockTransp) msgLen() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.msg)
}

type mockEndPoint struct {
	mu        sync.Mutex
	msg       []*sipmsg.Message
	tout      bool
	destroyID string
	err       error
}

func (e *mockEndPoint) TUConsume(msg *sipmsg.Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.msg = append(e.msg, msg)
}
func (e *mockEndPoint) Error(err error, msg *sipmsg.Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.err = err
}
func (e *mockEndPoint) TxnDestroy(id string) { e.destroyID = id }
func (e *mockEndPoint) isTout() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return errors.Is(ErrTimeout, e.err)
}
func (e *mockEndPoint) msgLen() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.msg)
}

func createMock(rawsip string) (*mockEndPoint, *mockTransp, *sipmsg.Message, netip.AddrPort) {
	ep := &mockEndPoint{msg: make([]*sipmsg.Message, 0)}
	tr := &mockTransp{msg: make([]*sipmsg.Message, 0)}
	msg, err := sipmsg.Parse(rawsip)
	if err != nil {
		panic(err)
	}
	addr, _ := netip.ParseAddrPort("127.0.0.1:5670")
	return ep, tr, msg, addr
}
