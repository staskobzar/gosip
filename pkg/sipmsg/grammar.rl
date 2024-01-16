%%{
# SIP rfc3261 grammar definitions with some constraints
machine grammar;

SP        = 0x20; # space
CRLF      = "\r\n";
HTAB      = 0x09;
DQUOTE    = 0x22;
WSP       = SP | HTAB;
LWS       = ( WSP* CRLF )? WSP+;  # linear whitespace
SWS       = LWS?;                 # sep whitespace
HCOLON    = WSP* ":" SWS;
SEMI      = SWS ";" SWS;
SLASH     = SWS "/" SWS;
COMMA     = SWS "," SWS;
COLON     = SWS ":" SWS;
EQUAL     = SWS "=" SWS;
LAQUOT    = SWS "<";
RAQUOT    = ">" SWS;
STAR      = SWS "*" SWS;

UTF8_CONT       = 0x80..0xBF;
UTF8_NONASCII   = ( 0xC0..0xDF UTF8_CONT{1} ) |
                  ( 0xE0..0xEF UTF8_CONT{2} ) |
                  ( 0xF0..0xF7 UTF8_CONT{3} ) |
                  ( 0xF8..0xFB UTF8_CONT{4} ) |
                  ( 0xFC..0xFD UTF8_CONT{5} );

unreserved      = alnum | [\-_.!~*'()]; # alpha | digit | mark
escaped         = "%" xdigit xdigit;
user_unreserved = [&=+$,;?/];
hexseq          = xdigit{1,4} ( ":" xdigit{1,4} )*;
hexpart         = hexseq | hexseq "::" hexseq? | "::" hexseq?;
token           = ( alnum | [\-.!%*_+`'~] )+;
word            = ( token | [()<>:\\/\[\]?{}] )+;

password        = ( unreserved | escaped | [&=+$,] )*;
user            = ( unreserved | escaped | user_unreserved )+;
paramchar       = [[\]/:&+$] | unreserved | escaped;
hdrchar         = [[\]/?:+$] | unreserved | escaped;
IPv4address     = digit{1,3} "." digit{1,3} "." digit{1,3} "." digit{1,3};
IPv6address     = hexpart ( ":" IPv4address )?;
IPv6reference   = "[" IPv6address "]";
domainlabel     = alnum | alnum ( alnum | "-" )* alnum;
toplabel        = alpha | alpha ( alnum | "-" )* alnum;
hostname        = ( domainlabel "." )* toplabel "."?;
port            = digit{1,5};
host            = hostname | IPv4address | IPv6reference;
uri_param       = paramchar+ ("=" paramchar+)?;
uri_header      = hdrchar+ "=" hdrchar*;

qdtext          =  LWS | 0x21 | 0x23..0x5B | 0x5D..0x7E | UTF8_NONASCII;
quoted_pair     = '\\' (0x00..0x09 | 0x0B..0x0C | 0x0E..0x7F);
quoted_string   = SWS DQUOTE ( qdtext | quoted_pair )* DQUOTE;
gen_value       = token | host | quoted_string;
generic_param   = token ( EQUAL gen_value)?;
qvalue          = [01] ("." digit{,3})?;

# CONSTRAINT: only sip/sips scheme limited
scheme          = ( "sip"i | "sips"i );
# CONSTRAINT: userinfo is missing telephone-subscriber for simplicity
userinfo        = user ( ":" password )?;
hostport        = host ( ":" port )?;
uri_params      = uri_param ( ";" uri_param )*;
uri_headers     = uri_header ( "&" uri_header)*;
}%%
