package backup

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type writer struct {
	file      *os.File
	bufWriter *bufio.Writer
	encoder   *gob.Encoder
}

func newWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	bufWriter := bufio.NewWriter(file)
	encoder := gob.NewEncoder(bufWriter)

	return &writer{
		file,
		bufWriter,
		encoder,
	}, nil
}

func (w *writer) writeURL(u model.StoreURL) error {
	gobU := newGobURL(u)
	if errEncode := w.encoder.Encode(gobU); errEncode != nil {
		return fmt.Errorf("cannot write to storage: %w", errEncode)
	}
	return nil
}

func (w *writer) close() error {
	if errFlush := w.bufWriter.Flush(); errFlush != nil {
		return fmt.Errorf("cannot write buffered data to file: %w", errFlush)
	}
	return w.file.Close()
}
