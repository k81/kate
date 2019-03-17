package utils

import (
	"fmt"
	"reflect"
)

// IsZeroValue return true if the value is zero value
func IsZeroValue(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	typ := value.Type()
	kind := typ.Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		return value.Len() == 0
	}

	if !typ.Comparable() {
		panic(fmt.Errorf("type is not comparable: %v", typ))
	}
	return reflect.Zero(typ).Interface() == value.Interface()
}

// IsType check type
func IsType(value reflect.Value, expected reflect.Type) bool {
	if !value.IsValid() {
		return false
	}

	typ := value.Type()
	kind := value.Kind()
	if kind == reflect.Ptr {
		typ = value.Type().Elem()
	}
	return typ == expected
}
