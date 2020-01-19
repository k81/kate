package sampler_test

import (
	"testing"
	"time"

	"github.com/k81/kate/debug/sampler"
)

func TestSamplerCheck(t *testing.T) {
	s := sampler.New(time.Second, 0, 10)
	count := 0
	start := time.Now()
	stop := start.Add(5 * time.Second)
	for {
		now := time.Now()
		if now.After(stop) {
			break
		}

		if s.Check(time.Now()) {
			count++
		}

		time.Sleep(time.Millisecond)
	}

	t.Logf("count=%d", count)
}
