// Package state provides SIP transactions' states as described in rfc3261
package state

import "sync"

// Type of SIP transaction state
type Type uint8

// SIP transaction states
const (
	Unknown Type = iota
	Calling
	Completed
	Confirmed
	Proceeding
	Terminated
	Trying
)

// State core structure
type State struct {
	mu sync.RWMutex
	v  Type
}

// New state create
func New() *State {
	return &State{}
}

// Set SIP transaction state
func (st *State) Set(v Type) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.v = v
}

// IsCalling returns true if state is Calling
func (st *State) IsCalling() bool { return st.is(Calling) }

// IsCompleted returns true if state is Completed
func (st *State) IsCompleted() bool { return st.is(Completed) }

// IsConfirmed returns true if state is Confirmed
func (st *State) IsConfirmed() bool { return st.is(Confirmed) }

// IsProceeding returns true if state is Proceeding
func (st *State) IsProceeding() bool { return st.is(Proceeding) }

// IsTerminated returns true if state is Terminated
func (st *State) IsTerminated() bool { return st.is(Terminated) }

// IsTrying returns true if state is Trying
func (st *State) IsTrying() bool { return st.is(Trying) }

// String is stringify interface for SIP state
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
