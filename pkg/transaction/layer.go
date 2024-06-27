package transaction

import (
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/transaction/pool"
	"gosip/pkg/transaction/timer"
	"time"
)

// transaction layer errors
var (
	Error      = errors.New("transaction")
	ErrTimeout = fmt.Errorf("%w: timeout", Error)
	ErrTxnFail = fmt.Errorf("%w: transaction failure: no ACK received", Error)
)

// Layer for SIP transactions
type Layer struct {
	SetupTimers func(*timer.Timer) *timer.Timer // callback function to setup transactions timers
	pool        *pool.Pool
	sndTransp   chan *sip.Packet // send to transport
	sndTU       chan *sip.Packet // send to TU
	err         chan error
}

// Init creates new layer
func Init() *Layer {
	return &Layer{
		pool:      pool.New(),
		sndTransp: make(chan *sip.Packet, 12),
		sndTU:     make(chan *sip.Packet, 12),
		err:       make(chan error, 1),
	}
}

// SendTransp channel for packets to send to transport
func (l *Layer) SendTransp() <-chan *sip.Packet {
	return l.sndTransp
}

// SendTU channel for packets to send to TU
func (l *Layer) SendTU() <-chan *sip.Packet {
	return l.sndTU
}

// Err channel for errors
func (l *Layer) Err() <-chan error {
	return l.err
}

// Destroy transaction for SIP packet on layer level
func (l *Layer) Destroy(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("layer Destroy got invalid SIP message")
		return
	}
	branchID := pack.Message.TopViaBranch()
	logger.Log("lookup transaction %q to destroy", branchID)

	txn, found := l.pool.Get(branchID)

	if !found {
		logger.Err("transaction %q not found", branchID)
		return
	}

	txn.Terminate()
}

// RecvTU process packets received from TU
func (l *Layer) RecvTU(pack *sip.Packet) {
	logger.Log("transaction received message from TU")
	if pack.Message == nil {
		logger.Err("empty Message received from TU")
		return
	}

	// if request then try to create a new client transaction
	if pack.Message.IsRequest() {
		logger.Log("create and add new client transaction")
		err := l.pool.Add(func() sip.Transaction {
			if pack.Message.IsInvite() {
				return initClientInvite(pack, l)
			}
			return initClientNonInvite(pack, l)
		}())
		if err != nil {
			logger.Err("%s", err)
		}
		return
	}

	if txn, ok := l.pool.Match(pack.Message); ok {
		logger.Log("found transaction and trying to consume message")
		txn.Consume(pack)
		return
	}

	logger.Err("transaction %q does not exists", pack.Message.FirstLine())
}

// RecvTransp to process packages received from network
func (l *Layer) RecvTransp(pack *sip.Packet) {
	if pack.Message == nil {
		logger.Err("empty Message received on %q from %q", pack.LocalSock, pack.RemoteSock)
		return
	}

	if txn, ok := l.pool.Match(pack.Message); ok {
		txn.Consume(pack)
		return
	}

	if pack.Message.IsRequest() {
		logger.Log("txn:layer: received request from %q", pack.RemoteSock)
		l.serverTxn(pack)
		return
	}
}

// process incoming SIP packet from network
func (l *Layer) serverTxn(pack *sip.Packet) {
	logger.Log("txn:layer: start server transaction processing")
	logger.Log("txn:layer: create and store new server transaction with branch %q",
		pack.Message.TopViaBranch())

	err := l.pool.Add(func() sip.Transaction {
		if pack.Message.IsInvite() {
			return initServerInvite(pack, l)
		}
		return initServerNonInvite(pack, l)
	}())
	if err != nil {
		logger.Err("%s", err)
	}
}

func (l *Layer) passToTU(pack *sip.Packet) {
	select {
	case l.sndTU <- pack:
		logger.Log("txn:layer: pass packet to TU")
	case <-time.After(500 * time.Millisecond):
		logger.Err("txn:layer: failed to send to TU")
	}
}

func (l *Layer) passToTransp(pack *sip.Packet) {
	logger.Log("txn:layer: to transport message %q", pack.Message.FirstLine())
	select {
	case l.sndTransp <- pack:
	case <-time.After(500 * time.Millisecond):
		logger.Err("txn:layer: failed to send message to transport. Channel blocked")
	}
}

func (l *Layer) passErr(err error) {
	select {
	case l.err <- err:
	case <-time.After(500 * time.Millisecond):
		logger.Err("txn:layer: failed to send error. Channel blocked")
	}
}
