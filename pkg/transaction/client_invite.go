package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
)

type ClientInvite struct {
	*Transaction
}

func initClientInvite(pack *sip.Packet, sndTransp chan *sip.Packet, sndTU chan *sip.Packet, err chan error) *ClientInvite {
	logger.Log("TODO: !! new ClientInvite transaction and init !!")
	return nil
}

func (txn *ClientInvite) Init()               {}
func (txn *ClientInvite) Consume(*sip.Packet) {}
func (txn *ClientInvite) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	if txn.MatchClient(msg) {
		return txn, true
	}
	return nil, false
}
