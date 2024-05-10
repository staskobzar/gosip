package transaction

import (
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
)

type ServerInvite struct {
	*Transaction
}

func initServerInvite(pack *sip.Packet, sndTransp chan *sip.Packet, sndTU chan *sip.Packet, err chan error) *ServerInvite {
	logger.Log("TODO: !! new ServerInvite transaction and init !!")
	return nil
}

func (txn *ServerInvite) Init()               {}
func (txn *ServerInvite) Consume(*sip.Packet) {}
func (txn *ServerInvite) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	if txn.MatchServer(msg) {
		return txn, true
	}
	return nil, false
}
