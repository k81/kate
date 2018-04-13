package utils

import (
	"encoding/json"
)

func ToJSON(v interface{}) string {
	var (
		data []byte
		err  error
	)
	if data, err = json.Marshal(v); err != nil {
		return "encoding failure"
	}
	return string(data)
}
