package csv

import (
	"context"
	"encoding/csv"
	"os"
	"path"

	"github.com/k81/log"
)

type Writer struct {
	FileName string
	file     *os.File
	writer   *csv.Writer
}

func NewWriter(fileName string) (writer *Writer, err error) {
	if err = os.MkdirAll(path.Dir(fileName), 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	writer = &Writer{
		FileName: fileName,
		file:     file,
		writer:   csv.NewWriter(file),
	}
	return writer, nil
}

func (writer *Writer) Write(ctx context.Context, record []string) (err error) {
	return writer.writer.Write(record)
}

func (writer *Writer) WriteAll(ctx context.Context, records [][]string) (err error) {
	return writer.writer.WriteAll(records)
}

func (writer *Writer) Close() error {
	writer.writer.Flush()
	if err := writer.writer.Error(); err != nil {
		log.Error(context.TODO(), "flushing csv writer", "file", writer.FileName, "error", err)
	}
	return writer.file.Close()
}
