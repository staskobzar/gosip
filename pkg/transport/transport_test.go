package transport

import (
	"context"
	"crypto/tls"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockTLSCert() *tls.Config {
	var (
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

	cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
	if err != nil {
		panic(err)
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12}
}

func mockInviteMsg() *sipmsg.Message {
	input := "INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP atlanta.com;branch=z9hG4bK776a\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@atlanta.com\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
		"Contact: <sip:alice@atlanta.com>\r\n\r\n" +
		"v=0\r\no=jdoe 3724394400 3724394405 IN IP4 198.51.100.1\r\n" +
		"s=Call to Bob\r\nc=IN IP4 198.51.100.1\r\nt=0 0\r\n" +
		"m=audio 49170 RTP/AVP 0\r\nc=IN IP6 2001:db8::2\r\na=sendrecv\r\n"
	msg, _ := sipmsg.Parse(input)

	return msg
}

func mgrStart(ln func(*Manager) error) *Manager {

	mgr := Init()

	if err := ln(mgr); err != nil {
		panic(err)
	}

	<-time.After(time.Millisecond) // let store listen address in goroutine

	return mgr
}

func TestManagerSendUDP(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	srvStart := func() (net.Addr, chan []byte) {
		ln, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		assert.Nil(t, err)

		ch := make(chan []byte)

		go func() {
			buf := make([]byte, 1024)
			n, err := ln.Read(buf)
			assert.Nil(t, err)
			ch <- buf[:n]
		}()

		return ln.LocalAddr(), ch
	}

	mgr := mgrStart(func(mgr *Manager) error {
		return mgr.ListenUDP(ctx, "127.0.0.1:0")
	})

	addr, ch := srvStart()
	pack := &sip.Packet{
		SendTo:  []net.Addr{addr},
		Message: mockInviteMsg(),
	}

	t.Run("send successfully", func(t *testing.T) {
		mgr.Send(pack)
		assert.Equal(t, "INVITE sip:bob@biloxi.com SIP/2.0\r\n", string(<-ch)[:35])
		assert.Equal(t, 1, mgr.sock.Len())
	})

	t.Run("fail when no local listener socket found", func(t *testing.T) {
		for name, _ := range mgr.sock.pool {
			mgr.sock.Del(name)
		}

		assert.Equal(t, 0, mgr.sock.Len())

		err := mgr.SendUDP(addr, nil, pack.Message)
		assert.ErrorContains(t, err, "was not found")
	})

	t.Run("fail when found listener is not UDP", func(t *testing.T) {
		assert.Equal(t, 0, mgr.sock.Len())
		mgr.sock.Put(sockName(addr), &TCPListener{})

		err := mgr.SendUDP(addr, nil, pack.Message)
		assert.ErrorContains(t, err, "type is not UDPConn")
	})
}

func TestManagerSendTCP(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	srvStart := func() (net.Addr, chan []byte) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		assert.Nil(t, err)

		ch := make(chan []byte)

		go func() {
			conn, err := ln.Accept()
			assert.Nil(t, err)

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			assert.Nil(t, err)
			ch <- buf[:n]
		}()

		return ln.Addr(), ch
	}

	startTCPMgr := func() *Manager {
		return mgrStart(func(mgr *Manager) error {
			return mgr.ListenTCP(ctx, "127.0.0.1:0")
		})
	}

	t.Run("successfully sent when no conn in store", func(t *testing.T) {
		mgr := startTCPMgr()
		addr, ch := srvStart()
		pack := &sip.Packet{
			SendTo:  []net.Addr{addr},
			Message: mockInviteMsg(),
		}

		assert.Equal(t, 0, mgr.conn.Len())

		mgr.Send(pack)
		select {
		case err := <-mgr.Err():
			t.Errorf("failed send: %q", err.Err)
		case msg := <-ch:
			assert.Equal(t, "INVITE sip:bob@biloxi.com SIP/2.0\r\n", string(msg[:35]))
			assert.Equal(t, 1, mgr.conn.Len())
		}
	})

	t.Run("fail when found not TCP in the store", func(t *testing.T) {
		mgr := startTCPMgr()
		for name, _ := range mgr.conn.pool {
			mgr.conn.Del(name)
		}

		assert.Equal(t, 0, mgr.conn.Len())
		mgr.conn.Put("tcp:foo:5050", &UDP{})

		pack := &sip.Packet{
			SendTo:  []net.Addr{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 45656}},
			Message: mockInviteMsg(),
		}
		mgr.Send(pack)
		assert.ErrorContains(t, (<-mgr.Err()).Err, "type is not TCPConn")
	})

	t.Run("send to existing conn", func(t *testing.T) {
		mgr := startTCPMgr()
		addr, ch := srvStart()

		conn, err := net.DialTCP("tcp", nil, addr.(*net.TCPAddr))
		assert.Nil(t, err)

		tcp := &TCP{conn: conn}
		mgr.conn.Put(tcp.key(), tcp)

		pack := &sip.Packet{
			SendTo:    []net.Addr{addr},
			LocalSock: conn.LocalAddr(),
			Message:   mockInviteMsg(),
		}

		mgr.Send(pack)
		select {
		case err := <-mgr.Err():
			t.Errorf("failed send: %q", err.Err)
		case msg := <-ch:
			assert.Equal(t, "INVITE sip:bob@biloxi.com SIP/2.0\r\n", string(msg[:35]))
		}
	})
}
