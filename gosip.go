package gosip

import "net"


type UDP struct{}
type TCP struct{}

type UA struct {
	tu *TU
}

type Transp interface{}

func mainBLOCK() {

	ua := BuildUA().
		WithTransport(&UDP{}, "0.0.0.0:5060").
		WithTransport(&TCP{}, "0.0.0.0:5061").
		WithResolver(&DNSResolver{})

	if err := ua.Start(); err != nil {
		panic(err)
	}

	ua.AddAccount("sip:alice@pbx.com", "pa55w0rd")

	ua.OnCall(func() {})

	ua.Call("sip:123@pbx.com", func() {})
}

func BuildUA() *UA                                     { return &UA{} }
func (ua *UA) WithTransport(_ Transp, addr string) *UA { return ua }

// --------------------------------------------------

// = transport start

type Transport interface {
	Send(addr net.Addr)
	IsReliable() bool
}

// --------------------------------------------------

type Txn struct {}
type TxnList map[string]Txn
type TU struct{
	client struct{
		inv TxnList
		noninv TxnList
	}
	srv struct {
		inv TxnList
		noninv TxnList
	}
}

ua := &UA{tu: &TU{}}

// With call back direct
ua.Call("sip:alice@1.0.1.1:5656").
	OnEarly(func(msg *MSG){}).
	OnConfirm(func(msg*MSG){}).
	OnRedirect(func(msg *MSG){}).
	OnError(func(msg *MSG){})

// With callback handler
type Handler interface{
	OnEarly(msg *MSG)
	OnConfirm(msg *MSG)
	OnRedirect(msg *MSG)
	OnError(msg*MSG)
}
type callHandler struct{}
func (h*callHandler) OnEarly(msg *MSG) {}
func (h*callHandler) OnConfirm(msg *MSG) {}
func (h*callHandler) OnRedirect(msg *MSG) {}
func (h*callHandler) OnError(msg*MSG) {}

ua.Call("sip:100@pbx.com", &callHandler{})

// --------------------------------------------------
