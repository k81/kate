package httpsrv

// Result define the handle result for http request
type Result struct {
	ErrNO  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}
