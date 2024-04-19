package timer

import "time"

type Timer struct {
	T1 time.Duration
	T2 time.Duration
	T4 time.Duration
	D  time.Duration
	J  time.Duration
}

func New() *Timer {
	t := &Timer{
		T1: 500 * time.Millisecond,
		T2: 4 * time.Second,
		T4: 5 * time.Second,
		D:  32 * time.Second,
	}

	t.J = t.T1 * 64

	return t
}

func (t *Timer) FireJ() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		<-time.After(t.J)
	}()
	return ch
}
