package utils

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Abs(v int64) int64 {
	if v >= 0 {
		return v
	}
	return -v
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Split(str string, sep string) (r []string) {
	p := strings.Split(str, sep)
	if len(p) == 1 && p[0] == "" {
		r = []string{}
		return
	}

	r = make([]string, 0, len(p))
	for _, v := range p {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			r = append(r, v)
		}
	}
	return
}

func TrimUntil(s string, stop string) string {
	p := strings.SplitN(s, stop, 2)
	n := len(p)
	return p[n-1]
}

func JoinN(c string, sep string, n int) string {
	return strings.Join(strings.Split(strings.Repeat(c, n), ""), sep)
}

func JoinSlice(slice interface{}, sep string) string {
	var (
		v = reflect.ValueOf(slice)
		b bytes.Buffer
	)
	if v.Kind() != reflect.Slice {
		panic(errors.New("not slice"))
	}

	for i := 0; i < v.Len(); i++ {
		b.WriteString(fmt.Sprint(v.Index(i)))
		b.WriteString(sep)
	}
	if b.Len() > 0 {
		b.Truncate(b.Len() - len(sep))
	}

	return b.String()
}
