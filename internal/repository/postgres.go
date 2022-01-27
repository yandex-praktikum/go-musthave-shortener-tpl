package repository

import (
	"context"
	"log"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/gorm"
)

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
