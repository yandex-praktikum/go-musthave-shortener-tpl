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

//this interface implements pgx.Conn, pgx.Pool and pgx.Mock
type Client interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Ping(ctx context.Context) error
}

func NewDBClient(ctx context.Context, conf *configs.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, conf.DBConnectionStr)
	if err != nil {
		return nil, err
	}
	log.Println("DB connection success")
	return pool, nil
}

func Migration(conf *configs.Config) error {
	db, err := gorm.Open("postgres", conf.DBConnectionStr)
	if err != nil {
		log.Fatal(err)
		return err

	}
	defer db.Close()

	db.AutoMigrate(&model.Shorten{}, &model.Session{})
	db.Model(&model.Shorten{}).AddForeignKey("session_ID", "Sessions(id)", "RESTRICT", "RESTRICT")
	log.Println("Migration success")

	return nil
}
