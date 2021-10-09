package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/config"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/backup"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/inmem"
)

type URLShortenerServer struct {
	http.Server
	Storage     storage.BulkStorage
	StorageFile string
}

// New makes an instance of HTTP server that is ready to run
func New(conf config.Config) *URLShortenerServer {
	storage := inmem.New()
	server := URLShortenerServer{
		Server: http.Server{
			Addr:    conf.ServerAddress,
			Handler: handler.New(storage, conf.BaseURL),
		},
		Storage:     storage,
		StorageFile: conf.StorageFile,
	}

	return &server
}

// ListenAndServe restores the server state from the backup file
// and starts the HTTP server
func (s *URLShortenerServer) ListenAndServe() error {
	if errRestore := backup.Restore(s.StorageFile, s.Storage); errRestore != nil {
		return fmt.Errorf("cannot restore URLs from storage file: %w", errRestore)
	}
	log.Printf("URL storage restored from [%s].", s.StorageFile)

	return s.Server.ListenAndServe()
}

// Shutdown backs-up the server state into the backup file
// and stops the HTTP server gracefully
func (s *URLShortenerServer) Shutdown(ctx context.Context) error {
	if errBackup := backup.Backup(s.StorageFile, s.Storage); errBackup != nil {
		return fmt.Errorf("cannot backup URLs to storage file: %w", errBackup)
	}
	log.Printf("URL storage backed up to [%s].", s.StorageFile)

	if errShutdown := s.Server.Shutdown(ctx); errShutdown != nil {
		return fmt.Errorf("cannot shutdown the server: %w", errShutdown)
	}

	return nil
}
