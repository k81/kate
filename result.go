package kate

type Result struct {
	Status int         `json:"errno"`
	Msg    string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}
