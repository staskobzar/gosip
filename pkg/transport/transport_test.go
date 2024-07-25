package transport

import (
	"context"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestTransportSendUDP(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	mgrStart := func(ln func(*Manager) error) *Manager {

		mgr := Init()

		err := ln(mgr)
		assert.Nil(t, err)

		<-time.After(time.Millisecond) // let store listen address in goroutine

		return mgr
	}

	t.Run("send using UDP", func(t *testing.T) {
		t.Parallel()

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

		mgr.Send(pack)
		assert.Equal(t, "INVITE sip:bob@biloxi.com SIP/2.0\r\n", string(<-ch)[:35])
	})

	t.Run("send creating TCP new connection", func(t *testing.T) {
		t.Parallel()

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

		mgr := mgrStart(func(mgr *Manager) error {
			return mgr.ListenTCP(ctx, "127.0.0.1:0")
		})

		addr, ch := srvStart()
		pack := &sip.Packet{
			SendTo:  []net.Addr{addr},
			Message: mockInviteMsg(),
		}

		// assert.Equal(t, 0, mgr.conn.Len())

		mgr.Send(pack)
		select {
		case err := <-mgr.Err():
			t.Errorf("failed send: %q", err.Err)
		case msg := <-ch:
			assert.Equal(t, "INVITE sip:bob@biloxi.com SIP/2.0\r\n", string(msg[:35]))
			// assert.Equal(t, 1, mgr.conn.Len())
		}
	})
}
