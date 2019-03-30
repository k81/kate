package utils

import (
	"reflect"
)

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

// IsEmptyValue return true if the value is zero value
func IsEmptyValue(v reflect.Value) bool {
	for {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			if v.IsNil() {
				return true
			}
			v = v.Elem()
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				if !IsEmptyValue(v.Index(i)) {
					return false
				}
			}
			return true
		case reflect.Map:
			iter := v.MapRange()
			for iter.Next() {
				if !IsEmptyValue(iter.Value()) {
					return false
				}
			}
			return true
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				if !IsEmptyValue(v.Field(i)) {
					return false
				}
			}
			return true
		default:
			return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
		}
	}
}
