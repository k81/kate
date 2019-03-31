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

type dynField struct {
	Name  string
	Value *reflect.Value
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

	if ind.Kind() != reflect.Struct {
		return json.Unmarshal(data, ptr)
	}

	typ := ind.Type()
	dynFields := []*dynField{}
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
			dynFields = append(dynFields, &dynField{
				Name:  sf.Name,
				Value: &field,
			})
		}
	}

	if err := json.Unmarshal(data, ptr); err != nil {
		return err
	}

	for _, dynField := range dynFields {
		rawMsg := dynField.Value.Interface().(*json.RawMessage)
		dynVal := ptr.(DynamicFielder).NewDynamicField(dynField.Name)
		if dynVal != nil && len(*rawMsg) > 0 {
			if err := ParseJSON([]byte(*rawMsg), dynVal); err != nil {
				return err
			}
			dynField.Value.Set(reflect.ValueOf(dynVal))
		} else {
			dynField.Value.Set(reflect.Zero(dynField.Value.Type())) // for json:",omitempty"
		}
	}

	return nil
}
