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

func TouchFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	return f.Close()
}

func IsFileExists(fileName string) (exists bool, err error) {
	if _, err = os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
