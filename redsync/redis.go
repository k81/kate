package redsync

import (
	"github.com/garyburd/redigo/redis"
	"github.com/k81/kate/redismgr"
)

// A Pool maintains a pool of Redis connections.
type Pool interface {
	GetConn() redis.Conn
}

type DefaultPool struct{}

func (p *DefaultPool) GetConn() redis.Conn {
	return redismgr.GetConn()
}
