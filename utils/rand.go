package utils

import (
	"math/rand"
	"sync"
	"time"
)

// Rand return an random number between [min, max)
func Rand(min, max int) int {
	return rand.Intn(max-min) + min
}

const (
	AlphaNumber = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Number      = "0123456789"
	Alpha       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randSrcPool *sync.Pool

// nolint:gochecknoinits
func init() {
	randSrcPool = &sync.Pool{
		New: func() interface{} {
			return rand.NewSource(time.Now().UnixNano())
		},
	}
}

// RandString return a random string of length n
func RandString(n int, letterBytes string) string {
	// nolint:errcheck
	src := randSrcPool.Get().(rand.Source)
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	randSrcPool.Put(src)

	return string(b)
}
