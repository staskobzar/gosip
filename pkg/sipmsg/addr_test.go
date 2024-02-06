package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseContactHeaders(t *testing.T) {
	t.Run("single header", func(t *testing.T) {
		hdr := "Contact: \"Bob C\" <sip:bob@192.0.2.4>;q=0.5;foo=bar;expires=1800"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)

		cnt := msg.Find(HContact).(*HeaderContact)
		assert.Equal(t, "Contact", cnt.HeaderName)
		assert.Equal(t, "\"Bob C\"", cnt.DisplayName)
		assert.Equal(t, "sip:bob@192.0.2.4", cnt.Addr.String())
		assert.Equal(t, "q=0.5;foo=bar;expires=1800", cnt.Params.str())
		assert.Equal(t, "0.5", cnt.Q)
		assert.Equal(t, "1800", cnt.Expires)
		assert.Nil(t, cnt.Next)
		assert.Equal(t, len(hdr), cnt.Len())
	})

	t.Run("start value", func(t *testing.T) {
		hdr := "Contact: *"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)
		cnt := msg.Find(HContact).(*HeaderContact)
		assert.Equal(t, "*", cnt.Params.str())
		assert.Equal(t, len(hdr), cnt.Len())
	})

	t.Run("linked header", func(t *testing.T) {
		hdr := "m: <sip:100@192.0.2.4>,<sip:100@10.0.0.4:4555>;q=0.5, <sip:100@[2041:0:140F::875B:131B]>;q=0.8"
		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)

		cnt := msg.Find(HContact).(*HeaderContact)
		assert.Equal(t, "m", cnt.HeaderName)
		assert.Equal(t, "sip:100@192.0.2.4", cnt.Addr.String())
		assert.Equal(t, "", cnt.Params.String())

		cnt = cnt.Next
		assert.Equal(t, "sip:100@10.0.0.4:4555", cnt.Addr.String())
		assert.Equal(t, "0.5", cnt.Q)

		cnt = cnt.Next
		assert.Equal(t, "sip:100@[2041:0:140F::875B:131B]", cnt.Addr.String())
		assert.Equal(t, "0.8", cnt.Q)
		assert.Nil(t, cnt.Next)
		assert.Equal(t, len(hdr)-1, // parse with space between contacts but stringify without
			msg.Find(HContact).(*HeaderContact).Len())
	})

	t.Run("multiple headers", func(t *testing.T) {
		hdr := "m: <sip:100@192.0.2.4>\r\n" +
			"Contact: <sip:caller@u1.space.com>;q=0.1\r\n" +
			"Contact: <sips:caller@u2.space.com>;q=0.3"

		msg, err := Parse(toMsg([]string{hdr}))
		assert.Nil(t, err)
		list := msg.FindAll(HContact)
		assert.Equal(t, 3, list.Len())
	})
}

func TestParseFromToHeader(t *testing.T) {
	t.Run("parse to member and headers list successfully", func(t *testing.T) {
		hdrs := "To: \"Alice Home\" <sip:alice@biloxi.com>;user=phone;tag=ff00aa\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456248;day=monday;free"
		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		assert.Equal(t, 1, msg.FindAll(HFrom).Len())
		assert.Equal(t, 1, msg.FindAll(HTo).Len())

		from := msg.Find(HFrom).(*NameAddr)
		assert.Same(t, msg.From, from)
		to := msg.Find(HTo).(*NameAddr)
		assert.Same(t, msg.To, to)

		assert.Equal(t, "To", to.HeaderName)
		assert.Equal(t, `"Alice Home"`, to.DisplayName)
		assert.Equal(t, "sip:alice@biloxi.com", to.Addr.String())
		assert.Equal(t, "ff00aa", to.Tag)
		assert.Equal(t, "user=phone;tag=ff00aa", to.Params.str())
		assert.Equal(t, len(to.String()), to.Len())

		assert.Equal(t, "From", from.HeaderName)
		assert.Equal(t, "Bob", from.DisplayName)
		assert.Equal(t, "sip:bob@biloxi.com", from.Addr.String())
		assert.Equal(t, "456248", from.Tag)
		assert.Equal(t, "tag=456248;day=monday;free", from.Params.str())
		assert.Equal(t, len(from.String()), from.Len())
	})

	t.Run("fail when more then one To header", func(t *testing.T) {
		hdrs := "To: \"Alice Home\" <sip:alice@biloxi.com>;tag=ff00aa\r\n" +
			"t: <sip:alice@biloxi.com>\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456249;day=monday;free"

		_, err := Parse(toMsg([]string{hdrs}))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "more then one To headers found")
	})

	t.Run("fail when more then one From header", func(t *testing.T) {
		hdrs := "To: \"Alice Home\" <sip:alice@biloxi.com>;tag=ff00aa\r\n" +
			"f: <sip:alice@biloxi.com>\r\n" +
			"From: Bob <sip:bob@biloxi.com>;tag=456249;day=monday;free"

		_, err := Parse(toMsg([]string{hdrs}))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "more then one From headers found")
	})
}

