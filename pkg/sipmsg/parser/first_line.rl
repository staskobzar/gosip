%%{
machine first_line;

action first_line_requ {
    msg.t      = HRequest
    msg.Method = data[:p]
    msg.RURI   = &URI{}
}
action first_line_resp {
    msg.t = HResponse
}
action first_line_code   { msg.Code   = data[m:p] }
action first_line_reason { msg.Reason = data[m:p] }
action ruri_scheme       { msg.RURI.Scheme   = data[m:p] }
action ruri_userinfo     { msg.RURI.Userinfo = data[m:p] }
action ruri_hostport     { msg.RURI.Hostport = data[m:p] }
action ruri_params       { msg.RURI.Params   = Params(data[m:p]).setup() }
action ruri_headers      { msg.RURI.Headers  = data[m:p] }
action uri_trp           { msg.RURI.Transport = data[m1:p] }

# CONSTRAINT: fixed version
sipver            = "SIP/2.0";
status_code       = [1-6][0-9]{2};
transp_param      = "transport="i token >sm1 %uri_trp;
param             = transp_param | uri_params;
params            = param ( ";" param )*;
request_line_ruri = scheme >sm %ruri_scheme ":" ( userinfo >sm %ruri_userinfo "@" )?
                    hostport >sm %ruri_hostport ( ";" params >sm %ruri_params )?
                    ( "?" uri_headers >sm %ruri_headers )?;
status_line       = sipver %first_line_resp SP
                    status_code >sm %first_line_code SP (any+ -- CRLF) >sm %first_line_reason;
request_line      = token %first_line_requ SP request_line_ruri SP sipver;

first_line = status_line | request_line;
}%%
