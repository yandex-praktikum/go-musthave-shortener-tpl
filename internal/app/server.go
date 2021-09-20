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
	GetUrlBy(id int) *url.URL
	SaveUrl(u *url.URL) int
}

type UrlShortener struct {
	*chi.Mux
	Repo Repository
}

func NewServer() *http.Server {
	return &http.Server{
		Addr:    ServiceAddr,
		Handler: NewUrlShortener(NewRepository()),
	}
}

func NewUrlShortener(repo Repository) http.Handler {
	shortener := &UrlShortener{
		Mux:  chi.NewMux(),
		Repo: repo,
	}
	shortener.Post("/", shortener.handlePostLongUrl)
	shortener.Get("/{id}", shortener.handleGetShortUrl)

	return shortener
}

func (s *UrlShortener) handlePostLongUrl(w http.ResponseWriter, r *http.Request) {
	rawUrl, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("Cannot read request: %v", errRead)
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}

	log.Printf("Got url to shorten: %v", rawUrl)
	url, errParse := url.Parse(string(rawUrl))
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	id := s.Repo.SaveUrl(url)
	log.Printf("Shortened: %v - %d", url, id)

	shortUrl := fmt.Sprintf("http://%s/%d", ServiceAddr, id)

	w.WriteHeader(http.StatusCreated)
	_, errWrite := w.Write([]byte(shortUrl))
	if errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}

func (s *UrlShortener) handleGetShortUrl(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid id [%v]", idStr)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	url := s.Repo.GetUrlBy(id)
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
		store:     make(map[int]*url.URL, 0),
		storeLock: sync.Mutex{},
	}
}

func (r *MemRepository) SaveUrl(u *url.URL) int {
	r.storeLock.Lock()
	defer r.storeLock.Unlock()

	id := len(r.store)
	r.store[id] = u

	return id
}

func (r *MemRepository) GetUrlBy(id int) *url.URL {
	r.storeLock.Lock()
	defer r.storeLock.Unlock()

	return r.store[id]
}
