package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type aContent struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type bContent struct {
	SubType int         `json:"sub_type"`
	Value   interface{} `json:"value" dynamic:"true"`
}

type bValue1 []int
type bValue2 int

func (b *bContent) NewDynamicField(fieldName string) interface{} {
	switch b.SubType {
	case 1:
		return new(bValue1)
	case 2:
		return new(bValue2)
	}
	return nil
}

type dynamicReq struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content" dynamic:"true"`
}

func (r *dynamicReq) NewDynamicField(fieldName string) interface{} {
	switch r.Type {
	case "A":
		return new(aContent)
	case "B":
		return new(bContent)
	}
	return nil
}

func TestParseJSON(t *testing.T) {
	ptrA := &dynamicReq{}
	inputA := []byte(`{"type":"A", "content":{"name":"zhangsan", "value": 11}}`)
	errA := ParseJSON(bytes.NewBuffer(inputA), ptrA)
	require.NoError(t, errA, "parse inputA")
	require.Equal(t, "A", ptrA.Type)
	require.IsType(t, new(aContent), ptrA.Content)
	require.Equal(t, "zhangsan", ptrA.Content.(*aContent).Name)
	require.Equal(t, 11, ptrA.Content.(*aContent).Value)

	ptrB1 := &dynamicReq{}
	inputB1 := []byte(`{"type":"B", "content":{"sub_type":1, "value":[1,3,5]}}`)
	errB1 := ParseJSON(bytes.NewBuffer(inputB1), ptrB1)
	require.NoError(t, errB1, "parse inputB1")
	require.Equal(t, "B", ptrB1.Type)
	require.IsType(t, new(bContent), ptrB1.Content)
	require.Equal(t, &bValue1{1, 3, 5}, ptrB1.Content.(*bContent).Value)

	ptrB2 := &dynamicReq{}
	inputB2 := []byte(`{"type":"B", "content":{"sub_type":2, "value":13}}`)
	errB2 := ParseJSON(bytes.NewBuffer(inputB2), ptrB2)
	require.NoError(t, errB2, "parse inputB2")
	require.Equal(t, "B", ptrB2.Type)
	require.IsType(t, new(bContent), ptrB2.Content)
	bValue2Expected := bValue2(13)
	require.Equal(t, &bValue2Expected, ptrB2.Content.(*bContent).Value)
}
