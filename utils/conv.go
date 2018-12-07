package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

// GetBool convert interface to bool.
func GetBool(v interface{}) bool {
	// nolint:errcheck
	b, _ := strconv.ParseBool(GetString(v))
	return b
}

// GetString convert interface to string.
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

// GetInt convert interface to int.
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
			// nolint:errcheck
			value, _ := strconv.Atoi(d)
			return value
		}
	}
	return 0
}

// GetInt8 convert interface to int8.
func GetInt8(v interface{}) int8 {
	// nolint:errcheck
	s, _ := strconv.ParseInt(GetString(v), 10, 8)
	return int8(s)
}

// GetInt16 convert interface to int16.
func GetInt16(v interface{}) int16 {
	// nolint:errcheck
	s, _ := strconv.ParseInt(GetString(v), 10, 16)
	return int16(s)
}

// GetInt32 convert interface to int32.
func GetInt32(v interface{}) int32 {
	// nolint:errcheck
	s, _ := strconv.ParseInt(GetString(v), 10, 32)
	return int32(s)
}

// GetInt64 convert interface to int64.
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
			// nolint:errcheck
			value, _ := strconv.ParseInt(d, 10, 64)
			return value
		}
	}
	return 0
}

// GetUint convert interface to uint.
func GetUint(v interface{}) uint {
	// nolint:errcheck
	s, _ := strconv.ParseUint(GetString(v), 10, 64)
	return uint(s)
}

// GetUint8 convert interface to uint8.
func GetUint8(v interface{}) uint8 {
	// nolint:errcheck
	s, _ := strconv.ParseUint(GetString(v), 10, 8)
	return uint8(s)
}

// GetUint16 convert interface to uint16.
func GetUint16(v interface{}) uint16 {
	// nolint:errcheck
	s, _ := strconv.ParseUint(GetString(v), 10, 16)
	return uint16(s)
}

// GetUint32 convert interface to uint32.
func GetUint32(v interface{}) uint32 {
	// nolint:errcheck
	s, _ := strconv.ParseUint(GetString(v), 10, 32)
	return uint32(s)
}

// GetUint64 convert interface to uint64.
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
			// nolint:errcheck
			value, _ := strconv.ParseUint(d, 10, 64)
			return value
		}
	}
	return 0
}

// GetFloat32 convert interface to float32.
func GetFloat32(v interface{}) float32 {
	// nolint:errcheck
	f, _ := strconv.ParseFloat(GetString(v), 32)
	return float32(f)
}

// GetFloat64 convert interface to float64.
func GetFloat64(v interface{}) float64 {
	// nolint:errcheck
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
