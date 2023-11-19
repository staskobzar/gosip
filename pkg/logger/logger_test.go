package logger

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisabled(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	Enable(false)

	Log("debug test")
	Err("error test")

	assert.Equal(t, 0, buf.Len())
}

func TestOutput(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	Enable(true)
	Log("debug test %d", 123)
	assert.Regexp(t, "^[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6} debug test 123\n$", buf.String())

	buf.Reset()
	Wrn("warn test %v", true)
	assert.Regexp(t, "^[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6} WRN warn test true\n$", buf.String())

	buf.Reset()
	Err("error test %s", "foo")
	assert.Regexp(t, "^[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6} ERR error test foo\n$", buf.String())
}

func BenchmarkLog(b *testing.B) {
	log.SetOutput(io.Discard)
	Enable(true)

	for i := 0; i < b.N; i++ {
		Log("benchmark log")
		Err("benchmark error %d", 0)
		Wrn("benchmark warn")
	}
}
