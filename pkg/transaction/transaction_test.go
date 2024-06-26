package transaction

import (
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockRegisterMsg() *sipmsg.Message {
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n\r\n"
	msg, _ := sipmsg.Parse(input)

	return msg
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

func mockNonInvite() (*Layer, *sip.Packet) {
	layer := Init()
	pack := &sip.Packet{
		Message: mockRegisterMsg(),
	}
	return layer, pack
}

func mockInvite() (*Layer, *sip.Packet) {
	layer := Init()
	pack := &sip.Packet{
		Message: mockInviteMsg(),
	}
	return layer, pack
}

func TestTransactionMatch(t *testing.T) {
	t.Run("client transaction", func(t *testing.T) {
		layer, pack := mockNonInvite()
		resp := pack.Message.Response(100, "Trying")
		txn := initClientNonInvite(pack, layer)

		t.Run("match", func(t *testing.T) {
			assert.True(t, txn.MatchClient(resp))
		})

		t.Run("not match", func(t *testing.T) {
			resp.Method = "CANCEL"
			assert.False(t, txn.MatchClient(resp))
		})
	})

	t.Run("server transaction match", func(t *testing.T) {
		setup := func() (*ServerNonInvite, *sipmsg.Message) {
			layer, pack := mockNonInvite()
			resp := pack.Message.Response(100, "Trying")
			txn := initServerNonInvite(pack, layer)
			return txn, resp
		}

		t.Run("match", func(t *testing.T) {
			txn, resp := setup()
			assert.True(t, txn.MatchServer(resp))
		})

		t.Run("false when no top via in txn request", func(t *testing.T) {
			txn, resp := setup()
			txn.req.Message.DelHeader("Via")
			assert.False(t, txn.MatchServer(resp))
		})

		t.Run("false when no top via in response", func(t *testing.T) {
			txn, resp := setup()
			resp.DelHeader("Via")
			assert.False(t, txn.MatchServer(resp))
		})

		t.Run("false when sentby not match", func(t *testing.T) {
			txn, resp := setup()
			via := resp.TopVia()
			via.Host = "example.com"
			assert.False(t, txn.MatchServer(resp))
		})

		t.Run("true when method is ACK", func(t *testing.T) {
			txn, resp := setup()
			resp.Method = "ACK"
			assert.True(t, txn.MatchServer(resp))
		})

		t.Run("false when branch id is invalid", func(t *testing.T) {
			txn, resp := setup()
			via := resp.TopVia()
			via.Branch = "foo-bar"
			assert.False(t, txn.MatchServer(resp))
		})
	})
}
