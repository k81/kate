package utils

import "testing"

func TestRandString(t *testing.T) {
	str := RandString(10, Number)
	t.Logf("str=%v", str)
	if len(str) != 10 {
		t.Fatalf("length not 10: len=%d, str=%v", len(str), str)
	}
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(10, Number)
	}
}
