package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/state"
	"net"
	"sync"
	"time"
)

// ServerInvite transaction
type ServerInvite struct {
	mu sync.Mutex
	*Transaction
	response *sip.Packet
}

func initServerInvite(pack *sip.Packet, layer *Layer) *ServerInvite {
	logger.Log("txn:server:inv: init new transaction")
	txn := &ServerInvite{
		Transaction: newTransaction(pack, layer),
	}
	txn.state.Set(state.Proceeding)
	txn.layer.passToTU(pack)

	go txn.fireEarly()
	return txn
}

// Consume invite requests and re-transactions and ACK
func (txn *ServerInvite) Consume(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("txn:server:inv: consume packet has <nil> SIP Message")
		return
	}

	msg := pack.Message

	if msg.IsRequest() {
		txn.consumeRequest(msg)
		return
	}

	// responses are expected only in the Proceeding state
	if txn.state.IsProceeding() {
		txn.setResponse(pack)
		txn.layer.passToTransp(pack)
		switch {
		case msg.IsSuccess():
			txn.Terminate()
		case msg.IsRedirOrError():
			txn.state.Set(state.Completed)
			go txn.fireH()
			go txn.fireG()
		}
	}
}

// Match server transaction
func (txn *ServerInvite) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	if txn.MatchServer(msg) {
		return txn, true
	}
	return nil, false
}

func (txn *ServerInvite) fireH() {
	select {
	case <-txn.timer.FireH():
		if txn.state.IsCompleted() {
			logger.Wrn("txn:server:inv: timer H fired in completed state. terminating txn")
			txn.layer.passErr(ErrTxnFail)
			txn.Terminate()
		}
	case <-txn.halt:
		logger.Wrn("txn:server:inv: timer H interupted by halt event")
	}
}

func (txn *ServerInvite) fireG() {
	tick := txn.timer.TickerG()
	timer := time.NewTimer(tick())
	for {
		select {
		case <-timer.C:
			if !txn.state.IsCompleted() {
				return
			}
			txn.layer.passToTransp(txn.response)
			timer.Reset(tick())
		case <-txn.halt:
			return
		}
	}
}

func (txn *ServerInvite) fireEarly() {
	select {
	case <-time.After(200 * time.Millisecond):
	case <-txn.halt:
		return
	}

	txn.mu.Lock()
	defer txn.mu.Unlock()
	if txn.response != nil {
		return
	}

	txn.layer.passToTransp(txn.gen100())
}

func (txn *ServerInvite) fireI() {
	if txn.IsReliable() {
		txn.Terminate()
		return
	}

	select {
	case <-txn.halt:
	case <-txn.timer.FireI():
		txn.Terminate()
	}
}

func (txn *ServerInvite) consumeRequest(msg *sipmsg.Message) {
	if msg.IsInvite() {
		if !(txn.state.IsProceeding() || txn.state.IsCompleted()) {
			logger.Wrn("txn:server:inv: ignore INVITE retransmission for state %q", txn.state)
			return
		}

		if txn.response == nil {
			txn.setResponse(txn.gen100())
		}
		txn.layer.passToTransp(txn.response)
		return
	}

	if msg.IsMethod("ACK") {
		switch {
		case txn.state.IsCompleted():
			txn.state.Set(state.Confirmed)
			go txn.fireI()
		case txn.state.IsConfirmed():
			logger.Log("txn:server:inv: absorb ACK in Confirmed state")
		}
		return
	}
	logger.Wrn("txn:server:inv: unacceptable request %q", msg.FirstLine())
}

func (txn *ServerInvite) setResponse(respPack *sip.Packet) {
	txn.mu.Lock()
	defer txn.mu.Unlock()
	txn.response = respPack
}

func (txn *ServerInvite) gen100() *sip.Packet {
	resp := txn.req.Message.Response(100, "Trying")
	pack := &sip.Packet{
		Message:    resp,
		SendTo:     []net.Addr{txn.req.RemoteSock},
		ReqAddrs:   txn.req.ReqAddrs,
		LocalSock:  txn.req.LocalSock,
		RemoteSock: txn.req.RemoteSock,
	}

	return pack
}
