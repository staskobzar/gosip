package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/internal/pool"
	"gosip/pkg/transaction/internal/state"
)

type ServerNonInvite struct {
	*Transaction
	response *sip.Packet
}

func initServerNonInvite(pack *sip.Packet, layer *Layer) *ServerNonInvite {
	txn := &ServerNonInvite{
		Transaction: newTransaction(pack, layer),
	}
	logger.Log("txn:srv:noninv: init new transaction in trying state")
	txn.state.Set(state.Trying)
	txn.layer.passToTU(pack) // TODO: ?? should it be done in goroutin ??
	return txn
}

// Consume accept and process SIP message from network
// on running transaction
func (txn *ServerNonInvite) Consume(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("txn:srv:noninv: consume packet has <nil> SIP Message")
		return
	}
	msg := pack.Message
	logger.Log("txn:srv:noninv: consume message %q while state is %q",
		msg.FirstLine(), txn.state.String())

	if msg.IsResponse() {
		logger.Log("txn:srv:noninv: recv response %s %s", msg.Code, msg.Reason)
		if txn.state.IsTrying() {
			logger.Log("txn:srv:noninv: is in trying state")
			if msg.IsProvisional() {
				logger.Log("txn:srv:noninv: moved to proceeding state")
				txn.state.Set(state.Proceeding)
				txn.response = pack
				txn.layer.passToTransp(pack)
				return
			}

			if msg.IsFinalResponse() {
				txn.completed(pack)
				return
			}
		}

		if txn.state.IsProceeding() {
			logger.Log("txn:srv:noninv: is in proceeding state")
			if msg.IsProvisional() {
				txn.response = pack
				txn.layer.passToTransp(pack)
				return
			}

			if msg.IsFinalResponse() {
				txn.completed(pack)
				return
			}
		}

		if txn.state.IsCompleted() {
			logger.Log("txn:srv:noninv: is in completed state")
			// Any other final responses passed by the TU to the server
			// transaction MUST be discarded while in the "Completed" state.
			return
		}
	} else {
		logger.Log("txn:srv:noninv: re-transmitting request %q", pack.Message.FirstLine())
		// Request re-transmissions
		// resend last response
		if txn.state.IsProceeding() || txn.state.IsCompleted() {
			txn.layer.passToTransp(txn.response)
			return
		}
		//
	}
}

func (txn *ServerNonInvite) completed(pack *sip.Packet) {
	logger.Log("txn:srv:noninv: moved to completed state")
	txn.response = pack
	txn.layer.passToTransp(pack)
	txn.state.Set(state.Completed)
	go func() {
		<-txn.timer.FireJ()
		txn.state.Set(state.Terminated)
		logger.Log("txn:srv:noninv: done timer J. set terminated state and destroy transaction")
		txn.layer.Destroy(txn)
	}()
}

func (txn *ServerNonInvite) Match(msg *sipmsg.Message) (pool.Transaction, bool) {
	if txn.MatchServer(msg) {
		return txn, true
	}
	return nil, false
}