func TestParseRoutingHeaders(t *testing.T) {
	t.Run("multiple headers", func(t *testing.T) {
		hdrs := "Record-Route: <sip:h1.domain.com;lr>;host=one\r\n" +
			"Record-Route: H2 <sip:h2.domain.com>\r\n" +
			"Record-Route: <sip:h3.domain.com>\r\n" +
			"Route: <sip:s1.pbx.com>\r\n" +
			"Route: <sip:s2.pbx.com>"

		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		assert.Equal(t, 3, msg.FindAll(HRecordRoute).Len())
		assert.Equal(t, 2, msg.FindAll(HRoute).Len())
	})

	t.Run("linked record-route headers", func(t *testing.T) {
		hdrs := "Record-Route: <sip:h1.domain.com;lr>;host=one\r\n" +
			"Record-Route: <sip:h2.domain.com;lr>, <sip:dd1.pbx.com>;user=pbx,<sips:dd2.pbx.com>\r\n" +
			"Record-Route: <sip:h3.domain.com>"
		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		list := msg.FindAll(HRecordRoute)
		assert.Equal(t, 3, list.Len())
		r := list[1].(*Route)
		assert.Equal(t, len(r.String()), r.Len())

		assert.Equal(t, "Record-Route", r.HeaderName)
		assert.Equal(t, "sip:h2.domain.com;lr", r.Addr.String())
		assert.Equal(t, "", r.Params.String())
		assert.NotNil(t, r.Next)

		r = r.Next
		assert.Equal(t, "sip:dd1.pbx.com", r.Addr.String())
		assert.Equal(t, "user=pbx", r.Params.str())
		assert.NotNil(t, r.Next)

		r = r.Next
		assert.Equal(t, "sips:dd2.pbx.com", r.Addr.String())
		assert.Equal(t, "", r.Params.String())
		assert.Nil(t, r.Next)
	})

	t.Run("linked route headers", func(t *testing.T) {
		hdrs := "Route: <sip:s1.pbx.com;lr>, <sip:h100.sip.com:5060>;now\r\n" +
			"Route: <sip:s2.pbx.com>"

		msg, err := Parse(toMsg([]string{hdrs}))
		assert.Nil(t, err)
		list := msg.FindAll(HRoute)
		assert.Equal(t, 2, list.Len())
		r := list[0].(*Route)
		assert.Equal(t, len(r.String()), r.Len())

		assert.Equal(t, "Route", r.HeaderName)
		assert.Equal(t, "sip:s1.pbx.com;lr", r.Addr.String())
		assert.Equal(t, "", r.Params.String())
		assert.NotNil(t, r.Next)

		r = r.Next
		assert.Equal(t, "sip:h100.sip.com:5060", r.Addr.String())
		assert.Equal(t, "now", r.Params.str())
		assert.Nil(t, r.Next)
	})
}

