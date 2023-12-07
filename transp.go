package main

import (
	"context"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/transport"
	"log"
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

	mgm := transport.InitManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := mgm.ListenUDP(ctx, "127.0.0.1:5060"); err != nil {
		log.Fatal(err)
	}

	// select {
	// case pack := <-mgm.Recv():
	for pack := range mgm.Recv() {
		fmt.Printf("%#v\n", mgm)
		fmt.Printf("%#v\n", pack)
		msg := &message{
			b: append([]byte("recv: "), pack.Payload...),
		}
		mgm.Send(pack.Laddr, pack.Raddr, msg)
		msg.b = append([]byte("any: "), msg.b...)
		mgm.Send(nil, pack.Raddr, msg)
	}
	// <-time.After(time.Second)
	// log.Println("========== end ==========")
}
