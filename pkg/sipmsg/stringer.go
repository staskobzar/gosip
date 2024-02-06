package sipmsg

import "bytes"

// Stringer provide helper structure to build SIP Message
// or its headers as a string
type Stringer struct {
	buf *bytes.Buffer
}

// NewStringer creates new stringer of the given length
func NewStringer(capacity int) *Stringer {
	return &Stringer{
		buf: bytes.NewBuffer(make([]byte, 0, capacity)),
	}
}

// Print prints into stringer buffer strings
func (str *Stringer) Print(val ...string) {
	for _, v := range val {
		str.buf.WriteString(v)
	}
}

// String returns stringer buffer as string
func (str *Stringer) String() string {
	return str.buf.String()
}
