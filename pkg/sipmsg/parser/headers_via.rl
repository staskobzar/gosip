%%{
machine headers_via;

action hdr_via_new {
    via = NewHeaderVia(data[m:p])
	msg.Headers = append(msg.Headers, via)
}

action hdr_via_link {
    linkVia = NewHeaderVia("") // no name for linked comma via
    via.Via = linkVia
    via = linkVia
}

action via_proto  { via.Proto  = data[m:p] }
action via_transp { via.Transp = data[m:p] }
action via_host   { via.Host   = data[m:p] }
action via_port   { via.Port   = data[m:p] }
action via_branch { via.Branch = data[m:p] }
action via_recvd  { via.Recvd  = data[m:p] }
action via_hprm   { via.Params = data[m1:p] }

proto_name    = "SIP" SLASH "2.0" SLASH;
sent_protocol = proto_name >sm %via_proto token >sm %via_transp;
sent_by       = host >sm %via_host ( COLON port >sm %via_port )?;
via_ttl       = "ttl"i EQUAL digit{3};
via_maddr     = "maddr"i EQUAL host;
via_recvd     = "received"i EQUAL (IPv4address | IPv6address) >sm %via_recvd;
via_branch    = "branch"i EQUAL token >sm %via_branch;
via_hprms     = via_ttl | via_maddr | via_recvd | via_branch | generic_param;
via_param     = sent_protocol LWS sent_by %sm1 ( SEMI via_hprms )* %via_hprm;

hdr_via = ("Via"i | "v"i) >sm %hdr_via_new HCOLON $(hdr,1) via_param ( COMMA via_param >hdr_via_link )*;
}%%
