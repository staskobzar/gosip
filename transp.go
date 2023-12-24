package main

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/transport"
	"log"
	"time"
)

type message struct {
	b []byte
}

func (m *message) Ack() sip.Message     { return m }
func (m *message) IsResponse() bool     { return false }
func (m *message) TopViaBranch() string { return "" }
func (m *message) Method() string       { return "INVITE" }
func (m *message) ResponseCode() int    { return 0 }
func (m *message) Byte() []byte         { return m.b }

func main() {
	fmt.Println("===================================================")

	logger.Enable(true)

	mgr := transport.Init()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := mgr.ListenTCP(ctx, "127.0.0.1:5060"); err != nil {
		log.Fatal(err)
	}

	if err := mgr.ListenUDP(ctx, "127.0.0.1:5060"); err != nil {
		log.Fatal(err)
	}

	for pack := range mgr.Recv() {
		logger.Log("[<] rcv: from %q to %q msg %q\n", pack.Raddr, pack.Laddr, pack.Payload)
		msg := &message{
			b: append([]byte("recv: "), pack.Payload...),
		}
		<-time.After(time.Millisecond * 2)
		if err := mgr.Send(pack.Laddr, pack.Raddr, msg); err != nil {
			logger.Err(err.Error())
		}
		// msg.b = append([]byte("any: "), msg.b...)
		// mgr.Send(nil, pack.Raddr, msg)
	}
}
