package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction/internal/pool"
)

type ClientNonInvite struct {
	*Transaction
}

func initClientNonInvite(pack *sip.Packet, sndTransp chan *sip.Packet, sndTU chan *sip.Packet, err chan error) *ClientNonInvite {
	logger.Log("TODO: !! new ClientNonInvite transaction and init !!")
	return nil
}

func (txn *ClientNonInvite) Init()               {}
func (txn *ClientNonInvite) Consume(*sip.Packet) {}
func (txn *ClientNonInvite) Match(msg *sipmsg.Message) (pool.Transaction, bool) {
	return nil, false
}
