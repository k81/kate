package csv

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/k81/kate/utils"
)

type Reader struct {
	FileName string
	file     *os.File
	scanner  *csv.Reader
}

func NewReader(fileName string, skipLine int) (reader *Reader, err error) {
	reader = &Reader{FileName: fileName}
	if reader.file, err = os.Open(fileName); err != nil {
		return nil, err
	}
	reader.scanner = csv.NewReader(reader.file)
	for i := 0; i < skipLine; i++ {
		if _, err = reader.scanner.Read(); err != nil {
			if err == io.EOF {
				break
			}
			// nolint:errcheck
			_ = reader.file.Close()
			return nil, err
		}
	}
	return reader, nil
}

func (reader *Reader) Read(ctx context.Context) (record []string, err error) {
	return reader.scanner.Read()
}

func (reader *Reader) Count() (int, error) {
	return utils.CountLine(reader.FileName)
}

func (reader *Reader) Close() (err error) {
	return reader.file.Close()
}
