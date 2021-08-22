package repository

import (
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
)

type Repository interface {
	GetShortenerBy(string) (*model.Shortener, error)
	SaveShortener(string, *model.Shortener) error
	IncludesCode(string) bool
}
