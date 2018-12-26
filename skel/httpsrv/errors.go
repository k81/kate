package httpsrv

import "github.com/k81/kate"

// 错误码列表，根据需要定义
const (
	ErrnoTimeout = 10405
)

// 错误列表
var (
	ErrTimeout = kate.NewError(ErrnoTimeout, "请求超时")
)
