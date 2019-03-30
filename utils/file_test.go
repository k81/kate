package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCountLine(t *testing.T) {
	fileName := "./file.go"
	tBegin := time.Now()
	count, err := CountLine(fileName)
	require.NoError(t, err, "CountLine")
	t.Logf("CountLine(%s) = %v, elapsed_ms=%v", fileName, count, time.Since(tBegin))
}
