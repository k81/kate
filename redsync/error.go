package redsync

import "errors"

// ErrFailed indicates error happened when acquire the lock
var ErrFailed = errors.New("redsync: failed to acquire lock")
