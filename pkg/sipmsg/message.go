package sipmsg

type HType uint8

const (
	HGeneric HType = iota
)

type Message[T Request | Response] struct {
	FirstLine T
	Headers   Headers
}

type Request struct {
	Method string
	RURI   string
}

type Response struct {
	Code   string
	Reason string
}

type Headers struct {
	list []any
}

type HeaderGeneric struct {
	Type HType
	Name string
}
