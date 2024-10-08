package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"gosip/pkg/dns"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"gosip/pkg/transaction"
	"gosip/pkg/transport"
	"net"
)

const (
	certPem = `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`
	keyPem = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`
)

func main() {
	fmt.Println("================= transport.ListernTLS ==============")
	logger.Enable(true)

	cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
	if err != nil {
		panic(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	mgr := transport.Init()
	err = mgr.ListenTLS(context.Background(), "127.0.0.1:5061", config)
	if err != nil {
		panic(err)
	}
	pack := <-mgr.Recv()
	fmt.Printf("%#v\n", pack)
	fmt.Printf("%s\n", pack.Payload)
}

func main1() {
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
		//Message: notifyReq(),
		Message: inviteReq(),
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

func inviteReq() *sipmsg.Message {
	domain := "alice@clusterpbx.xyz;transport=UDP"
	input := "INVITE sip:" + domain + " SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 199.182.135.220:5060;branch=z9hG4bK66746c85\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Alice <sip:" + domain + ">\r\n" +
		"From: Bob <sip:bob@clusterpbx.xyz>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@clusterpbx.xyz\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
		"Contact: <sip:bob@clusterpbx.xyz>\r\n\r\n" +
		"v=0\r\no=jdoe 3724394400 3724394405 IN IP4 198.51.100.1\r\n" +
		"s=Call to Bob\r\nc=IN IP4 198.51.100.1\r\nt=0 0\r\n" +
		"m=audio 49170 RTP/AVP 0\r\nc=IN IP6 2001:db8::2\r\na=sendrecv\r\n"
	msg, _ := sipmsg.Parse(input)
	return msg
}
