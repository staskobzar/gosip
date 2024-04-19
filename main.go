package main

import (
	"context"
	"fmt"
	"gosip/pkg/dns"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction"
	"gosip/pkg/transport"
	"net"
)

func main() {
	fmt.Println("====================== SIPUA ========================")
	logger.Enable(true)

	mgr := transport.Init()
	resolv, err := dns.NewResolver("/etc/resolv.conf")
	if err != nil {
		panic(err)
	}

	mgr.WithResolver(resolv)
	mgr.ListenUDP(context.Background(), "192.168.1.102:5060")

	decoder := sipmsg.NewDecoder()
	txn := transaction.Init()

	TURecv := func(pack *sip.Packet) {
		if pack.Message == nil {
			logger.Err("TU received invalid packet with nil SIP message")
			return
		}
		r100 := pack.Message.Response(200, "OK")
		txn.RecvTU(&sip.Packet{
			SendTo:     []net.Addr{pack.RemoteSock},
			ReqAddrs:   pack.ReqAddrs,
			LocalSock:  pack.LocalSock,
			RemoteSock: pack.RemoteSock,
			Message:    r100,
		})
	}

	for {
		select {
		case pack := <-mgr.Recv():
			fmt.Printf("==> remote addr: %q, local addr: %q\n", pack.Raddr, pack.Laddr)
			decoder.Decode(pack.Payload, pack.Laddr, pack.Raddr)
		case pack := <-decoder.Recv():
			fmt.Printf("==> decoded pack: %#v\n", pack)
			mgr.ResolveRURI(pack)
		case pack := <-mgr.Resolved():
			fmt.Printf("==> resolved pack: %#v\n", pack)
			txn.RecvTransp(&pack)
		case pack := <-txn.SendTransp():
			fmt.Printf("==> TXN TO TRANSP SEND: %#v\n", pack)
			mgr.Send(pack)
		// transport manager must send pack
		case pack := <-txn.SendTU():
			fmt.Println("==> TU received message", pack)
			// txn to TU
			TURecv(pack)
		case err := <-decoder.Err():
			fmt.Printf("==> ERR DECODER: %s\n", err)
		case err := <-mgr.Err():
			fmt.Printf("==> ERR TRANSPORT: %s\n", err)
		case err := <-txn.Err():
			fmt.Printf("==> ERR TXN: %s\n", err)
		}
	}
}
