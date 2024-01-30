package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageFind(t *testing.T) {
	//nolint:goconst
	input := "INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds.1\r\n" +
		"Via: SIP/2.0/UDP pc32.atlanta.com;branch=z9hG4bK776asdhds.2\r\n" +
		"Via: SIP/2.0/UDP pc31.atlanta.com;branch=z9hG4bK776asdhds.3\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
		"Contact: <sip:alice@pc33.atlanta.com>\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	t.Run("first available when single header", func(t *testing.T) {
		found := msg.Find(HAllow)
		hdr, ok := found.(*HeaderGeneric)
		assert.True(t, ok)
		assert.Equal(t, "Allow", hdr.HeaderName)
		assert.Equal(t, "INVITE, ACK, OPTIONS, CANCEL, BYE", hdr.Value)
	})

	t.Run("first available when multiple header", func(t *testing.T) {
		found := msg.Find(HVia)
		via, ok := found.(*HeaderVia)
		assert.True(t, ok)
		assert.Equal(t, "Via", via.HeaderName)
		assert.Equal(t, "z9hG4bK776asdhds.1", via.Branch)
		assert.Equal(t, "pc33.atlanta.com", via.Host)
	})

	t.Run("FindAll headers list", func(t *testing.T) {
		list := msg.FindAll(HVia)
		assert.Equal(t, 3, list.Len())
		via, ok := list[2].(*HeaderVia)
		assert.True(t, ok)
		assert.Equal(t, "z9hG4bK776asdhds.3", via.Branch)
	})

	t.Run("none existing headers", func(t *testing.T) {
		assert.Nil(t, msg.Find(HWWWAuthenticate))

		list := msg.FindAll(HWWWAuthenticate)
		assert.Equal(t, 0, list.Len())
	})
}

func TestMessageFindByName(t *testing.T) {
	//nolint:goconst
	input := "INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds.1\r\n" +
		"Via: SIP/2.0/UDP pc32.atlanta.com;branch=z9hG4bK776asdhds.2\r\n" +
		"Via: SIP/2.0/UDP pc31.atlanta.com;branch=z9hG4bK776asdhds.3\r\n" +
		"X-Foo: v1\r\n" +
		"Max-Forwards: 70\r\n" +
		"X-Foo: v2\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
		"X-Foo: v3\r\n" +
		"Contact: <sip:alice@pc33.atlanta.com>\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	t.Run("find first parsed", func(t *testing.T) {
		hdr := msg.FindByName("X-Foo").(*HeaderGeneric)
		assert.Equal(t, "v1", hdr.Value)
	})

	t.Run("find all as list", func(t *testing.T) {
		list := msg.FindByNameAll("X-Foo")
		assert.Equal(t, 3, list.Len())
	})

	t.Run("no headres found", func(t *testing.T) {
		hdr := msg.FindByName("X-Unknown")
		assert.Nil(t, hdr)

		list := msg.FindByNameAll("X-Unknown")
		assert.Equal(t, 0, list.Len())
	})
}

