package backup

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type reader struct {
	file    *os.File
	decoder *gob.Decoder
}

func newReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &reader{
		file:    file,
		decoder: gob.NewDecoder(bufio.NewReader(file)),
	}, nil
}

func (r *reader) readURL() (*model.StoreURL, error) {
	u := &gobURL{}
	errDecode := r.decoder.Decode(u)
	if errDecode == io.EOF {
		return nil, nil
	}
	if errDecode != nil {
		return nil, errDecode
	}
	storeURL, errParse := u.ToStoreURL()
	if errParse != nil {
		return nil, fmt.Errorf("cannot read StoreURL: %w", errParse)
	}

	return storeURL, nil
}

func (r *reader) close() error {
	return r.file.Close()
}
