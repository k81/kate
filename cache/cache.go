package cache

import (
	"context"

	"github.com/garyburd/redigo/redis"
	"github.com/k81/kate/redismgr"
	"github.com/k81/log"
)

// Load fetch the cache by key
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

// Store save the value to cache under key
func Store(ctx context.Context, key string, value string) (err error) {
	if _, err = redismgr.TryRedisCmd(ctx, nil, "SET", key, value); err != nil {
		log.Error(ctx, "store cache", "key", key, "value", value, "error", err)
	}
	return
}

// StoreWithTimeout save the value to cache with timeout
func StoreWithTimeout(ctx context.Context, key string, value string, timeout int64) (err error) {
	if _, err = redismgr.TryRedisCmd(ctx, nil, "SET", key, value, "EX", timeout); err != nil {
		log.Error(ctx, "store cache", "key", key, "value", value, "timeout", timeout, "error", err)
	}
	return
}

// Delete remove key from cache
func Delete(ctx context.Context, key string) (err error) {
	if _, err = redismgr.TryRedisCmd(ctx, nil, "DEL", key); err != nil {
		log.Error(ctx, "delete cache", "key", key, "error", err)
	}
	return
}
