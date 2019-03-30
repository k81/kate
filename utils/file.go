package utils

import (
	"bytes"
	"io"
	"os"
)

func CountLine(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], []byte{'\n'})
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
