package storage

import "github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"

type Storage interface {
	GetByID(id int) *model.StoreURL
	Save(model.StorableURL) model.StoreURL
}

type BulkStorage interface {
	Storage
	GetAll() []model.StoreURL
	Load(u model.StoreURL)
}
