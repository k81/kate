package kate

type ErrorInfo interface {
	error
	Code() int
}

type errSimple struct {
	ErrCode    int
	ErrMessage string
}

func NewError(code int, message string) ErrorInfo {
	return &errSimple{code, message}
}

func (e *errSimple) Code() int {
	return e.ErrCode
}

func (e *errSimple) Error() string {
	return e.ErrMessage
}

var ErrSuccess = NewError(10000, "success")

type ErrorInfoWithData interface {
	error
	Code() int
	Data() interface{}
}

type errWithData struct {
	*errSimple
	ErrData interface{}
}

func (e *errWithData) Data() interface{} {
	return e.ErrData
}

func NewErrorWithData(code int, message string, data interface{}) ErrorInfoWithData {
	return &errWithData{
		errSimple: &errSimple{code, message},
		ErrData:   data,
	}
}
