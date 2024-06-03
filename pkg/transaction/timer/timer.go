package timer

import (
	"time"
)

// Timer contains SIP timers rfc3261
type Timer struct {
	T1 time.Duration
	T2 time.Duration
	T4 time.Duration
	B  time.Duration
	D  time.Duration
	F  time.Duration
	J  time.Duration
}

// New Timer create with default values as in rfc3261#appendix-A
func New() *Timer {
	t := &Timer{
		T1: 500 * time.Millisecond,
		T2: 4 * time.Second,
		T4: 5 * time.Second,
		D:  32 * time.Second,
	}

	t.B = t.T1 * 64
	t.J = t.T1 * 64
	t.F = t.T1 * 64

	return t
}

// FireB blocks until timer B
func (t *Timer) FireB() <-chan struct{} {
	return fire(t.B)
}

// FireD blocks until timer B
func (t *Timer) FireD() <-chan struct{} {
	return fire(t.D)
}

// FireJ blocks until timer J
func (t *Timer) FireJ() <-chan struct{} {
	return fire(t.J)
}

// FireF blocks until timer F
func (t *Timer) FireF() <-chan struct{} {
	return fire(t.F)
}

// FireK blocks until timer K
func (t *Timer) FireK() <-chan struct{} {
	return fire(t.T4)
}

// TickerA returns the timer A value
// with T = T1 and for each next call T = 2*T
func (t *Timer) TickerA() func() time.Duration {
	dur := t.T1

	return func() time.Duration {
		t := dur
		dur *= 2
		return t
	}
}

// TickerE return function that when called returns
// next duration as described in rfc3261#17.1.2.2
func (t *Timer) TickerE() func(bool) time.Duration {
	t1, t2 := t.T1, t.T2
	return func(isProcessing bool) time.Duration {
		if isProcessing {
			return t2
		}
		t := min(t1, t2)
		if t1 < t2 {
			t1 *= 2
		}
		return t
	}
}

func fire(dur time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		<-time.After(dur)
	}()
	return ch
}
