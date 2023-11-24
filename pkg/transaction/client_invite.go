package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"net/netip"
	"time"
)

type ClientInvite struct {
	Basic
}

func createClientInvite(transp sip.Transport, endpoint EndPoint, msg sip.Message) *ClientInvite {
	return &ClientInvite{
		Basic: initBasicTxn(transp, endpoint, msg),
	}
}

func (txn *ClientInvite) Init(msg sip.Message, addr netip.AddrPort) {
	txn.addr = addr
	txn.calling(msg)
	txn.fireTimerB(msg)
}

func (txn *ClientInvite) Consume(msg sip.Message) {
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
			// 17.1.1.3 Construction of the ACK Request
			txn.Send(msg.Ack())
		}
	default:
		logger.Err("client invite txn should not be in this state %d", txn.state.Load())
	}
}

func (txn *ClientInvite) calling(msg sip.Message) {
	txn.state.Store(Calling)
	txn.Send(msg)
	if txn.transp.IsReliable() {
		return // no timer A for reliable transport
	}

	// for unrelables transport control request retransmissions with timer A
	txn.fireTimerA(msg)
}

func (txn *ClientInvite) fireTimerA(msg sip.Message) {
	go func() {
		t1 := txn.timer.T1
		timer := time.NewTimer(0)
		state := txn.state.Load()
		for t := t1; state == Calling; t *= 2 {
			timer.Reset(t)
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

func (txn *ClientInvite) fireTimerB(msg sip.Message) {
	go func() {
		select {
		case <-time.After(txn.timer.T1 * 64):
			if txn.state.Load() == Calling {
				txn.endpoint.Error(ErrTimeout, msg)
				txn.terminate()
			}
		case <-txn.halt:
			// transaction destroyed
		}
	}()
}

func (txn *ClientInvite) fireTimerD() {
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
			// transaction destroyed
		}
	}()
}

func (txn *ClientInvite) proceed(code int, msg sip.Message) {
	switch {
	case code >= 100 && code < 200:
		txn.state.Store(Proceeding)
	case code >= 200 && code < 300:
		txn.terminate()
	case code >= 300 && code <= 699:
		txn.state.Store(Completed)
		txn.Send(msg.Ack())
		txn.fireTimerD()
	default:
		logger.Err("invalid proceed code %d", code)
	}
	txn.endpoint.TUConsume(msg)
}
