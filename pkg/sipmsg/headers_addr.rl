%%{
# parsing grammar rules for a family of headers
# which contain named address value like From, To etc
machine headers_addr;

action hdr_from_init {
    msg.From = NewNameAddr(data[m:p])
    naddr = msg.From
}
action hdr_to_init {
    msg.To = NewNameAddr(data[m:p])
    naddr = msg.To
}
action hdr_cnt_init {
    cnt = NewHeaderContact(data[m:p])
    naddr = cnt
    msg.Contact = append(msg.Contact, cnt)
}
action hdr_cnt_link {
    cnt.Next = NewHeaderContact("")
    cnt = cnt.Next
    naddr = cnt
}
action hdr_rroute_init {
    route = NewRoute(data[m:p])
    msg.RecRoute = append(msg.RecRoute, route)
    naddr = route
}
action hdr_route_init {
    route = NewRoute(data[m:p])
    msg.Route = append(msg.Route, route)
    naddr = route
}
action hdr_route_link {
    route.Next = NewRoute("")
    route = route.Next
    naddr = route
}

action display_name   { naddr.setDisplayName(data[m:p]) }
action naddr_scheme   { naddr.setURIScheme(data[m:p]) }
action naddr_userinfo { naddr.setURIUserinfo(data[m:p]) }
action naddr_hostport { naddr.setURIHostport(data[m:p]) }
action naddr_params   { naddr.setURIParams(data[m:p]) }
action naddr_headers  { naddr.setURIHeaders(data[m:p]) }
action hdr_naddr_tag  { naddr.setParam("tag", data[m:p]) }

action hdr_naddr_prms { naddr.setParam("params", data[m1:p]) }

action contact_star   { naddr.setParam("params", "*") }
action hdr_cnt_q      { naddr.setParam("q", data[m:p]) }
action hdr_cnt_expr   { naddr.setParam("expires", data[m:p]) }

display_name   = (token LWS)* | quoted_string;
name_addr_tag  = "tag"i EQUAL token >sm %hdr_naddr_tag;
name_addr_prm  = name_addr_tag | generic_param;
addr_spec      = scheme >sm %naddr_scheme ":" ( userinfo >sm %naddr_userinfo "@" )?
                 hostport >sm %naddr_hostport ( ";" uri_params >sm %naddr_params )?
                 ( "?" uri_headers >sm %naddr_headers )?;
name_addr      = ( display_name >sm %display_name)? LAQUOT addr_spec RAQUOT;

name_addr_spec = name_addr | addr_spec;

cnt_prm_expr   = "expires"i EQUAL digit+ >sm %hdr_cnt_expr;
cnt_prm_q      = "q"i EQUAL qvalue >sm %hdr_cnt_q;

contact_prms   = cnt_prm_q | cnt_prm_expr | generic_param;
hdr_cnt_prm    = name_addr_spec %sm1 (SEMI contact_prms)* %hdr_naddr_prms;
hdr_cnts_list  = hdr_cnt_prm (COMMA %hdr_cnt_link hdr_cnt_prm)*;
contact_value  = STAR %contact_star | hdr_cnts_list;

route_prm      = name_addr %sm1 (SEMI generic_param)* %hdr_naddr_prms;

# headers machines
hdr_from    = ("From"i | "f"i) >sm %hdr_from_init HCOLON $(hdr,1) name_addr_spec %sm1
              (SEMI name_addr_prm)* %hdr_naddr_prms;

hdr_to      = ("To"i | "t"i) >sm %hdr_to_init HCOLON $(hdr,1) name_addr_spec %sm1
              (SEMI name_addr_prm)* %hdr_naddr_prms;

hdr_contact = ("Contact"i | "m"i) >sm %hdr_cnt_init HCOLON $(hdr,1) contact_value;

hdr_rroute  = "Record-Route"i >sm %hdr_rroute_init HCOLON $(hdr,1) route_prm (COMMA %hdr_route_link route_prm)*;

hdr_route   = "Route"i >sm %hdr_route_init HCOLON $(hdr,1) route_prm (COMMA %hdr_route_link route_prm)*;
}%%
