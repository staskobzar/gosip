package transaction

import (
	"errors"
	"gosip/pkg/sip"
	"net/netip"
	"sync/atomic"
	"time"
)

const (
	Unknown uint32 = iota
	Calling
	Trying
	Proceeding
	Completed
	Terminated
)

var (
	ErrTimeout = errors.New("SIP Timeout")
)

type EndPoint interface {
	TUConsume(msg sip.Message)
	Error(err error, msg sip.Message)
	TxnDestroy(ID string)
}

type Basic struct {
	transp   sip.Transport
	endpoint EndPoint
	branch   string
	method   string
	state    *atomic.Uint32
	timer    Timer
	addr     netip.AddrPort
	halt     chan struct{}
}

type Timer struct {
	T1 time.Duration
	T2 time.Duration
	T4 time.Duration
	D  time.Duration
}

func initTimer() Timer {
	return Timer{
		T1: 500 * time.Millisecond,
		T2: 4 * time.Second,
		T4: 5 * time.Second,
		D:  32 * time.Second,
	}
}

func initBasicTxn(transp sip.Transport, endpoint EndPoint, msg sip.Message) Basic {
	state := new(atomic.Uint32)
	state.Store(Unknown)
	return Basic{
		transp:   transp,
		endpoint: endpoint,
		method:   msg.Method(),
		branch:   msg.TopViaBranch(),
		state:    state,
		timer:    initTimer(),
		halt:     make(chan struct{}),
	}
}

func (txn Basic) ID() string {
	return txn.branch
}

func (txn Basic) Send(msg sip.Message) {
	err := txn.transp.Send(txn.addr, msg)
	if err != nil {
		txn.state.Store(Terminated)
		txn.endpoint.Error(err, msg)
	}
}

func (txn Basic) terminate() {
	// channel halt is only to be closed and not supposed to receive anything
	select {
	case <-txn.halt:
		// channel is already closed
	default:
		close(txn.halt)
	}
	txn.state.Store(Terminated)
	txn.endpoint.TxnDestroy(txn.ID())
}

type Transaction interface {
	Consume(msg sip.Message)
	ID() string
	Init(msg sip.Message, addr netip.AddrPort)
}

type pool map[string]Transaction

type Layer struct {
	pool pool
}

// EndPoint is actually an interface to TU
func New(endpoint EndPoint) *Layer {
	return &Layer{
		pool: make(pool),
	}
}

// client transaction create
// used by TU to start new transaction
func (txl *Layer) Client(msg sip.Message, transp sip.Transport, addr netip.AddrPort) {
	var txn Transaction
	if msg.Method() == "INVITE" {
		txn = createClientInvite(transp, txl, msg)
	} else {
		txn = createClientNonInvite(transp, txl, msg)
	}

	txn.Init(msg, addr)

	txl.push(txn)
}

// consume new message from transport
// match existing or create new Server transaction
func (txl *Layer) Consume(msg sip.Message, transp sip.Transport, addr netip.AddrPort) {
	if txn, exists := txl.pool[msg.TopViaBranch()]; exists {
		txn.Consume(msg)
	} else {
		// new server txn
	}
}

func (txl *Layer) TxnDestroy(txnID string) {
	delete(txl.pool, txnID)
}

func (txl *Layer) TUConsume(msg sip.Message) {
	// TODO send to chan msg
}

func (txl *Layer) Error(err error, msg sip.Message) {
	// TODO send err to upstream chan
}

func (txl *Layer) push(txn Transaction) {
	txl.pool[txn.ID()] = txn
}
