package utils

import "testing"

func BenchmarkFastUUID(t *testing.B) {
	for i := 0; i < t.N; i++ {
		FastUUID()
	}
}

func BenchmarkFastUUIDStr(t *testing.B) {
	for i := 0; i < t.N; i++ {
		FastUUIDStr()
	}
}
