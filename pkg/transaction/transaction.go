package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/state"
	"gosip/pkg/transaction/timer"
)

type Transaction struct {
	req   *sip.Packet
	layer *Layer
	state *state.State
	timer *timer.Timer
	halt  chan struct{}
}

func newTransaction(pack *sip.Packet, layer *Layer) *Transaction {
	timerSetup := func() *timer.Timer {
		if layer.SetupTimers == nil {
			return timer.New()
		}
		return layer.SetupTimers(timer.New())
	}

	return &Transaction{
		req:   pack,
		layer: layer,
		state: state.New(),
		timer: timerSetup(),
		halt:  make(chan struct{}),
	}
}

func (txn *Transaction) BranchID() string {
	if via := txn.reqTopVia(); via != nil {
		return via.Branch
	}
	return ""
}

// IsTranspReliable returns true if transport is not UDP
func (txn *Transaction) IsReliable() bool {
	if txn.req == nil || len(txn.req.SendTo) == 0 {
		return false
	}
	switch txn.req.SendTo[0].Network() {
	case "udp", "udp4", "udp6":
		return false
	default:
		return true
	}
}

func (txn *Transaction) MatchClient(msg *sipmsg.Message) bool {
	// TODO: implement match client rfc3261#17.1.3
	return false
}

// MatchServer tries to match request to transaction
// following rfc3261#17.2.3 Matching Requests to Server Transactions
func (txn *Transaction) MatchServer(msg *sipmsg.Message) bool {
	reqvia := txn.reqTopVia()
	if reqvia == nil {
		logger.Wrn("initial request message top Via header not found")
		return false
	}

	incomevia := msg.TopVia()
	if incomevia == nil {
		logger.Wrn("incoming message top Via header not found for %s %s",
			msg.Method, msg.RURI)
		return false
	}
	// TODO: it should have a check if branch starts with "z9hG4bK"
	// and if not handle backwards compatibility with RFC 2543

	// 1. the branch parameter in the request is equal to the one in the
	// top Via header field of the request that created the transaction, and
	if reqvia.Branch != incomevia.Branch {
		logger.Log("incoming message top Via branch %q not matching transaction branch %q",
			incomevia.Branch, reqvia.Branch)
		return false
	}

	// 2. the sent-by value in the top Via of the request is equal to the
	// one in the request that created the transaction, and
	if len(reqvia.Host) == 0 || (reqvia.Host != incomevia.Host && reqvia.Port != incomevia.Port) {
		logger.Wrn("request top via branch host %q port %q does not match incoming host %q port %q",
			reqvia.Host, reqvia.Port, incomevia.Host, incomevia.Port)
		return false
	}

	// 3. the method of the request matches the one that created the
	// transaction, except for ACK, where the method of the request
	// that created the transaction is INVITE.
	if msg.IsMethod("ACK") {
		logger.Wrn("unexpected ACK in server transaction match")
		return false
	}

	return txn.req.Message.IsMethod(msg.Method)
}

// stop all running background timers and actions
// remove transaction from the store
func (txn *Transaction) terminate() {
	select {
	case <-txn.halt:
		// channel is already closed
	default:
		close(txn.halt)
	}
	txn.state.Set(state.Terminated)
	txn.layer.Destroy(txn.BranchID())
}

func (txn *Transaction) reqTopVia() *sipmsg.HeaderVia {
	msg := txn.reqMessage()
	if msg == nil {
		return nil
	}
	return txn.req.Message.TopVia()
}

func (txn *Transaction) reqMessage() *sipmsg.Message {
	if txn.req == nil {
		return nil
	}
	return txn.req.Message
}
