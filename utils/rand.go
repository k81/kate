package utils

import (
	"math/rand"
	"sync"
	"time"
)

func Rand(min, max int) int {
	return rand.Intn(max-min) + min
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	randSrcPool *sync.Pool
)

func init() {
	randSrcPool = &sync.Pool{
		New: func() interface{} {
			return rand.NewSource(time.Now().UnixNano())
		},
	}
}

func RandString(n int) string {
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
