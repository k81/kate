package regexpool

import (
	"regexp"
	"sync"
)

type Pool struct {
	pool    *sync.Pool
	Pattern string
}

func New(pattern string) *Pool {
	_ = regexp.MustCompile(pattern)

	return &Pool{
		Pattern: pattern,
		pool: &sync.Pool{
			New: func() interface{} {
				return regexp.MustCompile(pattern)
			},
		},
	}
}

func (p *Pool) GetMatcher() (m *regexp.Regexp) {
	return p.pool.Get().(*regexp.Regexp)
}

func (p *Pool) PutMatcher(m *regexp.Regexp) {
	p.pool.Put(m)
}