func TestMessageAppendInsert(t *testing.T) {
	createMsg := func() *Message {
		//nolint:goconst
		input := "INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
			"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds.1\r\n" +
			"Max-Forwards: 70\r\n" +
			"To: Bob <sip:bob@biloxi.com>\r\n" +
			"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
			"Via: SIP/2.0/UDP pc31.atlanta.com;branch=z9hG4bK776asdhds.3\r\n" +
			"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
			"CSeq: 314159 INVITE\r\n" +
			"Allow: INVITE, ACK, OPTIONS, CANCEL, BYE\r\n" +
			"X-Foo-Bar: v1-ff-foo-bar\r\n" +
			"Contact: <sip:alice@pc33.atlanta.com>\r\n\r\n"

		msg, err := Parse(input)
		if err != nil {
			panic(err)
		}
		return msg
	}

	t.Run("append via to via", func(t *testing.T) {
		msg := createMsg()
		assert.Equal(t, 10, msg.HLen())
		assert.Equal(t, 2, msg.FindAll(HVia).Len())
		msg.Append(HVia, &HeaderVia{HeaderName: "Via", Proto: "SIP/2.0/", Transp: "UDP", Host: "append.com"})
		assert.Equal(t, 11, msg.HLen())
		assert.Equal(t, 3, msg.FindAll(HVia).Len())
		hdr := msg.Headers[1].(*HeaderVia)
		assert.Equal(t, "append.com", hdr.Host)
	})

	t.Run("append generic header to via", func(t *testing.T) {
		msg := createMsg()
		assert.Equal(t, 10, msg.HLen())
		assert.Equal(t, 1, msg.FindByNameAll("X-Foo-Bar").Len())
		msg.Append(HVia, &HeaderGeneric{HeaderName: "X-Foo-Bar", Value: "x-foo-value", T: HGeneric})
		assert.Equal(t, 11, msg.HLen())
		assert.Equal(t, 2, msg.FindByNameAll("X-Foo-Bar").Len())
		hdr := msg.Headers[1].(*HeaderGeneric)
		assert.Equal(t, "x-foo-value", hdr.Value)
	})

	t.Run("append generic to the end of a message", func(t *testing.T) {
		msg := createMsg()
		assert.Equal(t, 10, msg.HLen())
		assert.Equal(t, 0, msg.FindByNameAll("X-Bar").Len())
		msg.Append(HGeneric, &HeaderGeneric{HeaderName: "X-Bar", Value: "x-bar-value", T: HGeneric})
		assert.Equal(t, 11, msg.HLen())
		assert.Equal(t, 1, msg.FindByNameAll("X-Bar").Len())
		hdr := msg.Headers[10].(*HeaderGeneric)
		assert.Equal(t, "x-bar-value", hdr.Value)
	})

	t.Run("insert before CallID header", func(t *testing.T) {
		msg := createMsg()
		assert.Equal(t, 10, msg.HLen())
		assert.Equal(t, 0, msg.FindByNameAll("X-Bar").Len())
		msg.Insert(HCallID, &HeaderGeneric{HeaderName: "X-Bar", Value: "xx-bar-value", T: HGeneric})
		assert.Equal(t, 11, msg.HLen())
		assert.Equal(t, 1, msg.FindByNameAll("X-Bar").Len())
		hdr := msg.Headers[5].(*HeaderGeneric)
		assert.Equal(t, "xx-bar-value", hdr.Value)
	})

	t.Run("insert header to the very top", func(t *testing.T) {
		msg := createMsg()
		assert.Equal(t, 10, msg.HLen())
		assert.Equal(t, 0, msg.FindByNameAll("X-Foo").Len())
		msg.Insert(HGeneric, &HeaderGeneric{HeaderName: "X-Foo", Value: "x-foo-value", T: HGeneric})
		assert.Equal(t, 11, msg.HLen())
		assert.Equal(t, 1, msg.FindByNameAll("X-Foo").Len())
		hdr := msg.Headers[0].(*HeaderGeneric)
		assert.Equal(t, "x-foo-value", hdr.Value)
	})
}

func TestMessageRequestString(t *testing.T) {
	//nolint:goconst
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)

	assert.Equal(t, input, msg.String())
}

