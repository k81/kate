package utils

import (
	"bytes"
	"encoding/json"
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

// getInnerPtrValue
// val: *struct{}
// ind: struct{}
func getInnerPtrValue(ptr interface{}) (val, ind reflect.Value) {
	val = reflect.ValueOf(ptr)
	ind = reflect.Indirect(val)

	for {
		switch ind.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			if ind.IsNil() {
				ind.Set(reflect.New(ind.Type().Elem()))
			}
			val = ind
			ind = val.Elem()
		default:
			return val, ind
		}
	}
}

// ParseJSON parse json with dynamic field parse support
func ParseJSON(data []byte, ptr interface{}) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		// ignore empty field
		return nil
	}

	val, ind := getInnerPtrValue(ptr)
	ptr = val.Interface()
	_, ok := ptr.(DynamicFielder)
	if !ok {
		return json.Unmarshal(data, ptr)
	}

	dynFieldMap := make(map[string]*json.RawMessage)
	typ := ind.Type()
	for i := 0; i < ind.NumField(); i++ {
		sf := typ.Field(i)
		field := ind.Field(i)

		if !field.CanSet() {
			continue
		}

		dynamic := sf.Tag.Get("dynamic")
		if dynamic == "true" {
			rawMsg := new(json.RawMessage)
			field.Set(reflect.ValueOf(rawMsg))
			dynFieldMap[sf.Name] = rawMsg
		}
	}

	if err := json.Unmarshal(data, ptr); err != nil {
		return err
	}

	for name, rawMsg := range dynFieldMap {
		field := ind.FieldByName(name)
		dynVal := ptr.(DynamicFielder).NewDynamicField(name)
		if dynVal != nil && len(*rawMsg) > 0 {
			if err := ParseJSON([]byte(*rawMsg), dynVal); err != nil {
				return err
			}
			field.Set(reflect.ValueOf(dynVal))
		} else {
			field.Set(reflect.Zero(field.Type())) // for json:",omitempty"
		}
	}

	return nil
}
