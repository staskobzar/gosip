// Package sipmsg provides SIP Message parse and generate
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

// Decode parse received SIP message into Message
func Decode(_ []byte) (sip.Message, error) {
	return nil, nil
}
