package logger

import (
	"log"
	"sync/atomic"
)

var enabled atomic.Bool

func Enable(val bool) {
	enabled.Store(val)

	if val {
		log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	}
}

func Log(format string, args ...any) {
	if enabled.Load() {
		log.Printf(format, args...)
	}
}

func Wrn(format string, args ...any) {
	if enabled.Load() {
		log.Printf("WRN "+format, args...)
	}
}

func Err(format string, args ...any) {
	if enabled.Load() {
		log.Printf("ERR "+format, args...)
	}
}
