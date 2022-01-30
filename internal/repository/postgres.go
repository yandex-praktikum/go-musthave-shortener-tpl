package repository

import (
	"context"
	"log"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/gorm"
)

type Client interface {
	//this interface need to mock BD in test(pgxmock have the same metods)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Ping(ctx context.Context) error
}

func NewDBClient(ctx context.Context, conf *configs.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, conf.DBConnectionStr)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
		return nil, err
	}
	log.Println("DB connection succes")
	return pool, nil
}

func Migration(conf *configs.Config) {
	db, err := gorm.Open("postgres", conf.DBConnectionStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&model.Shorten{}, &model.Session{})
	db.Model(&model.Shorten{}).AddForeignKey("session_ID", "Sessions(id)", "RESTRICT", "RESTRICT")
	log.Println("Migration succes")
}
