package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type exampleVal struct {
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
	IntPtr      *int     `query:"val_int"`
}

func TestBind(t *testing.T) {
	v := &exampleVal{}

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
	require.Equal(t, true, v.Bool, "bind bool")
	require.Equal(t, int(10), v.Int, "bind int")
	require.Equal(t, int8(10), v.Int8, "bind int8")
	require.Equal(t, int16(10), v.Int16, "bind int16")
	require.Equal(t, int32(10), v.Int32, "bind int32")
	require.Equal(t, int64(10), v.Int64, "bind int64")
	require.Equal(t, uint(100), v.Uint, "bind uint")
	require.Equal(t, uint8(100), v.Uint8, "bind uint8")
	require.Equal(t, uint16(100), v.Uint16, "bind uint16")
	require.Equal(t, uint32(100), v.Uint32, "bind uint32")
	require.Equal(t, uint64(100), v.Uint64, "bind uint64")
	require.Equal(t, float32(100.1), v.Float32, "bind float32")
	require.Equal(t, float64(100.1), v.Float64, "bind float64")
	require.Equal(t, "zhangsan", v.String, "bind string")
	require.Equal(t, []uint64{1, 2, 3, 4}, v.Uint64Slice, "bind uint64 slice")
	require.NotNil(t, v.IntPtr)
	require.Equal(t, int(10), *v.IntPtr, "bind int ptr")
}

type taggedStruct struct {
	ID    int
	Name  string `field:"create"`
	Value string `field:"create,update"`
}

func TestFillStructByTag(t *testing.T) {
	createVal := &taggedStruct{}
	createParams := map[string]interface{}{
		"ID":    100,
		"Name":  "zhangsan",
		"Value": "hello",
	}
	filledByCreate, err := FillStructByTag(createVal, "create", createParams)
	require.NoError(t, err, "fill struct by tag create")
	require.Equal(t, 0, createVal.ID)
	require.Equal(t, "zhangsan", createVal.Name)
	require.Equal(t, "hello", createVal.Value)
	require.Equal(t, []string{"Name", "Value"}, filledByCreate)

	updateVal := &taggedStruct{}
	updateParams := map[string]interface{}{
		"Name":  "lisi",
		"Value": "world",
	}
	filledByUpdate, err := FillStructByTag(updateVal, "update", updateParams)
	require.NoError(t, err, "fill struct by tag update")
	require.Equal(t, "", updateVal.Name)
	require.Equal(t, "world", updateVal.Value)
	require.Equal(t, []string{"Value"}, filledByUpdate)
}
