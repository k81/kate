package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// BindSliceSep the separator for parsing slice field
var BindSliceSep = ","

// BindUnmarshaler the bind unmarshal interface
type BindUnmarshaler interface {
	UnmarshalBind(value string) error
}

// FillStruct set the field value of s according data kv map.
func FillStruct(s interface{}, data map[string]interface{}) {
	err := Bind(s, "", data)
	if err != nil {
		panic(err)
	}
}

// Bind bind values to struct ptr
func Bind(ptr interface{}, tag string, input map[string]interface{}) error {
	val := reflect.ValueOf(ptr)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	fullName := typ.PkgPath() + "." + typ.Name()

	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("bind: cannot use non-ptr struct `%s`", fullName))
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("bind: only allow ptr of struct"))
	}

	for i := 0; i < ind.NumField(); i++ {
		structField := ind.Type().Field(i)
		field := ind.Field(i)

		if !field.CanSet() {
			continue
		}

		name := ""
		if tag != "" {
			name = structField.Tag.Get(tag)
			if name == "" {
				continue
			}
		} else {
			name = structField.Name
		}

		value, ok := input[name]
		if !ok {
			continue
		}

		if err := bindValue(field, value); err != nil {
			return err
		}
	}
	return nil
}

func bindSlice(field reflect.Value, value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		field.Set(reflect.ValueOf(value))
		return nil
	}

	vals := strings.Split(strValue, BindSliceSep)
	if len(vals) == 0 {
		return nil
	}

	ind := reflect.Indirect(field)
	typ := ind.Type().Elem()
	isPtr := typ.Kind() == reflect.Ptr

	if isPtr {
		typ = typ.Elem()
	}

	slice := reflect.New(ind.Type()).Elem()
	for _, val := range vals {
		elem := reflect.New(typ)
		elemInd := reflect.Indirect(elem)

		if err := bindValue(elemInd, val); err != nil {
			return err
		}

		if isPtr {
			slice = reflect.Append(slice, elemInd.Addr())
		} else {
			slice = reflect.Append(slice, elemInd)
		}
	}

	ind.Set(slice)

	return nil
}

// nolint:gocyclo
func bindValue(field reflect.Value, value interface{}) error {
	ok, err := unmarshalBind(field, value)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(GetBool(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(GetInt64(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(GetUint64(value))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(GetFloat64(value))
	case reflect.String:
		field.SetString(GetString(value))
	case reflect.Slice:
		if err := bindSlice(field, value); err != nil {
			return err
		}
	default:
		field.Set(reflect.ValueOf(value))
	}
	return nil
}

func unmarshalBind(field reflect.Value, value interface{}) (ok bool, err error) {
	strValue, ok := value.(string)
	if !ok {
		return false, nil
	}

	ptr := reflect.New(field.Type())
	if !ptr.CanInterface() {
		return false, nil
	}

	unmarshaler, ok := ptr.Interface().(BindUnmarshaler)
	if !ok {
		return false, nil
	}

	if err = unmarshaler.UnmarshalBind(strValue); err != nil {
		return false, err
	}

	field.Set(reflect.Indirect(ptr))
	return true, nil
}
