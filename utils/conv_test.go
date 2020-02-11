package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type A struct {
	Name     string `json:"name"`
	AgeValue int    `json:"value"`
}

func BenchmarkFillStruct(b *testing.B) {
	a := &A{}
	m := map[string]interface{}{
		"name":     "zhangsan",
		"AgeValue": "99",
	}

	for i := 0; i < b.N; i++ {
		FillStruct(a, m)
	}
}

func TestBytes2Str(t *testing.T) {
	b := []byte("abc123")
	s := Bytes2Str(b)
	require.Equal(t, "abc123", s)
}

func TestStr2Bytes(t *testing.T) {
	s := "abc123"
	b := Str2Bytes(s)
	require.Equal(t, []byte("abc123"), b)
}
