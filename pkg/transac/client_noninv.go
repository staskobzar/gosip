package txnlayer

import (
	"net/netip"
	"time"
)

type TxnClientNonInvite struct {
	TxnBasic
}

func createClientNonInvTxn(transp Transport, endpoint EndPoint, msg Message) *TxnClientNonInvite {
	return &TxnClientNonInvite{
		TxnBasic: initBasicTxn(transp, endpoint, msg),
	}
}

func (txn *TxnClientNonInvite) Init(msg Message, addr netip.AddrPort) {
	txn.addr = addr
	txn.trying(msg)
	txn.fireTimerF(msg)
}

func (txn *TxnClientNonInvite) Consume(msg Message) {
	if !msg.IsResponse() {
		return
	}

	code := msg.ResponseCode()

	switch txn.state.Load() {
	case Trying, Proceeding:
		if code >= 100 && code <= 199 {
			txn.state.Store(Proceeding)
		} else if code >= 200 && code <= 699 {
			txn.complete()
		} else {
			// TODO log invalid code
		}
	case Completed:
		// TODO log absorbing response retransmissions
	default:
		// TODO log unknown state
	}
}

func (txn *TxnClientNonInvite) trying(msg Message) {
	txn.state.Store(Trying)
	txn.Send(msg)

	if txn.transp.IsReliable() {
		return // no E timer for reliable transport
	}
	txn.fireTimerE(msg)
}

func (txn *TxnClientNonInvite) complete() {
	txn.state.Store(Completed)

	if txn.transp.IsReliable() {
		txn.terminate()
		return // no K timer for reliable transport
	}

	txn.fireTimerK()
}

func (txn *TxnClientNonInvite) fireTimerE(msg Message) {
	go func() {
		tick := tickTimerE(txn.timer.T1, txn.timer.T2)
		timer := time.NewTimer(0)
		state := txn.state.Load()
		for state == Trying || state == Proceeding {
			timer.Reset(tick(state))
			select {
			case <-timer.C:
				state = txn.state.Load()
				txn.Send(msg) // resend SIP message
			case <-txn.halt:
				return
			}
		}
	}()
}

func tickTimerE(t1, t2 time.Duration) func(uint32) time.Duration {
	return func(state uint32) time.Duration {
		if state == Proceeding {
			return t2
		}
		t := min(t1, t2)
		if t1 < t2 {
			t1 *= 2
		}
		return t
	}
}

func (txn *TxnClientNonInvite) fireTimerF(msg Message) {
	go func() {
		select {
		case <-time.After(txn.timer.T1 * 64):
			state := txn.state.Load()
			if state == Trying || state == Proceeding {
				txn.endpoint.TimeoutError(msg)
				txn.terminate()
			}
		case <-txn.halt:
		}
	}()
}

// just buffer any additional response retransmissions
func (txn *TxnClientNonInvite) fireTimerK() {
	go func() {
		select {
		case <-time.After(txn.timer.T4):
			txn.terminate()
		case <-txn.halt:
		}
	}()
}
