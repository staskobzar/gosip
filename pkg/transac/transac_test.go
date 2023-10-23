package txnlayer

import (
	"net/netip"
	"sync"
)

type mockMsg struct {
	isInv  bool
	isResp bool
	isAck  bool
	code   int
	branch string
}

func (m *mockMsg) Ack() Message {
	return &mockMsg{isAck: true, branch: m.branch}
}
func (m *mockMsg) IsInvite() bool       { return m.isInv }
func (m *mockMsg) IsResponse() bool     { return m.isResp }
func (m *mockMsg) TopViaBranch() string { return m.branch }
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
	mu   sync.Mutex
	msg  []Message
	tout bool
	err  error
}

func (e *mockEndPoint) Consume(msg Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.msg = append(e.msg, msg)
}
func (e *mockEndPoint) Error(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.err = err
}
func (e *mockEndPoint) TimeoutError(_ Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.tout = true
}
