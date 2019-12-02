package csv

import (
	"context"
	"encoding/csv"
	"os"
	"path"

	"go.uber.org/zap"
)

type Writer struct {
	FileName string
	file     *os.File
	writer   *csv.Writer
	logger   *zap.Logger
}

func NewWriter(fileName string, logger *zap.Logger) (writer *Writer, err error) {
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
		logger:   logger,
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
		writer.logger.Error("flushing csv writer", zap.String("file", writer.FileName), zap.Error(err))
	}
	return writer.file.Close()
}
