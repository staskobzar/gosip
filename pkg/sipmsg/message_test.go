package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageFind(t *testing.T) {
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
		assert.Equal(t, "Allow", hdr.Name)
		assert.Equal(t, "INVITE, ACK, OPTIONS, CANCEL, BYE", hdr.Value)
	})

	t.Run("first available when multiple header", func(t *testing.T) {
		found := msg.Find(HVia)
		via, ok := found.(*HeaderVia)
		assert.True(t, ok)
		assert.Equal(t, "Via", via.Name)
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

func TestMessageRequestString(t *testing.T) {
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

func TestMessageResponseString(t *testing.T) {
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
}

func BenchmarkMessageToString(b *testing.B) {
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
