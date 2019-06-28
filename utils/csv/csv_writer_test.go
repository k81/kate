package csv

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCSVWriter(t *testing.T) {
	csvFileName := "./tests/writer.csv"
	csvWriter, err := NewCSVWriter(csvFileName)
	if err == nil {
		defer func() {
			_ = csvWriter.Close()
			_ = os.Remove(csvFileName)
		}()
	}
	require.NoError(t, err, "NewCSVWriter(%v)", csvFileName)
	err = csvWriter.Write(context.TODO(), []string{"c1", "c2", "c3", "c4"})
	require.NoError(t, err, "write single record")
	err = csvWriter.WriteAll(context.TODO(), [][]string{{"a1", "a2", "a3", "a4"}, {"b1", "b2", "b3", "b4"}})
	require.NoError(t, err, "write batch records")
}
