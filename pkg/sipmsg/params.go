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

	buf := bytes.NewBufferString(p.str())
	p.bufAdd(buf, key, val)

	return Params(buf.String())
}

// Set sets the parameter entries associated with key to
// the single element value. It replaces any existing
// values associated with key.
func (p Params) Set(key, val string) Params {
	if p.Len() == 0 || len(key) == 0 {
		return p
	}

	buf := p.makeBuf(key, val)

	for _, pt := range p.split() {
		prm := p.sub(pt[0], pt[1])
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
	if p.Len() == 0 || len(key) == 0 {
		return p
	}
	buf := p.makeBuf("", "")

	for _, pt := range p.split() {
		prm := p.sub(pt[0], pt[1])
		if prm == key || strings.HasPrefix(prm, key+"=") {
			continue
		}
		p.bufAdd(buf, prm, "")
	}

	return Params(buf.String())
}

// Get gets the first value associated with the given key
// case-insensitive match
func (p Params) Get(key string) (string, bool) {
	if p.Len() == 0 || len(key) == 0 {
		return "", false
	}
	for _, pt := range p.split() {
		prm := p.sub(pt[0], pt[1])
		if strings.EqualFold(prm, key) {
			return "", true
		}

		if len(prm) < len(key)+1 {
			continue
		}

		if strings.EqualFold(prm[:len(key)+1], key+"=") {
			return prm[len(key)+1:], true
		}
	}
	return "", false
}

// Len returns length of params in bytes
func (p Params) Len() int { return len(p) }

// String representation of the parameters
func (p Params) String() string {
	if p.Len() == 0 {
		return ""
	}
	return ";" + p.str()
}

func (p Params) str() string { return string(p) }

func (p Params) sub(x, y int) string { return string(p[x:y]) }

// setup prepares parameters string
// it will trim any space or semicolon on the left
// it is a bit faster then built-in go func
// strings.TrimLeft(input, " ;")
func (p Params) setup() Params {
	if p.Len() == 0 {
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

func (p Params) split() [][2]int {
	pt := make([][2]int, 0, 4)
	s, e := 0, 0
	for i, c := range p {
		switch c {
		case ';':
			pt = append(pt, [2]int{s, e + 1})
			s = i + 1
		case ' ':
			if s >= e {
				s = i + 1
			}
		default:
			e = i
		}
	}
	pt = append(pt, [2]int{s, e + 1})
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
