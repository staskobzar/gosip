package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/internal/state"
	"gosip/pkg/transaction/internal/timer"
)

type Transaction struct {
	req   *sip.Packet
	layer *Layer
	state *state.State
	timer *timer.Timer
}

func newTransaction(pack *sip.Packet, layer *Layer) *Transaction {
	return &Transaction{
		req:   pack,
		layer: layer,
		state: state.New(),
		timer: timer.New(),
	}
}

func (txn *Transaction) BranchID() string {
	if via := txn.reqTopVia(); via != nil {
		return via.Branch
	}
	return ""
}

func (txn *Transaction) MatchClient(msg *sipmsg.Message) bool {
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
