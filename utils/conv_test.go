package utils

import "testing"

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
