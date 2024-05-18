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
	// mgr.ListenUDP(context.Background(), "192.168.1.102:5060")
	mgr.ListenUDP(context.Background(), "10.54.197.36:5060")

	decoder := sipmsg.NewDecoder()
	txn := transaction.Init()

	TURecv := func(pack *sip.Packet) {
		if pack.Message == nil {
			logger.Err("TU received invalid packet with nil SIP message")
			return
		}
		if pack.Message.IsResponse() {
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

	// send notify as client
	txn.RecvTU(&sip.Packet{
		// 8.1.1.7 When the UAC creates a request, it MUST insert a Via
		Message: notifyReq(),
	})

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

func notifyReq() *sipmsg.Message {
	domain := "alice@clusterpbx.xyz;transport=UDP"
	input := "OPTIONS sip:" + domain + " SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 199.182.135.220:5060;branch=z9hG4bK66746c85;rport\r\n" +
		"Max-Forwards: 70\r\n" +
		"From: <sip:" + domain + ">;tag=as01e75d0c\r\n" +
		"To: <sip:" + domain + ">\r\n" +
		"Contact: <sip:asterisk@192.168.1.102>\r\n" +
		"Call-ID: 0190718f4d2fcfd514f931d359586c24@192.168.1.102\r\n" +
		"CSeq: 102 OPTIONS\r\n" +
		"User-Agent: ClearlyIP PBX\r\n" +
		"Allow: INVITE, ACK, CANCEL, OPTIONS, BYE, REFER, SUBSCRIBE, NOTIFY, INFO, PUBLISH, MESSAGE\r\n" +
		"Supported: replaces, timer\r\n" +
		"Content-Length: 0\r\n\r\n"

	msg, _ := sipmsg.Parse(input)
	via := msg.TopVia()
	via.Transp = ""
	via.Host = ""
	via.Port = ""

	return msg
}
