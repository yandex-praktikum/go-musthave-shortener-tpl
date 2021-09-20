package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

const ServiceAddr = "localhost:8080"

type Repository interface {
	GetURLBy(id int) *url.URL
	SaveURL(u *url.URL) int
}

type URLShortener struct {
	*chi.Mux
	Repo Repository
}

func NewServer() *http.Server {
	return &http.Server{
		Addr:    ServiceAddr,
		Handler: NewURLShortener(NewRepository()),
	}
}

func NewURLShortener(repo Repository) http.Handler {
	shortener := &URLShortener{
		Mux:  chi.NewMux(),
		Repo: repo,
	}
	shortener.Post("/", shortener.handlePostLongURL)
	shortener.Get("/{id}", shortener.handleGetShortURL)

	return shortener
}

func (s *URLShortener) handlePostLongURL(w http.ResponseWriter, r *http.Request) {
	rawURL, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("Cannot read request: %v", errRead)
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}

	log.Printf("Got url to shorten: %s", rawURL)
	url, errParse := url.Parse(string(rawURL))
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	id := s.Repo.SaveURL(url)
	log.Printf("Shortened: %v - %d", url, id)

	shortURL := fmt.Sprintf("http://%s/%d", ServiceAddr, id)

	w.WriteHeader(http.StatusCreated)
	_, errWrite := w.Write([]byte(shortURL))
	if errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}

func (s *URLShortener) handleGetShortURL(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid id [%v]", idStr)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	url := s.Repo.GetURLBy(id)
	if url == nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("Found: %d - %v", id, url)

	w.Header().Add("Location", url.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type MemRepository struct {
	store     map[int]*url.URL
	storeLock sync.Mutex
}

func NewRepository() Repository {
	return &MemRepository{
		store:     make(map[int]*url.URL),
		storeLock: sync.Mutex{},
	}
}

func (r *MemRepository) SaveURL(u *url.URL) int {
	r.storeLock.Lock()
	defer r.storeLock.Unlock()

	id := len(r.store)
	r.store[id] = u

	return id
}

func (r *MemRepository) GetURLBy(id int) *url.URL {
	r.storeLock.Lock()
	defer r.storeLock.Unlock()

	return r.store[id]
}
