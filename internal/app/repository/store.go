package repository

import (
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/repository/web"
)

type Store struct {
	MemoryDB  map[string]*model.Shortener
	Shortener Repository
}

func New() (*Store, error) {
	db := make(map[string]*model.Shortener)
	var store Store
	store.MemoryDB = db
	store.Shortener = web.NewShortenerRepo(db)
	return &store, nil
}
