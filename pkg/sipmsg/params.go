package sipmsg

import (
	"bytes"
	"strings"
)

// Params for URIs and SIP headers
type Params string

// Add appends parameter as key=val. If val
// is empty then adds only key
func (p Params) Add(key, val string) Params {
	if len(key) == 0 {
		return p
	}

	buf := bytes.NewBufferString(string(p))
	p.bufAdd(buf, key, val)

	return Params(buf.String())
}

// Set sets the parameter entries associated with key to
// the single element value. It replaces any existing
// values associated with key.
func (p Params) Set(key, val string) Params {
	if len(p) == 0 || len(key) == 0 {
		return p
	}

	buf := p.makeBuf(key, val)

	for _, pt := range p.split() {
		prm := string(p[pt[0]:pt[1]])
		if prm == key && len(val) > 0 {
			p.bufAdd(buf, key, val)
			continue
		}

		if !strings.HasPrefix(prm, key+"=") {
			p.bufAdd(buf, prm, "")
			continue
		}

		p.bufAdd(buf, key, val)
	}
	return Params(buf.String())
}

// Del deletes parameters by key
func (p Params) Del(key string) Params {
	if len(p) == 0 || len(key) == 0 {
		return p
	}
	buf := p.makeBuf("", "")

	for _, pt := range p.split() {
		prm := string(p[pt[0]:pt[1]])
		if prm == key || strings.HasPrefix(prm, key+"=") {
			continue
		}
		p.bufAdd(buf, prm, "")
	}

	return Params(buf.String())
}

// Get gets the first value associated with the given key
// TODO: case insensitive
func (p Params) Get(key string) (string, bool) {
	if len(p) == 0 || len(key) == 0 {
		return "", false
	}
	for _, pt := range p.split() {
		prm := string(p[pt[0]:pt[1]])
		if prm == key {
			return "", true
		}

		if strings.HasPrefix(prm, key+"=") {
			return prm[len(key)+1:], true
		}
	}
	return "", false
}

// Len returns lenght of params in bytes
func (p Params) Len() int { return len(p) }

// String representation of the parameters
func (p Params) String() string {
	if p.Len() == 0 {
		return ""
	}
	return ";" + p.str()
}

func (p Params) str() string { return string(p) }

// setup prepares parameters string
// it will trim any space or semicolon on the left
// it is a bit faster then built-in go func
// strings.TrimLeft(input, " ;")
func (p Params) setup() Params {
	if len(p) == 0 {
		return p
	}
	idx := 0
	for i, c := range p {
		if c == ' ' || c == ';' {
			idx = i + 1
			continue
		}
		break
	}
	return p[idx:]
}

// TODO: handle spaces in the beginning
// ; branch=z9hG4bKna; maddr=10.0.0.1;received=10.0.0.100 ;ttl=120
func (p Params) split() [][2]int {
	pt := make([][2]int, 0, 4)
	start := 0
	for i, s := range p {
		if s == ';' {
			pt = append(pt, [2]int{start, i})
			start = i + 1
		}
	}
	pt = append(pt, [2]int{start, len(p)})
	return pt
}

func (p Params) makeBuf(k, v string) *bytes.Buffer {
	buf := make([]byte, 0, len(p)+len(k)+len(v)+2)
	return bytes.NewBuffer(buf)
}

func (p Params) bufAdd(buf *bytes.Buffer, k, v string) {
	if buf.Len() > 0 {
		buf.WriteByte(';')
	}
	buf.WriteString(k)
	if len(v) > 0 {
		buf.WriteByte('=')
		buf.WriteString(v)
	}
}
