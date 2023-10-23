package txnlayer

import (
	"net/netip"
	"time"
)

type TxnClientInvite struct {
	TxnBasic
}

func createClientInvTxn(transp Transport, endpoint EndPoint, msg Message) *TxnClientInvite {
	return &TxnClientInvite{
		TxnBasic: initBasicTxn(transp, endpoint, msg.TopViaBranch()),
	}
}

func (txn *TxnClientInvite) Init(msg Message, addr netip.AddrPort) {
	txn.addr = addr
	txn.calling(msg)
	txn.fireTimerB(msg)
}

func (txn *TxnClientInvite) Consume(msg Message) {
	if !msg.IsResponse() {
		return
	}

	code := msg.ResponseCode()

	switch txn.state.Load() {
	case Calling, Proceeding:
		txn.proceed(code, msg)
	case Completed:
		// absorb re-transactions
		if code >= 300 && code <= 699 {
			txn.Send(msg.Ack())
		}
	default:
		// TODO: log invalid state
	}
}

func (txn *TxnClientInvite) calling(msg Message) {
	txn.state.Store(Calling)
	if txn.transp.IsReliable() {
		txn.Send(msg)
		return // no timer A for reliable transport
	}

	// for unrelables transport control request retransmissions with timer A
	txn.fireTimerA(msg)
}

func (txn *TxnClientInvite) fireTimerA(msg Message) {
	go func() {
		t1 := txn.timer.T1
		timer := time.NewTimer(t1)
		txn.Send(msg)
		// in case if state is failed to change use the same
		// time limit as timer B (64*T1) INVITE transaction limit
		for t := t1; t <= (64 * t1); t *= 2 {
			timer.Reset(t)
			select {
			case <-timer.C:
				if txn.state.Load() != Calling {
					return
				}
			case <-txn.halt:
				return
			}
			// resend SIP message
			txn.Send(msg)
		}
	}()
}

func (txn *TxnClientInvite) fireTimerB(msg Message) {
	go func() {
		select {
		case <-time.After(txn.timer.T1 * 64):
			if txn.state.Load() == Calling {
				txn.endpoint.TimeoutError(msg)
				txn.terminate()
			}
		case <-txn.halt:
			return
		}
	}()
}

func (txn *TxnClientInvite) fireTimerD() {
	if txn.transp.IsReliable() {
		// for reliable transport timer D is 0
		txn.terminate()
		return
	}
	go func() {
		select {
		case <-time.After(txn.timer.D):
			txn.terminate()
		case <-txn.halt:
			return
		}
	}()
}

func (txn *TxnClientInvite) proceed(code int, msg Message) {
	if code >= 100 && code < 200 {
		txn.state.Store(Proceeding)
	} else if code >= 200 && code < 300 {
		txn.terminate()
	} else if code >= 300 && code <= 699 {
		txn.state.Store(Completed)
		txn.Send(msg.Ack())
		txn.fireTimerD()
	} else {
		// TODO: log invalid code error
	}
	txn.endpoint.Consume(msg)
}
