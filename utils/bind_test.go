package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExampleVal struct {
	Bool        bool     `rest:"val_bool"`
	Int         int      `rest:"val_int"`
	Int8        int8     `rest:"val_int"`
	Int16       int16    `rest:"val_int"`
	Int32       int32    `rest:"val_int"`
	Int64       int64    `rest:"val_int"`
	Uint        uint     `query:"val_uint"`
	Uint8       uint8    `query:"val_uint"`
	Uint16      uint16   `query:"val_uint"`
	Uint32      uint32   `query:"val_uint"`
	Uint64      uint64   `query:"val_uint"`
	Uint64Slice []uint64 `query:"val_uints"`
	Float32     float32  `query:"val_float"`
	Float64     float64  `query:"val_float"`
	String      string   `query:"val_string"`
}

func TestBind(t *testing.T) {
	v := &ExampleVal{}

	data := map[string]interface{}{
		"val_bool":   "true",
		"val_int":    "10",
		"val_uint":   "100",
		"val_float":  "100.1",
		"val_string": "zhangsan",
		"val_uints":  "1,2,3,4",
	}

	Bind(v, "query", data)
	Bind(v, "rest", data)
	assert.Equal(t, true, v.Bool, "bind bool")
	assert.Equal(t, int(10), v.Int, "bind int")
	assert.Equal(t, int8(10), v.Int8, "bind int8")
	assert.Equal(t, int16(10), v.Int16, "bind int16")
	assert.Equal(t, int32(10), v.Int32, "bind int32")
	assert.Equal(t, int64(10), v.Int64, "bind int64")
	assert.Equal(t, uint(100), v.Uint, "bind uint")
	assert.Equal(t, uint8(100), v.Uint8, "bind uint8")
	assert.Equal(t, uint16(100), v.Uint16, "bind uint16")
	assert.Equal(t, uint32(100), v.Uint32, "bind uint32")
	assert.Equal(t, uint64(100), v.Uint64, "bind uint64")
	assert.Equal(t, float32(100.1), v.Float32, "bind float32")
	assert.Equal(t, float64(100.1), v.Float64, "bind float64")
	assert.Equal(t, "zhangsan", v.String, "bind string")
	assert.Equal(t, []uint64{1, 2, 3, 4}, v.Uint64Slice, "bind uint64 slice")
}
