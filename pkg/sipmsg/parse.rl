// -*-go-*-
//

package sipmsg

import (
	"fmt"
	"strconv"
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
	var via *HeaderVia
	var naddr nameAddr
	var cnt *HeaderContact
	var route *HeaderRoute

	%%{
		action sm     { m = p }
		action sm1    { m1 = p }
		action body   { msg.Body = data[m:p] }

		include grammar    "parser/grammar.rl";
		include first_line "parser/first_line.rl";
		include headers    "parser/headers.rl";

		body = extend+ >sm %body;

		main := first_line CRLF headers* CRLF body?;
	}%%

	%% write init;
	%% write exec;

	if cs >= sipmsg_first_final {
		return msg, nil
	}

	if p == pe {
		return nil, fmt.Errorf("%w: unexpected eof: %s...", ErrMsgParse, data[:p])
	}

	return nil, fmt.Errorf("%w: error in uri at pos %d: %q>>>%q", ErrMsgParse, p, data[:p],data[p:])
}

// simplified string to number converter for parser
// not checking convert errors because parser already
// essures that string is a number
func atoi(num string) int {
	n, _ := strconv.Atoi(num)
	return n
}
