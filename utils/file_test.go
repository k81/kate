package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountLine(t *testing.T) {
	fileName := "./file.go"
	count, err := CountLine(fileName)
	require.NoError(t, err, "CountLine")
	t.Logf("CountLine(%s) = %v", fileName, count)
}
