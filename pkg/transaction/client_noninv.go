package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"net/netip"
	"time"
)

type ClientNonInvite struct {
	Basic
}

func createClientNonInvite(transp sip.Transport, endpoint EndPoint, msg *sipmsg.Message) *ClientNonInvite {
	return &ClientNonInvite{
		Basic: initBasicTxn(transp, endpoint, msg),
	}
}

func (txn *ClientNonInvite) Init(msg *sipmsg.Message, addr netip.AddrPort) {
	txn.addr = addr
	txn.trying(msg)
	txn.fireTimerF(msg)
}

func (txn *ClientNonInvite) Consume(msg *sipmsg.Message) {
	if !msg.IsResponse() {
		return
	}

	code := msg.ResponseCode()
	logger.Log("client non invite txn consume response %d", code)

	switch txn.state.Load() {
	case Trying, Proceeding:
		if code >= 100 && code <= 199 {
			txn.state.Store(Proceeding)
		} else {
			txn.complete()
		}
	case Completed:
		logger.Log("absorbing response %d retransmission", code)
	default:
		logger.Err("unknown state %d", txn.state.Load())
	}
}

func (txn *ClientNonInvite) trying(msg *sipmsg.Message) {
	txn.state.Store(Trying)
	txn.Send(msg)

	if txn.transp.IsReliable() {
		return // no E timer for reliable transport
	}
	txn.fireTimerE(msg)
}

func (txn *ClientNonInvite) complete() {
	txn.state.Store(Completed)

	if txn.transp.IsReliable() {
		txn.terminate()
		return // no K timer for reliable transport
	}

	txn.fireTimerK()
}

func (txn *ClientNonInvite) fireTimerE(msg *sipmsg.Message) {
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

func (txn *ClientNonInvite) fireTimerF(msg *sipmsg.Message) {
	go func() {
		select {
		case <-time.After(txn.timer.T1 * 64):
			state := txn.state.Load()
			if state == Trying || state == Proceeding {
				txn.endpoint.Error(ErrTimeout, msg)
				txn.terminate()
			}
		case <-txn.halt:
		}
	}()
}

// just buffer any additional response retransmissions
func (txn *ClientNonInvite) fireTimerK() {
	go func() {
		select {
		case <-time.After(txn.timer.T4):
			txn.terminate()
		case <-txn.halt:
		}
	}()
}
