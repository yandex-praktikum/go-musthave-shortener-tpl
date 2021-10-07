package backup

import (
	"log"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

// Backup stores all the BulkStorage state into a file
func Backup(fileName string, s storage.BulkStorage) error {
	writer, errOpen := newWriter(fileName)
	if errOpen != nil {
		return errOpen
	}
	defer writer.close()

	for _, url := range s.GetAll() {
		if errWrite := writer.writeURL(url); errWrite != nil {
			return errWrite
		}
		log.Printf("Url backed up [%s]", url)
	}

	return nil
}
