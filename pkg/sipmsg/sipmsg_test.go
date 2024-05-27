package sipmsg

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoderDercode(t *testing.T) {
	makeSIPMsg := func(vh []string) []byte {
		vias := ""
		if len(vh) > 0 {
			vias = strings.Join(vh, "")
		}
		input := "INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
			vias +
			"Max-Forwards: 70\r\n" +
			"To: Bob <sip:bob@biloxi.com>\r\n" +
			"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
			"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
			"CSeq: 314159 INVITE\r\n" +
			"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
			"Contact: <sip:alice@pc33.atlanta.com>\r\n\r\n"

		return []byte(input)
	}
	laddr := &net.UDPAddr{IP: net.ParseIP("10.100.100.1"), Port: 5060}
	raddr := &net.UDPAddr{IP: net.ParseIP("10.200.200.2"), Port: 8060}
	dec := NewDecoder()
	// "Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds.1\r\n" +
	// "Via: SIP/2.0/UDP pc32.atlanta.com;branch=z9hG4bK776asdhds.2\r\n" +
	// "Via: SIP/2.0/UDP pc31.atlanta.com;branch=z9hG4bK776asdhds.3\r\n" +
	t.Run("fail on invalid SIP message", func(t *testing.T) {
		input := makeSIPMsg([]string{"foobar invalid packet\r\n"})
		dec.Decode(input, laddr, raddr)
		err := <-dec.Err()
		assert.ErrorIs(t, err, ErrMsgParse)
		assert.ErrorContains(t, err, "SIP Message parse")
		assert.Zero(t, len(dec.rcv))
	})

	t.Run("fail on missing top Via header", func(t *testing.T) {
		input := makeSIPMsg(nil)
		dec.Decode(input, laddr, raddr)
		err := <-dec.Err()
		assert.ErrorIs(t, err, ErrDecoder)
		assert.ErrorContains(t, err, "invalid top Via")
		assert.Zero(t, len(dec.rcv))
	})

	t.Run("failed when remote address is nil", func(t *testing.T) {
		input := makeSIPMsg([]string{"Via: SIP/2.0/UDP atlanta.com;branch=z9hG4bK776\r\n"})
		dec.Decode(input, laddr, nil)
		err := <-dec.Err()
		assert.ErrorIs(t, err, ErrDecoder)
		assert.ErrorContains(t, err, "invalid remote address")
		assert.Zero(t, len(dec.rcv))
	})

	t.Run("add received when Via sent-by is domain name", func(t *testing.T) {
		input := makeSIPMsg([]string{"Via: SIP/2.0/UDP atlanta.com;branch=z9hG4bK776\r\n"})
		dec.Decode(input, laddr, raddr)
		assert.Zero(t, len(dec.err))
		pack := <-dec.Recv()

		assert.Equal(t, pack.Laddr, laddr)
		assert.Equal(t, pack.Raddr, raddr)
		via := pack.Message.TopVia()
		assert.Equal(t, "10.200.200.2", via.Recvd)
	})

	t.Run("add received when Via sent-by IP is defferent from source address", func(t *testing.T) {
		input := makeSIPMsg([]string{"Via: SIP/2.0/UDP 172.0.0.1;branch=z9hG4bK776\r\n"})
		dec.Decode(input, laddr, raddr)
		assert.Zero(t, len(dec.err))
		pack := <-dec.Recv()

		assert.Equal(t, pack.Laddr, laddr)
		assert.Equal(t, pack.Raddr, raddr)
		via := pack.Message.TopVia()
		assert.Equal(t, "10.200.200.2", via.Recvd)
	})

	t.Run("do not add received when Via sent-by is IP that match source address", func(t *testing.T) {
		input := makeSIPMsg([]string{"Via: SIP/2.0/UDP 10.200.200.2;branch=z9hG4bK776\r\n"})
		dec.Decode(input, laddr, raddr)
		assert.Zero(t, len(dec.err))
		pack := <-dec.Recv()

		assert.Equal(t, pack.Laddr, laddr)
		assert.Equal(t, pack.Raddr, raddr)
		via := pack.Message.TopVia()
		assert.Empty(t, via.Recvd)
	})
}
