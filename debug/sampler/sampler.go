package sampler

import (
	"time"
)

type Sampler struct {
	counter           counter
	tick              time.Duration
	first, thereafter uint64
}

// 每tick间隔内:
// 1. 对前first个请求返回true
// 2. 之后开始采样，采样率为1/thereafter
func New(tick time.Duration, first, thereafter int) *Sampler {
	return &Sampler{
		tick:       tick,
		first:      uint64(first),
		thereafter: uint64(thereafter),
	}
}

func (s *Sampler) Check(now time.Time) bool {
	n := s.counter.IncCheckReset(now, s.tick)
	if n > s.first && (n-s.first)%s.thereafter != 0 {
		return false
	}
	return true
}
