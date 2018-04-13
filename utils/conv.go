package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func GetBool(v interface{}) bool {
	b, _ := strconv.ParseBool(GetString(v))
	return b
}

// convert interface to string.
func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprintf("%v", result)
		}
	}
	return ""
}

// convert interface to int.
func GetInt(v interface{}) int {
	switch result := v.(type) {
	case int:
		return result
	case int32:
		return int(result)
	case int64:
		return int(result)
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.Atoi(d)
			return value
		}
	}
	return 0
}

func GetInt8(v interface{}) int8 {
	s, _ := strconv.ParseInt(GetString(v), 10, 8)
	return int8(s)
}

func GetInt16(v interface{}) int16 {
	s, _ := strconv.ParseInt(GetString(v), 10, 16)
	return int16(s)
}

func GetInt32(v interface{}) int32 {
	s, _ := strconv.ParseInt(GetString(v), 10, 32)
	return int32(s)
}

// convert interface to int64.
func GetInt64(v interface{}) int64 {
	switch result := v.(type) {
	case int:
		return int64(result)
	case int32:
		return int64(result)
	case int64:
		return result
	default:

		if d := GetString(v); d != "" {
			value, _ := strconv.ParseInt(d, 10, 64)
			return value
		}
	}
	return 0
}

func GetUint(v interface{}) uint {
	s, _ := strconv.ParseUint(GetString(v), 10, 64)
	return uint(s)
}

func GetUint8(v interface{}) uint8 {
	s, _ := strconv.ParseUint(GetString(v), 10, 8)
	return uint8(s)
}

func GetUint16(v interface{}) uint16 {
	s, _ := strconv.ParseUint(GetString(v), 10, 16)
	return uint16(s)
}

func GetUint32(v interface{}) uint32 {
	s, _ := strconv.ParseUint(GetString(v), 10, 32)
	return uint32(s)
}

// convert interface to uint64.
func GetUint64(v interface{}) uint64 {
	switch result := v.(type) {
	case int:
		return uint64(result)
	case int32:
		return uint64(result)
	case int64:
		return uint64(result)
	case uint64:
		return result
	default:

		if d := GetString(v); d != "" {
			value, _ := strconv.ParseUint(d, 10, 64)
			return value
		}
	}
	return 0
}

func GetFloat32(v interface{}) float32 {
	f, _ := strconv.ParseFloat(GetString(v), 32)
	return float32(f)
}

func GetFloat64(v interface{}) float64 {
	f, _ := strconv.ParseFloat(GetString(v), 64)
	return f
}

func StringJoin(params ...interface{}) string {
	var buffer bytes.Buffer

	for _, para := range params {
		buffer.WriteString(GetString(para))
	}

	return buffer.String()
}

func GetIntSlices(v interface{}) []int {

	switch result := v.(type) {

	case []int:
		return []int(result)
	default:
		return nil
	}
}

func GetInt64Slices(v interface{}) []int64 {

	switch result := v.(type) {

	case []int64:
		return []int64(result)
	default:
		return nil
	}
}

func GetUint64Slices(v interface{}) []uint64 {

	switch result := v.(type) {

	case []uint64:
		return []uint64(result)
	default:
		return nil
	}
}

// convert interface to byte slice.
func GetByteArray(v interface{}) []byte {
	switch result := v.(type) {
	case []byte:
		return result
	case string:
		return []byte(result)
	default:
		return nil
	}
}

func StringsToInterfaces(keys []string) []interface{} {
	result := make([]interface{}, len(keys))
	for i, k := range keys {
		result[i] = k
	}
	return result
}

func GetByKind(kind reflect.Kind, v interface{}) (result interface{}) {
	switch kind {
	case reflect.Bool:
		result = GetBool(v)
	case reflect.Int:
		result = GetInt(v)
	case reflect.Int8:
		result = GetInt8(v)
	case reflect.Int16:
		result = GetInt16(v)
	case reflect.Int32:
		result = GetInt32(v)
	case reflect.Int64:
		result = GetInt64(v)
	case reflect.Uint:
		result = GetUint(v)
	case reflect.Uint8:
		result = GetUint8(v)
	case reflect.Uint16:
		result = GetUint16(v)
	case reflect.Uint32:
		result = GetUint32(v)
	case reflect.Uint64:
		result = GetUint64(v)
	case reflect.Float32:
		result = GetFloat32(v)
	case reflect.Float64:
		result = GetFloat64(v)
	default:
		result = v
	}
	return
}

func T(name string) func(string) bool {
	return func(field string) bool {
		if field == name {
			return true
		}
		if ToSnake(field) == name {
			return true
		}

		if ToCamel(field) == name {
			return true
		}

		if ToCamelLower(field) == name {
			return true
		}
		return false
	}
}

func FillStruct(s interface{}, data map[string]interface{}) {
	t := reflect.ValueOf(s).Elem()

	for k, v := range data {
		match := T(k)

		field := t.FieldByNameFunc(match)

		if !field.IsValid() {
			panic(fmt.Errorf("no such field %s", k))
		}

		if !field.CanSet() {
			panic(fmt.Errorf("field %s can not be set", k))
		}

		val := GetByKind(field.Kind(), v)
		field.Set(reflect.ValueOf(val))
	}
}
