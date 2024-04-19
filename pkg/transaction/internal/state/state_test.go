package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New()
	assert.Equal(t, Unknown, s.v)
}

func TestString(t *testing.T) {
	tests := []struct {
		input Type
		want  string
	}{
		{Unknown, "Unknown"},
		{Calling, "Calling"},
		{Completed, "Completed"},
		{Confirmed, "Confirmed"},
		{Proceeding, "Proceeding"},
		{Terminated, "Terminated"},
		{Trying, "Trying"},
	}

	st := New()
	for _, tc := range tests {
		st.Set(tc.input)
		assert.Equal(t, tc.want, st.String())
	}
}
