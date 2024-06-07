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

func (txn *ServerInvite) Consume(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("txn:server:inv: consume packet has <nil> SIP Message")
		return
	}

	msg := pack.Message

	if msg.IsRequest() {
		if txn.response != nil {
			txn.layer.passToTransp(txn.response)
		} else {
			txn.mu.Lock()
			defer txn.mu.Unlock()
			txn.response = txn.gen100()
			txn.layer.passToTransp(txn.response)
		}
		return
	}

	if txn.state.IsProceeding() {
		if msg.IsProvisional() {
			txn.mu.Lock()
			txn.response = pack
			txn.mu.Unlock()
			txn.layer.passToTransp(pack)
			return
		}
	}
}

func (txn *ServerInvite) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	if txn.MatchServer(msg) {
		return txn, true
	}
	return nil, false
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
