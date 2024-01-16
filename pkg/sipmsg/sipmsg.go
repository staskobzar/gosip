package sipmsg

import (
	"errors"
	"gosip/pkg/sip"
)

// errors
var (
	ErrURIParse = errors.New("SIP URI parse")
	ErrMsgParse = errors.New("SIP Message parse")
)

func Decode(input []byte) (sip.Message, error) {
	return nil, nil
}
