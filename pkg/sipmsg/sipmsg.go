// Package sipmsg provides SIP Message parse and generate
package sipmsg

import (
	"errors"
)

// errors
var (
	ErrURIParse = errors.New("SIP URI parse")
	ErrMsgParse = errors.New("SIP Message parse")
)

// Decode parse received SIP message into Message
func Decode(input []byte) (*Message, error) {
	return Parse(string(input))
}
