package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/config"
	auth "github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth/v1"
	shortener "github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener/v1"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/pg"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	conf, errConf := config.Load()
	if errConf != nil {
		log.Fatalf("Cannot load config: %s", errConf.Error())
	}

	if errMigrate := migrateDB(conf.DatabaseDSN); errMigrate != nil {
		log.Fatalf("Cannot migrate DB: %s", errMigrate.Error())
	}

	db, errDB := newDataSource(conf.DatabaseDSN)
	if errDB != nil {
		log.Fatalf("Cannot start DB: %s", errDB.Error())
	}

	authStorage := pg.NewAuthStorage(db)
	shortenerStorage := pg.NewShortenerStorage(db)

	authSrv := auth.New(authStorage)
	shortenerSrv := shortener.New(shortenerStorage, conf.BaseURL)

	server := api.New(shortenerSrv, authSrv, conf.ServerAddress, conf.BaseURL)
	log.Println("Starting server...")

	go start(server)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	if errShutdown := server.Shutdown(context.Background()); errShutdown != nil {
		log.Fatalf("Could not gracefully stop the server: %s", errShutdown.Error())
	}
	log.Println("Server stopped.")
}

func start(s *api.URLShortenerServer) {
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Cannot start the server: %v", err.Error())
	}
}

func migrateDB(databaseURL string) error {
	m, errMigrations := migrate.New("file://db/migrations", databaseURL)
	if errMigrations != nil {
		return fmt.Errorf("cannot init DB migrations: %w", errMigrations)
	}
	if errUp := m.Up(); errUp != nil && errUp != migrate.ErrNoChange {
		return fmt.Errorf("cannot apply migrations: %w", errUp)
	}

	return nil
}

func newDataSource(databaseURL string) (*sql.DB, error) {
	db, errOpen := sql.Open("pgx", databaseURL)
	if errOpen != nil {
		return nil, fmt.Errorf("cannot connect to DB: %w", errOpen)
	}

	if errPing := db.Ping(); errPing != nil {
		return nil, fmt.Errorf("cannot verify that DB connection is alive: %w", errPing)
	}

	return db, nil
}
