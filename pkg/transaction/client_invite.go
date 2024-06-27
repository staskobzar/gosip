package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/state"
	"time"
)

// ClientInvite transaction
type ClientInvite struct {
	*Transaction
}

func initClientInvite(pack *sip.Packet, layer *Layer) *ClientInvite {
	logger.Log("txn:client:inv init new transaction")
	txn := &ClientInvite{
		Transaction: newTransaction(pack, layer),
	}
	txn.state.Set(state.Calling)
	txn.layer.passToTransp(pack)

	go txn.timerA()
	go txn.timerB()

	return txn
}

// Consume process SIP Responses from transport
func (txn *ClientInvite) Consume(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("txn.client:inv: consume packet has <nil> SIP Message")
		return
	}

	if pack.Message.IsRequest() {
		return
	}

	resp := pack.Message
	logger.Log("txn:client:inv: consume message %q while state is %q", resp.FirstLine(), txn.state)

	if txn.state.IsCalling() || txn.state.IsProceeding() {
		switch {
		case resp.IsProvisional():
			txn.state.Set(state.Proceeding)
			txn.layer.passToTU(pack)
		case resp.IsSuccess():
			txn.layer.passToTU(pack)
			txn.Terminate()
		case resp.IsRedirOrError(): // 300-699
			txn.layer.passToTU(pack)
			txn.sendACK(pack)
			txn.complete()
		default:
			logger.Err("txn:client:inv: invalid response code %q", resp.Code)
		}
		return
	}

	if txn.state.IsCompleted() {
		txn.sendACK(pack)
		return
	}
}

// Match returns true if SIP message match transaction
func (txn *ClientInvite) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	if txn.MatchClient(msg) {
		return txn, true
	}
	return nil, false
}

func (txn *ClientInvite) sendACK(pack *sip.Packet) {
	ackMsg := txn.req.Message.ACK(pack.Message)
	ack := &sip.Packet{
		SendTo:     txn.req.SendTo,
		ReqAddrs:   txn.req.ReqAddrs,
		LocalSock:  txn.req.LocalSock,
		RemoteSock: txn.req.RemoteSock,
		Message:    ackMsg,
	}
	txn.layer.passToTransp(ack)
}

func (txn *ClientInvite) complete() {
	if txn.IsReliable() {
		txn.Terminate()
		logger.Log("txn:client:inv: skip timer D for reliable transport")
		return
	}
	txn.state.Set(state.Completed)
	go func() {
		select {
		case <-txn.timer.FireD():
			txn.Terminate()
		case <-txn.halt:
		}
	}()
}

func (txn *ClientInvite) timerA() {
	if txn.IsReliable() {
		logger.Log("txn:client:inv: skip timer A for reliable transport")
		return
	}
	timer := time.NewTimer(0)
	tick := txn.timer.TickerA()
	for txn.state.IsCalling() {
		timer.Reset(tick())
		select {
		case <-txn.halt:
			return
		case <-timer.C:
			txn.layer.passToTransp(txn.req)
		}
	}
}

func (txn *ClientInvite) timerB() {
	defer txn.Terminate()
	select {
	case <-txn.timer.FireB():
		if txn.state.IsCalling() {
			txn.layer.passErr(ErrTimeout)
		}
	case <-txn.halt:
	}
}
