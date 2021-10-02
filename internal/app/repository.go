package app

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sync"
)

type StorableURL struct {
	URL string
	ID  int
}

type Repository interface {
	GetURLBy(id int) *url.URL
	SaveURL(u url.URL) int
	Restore(fileName string) error
	Backup(fineName string) error
}

type MemRepository struct {
	sync.RWMutex

	store map[int]url.URL
}

func NewMemRepository() Repository {
	return &MemRepository{
		RWMutex: sync.RWMutex{},
		store:   make(map[int]url.URL),
	}
}

func (r *MemRepository) Restore(fileName string) error {
	reader, errOpen := NewReader(fileName)
	if errOpen != nil {
		return errOpen
	}
	defer reader.Close()

	for {
		storableURL, errDecode := reader.ReadURL()
		if errDecode != nil {
			return fmt.Errorf("cannot decode backed-up URL: %w", errDecode)
		}
		if storableURL == nil {
			break
		}

		url, errParse := url.Parse(storableURL.URL)
		if errParse != nil {
			return fmt.Errorf("cannot parse backed-up URL [%s]: %w", storableURL.URL, errParse)
		}
		r.SaveURL(*url)
		log.Printf("Url restored [%s]", url)
	}

	return nil
}

func (r *MemRepository) Backup(fileName string) error {
	writer, errOpen := NewWriter(fileName)
	if errOpen != nil {
		return errOpen
	}
	defer writer.Close()

	for id, shortURL := range r.store {
		errWrite := writer.WriteURL(StorableURL{
			ID:  id,
			URL: shortURL.String(),
		})
		if errWrite != nil {
			return errWrite
		}
		log.Printf("Url backed up [%s]", &shortURL)
	}
	return nil
}

func (r *MemRepository) SaveURL(u url.URL) int {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	id := len(r.store)
	r.store[id] = u

	return id
}

func (r *MemRepository) GetURLBy(id int) *url.URL {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	longURL, ok := r.store[id]
	if !ok {
		return nil
	}
	return &longURL
}

type reader struct {
	file    *os.File
	decoder *gob.Decoder
}

func NewReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &reader{
		file:    file,
		decoder: gob.NewDecoder(bufio.NewReader(file)),
	}, nil
}

func (r *reader) ReadURL() (*StorableURL, error) {
	url := &StorableURL{}
	errDecode := r.decoder.Decode(url)
	if errDecode == io.EOF {
		return nil, nil
	}
	if errDecode != nil {
		return nil, errDecode
	}
	return url, nil
}

func (r *reader) Close() error {
	return r.file.Close()
}

type writer struct {
	file      *os.File
	bufWriter *bufio.Writer
	encoder   *gob.Encoder
}

func NewWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_TRUNC|os.O_CREATE, 0777)
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

func (w *writer) WriteURL(u StorableURL) error {
	errEncode := w.encoder.Encode(u)
	if errEncode != nil {
		return fmt.Errorf("cannot write to storage: %w", errEncode)
	}
	return nil
}

func (w *writer) Close() error {
	errFlush := w.bufWriter.Flush()
	if errFlush != nil {
		return fmt.Errorf("cannot write buffered data to file: %w", errFlush)
	}
	return w.file.Close()
}
