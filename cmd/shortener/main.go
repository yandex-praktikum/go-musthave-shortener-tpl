package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

	authStorage, errAuthStorage := pg.NewAuthStorage(db)
	if errAuthStorage != nil {
		log.Fatalf("Cannot instantiate auth storage: %s", errAuthStorage.Error())
	}

	shortenerStorage, errShortenerStorage := pg.NewShortenerStorage(db)
	if errShortenerStorage != nil {
		log.Fatalf("Cannot instantiate shortener storage: %s", errShortenerStorage.Error())
	}

	authSrv, errAuthSrv := auth.New(authStorage)
	if errAuthSrv != nil {
		log.Fatalf("Cannot instantiate auth service: %s", errAuthSrv.Error())
	}

	shortenerSrv, errShortenerSrv := shortener.New(shortenerStorage, conf.BaseURL)
	if errShortenerSrv != nil {
		log.Fatalf("Cannot instantiate shortener service: %s", errShortenerSrv.Error())
	}

	server := api.New(shortenerSrv, authSrv, db, conf.ServerAddress, conf.BaseURL)

	awaitTermination()

	if errShutdown := server.Shutdown(context.Background()); errShutdown != nil {
		log.Fatalf("Could not gracefully stop the server: %s", errShutdown.Error())
	}
}

func awaitTermination() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func migrateDB(databaseURL string) error {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("cannot init DB migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("cannot apply migrations: %w", err)
	}

	return nil
}

func newDataSource(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot verify that DB connection is alive: %w", err)
	}

	return db, nil
}
