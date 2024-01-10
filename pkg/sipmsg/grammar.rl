%%{
# SIP rfc3261 grammar definitions with some constraints
machine grammar;

unreserved      = alnum | [\-_.!~*'()]; # alpha | digit | mark
escaped         = "%" xdigit xdigit;
user_unreserved = [&=+$,;?/];
hexseq          = xdigit{1,4} ( ":" xdigit{1,4} )*;
hexpart         = hexseq | hexseq "::" hexseq? | "::" hexseq?;

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

# CONSTRAINT: only sip/sips scheme limited
scheme          = ( "sip"i | "sips"i );
# CONSTRAINT: userinfo is missing telephone-subscriber for simplicity
userinfo        = user ( ":" password )?;
hostport        = host ( ":" port )?;
uri_params      = uri_param ( ";" uri_param )*;
uri_headers     = uri_header ( "&" uri_header)*;

}%%
