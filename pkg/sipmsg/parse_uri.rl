// -*-go-*-
//
// SIP URI parser

package sipmsg

import (
	"fmt"
)

%% machine uri;
%% write data;

func ParseURI(data string) (*URI, error) {
	uri := &URI{}
	m   := 0 // marker
	cs  := 0 // current state
	p   := 0 // data pointer
	pe  := len(data) // data end pointer
	eof := len(data)

	%%{
		action sm       { m = p }
		action scheme   { uri.Scheme   = data[:p] }
		action userinfo { uri.Userinfo = data[m:p] }
		action hostport { uri.Hostport = data[m:p] }
		action params   { uri.Params   = data[m:p] }
		action headers  { uri.Headers  = data[m:p] }

		include grammar "grammar.rl";

		main := scheme %scheme ":" ( userinfo >sm %userinfo "@" )?
		        hostport >sm %hostport ( ";" uri_params >sm %params )?
		        ( "?" uri_headers >sm %headers )?;
	}%%

	%% write init;
	%% write exec;

	if cs >= uri_first_final {
		return uri, nil
	}

	if p == pe {
		return nil, fmt.Errorf("%w: unexpected eof: %q", ErrURIParse, data)
	}

	return nil, fmt.Errorf("%w: error in uri at pos %d: %q>>%q<<", ErrURIParse, p, data[:p],data[p:])
}
