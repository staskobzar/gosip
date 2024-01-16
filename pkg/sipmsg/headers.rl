%%{
machine headers;

action hdr_name  { hdrname = data[m:p] }
action hdr_value { msg.pushHeader(htype, hdrname, data[m:p]) }

action hdr_callid   { msg.CallID = data[m:p]; msg.pushHeader(HCallID, hdrname, data[m:p]) }
action hdr_clen     { msg.ContentLen = data[m:p] }
action hdr_ctyp     { msg.ContentType = data[m:p] }
action hdr_cseq     { msg.CSeq = data[m:p] }
action hdr_maxfwd   { msg.MaxFwd = data[m:p] }
action hdr_expires  { msg.Expires = data[m:p] }

callid        = word ( "@" word )?;
generic_value = ( extend+ -- CRLF );

include headers_via  "headers_via.rl";
include headers_addr "headers_addr.rl";

# CONSTRAINT: simplified generic header without new lines
hdr_generic   = token >sm %hdr_name %{htype = HGeneric} HCOLON $(hdr,0)
                generic_value >sm %hdr_value;
hdr_accept    = "Accept"i >sm %hdr_name %{htype = HAccept} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_acc_enc   = "Accept-Encoding"i >sm %hdr_name %{htype = HAcceptEncoding}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_acc_lang  = "Accept-Language"i >sm %hdr_name %{htype = HAcceptLanguage}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_alertinf  = "Alert-Info"i >sm %hdr_name %{htype = HAlertInfo}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_allow     = "Allow"i >sm %hdr_name %{htype = HAllow} HCOLON $(hdr,1)
                %sm token (COMMA token)* %hdr_value;
hdr_auth_info = "Authentication-Info"i >sm %hdr_name %{htype = HAuthenticationInfo}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_authz     = "Authorization"i >sm %hdr_name %{htype = HAuthorization}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_callid    = ("Call-ID"i | "i"i) >sm %hdr_name HCOLON $(hdr,1) callid >sm %hdr_callid;
hdr_callinfo  = "Call-Info"i >sm %hdr_name %{htype = HCallInfo} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_cnt_disp  = "Content-Disposition"i >sm %hdr_name %{htype = HContentDisposition}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_cnt_enc   = ("Content-Encoding"i | "e"i) >sm %hdr_name %{htype = HContentEncoding}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_cnt_lang  = "Content-Language"i >sm %hdr_name %{htype = HContentLanguage}
                HCOLON $(hdr,1) generic_value >sm %hdr_value;
hdr_clen      = ("Content-Length"i | "l"i) HCOLON $(hdr,1) digit+ >sm %hdr_clen;
hdr_cseq      = "CSeq"i HCOLON $(hdr,1) digit+ >sm %hdr_cseq LWS token;
hdr_ctyp      = ("Content-Type"i | "c"i) HCOLON $(hdr,1)
                (token SLASH token (SEMI token SLASH token)*) >sm %hdr_ctyp;
hdr_date      = "Date"i >sm %hdr_name %{htype = HDate} HCOLON $(hdr,1)
                (alnum | [,:] | SP)+ >sm %hdr_value;
hdr_errinfo   = "Error-Info"i >sm %hdr_name %{htype = HErrorInfo} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_expires   = "Expires"i HCOLON $(hdr,1) digit+ >sm %hdr_expires;
hdr_inreply   = "In-Reply-To"i >sm %hdr_name %{htype = HInReplyTo} HCOLON $(hdr,1)
                callid >sm (COMMA callid)* %hdr_value;
hdr_maxfwd    = "Max-Forwards"i HCOLON $(hdr,1) digit+ >sm %hdr_maxfwd;
hdr_mimever   = "MIME-Version"i >sm %hdr_name %{htype = HMIMEVersion} HCOLON $(hdr,1)
                (digit+ "." digit+) >sm %hdr_value;
hdr_minexpr   = "Min-Expires"i >sm %hdr_name %{htype = HMinExpires} HCOLON $(hdr,1)
                digit+ >sm %hdr_value;
hdr_org       = "Organization"i >sm %hdr_name %{htype = HOrganization} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_priority  = "Priority"i >sm %hdr_name %{htype = HPriority} HCOLON $(hdr,1)
                token >sm %hdr_value;
hdr_pxy_authn = "Proxy-Authenticate"i >sm %hdr_name %{htype = HProxyAuthenticate} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_pxy_authz = "Proxy-Authorization"i >sm %hdr_name %{htype = HProxyAuthorization} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_pxy_req   = "Proxy-Require"i >sm %hdr_name %{htype = HProxyRequire} HCOLON $(hdr,1)
                token >sm (COMMA token)* %hdr_value;
hdr_replyto   = "Reply-To"i >sm %hdr_name %{htype = HReplyTo} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_require   = "Require"i >sm %hdr_name %{htype = HRequire} HCOLON $(hdr,1)
                token >sm (COMMA token)* %hdr_value;
hdr_retry     = "Retry-After"i >sm %hdr_name %{htype = HRetryAfter} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_server    = "Server"i >sm %hdr_name %{htype = HServer} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_subject   = ("Subject"i | "s"i) >sm %hdr_name %{htype = HSubject} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_supported = ("Supported"i | "k"i) >sm %hdr_name %{htype = HSupported} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_timestamp = "Timestamp"i >sm %hdr_name %{htype = HTimestamp} HCOLON $(hdr,1)
                (digit | "." | LWS)+ >sm %hdr_value;
hdr_unsupported = "Unsupported"i >sm %hdr_name %{htype = HUnsupported} HCOLON $(hdr,1)
                token >sm (COMMA token)* %hdr_value;
hdr_useragent = "User-Agent"i >sm %hdr_name %{htype = HUserAgent} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_warning   = "Warning"i >sm %hdr_name %{htype = HWarning} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;
hdr_www_auth  = "WWW-Authenticate"i >sm %hdr_name %{htype = HWWWAuthenticate} HCOLON $(hdr,1)
                generic_value >sm %hdr_value;

headers = (
  hdr_generic
  | hdr_accept
  | hdr_acc_enc
  | hdr_acc_lang
  | hdr_alertinf
  | hdr_allow
  | hdr_auth_info
  | hdr_authz
  | hdr_callid
  | hdr_callinfo
  | hdr_clen
  | hdr_contact
  | hdr_cnt_disp
  | hdr_cnt_enc
  | hdr_cnt_lang
  | hdr_cseq
  | hdr_ctyp
  | hdr_date
  | hdr_errinfo
  | hdr_expires
  | hdr_inreply
  | hdr_from
  | hdr_maxfwd
  | hdr_mimever
  | hdr_minexpr
  | hdr_org
  | hdr_priority
  | hdr_pxy_authn
  | hdr_pxy_authz
  | hdr_pxy_req
  | hdr_rroute
  | hdr_replyto
  | hdr_require
  | hdr_retry
  | hdr_route
  | hdr_server
  | hdr_subject
  | hdr_supported
  | hdr_timestamp
  | hdr_unsupported
  | hdr_useragent
  | hdr_to
  | hdr_via
  | hdr_warning
  | hdr_www_auth
) CRLF;
}%%
