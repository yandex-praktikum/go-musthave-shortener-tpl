package web

import (
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
)

// ShortenerRepo ...
type ShortenerRepo struct {
	Memory map[string]*model.Shortener
}

// NewShortenerRepo ...
func NewShortenerRepo(db map[string]*model.Shortener) *ShortenerRepo {
	var repo ShortenerRepo
	repo.Memory = db
	return &repo
}

func (repo *ShortenerRepo) GetShortenerBy(id string) (*model.Shortener, error) {
	shortener := repo.Memory[id]
	return shortener, nil
}

func (repo *ShortenerRepo) SaveShortener(code string, shortener *model.Shortener) error {
	repo.Memory[code] = shortener
	return nil
}

func (repo *ShortenerRepo) IncludesCode(code string) bool {
	_, result := repo.Memory[code]
	return result
}
