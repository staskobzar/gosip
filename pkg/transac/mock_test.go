package transac

import (
	"errors"
	"net/netip"
	"sync"
)

type mockMsg struct {
	method string
	code   int
	branch string
}

func (m *mockMsg) Ack() Message {
	return &mockMsg{method: "ACK", branch: m.branch}
}
func (m *mockMsg) IsResponse() bool     { return m.code >= 100 }
func (m *mockMsg) TopViaBranch() string { return m.branch }
func (m *mockMsg) Method() string       { return m.method }
func (m *mockMsg) ResponseCode() int    { return m.code }

type mockTransp struct {
	isReliable bool
	addr       netip.AddrPort
	msg        []Message
	mu         sync.Mutex
	senderr    error
}

func (t *mockTransp) Send(addr netip.AddrPort, msg Message) error {
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
	msg       []Message
	tout      bool
	destroyID string
	err       error
}

func (e *mockEndPoint) TUConsume(msg Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.msg = append(e.msg, msg)
}
func (e *mockEndPoint) Error(err error, msg Message) {
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

func createMock() (*mockEndPoint, *mockTransp, *mockMsg, netip.AddrPort) {
	ep := &mockEndPoint{msg: make([]Message, 0)}
	tr := &mockTransp{msg: make([]Message, 0)}
	msg := &mockMsg{method: "INVITE", branch: "z9hG4bK-f00"}
	addr, _ := netip.ParseAddrPort("127.0.0.1:5670")
	return ep, tr, msg, addr
}
