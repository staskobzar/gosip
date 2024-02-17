// -*-go-*-
//
// SIP URI parser

package sipmsg

import (
	"fmt"
)

func ParseURI(data string) (*URI, error) {
	%% machine uri;
	%% write data;
	uri := &URI{}
	m   := 0 // marker
	m1  := 0 // additional marker
	cs  := 0 // current state
	p   := 0 // data pointer
	pe  := len(data) // data end pointer
	eof := len(data)

	%%{
		action sm       { m = p }
		action sm1      { m1 = p }
		action scheme   { uri.Scheme    = data[:p] }
		action userinfo { uri.Userinfo  = data[m1:p] }
		action hostport { uri.Hostport  = data[m:p] }
		action params   { uri.Params    = Params(data[m:p]).setup() }
		action headers  { uri.Headers   = data[m:p] }
		action uri_trp  { uri.Transport = data[m1:p] }

		include grammar "parser/grammar.rl";

		# // set special marker m1 for userinfo because its machine contains ";"
		# // and same char is used in uri_params as a boarder condition
		# // it helps to handle correctly usernames in urls like "sip:alice;day=tuesday@sip.com"
		transp_param = "transport="i token >sm1 %uri_trp;
		param        = transp_param | uri_params;
		params       = param ( ";" param )*;
		main := scheme %scheme ":" ( userinfo >sm1 %userinfo "@" )?
				hostport >sm %hostport ( ";" params >sm %params )?
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
