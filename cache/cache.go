package cache

import (
	"context"

	"github.com/k81/kate/log"
	"github.com/k81/kate/redismgr"

	"github.com/garyburd/redigo/redis"
)

func Load(ctx context.Context, key string) (value string, err error) {
	value, err = redis.String(redismgr.TryRedisCmd(ctx, nil, "GET", key))
	switch {
	case err == redis.ErrNil:
		err = nil
	case err != nil:
		log.Error(ctx, "load cache", "key", key, "error", err)
	}
	return
}

func Store(ctx context.Context, key string, value string) (err error) {
	if _, err = redismgr.TryRedisCmd(ctx, nil, "SET", key, value); err != nil {
		log.Error(ctx, "store cache", "key", key, "value", value, "error", err)
	}
	return
}

func StoreWithTimeout(ctx context.Context, key string, value string, timeout int64) (err error) {
	if _, err = redismgr.TryRedisCmd(ctx, nil, "SET", key, value, "EX", timeout); err != nil {
		log.Error(ctx, "store cache", "key", key, "value", value, "timeout", timeout, "error", err)
	}
	return
}

func Delete(ctx context.Context, key string) (err error) {
	if _, err = redismgr.TryRedisCmd(ctx, nil, "DEL", key); err != nil {
		log.Error(ctx, "delete cache", "key", key, "error", err)
	}
	return
}
