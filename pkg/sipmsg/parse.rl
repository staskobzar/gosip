// -*-go-*-
//

package sipmsg

import (
	"fmt"
)

func Parse(data string) (*Message, error) {
	%% machine sipmsg;
	%% write data;
	msg := NewMessage()
	m   := 0 // marker
	m1  := 0 // additional marker
	cs  := 0 // current state
	p   := 0 // data pointer
	pe  := len(data) // data end pointer
	eof := len(data)
	var htype HType
	var hdrname string
	var via, linkVia *HeaderVia
	var naddr nameAddr
	var cnt *HeaderContact
	var route *Route

	%%{
		action sm     { m = p }
		action sm1    { m1 = p }
		action body   { msg.Body = data[m:p] }

		include grammar    "grammar.rl";
		include first_line "first_line.rl";
		include headers    "headers.rl";

		body = extend+ >sm %body;

		main := first_line CRLF headers* CRLF body?;
	}%%

	%% write init;
	%% write exec;

	if cs >= sipmsg_first_final {
		return msg, nil
	}

	if p == pe {
		return nil, fmt.Errorf("%w: unexpected eof: %q", ErrMsgParse, data)
	}

	return nil, fmt.Errorf("%w: error in uri at pos %d: %q>>%q<<", ErrMsgParse, p, data[:p],data[p:])
}
