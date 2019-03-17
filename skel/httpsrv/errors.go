package httpsrv

var (
	// ErrNoSuccess the error number for success
	ErrNoSuccess = 0 // success
)

var (
	// ErrSuccess indicates api success
	ErrSuccess = NewError(ErrNoSuccess, "success")
)
