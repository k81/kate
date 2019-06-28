package utils

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Abs return the absolute value
func Abs(v int64) int64 {
	if v >= 0 {
		return v
	}
	return -v
}

// Min return the minimum value
func Min(vals ...int) int {
	min := vals[0]
	for _, val := range vals[1:] {
		if val < min {
			min = val
		}
	}
	return min
}

// Max return the maximum value
func Max(vals ...int) int {
	max := vals[0]
	for _, val := range vals[1:] {
		if val > max {
			max = val
		}
	}
	return max
}

// Split split the string.
// if empty return an empty slice instead of nil
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

// TrimUntil trim the prefix until find the stop character
func TrimUntil(s string, stop string) string {
	p := strings.SplitN(s, stop, 2)
	n := len(p)
	return p[n-1]
}

// RepeatWithSep repeat the `s` `count` times, and separated by `sep`
func RepeatWithSep(s string, sep string, count int) string {
	if count < 0 {
		panic("negative RepeatWithSep count")
	} else if count > 0 && len(s)*count/count != len(s) {
		panic("RepeatWithSep count causes overflow")
	}

	b := make([]byte, len(s)*count+len(sep)*(count-1))
	bp := copy(b, s)
	bp += copy(b[bp:], sep)
	for bp < len(b) {
		copy(b[bp:], b[:bp])
		bp *= 2
	}
	return string(b)
}

// JoinSlice join the slice to string with separator `sep`
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
