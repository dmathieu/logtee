package handler

import (
	"testing"

	"github.com/dmathieu/logtee/logtee"
)

var logLine = logtee.LogLine{Path: "/report"}

func BenchmarkHandler(b *testing.B) {
	h := Handler{}

	for i := 0; i < b.N; i++ {
		_ = h.Handle(logLine)
	}
}
