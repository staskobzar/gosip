package state

import "sync"

type Type uint8

const (
	Unknown Type = iota
	Calling
	Completed
	Confirmed
	Proceeding
	Terminated
	Trying
)

type State struct {
	mu sync.RWMutex
	v  Type
}

func New() *State {
	return &State{}
}

func (st *State) Set(v Type) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.v = v
}

func (st *State) IsCalling() bool    { return st.is(Calling) }
func (st *State) IsCompleted() bool  { return st.is(Completed) }
func (st *State) IsConfirmed() bool  { return st.is(Confirmed) }
func (st *State) IsProceeding() bool { return st.is(Proceeding) }
func (st *State) IsTrying() bool     { return st.is(Trying) }
func (st *State) IsTerminated() bool { return st.is(Terminated) }

func (st *State) String() string {
	st.mu.RLock()
	defer st.mu.RUnlock()
	switch st.v {
	case Calling:
		return "Calling"
	case Completed:
		return "Completed"
	case Confirmed:
		return "Confirmed"
	case Proceeding:
		return "Proceeding"
	case Terminated:
		return "Terminated"
	case Trying:
		return "Trying"
	default:
		return "Unknown"
	}
}

func (st *State) is(t Type) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.v == t
}
