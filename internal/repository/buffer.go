package repository

import (
	"sync"
	"time"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
)

type DeleteBuffer struct {
	Buffer     []model.URL
	Mutex      sync.Mutex
	Signal     chan struct{}
	Msg        chan string
	LastUpdate time.Duration
}

func NewDeleteBuffer() *DeleteBuffer {
	return &DeleteBuffer{
		Buffer: make([]model.URL, 0, 10),
		Signal: make(chan struct{}),
		Msg:    make(chan string),
	}
}
func (buf *DeleteBuffer) ClearBuffer() {
	buf.Buffer = buf.Buffer[:0]
}
