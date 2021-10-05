package backup

import (
	"fmt"
	"log"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

func Restore(fileName string, s storage.BulkStorage) error {
	reader, errOpen := newReader(fileName)
	if errOpen != nil {
		return errOpen
	}
	defer reader.close()

	for {
		url, errDecode := reader.readURL()
		if errDecode != nil {
			return fmt.Errorf("cannot decode backed-up URL: %w", errDecode)
		}
		if url == nil {
			return nil
		}

		s.Load(*url)
		log.Printf("Url restored [%v]", url)
	}
}
