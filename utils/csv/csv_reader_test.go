package csv

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCSVReader(t *testing.T) {
	csvFileName := "./tests/data.csv"
	csvReader, err := NewCSVReader(csvFileName, 1)
	if err != nil {
		defer csvReader.Close()
	}
	require.NoError(t, err, "NewCSVReader(%v, 1)", csvFileName)
	record, err := csvReader.Read(context.TODO())
	require.NoError(t, err, "CSV Read()")
	require.Equal(t, "c1", record[0])
	require.Equal(t, "c2", record[1])
	require.Equal(t, "c3", record[2])
	require.Equal(t, "c4", record[3])
	require.Equal(t, "c5", record[4])
	require.Equal(t, "c6", record[5])
	require.Equal(t, "c7", record[6])
}
