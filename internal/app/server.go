package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type Repository interface {
	GetUrlBy(id int) *url.URL
	SaveUrl(u *url.URL) int
}

type UrlShortener struct {
	repo Repository
}

func NewServer() *http.Server {
	shortener := UrlShortener{
		repo: NewRepository(),
	}
	return &http.Server{
		Addr:    "localhost:8080",
		Handler: http.HandlerFunc(shortener.handler),
	}
}

func (s *UrlShortener) handler(w http.ResponseWriter, r *http.Request) {
	if matchesPostLongUrl(r) {
		s.handlePostLongUrl(w, r)
		return
	}

	if matchesGetShortUrl(r) {
		s.handleGetShortUrl(w, r)
		return
	}

	http.NotFound(w, r)
}

func matchesPostLongUrl(r *http.Request) bool {
	return r.URL.Path == "/" && r.Method == http.MethodPost
}

func matchesGetShortUrl(r *http.Request) bool {
	pathParts := strings.Split(r.URL.Path, "/")
	return len(pathParts) == 2
}

func (s *UrlShortener) handlePostLongUrl(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, errRead := buf.ReadFrom(r.Body)
	if errRead != nil {
		log.Printf("Cannot read request: %v", errRead)
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}

	rawUrl := buf.String()
	log.Printf("Got url to shorten: %v", rawUrl)
	url, errParse := url.Parse(rawUrl)
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	log.Printf("Shortened: %v", url)
	id := s.repo.SaveUrl(url)

	shortUrl := fmt.Sprintf("http://localhost:8080/%d", id)

	w.WriteHeader(http.StatusCreated)
	_, errWrite := w.Write([]byte(shortUrl))
	if errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}

func (s *UrlShortener) handleGetShortUrl(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		log.Printf("Invalid id: %v", pathParts[1])
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	url := s.repo.GetUrlBy(id)
	if url == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Add("Location", url.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type MemRepository struct {
	store     map[int]*url.URL
	storeLock *sync.Mutex
}

func NewRepository() Repository {
	return &MemRepository{
		store:     make(map[int]*url.URL, 0),
		storeLock: &sync.Mutex{},
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
	defer r.storeLock.Unlock()

	url, found := r.store[id]
	log.Printf("Found: %v", r.store)
	if !found {
		return nil
	}

	return url
}
