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
	assert.Equal(t, 32*time.Second, tm.D)
}
