package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/state"
	"time"
)

type ClientNonInvite struct {
	*Transaction
}

func initClientNonInvite(pack *sip.Packet, layer *Layer) *ClientNonInvite {
	logger.Log("txn:client:noninv: init new transaction in trying state")
	txn := &ClientNonInvite{
		Transaction: newTransaction(pack, layer),
	}

	txn.layer.passToTransp(txn.req) // send first messag
	txn.state.Set(state.Trying)

	if txn.IsReliable() {
		logger.Log("txn:client:noninv do not start timer F and E for reliable transport")
		return txn
	}

	go txn.timerF()
	go txn.timerE()

	return txn
}

func (txn *ClientNonInvite) Consume(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("txn:client:noninv: consume packet has <nil> SIP Message")
		return
	}
	if pack.Message.IsRequest() {
		// ignore requests
		return
	}
	msg := pack.Message
	logger.Log("txn:client:noninv: consume message %q while state is %q",
		msg.FirstLine(), txn.state.String())

	if txn.state.IsTrying() || txn.state.IsProceeding() {
		txn.layer.passToTU(pack)

		switch {
		case msg.IsProvisional():
			txn.state.Set(state.Proceeding)
		case msg.IsFinalResponse():
			txn.completed()
		default:
			logger.Wrn("txn:client:noninv: invalid response: %s", msg.FirstLine())
		}
		return
	}

	logger.Log("txn:client:noninv: absorb response %q", msg.FirstLine())
}

// Match returns transaction and true if match to the SIP Message
func (txn *ClientNonInvite) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	if txn.MatchClient(msg) {
		return txn, true
	}
	return nil, false
}

func (txn *ClientNonInvite) completed() {
	logger.Log("txn:client:noninv: enter completed state")
	txn.state.Set(state.Completed)
	if txn.IsReliable() {
		logger.Log("txn:client:noninv: terminate for reliable transport")
		txn.Terminate()
		return
	}
	// unreliable transport
	logger.Log("txn:client:noninv: fire timer K")
	go func() {
		select {
		case <-txn.timer.FireK():
		case <-txn.halt:
			return
		}
		txn.Terminate()
	}()
}

func (txn *ClientNonInvite) timerF() {
	logger.Log("txn:client:noninv: fire timer F")
	select {
	case <-txn.timer.FireF():
	case <-txn.halt:
		return
	}
	if txn.state.IsTrying() || txn.state.IsProceeding() {
		txn.layer.passErr(ErrTimeout)
	}
	logger.Log("txn:client:noninv: done timer F. set terminated state and destroy transaction")
	txn.Terminate()
}

func (txn *ClientNonInvite) timerE() {
	logger.Log("txn:client:noninv: fire timer E")
	tick := txn.timer.TickerE()
	timer := time.NewTimer(0)
	for {
		next := tick(txn.state.IsProceeding())
		timer.Reset(next)
		select {
		case <-timer.C:
			if txn.state.IsTrying() || txn.state.IsProceeding() {
				logger.Log("txn:client:noninv timer E fired after %v. resend initial request", next)
				txn.layer.passToTransp(txn.req)
				continue
			}
			logger.Log("txn:client:noninv state is %q. stop timer E", txn.state)
			return
		case <-txn.halt:
			logger.Log("txn:client:noninv timer E is halted")
			return
		}
	}
}
