package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tm := New()
	assert.Equal(t, 500*time.Millisecond, tm.T1)
	assert.Equal(t, 4*time.Second, tm.T2)
	assert.Equal(t, 5*time.Second, tm.T4)
	assert.Equal(t, 32*time.Second, tm.B)
	assert.Equal(t, 32*time.Second, tm.D)
	assert.Equal(t, 32*time.Second, tm.J)
	assert.Equal(t, 32*time.Second, tm.F)
}

func TestTimers(t *testing.T) {
	tm := New()
	fire := func(ch <-chan struct{}) bool {
		select {
		case <-ch:
			return false
		default:
			return true
		}
	}

	t.Run("timer B", func(t *testing.T) {
		tm.B = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireB()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireB()
		assert.False(t, ok)
	})

	t.Run("timer D", func(t *testing.T) {
		tm.D = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireD()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireD()
		assert.False(t, ok)
	})

	t.Run("timer H", func(t *testing.T) {
		tm.H = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireH()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireH()
		assert.False(t, ok)
	})

	t.Run("timer J", func(t *testing.T) {
		tm.J = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireJ()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireJ()
		assert.False(t, ok)
	})

	t.Run("timer I", func(t *testing.T) {
		tm.T4 = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireI()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireI()
		assert.False(t, ok)
	})

	t.Run("timer F", func(t *testing.T) {
		tm.F = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireF()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireF()
		assert.False(t, ok)
	})

	t.Run("timer K", func(t *testing.T) {
		tm.T4 = 1 * time.Millisecond
		assert.Eventually(t, func() bool { return fire(tm.FireK()) }, 2*time.Millisecond, 100*time.Microsecond)
		_, ok := <-tm.FireK()
		assert.False(t, ok)
	})
}

func TestTickerE(t *testing.T) {
	tests := map[string]struct {
		isProgress []bool
		t1, t2     time.Duration
		want       []time.Duration
	}{
		`for state trying only`: {
			isProgress: []bool{false, false, false, false, false, false},
			t1:         500 * time.Millisecond,
			t2:         4 * time.Second,
			want: []time.Duration{
				500 * time.Millisecond,
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
		`with short step`: {
			isProgress: []bool{false, false, false, false},
			t1:         1 * time.Second,
			t2:         4 * time.Second,
			want: []time.Duration{
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
		`next state is proceed`: {
			isProgress: []bool{false, true, true, true},
			t1:         500 * time.Millisecond,
			t2:         4 * time.Second,
			want: []time.Duration{
				500 * time.Millisecond,
				4 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tm := New()
			tm.T1 = tc.t1
			tm.T2 = tc.t2
			tick := tm.TickerE()
			for i := 0; i < len(tc.want); i++ {
				assert.Equal(t, tc.want[i], tick(tc.isProgress[i]))
			}
		})
	}
}

func TestTickerG(t *testing.T) {
	tests := []struct {
		t1, t2 time.Duration
		want   []time.Duration
	}{
		{
			500 * time.Millisecond, 4 * time.Second,
			[]time.Duration{
				500 * time.Millisecond,
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
		{
			1 * time.Second, 4 * time.Second,
			[]time.Duration{
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				4 * time.Second,
			},
		},
	}

	for _, tc := range tests {
		tm := New()
		tm.T1 = tc.t1
		tm.T2 = tc.t2
		tick := tm.TickerG()
		for i := 0; i < len(tc.want); i++ {
			assert.Equal(t, tc.want[i], tick())
		}
	}
}

func TestTickerA(t *testing.T) {
	tm := New()
	tick := tm.TickerA()

	assert.Equal(t, 500*time.Millisecond, tm.T1)
	assert.Equal(t, 500*time.Millisecond, tick())
	assert.Equal(t, 1*time.Second, tick())
	assert.Equal(t, 2*time.Second, tick())
	assert.Equal(t, 4*time.Second, tick())
	assert.Equal(t, 8*time.Second, tick())
	assert.Equal(t, 16*time.Second, tick())
	assert.Equal(t, 32*time.Second, tick())
}
