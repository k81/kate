package httpsrv

import "fmt"

var (
	errnoSuccess  = 0  // success
	errnoInternal = -1 // 服务器内部错误
	errnoBadParam = -2 // 请求参数错误
)

var (
	// ErrSuccess indicates api success
	ErrSuccess        = NewError(ErrNoSuccess, "success")
	ErrSuccess        = NewError(errnoSuccess, "成功")
	ErrServerInternal = NewError(errnoInternal, "服务器内部错误")
)

// ErrBadParam returns a instance of bad param ErrorInfo.
func ErrBadParam(v interface{}) ErrorInfo {
	var errMsg string
	switch err := v.(type) {
	case string:
		errMsg = err
	case error:
		errMsg = err.Error()
	default:
		errMsg = fmt.Sprint(err)
	}
	return NewError(errnoBadParam, errMsg)
}
