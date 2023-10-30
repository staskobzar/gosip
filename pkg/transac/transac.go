package txnlayer

import (
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

type Message interface {
	Ack() Message
	IsResponse() bool
	TopViaBranch() string
	Method() string
	ResponseCode() int
}

type Transport interface {
	Send(addr netip.AddrPort, msg Message) error
	IsReliable() bool
}

type EndPoint interface {
	Consume(msg Message)
	Error(err error)
	TimeoutError(msg Message)
	TxnDestroy(ID string)
}

type TxnBasic struct {
	transp   Transport
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

func initBasicTxn(transp Transport, endpoint EndPoint, msg Message) TxnBasic {
	state := new(atomic.Uint32)
	state.Store(Unknown)
	return TxnBasic{
		transp:   transp,
		endpoint: endpoint,
		method:   msg.Method(),
		branch:   msg.TopViaBranch(),
		state:    state,
		timer:    initTimer(),
		halt:     make(chan struct{}),
	}
}

func (txn TxnBasic) ID() string {
	return txn.branch
}

func (txn TxnBasic) Send(msg Message) {
	err := txn.transp.Send(txn.addr, msg)
	if err != nil {
		txn.state.Store(Terminated)
		txn.endpoint.Error(err)
	}
}

func (txn TxnBasic) terminate() {
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
	Consume(msg Message)
	ID() string
	Init(msg Message, addr netip.AddrPort)
}

type pool map[string]Transaction

type TxnLayer struct {
	pool     pool
	endpoint EndPoint
}

// EndPoint is actually an interface to TU
func New(endpoint EndPoint) *TxnLayer {
	return &TxnLayer{
		pool:     make(pool),
		endpoint: endpoint,
	}
}

// client transaction create
// used by TU to start new transaction
func (txl *TxnLayer) Client(msg Message, transp Transport, addr netip.AddrPort) {
	var txn Transaction
	if msg.Method() == "INVITE" {
		txn = createClientInvTxn(transp, txl.endpoint, msg)
	} else {
		txn = createClientNonInvTxn(transp, txl.endpoint, msg)
	}

	txn.Init(msg, addr)

	txl.push(txn)
}

// consume new message from transport
// match existing or create new Server transaction
func (txl *TxnLayer) Consume(msg Message, transp Transport, addr netip.AddrPort) {
	if txn, exists := txl.pool[msg.TopViaBranch()]; exists {
		txn.Consume(msg)
	} else {
		// new server txn
	}
}

func (txl *TxnLayer) Destroy(txnID string) {
	delete(txl.pool, txnID)
}

func (txl *TxnLayer) push(txn Transaction) {
	txl.pool[txn.ID()] = txn
}