func TestMessageRequestNone200Ack(t *testing.T) {
	//nolint:goconst
	inputReq := "INVITE sip:bob@biloxi.example.com SIP/2.0\r\n" +
		"Via: SIP/2.0/TCP client.atlanta.example.com:5060;branch=z9hG4bK74b43\r\n" +
		"Max-Forwards: 70\r\n" +
		"Route: <sip:ss1.atlanta.example.com;lr>\r\n" +
		"Route: <sip:ss2.atlanta.example.com;lr>\r\n" +
		"From: Alice <sip:alice@atlanta.example.com>;tag=9fxced76sl\r\n" +
		"To: Bob <sip:bob@biloxi.example.com>\r\n" +
		"Call-ID: 3848276298220188511@atlanta.example.com\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Contact: <sip:alice@client.atlanta.example.com;transport=tcp>\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 151\r\n\r\n"

	//nolint:goconst
	inputResp := "SIP/2.0 407 Proxy Authorization Required\r\n" +
		"Via: SIP/2.0/TCP client.atlanta.example.com:5060;branch=z9hG4bK74b43;received=192.0.2.101\r\n" +
		"From: Alice <sip:alice@atlanta.example.com>;tag=9fxced76sl\r\n" +
		"To: Bob <sip:bob@biloxi.example.com>;tag=3flal12sf\r\n" +
		"Call-ID: 3848276298220188511@atlanta.example.com\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Proxy-Authenticate: Digest realm=\"atlanta.example.com\", qop=\"auth\"," +
		"nonce=\"f84f1cec41e6cbe5aea9c8e88d359\", opaque=\"\", stale=FALSE, algorithm=MD5\r\n" +
		"Content-Length: 0\r\n\r\n"

	wantAck := "ACK sip:bob@biloxi.example.com SIP/2.0\r\n" +
		"Via: SIP/2.0/TCP client.atlanta.example.com:5060;branch=z9hG4bK74b43\r\n" +
		"Max-Forwards: 70\r\n" +
		"Route: <sip:ss1.atlanta.example.com;lr>\r\n" +
		"Route: <sip:ss2.atlanta.example.com;lr>\r\n" +
		"From: Alice <sip:alice@atlanta.example.com>;tag=9fxced76sl\r\n" +
		"To: Bob <sip:bob@biloxi.example.com>;tag=3flal12sf\r\n" +
		"Call-ID: 3848276298220188511@atlanta.example.com\r\n" +
		"CSeq: 1 ACK\r\n" +
		"Content-Length: 0\r\n\r\n"

	req, err := Parse(inputReq)
	assert.Nil(t, err)

	resp, err := Parse(inputResp)
	assert.Nil(t, err)

	ack := req.ACK(resp)
	assert.False(t, ack.IsResponse())
	assert.Equal(t, "ACK", ack.Method)
	assert.Equal(t, "sip:bob@biloxi.example.com", ack.RURI.String())
	assert.Equal(t, "3848276298220188511@atlanta.example.com", ack.CallID)
	assert.Equal(t, req.From.String(), ack.From.String())
	assert.Equal(t, resp.To.String(), ack.To.String())
	assert.Equal(t, req.Find(HVia).String(), ack.Find(HVia).String())
	assert.Equal(t, req.CSeq, ack.CSeq)
	assert.Equal(t, 2, ack.FindAll(HRoute).Len())
	assert.Equal(t, req.Find(HRoute).String(), ack.Find(HRoute).String())

	assert.Equal(t, wantAck, ack.String())
}

func TestMessageResponseString(t *testing.T) {
	//nolint:goconst
	input := "SIP/2.0 200 OK\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7 ;received=192.0.2.4\r\n" +
		"Route: <sip:alice@atlanta.com>,<sip:bob@biloxi.com>\r\n" +
		"To: Bob <sip:bob@biloxi.com>;tag=2493k59kd\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n" +
		"Content-Length: 0\r\n\r\n"

	msg, err := Parse(input)
	assert.Nil(t, err)

	assert.Equal(t, input, msg.String())
	assert.Equal(t, []byte(input), msg.Byte())
}

func BenchmarkMessageToString(b *testing.B) {
	//nolint:goconst
	input := "REGISTER sip:registrar.biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP bobspc.biloxi.com:5060;branch=z9hG4bKnashds7\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Bob <sip:bob@biloxi.com>;tag=456248\r\n" +
		"Call-ID: 843817637684230@998sdasdh09\r\n" +
		"CSeq: 1826 REGISTER\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Expires: 7200\r\n\r\n"

	msg, _ := Parse(input)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = msg.String()
	}
}
