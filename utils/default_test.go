package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type defaultStruct struct {
	Int         int      `default:"100"`
	String      string   `default:"2018sigh"`
	PtrInt      *int     `default:"100"`
	SliceInt    []int    `default:"1,2,3"`
	SliceString []string `default:"a,b,c"`
	PtrSliceInt *[]int   `default:"1,2,3"`
}

func TestSetDefaults(t *testing.T) {
	v := &defaultStruct{}
	err := SetDefaults(v)
	require.NoError(t, err)
	iVal := 100
	strVal := "2018sigh"
	sliceInt := []int{1, 2, 3}
	sliceStr := []string{"a", "b", "c"}
	require.Equal(t, iVal, v.Int)
	require.Equal(t, strVal, v.String)
	require.Equal(t, &iVal, v.PtrInt)
	require.Equal(t, sliceInt, v.SliceInt)
	require.Equal(t, sliceStr, v.SliceString)
	require.Equal(t, &sliceInt, v.PtrSliceInt)
}
