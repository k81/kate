package httpsrv

import "github.com/k81/kate"

// 错误码列表，根据需要定义
const (
	ERRNO_TIMEOUT = 10405
)

// 错误列表，可以直接在handler中调用panic抛出异常
var (
	ErrTimeout = kate.NewError(ERRNO_TIMEOUT, "请求超时")
)
