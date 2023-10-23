package txnlayer

import (
	"errors"
	"net/netip"
	"sync/atomic"
	"time"
)

const (
	Unknown uint32 = iota
	Calling
	Proceeding
	Completed
	Terminated
)

type Message interface {
	Ack() Message
	IsInvite() bool
	IsResponse() bool
	TopViaBranch() string
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
}

type TxnBasic struct {
	transp   Transport
	endpoint EndPoint
	branch   string
	state    *atomic.Uint32
	timer    Timer
	addr     netip.AddrPort
	halt     chan struct{}
}

type Timer struct {
	T1 time.Duration
	D  time.Duration
}

func initTimer() Timer {
	return Timer{
		T1: 500 * time.Millisecond,
		D:  32 * time.Second,
	}
}

func initBasicTxn(transp Transport, endpoint EndPoint, branch string) TxnBasic {
	state := new(atomic.Uint32)
	state.Store(Unknown)
	return TxnBasic{
		transp:   transp,
		endpoint: endpoint,
		branch:   branch,
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
}

type TxnClientNonInvite struct {
	TxnBasic
}

type Transaction interface {
	Consume(msg Message)
	ID() string
	Init(msg Message, addr netip.AddrPort)
}

type pool map[string]Transaction

type TxnLayer struct {
	client   pool
	server   pool
	endpoint EndPoint
}

// EndPoint is actually an interface to TU
func New(endpoint EndPoint) *TxnLayer {
	return &TxnLayer{
		client:   make(pool),
		server:   make(pool),
		endpoint: endpoint,
	}
}

// client transaction create
// used by TU to start new transaction
func (tl *TxnLayer) Client(msg Message, transp Transport, addr netip.AddrPort) error {
	var txn Transaction
	if msg.IsInvite() {
		txn = createClientInvTxn(transp, tl.endpoint, msg)
	} else {
		return errors.New("not implemented client non-invite txn")
		// txn = &TxnClientNonInvite{}
	}

	txn.Init(msg, addr)

	tl.client[txn.ID()] = txn
	return nil
}

// consume new message from transport
// match existing or create new transaction
func (s *TxnLayer) Consume(msg Message, transp Transport, addr netip.AddrPort) {
}
