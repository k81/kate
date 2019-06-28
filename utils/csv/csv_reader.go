package csv

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/k81/kate/utils"
)

type CSVReader struct {
	FileName string
	file     *os.File
	scanner  *csv.Reader
}

func NewCSVReader(fileName string, skipLine int) (reader *CSVReader, err error) {
	csvReader := &CSVReader{FileName: fileName}
	if csvReader.file, err = os.Open(fileName); err != nil {
		return nil, err
	}
	csvReader.scanner = csv.NewReader(csvReader.file)
	for i := 0; i < skipLine; i++ {
		if _, err = csvReader.scanner.Read(); err != nil {
			if err == io.EOF {
				break
			}
			// nolint:errcheck
			_ = csvReader.file.Close()
			return nil, err
		}
	}
	return csvReader, nil
}

func (reader *CSVReader) Read(ctx context.Context) (record []string, err error) {
	return reader.scanner.Read()
}

func (reader *CSVReader) Count() (int, error) {
	return utils.CountLine(reader.FileName)
}

func (reader *CSVReader) Close() (err error) {
	return reader.file.Close()
}
