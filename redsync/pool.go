package redsync

import (
	"github.com/garyburd/redigo/redis"
	"github.com/k81/kate/redismgr"
)

// Pool maintains a pool of Redis connections.
type Pool interface {
	GetConn() redis.Conn
}

type defaultPool struct{}

func (p *defaultPool) GetConn() redis.Conn {
	return redismgr.GetConn()
}
