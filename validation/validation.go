package validation

import (
	"fmt"
	"net/url"
)

const (
	Required = 1
	Optional = 2
)

type ErrMissingField struct {
	Field string
}

func (e *ErrMissingField) Error() string {
	return fmt.Sprintf("missing field `%s`", e.Field)
}

func NewErrMissingField(field string) *ErrMissingField {
	return &ErrMissingField{
		Field: field,
	}
}

func Validate(form url.Values, schema map[string]int) (params map[string]interface{}, err error) {
	params = make(map[string]interface{}, len(form))
	for k := range form {
		//filter out unknown fields
		if _, ok := schema[k]; ok {
			if v := form.Get(k); v != "" {
				params[k] = v
			}
		}
	}

	for k, v := range schema {
		if v == Required {
			if _, ok := params[k]; !ok {
				err = NewErrMissingField(k)
				return
			}
		}
	}
	return
}