func TestRouteString(t *testing.T) {
	tests := map[string]struct {
		route *Route
		want  string
	}{
		`simple route header`: {&Route{NameAddrSpec: NameAddrSpec{
			HeaderName: "Route",
			Addr:       &URI{Scheme: "sip", Hostport: "10.0.0.1"},
		}}, "Route: <sip:10.0.0.1>"},
		`with display name and params`: {&Route{NameAddrSpec: NameAddrSpec{
			HeaderName:  "Record-Route",
			DisplayName: "\"PBX f1\"",
			Addr:        &URI{Scheme: "sip", Hostport: "10.0.0.1"},
			Params:      Params("replica=true"),
		}}, "Record-Route: \"PBX f1\" <sip:10.0.0.1>;replica=true"},
		`route with linked header`: {
			&Route{
				NameAddrSpec: NameAddrSpec{
					HeaderName: "Record-Route",
					Addr:       &URI{Scheme: "sip", Hostport: "p1.sip.com", Params: Params("lr")},
				},
				Next: &Route{NameAddrSpec: NameAddrSpec{Addr: &URI{Scheme: "sips", Hostport: "p2.sip.com", Params: Params("lr")}}},
			}, "Record-Route: <sip:p1.sip.com;lr>,<sips:p2.sip.com;lr>",
		},
		`route with two linked header`: {
			&Route{
				NameAddrSpec: NameAddrSpec{
					HeaderName: "Route",
					Addr:       &URI{Scheme: "sip", Hostport: "p1.sip.com"},
				},
				Next: &Route{
					NameAddrSpec: NameAddrSpec{Addr: &URI{Scheme: "sips", Hostport: "p2.sip.com"}},
					Next:         &Route{NameAddrSpec: NameAddrSpec{Addr: &URI{Scheme: "sip", Hostport: "p3.sip.com"}}},
				},
			}, "Route: <sip:p1.sip.com>,<sips:p2.sip.com>,<sip:p3.sip.com>",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.route.String())
		})
	}
}

func TestFromToString(t *testing.T) {
	tests := []struct {
		addr *NameAddr
		want string
	}{
		{
			&NameAddr{NameAddrSpec: NameAddrSpec{
				HeaderName:  "From",
				DisplayName: "Alice",
				Addr:        &URI{Scheme: "sip", Userinfo: "alice", Hostport: "atlanta.com"},
				Params:      Params("tag=88sja8x"),
			}},
			"From: Alice <sip:alice@atlanta.com>;tag=88sja8x",
		},
		{
			&NameAddr{NameAddrSpec: NameAddrSpec{
				HeaderName: "f",
				Addr:       &URI{Scheme: "sip", Userinfo: "+12125551212", Hostport: "server.phone2net.com"},
				Params:     Params("tag=887s;user=bob"),
			}},
			"f: <sip:+12125551212@server.phone2net.com>;tag=887s;user=bob",
		},
		{
			&NameAddr{NameAddrSpec: NameAddrSpec{
				HeaderName:  "To",
				DisplayName: "\"Carol Chicago\"",
				Addr:        &URI{Scheme: "sip", Userinfo: "carol", Hostport: "chicago.com"},
			}},
			"To: \"Carol Chicago\" <sip:carol@chicago.com>",
		},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.want, tc.addr.String())
	}
}

func TestHeaderContactString(t *testing.T) {
	tests := []struct {
		cnt  *HeaderContact
		want string
	}{
		{
			&HeaderContact{NameAddrSpec: NameAddrSpec{
				HeaderName: "Contact",
				Addr:       &URI{Scheme: "sip", Userinfo: "alice", Hostport: "atlanta.com"},
				Params:     Params("expires=3600"),
			}},
			"Contact: <sip:alice@atlanta.com>;expires=3600",
		},
		{
			&HeaderContact{NameAddrSpec: NameAddrSpec{
				HeaderName: "Contact",
				Params:     Params("*"),
			}},
			"Contact: *",
		},
		{
			&HeaderContact{NameAddrSpec: NameAddrSpec{
				HeaderName:  "m",
				DisplayName: "Caller",
				Addr:        &URI{Scheme: "sip", Userinfo: "caller", Hostport: "u1.privspace.com", Params: Params("transport=UDP")},
			}},
			"m: Caller <sip:caller@u1.privspace.com;transport=UDP>",
		},
		{
			&HeaderContact{
				NameAddrSpec: NameAddrSpec{
					HeaderName:  "Contact",
					DisplayName: "\"Mr. Watson\"",
					Addr:        &URI{Scheme: "sip", Userinfo: "watson", Hostport: "ch.bell.com"},
					Params:      Params("q=0.7; expires=3600"),
				},
				Next: &HeaderContact{
					NameAddrSpec: NameAddrSpec{
						DisplayName: "Watson",
						Addr:        &URI{Scheme: "sips", Userinfo: "watson", Hostport: "bell.com"},
						Params:      Params("q=0.1"),
					},
				},
			},
			"Contact: \"Mr. Watson\" <sip:watson@ch.bell.com>;q=0.7; expires=3600,Watson <sips:watson@bell.com>;q=0.1",
		},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.want, tc.cnt.String())
	}
}
