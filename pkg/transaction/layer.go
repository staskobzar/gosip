package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/transaction/internal/pool"
)

type LayerError struct{}

type Layer struct {
	pool      *pool.Pool
	sndTransp chan *sip.Packet // send to transport
	sndTU     chan *sip.Packet // send to TU
	err       chan error
}

func Init() *Layer {
	return &Layer{
		pool:      pool.New(),
		sndTransp: make(chan *sip.Packet, 12),
		sndTU:     make(chan *sip.Packet, 12),
		err:       make(chan error),
	}
}

func (l *Layer) SendTransp() <-chan *sip.Packet {
	return l.sndTransp
}

func (l *Layer) SendTU() <-chan *sip.Packet {
	return l.sndTU
}

func (l *Layer) Err() <-chan error {
	return l.err
}

func (l *Layer) passToTU(pack *sip.Packet) {
	select {
	case l.sndTU <- pack:
		logger.Log("txn:layer: pass packet to TU")
	default:
		logger.Err("txn:layer: failed to send to TU")
	}
}

func (l *Layer) passToTransp(pack *sip.Packet) {
	// TODO: control blocking
	logger.Log("txn:layer: to transport message %q", pack.Message.FirstLine())
	l.sndTransp <- pack
}

func (l *Layer) RecvTU(pack *sip.Packet) {
	logger.Log("transaction received message from TU")
	if pack.Message == nil {
		logger.Err("empty Message received from TU")
		return
	}

	// if request then try to create a new client transaction
	if pack.Message.IsRequest() {
		logger.Log("create and add new client transaction")
		l.pool.Add(func() pool.Transaction {
			if pack.Message.IsInvite() {
				return initClientInvite(pack, l.sndTransp, l.sndTU, l.err)
			}
			return initClientNonInvite(pack, l.sndTransp, l.sndTU, l.err)
		}())
		return
	}

	if txn, ok := l.pool.Match(pack.Message); ok {
		logger.Log("found transaction and trying to consume message")
		txn.Consume(pack)
		return
	}

	logger.Err("transaction %q does not exists", pack.Message.FirstLine())
}

func (l *Layer) Destroy(txn pool.Transaction) {
	logger.Log("txn:layer:pool: destroy transaction %q", txn.BranchID())
	l.pool.Delete(txn)
}

func (l *Layer) RecvTransp(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("empty Message received on %q from %q", pack.LocalSock, pack.RemoteSock)
		return
	}
	if pack.Message.IsRequest() {
		logger.Log("txn:layer: received request from %q", pack.RemoteSock)
		l.serverTxn(pack)
		return
	}

	l.clientTxn(pack)
}

// process incoming SIP packet from network
func (l *Layer) serverTxn(pack *sip.Packet) {
	logger.Log("txn:layer: start server transaction processing")
	if txn, ok := l.pool.Match(pack.Message); ok {
		txn.Consume(pack)
		return
	}

	logger.Log("txn:layer: create and store new server transaction with branch %q",
		pack.Message.TopViaBranch())
	l.pool.Add(func() pool.Transaction {
		if pack.Message.IsInvite() {
			return initServerInvite(pack, l.sndTransp, l.sndTU, l.err)
		}
		return initServerNonInvite(pack, l)
	}())
}

func (l *Layer) clientTxn(pack *sip.Packet) {
	logger.Wrn("TODO: client transaction is not yet implemented")
}
