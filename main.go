package main

import (
	"context"
	"fmt"
	"gosip/pkg/dns"
	"gosip/pkg/logger"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transport"
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

	for {
		select {
		case pack := <-mgr.Recv():
			fmt.Printf("remote addr: %#v\n", pack.Raddr)
			fmt.Printf("local  addr: %#v\n", pack.Laddr)
			decoder.Decode(pack.Payload, pack.Laddr, pack.Raddr)
		case pack := <-decoder.Recv():
			fmt.Printf("decoded pack: %#v\n", pack)
			mgr.ResolveRURI(pack)
		case pack := <-mgr.Resolved():
			fmt.Printf("resolved pack: %#v\n", pack)
			// txn.Consume(pack)
		// case pack := <-txn.Rcv():
		// 	fmt.Printf("TXN RCV: %#v\n", pack)
		case err := <-decoder.Err():
			fmt.Printf("ERR DECODER: %s\n", err)
		case err := <-mgr.Err():
			fmt.Printf("ERR TRANSPORT: %s\n", err)
			// case err := <-txn.Err():
			// 	fmt.Printf("ERR TXN: %s\n", err)
		}
	}
}
