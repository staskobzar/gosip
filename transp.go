package main

import (
	"context"
	"fmt"
	"gosip/pkg/dns"
	"gosip/pkg/logger"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transport"
	"net"
	"time"
)

// type message struct {
// 	b []byte
// }

// func (m *message) Ack() sip.Message     { return m }
// func (m *message) IsResponse() bool     { return false }
// func (m *message) TopViaBranch() string { return "" }
// func (m *message) Method() string       { return "INVITE" }
// func (m *message) ResponseCode() int    { return 0 }
// func (m *message) Byte() []byte         { return m.b }

func mainDISABLE() {
	fmt.Println("===================================================")

	logger.Enable(true)

	mgr := transport.Init()
	resolv, err := dns.NewResolver("/etc/resolv.conf")
	if err != nil {
		panic(err)
	}
	mgr.WithResolver(resolv)
	mgr.ListenTCP(context.Background(), ":0")
	mgr.ListenUDP(context.Background(), ":0")

	<-time.After(time.Second)

	for _, sipuri := range []string{
		"sip:alice@apple.com", "sip:alice@apple.com;transport=UDP", "sips:bob@10.0.0.1",
		"sip:100@localhost:5858", "sips:1000@172.19.1.2:5858", "sip:200@clusterpbx.com",
		"sips:300@okon.ferry.clusterpbx.xyz", "sip:foobar.com", "sip:400@sip.google.com",
		"sip:700@amazon.com", "sips:xzxzy.dfg", "sip:222@250.250.250.255:5858",
	} {
		fmt.Printf("-------------------------------------------------- %q\n", sipuri)
		uri, err := sipmsg.ParseURI(sipuri)
		if err != nil {
			fmt.Printf("[!] ERR: failed parse uri: %s\n", err)
			continue
		}
		fmt.Println(mgr.Resolve(uri))
	}

	list, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, eth := range list {
		addr, err := eth.Addrs()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %s\n", eth.Name, addr)
	}
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// if err := mgr.ListenTCP(ctx, "127.0.0.1:5060"); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := mgr.ListenUDP(ctx, "127.0.0.1:5060"); err != nil {
	// 	log.Fatal(err)
	// }

	// for pack := range mgr.Recv() {
	// 	logger.Log("[<] rcv: from %q to %q msg %q\n", pack.Raddr, pack.Laddr, pack.Payload)
	// 	<-time.After(time.Millisecond * 1)
	// 	msg := &message{
	// 		b: append([]byte("recv: "), pack.Payload...),
	// 	}
	// 	if err := mgr.Send(pack.Laddr, pack.Raddr, msg); err != nil {
	// 		logger.Err(err.Error())
	// 	}
	// 	// msg.b = append([]byte("any: "), msg.b...)
	// 	// mgr.Send(nil, pack.Raddr, msg)
	// }
}
