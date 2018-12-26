package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
)

// ToJSON return json encoded string of interface.
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

// DynamicFielder is the dynamic fielder interface
// Struct which implement this interface will have dynamic field support
type DynamicFielder interface {
	NewDynamicField(fieldName string) interface{}
}

// ParseJSON parse json with dynamic field parse support
func ParseJSON(r io.Reader, ptr interface{}) error {
	_, ok := ptr.(DynamicFielder)
	if !ok {
		return json.NewDecoder(r).Decode(ptr)
	}

	val := reflect.ValueOf(ptr)
	ind := reflect.Indirect(val)
	dynFieldMap := make(map[string]*json.RawMessage)

	for i := 0; i < ind.NumField(); i++ {
		sf := ind.Type().Field(i)
		field := ind.Field(i)

		dynamic := sf.Tag.Get("dynamic")
		if dynamic == "true" {
			rawMsg := new(json.RawMessage)
			field.Set(reflect.ValueOf(rawMsg))
			dynFieldMap[sf.Name] = rawMsg
		}
	}

	if err := json.NewDecoder(r).Decode(ptr); err != nil {
		return err
	}

	for name, rawMsg := range dynFieldMap {
		field := ind.FieldByName(name)
		dynVal := ptr.(DynamicFielder).NewDynamicField(name)
		if dynVal != nil {
			if err := ParseJSON(bytes.NewReader([]byte(*rawMsg)), dynVal); err != nil {
				return err
			}
			field.Set(reflect.ValueOf(dynVal))
		}
	}

	return nil
}
