package app

import (
	"encoding/json"
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
	SaveURL(u url.URL) int
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
	shortener.Post("/api/shorten", shortener.handlePostApiShorten)
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
	longURL, errParse := url.Parse(string(rawURL))
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	shortURL := s.ShortenURL(*longURL)

	w.WriteHeader(http.StatusCreated)
	_, errWrite := fmt.Fprint(w, shortURL.String())
	if errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}

func (s *URLShortener) ShortenURL(longURL url.URL) url.URL {
	id := s.Repo.SaveURL(longURL)
	log.Printf("Shortened: %s - %d", longURL.String(), id)
	rawShortURL := fmt.Sprintf("http://%s/%d", ServiceAddr, id)
	shortURL, err := url.Parse(rawShortURL)
	if err != nil {
		log.Panicf("Cannot make URL [%s]", rawShortURL)
	}
	return *shortURL
}

type LongURLJson struct {
	URL string `json:"url"`
}

type ShortURLJson struct {
	Result string `json:"result"`
}

func (s *URLShortener) handlePostApiShorten(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	longURLJson := LongURLJson{}
	errDec := dec.Decode(&longURLJson)
	if errDec != nil {
		msg := fmt.Sprintf("Cannot decode request body: %v", errDec)
		http.Error(w, msg, http.StatusBadRequest)
		log.Println(msg)
		return
	}
	longURL, errParse := url.Parse(longURLJson.URL)
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	log.Printf("longURLJson.Url: [%v]", longURL)
	shortURL := s.ShortenURL(*longURL)
	shortURLJson := ShortURLJson{Result: shortURL.String()}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	errEnc := enc.Encode(shortURLJson)
	if errEnc != nil {
		log.Printf("Cannot write response: %v", errEnc)
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
	sync.RWMutex

	store map[int]url.URL
}

func NewRepository() Repository {
	return &MemRepository{
		RWMutex: sync.RWMutex{},
		store:   make(map[int]url.URL),
	}
}

func (r *MemRepository) SaveURL(u url.URL) int {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	id := len(r.store)
	r.store[id] = u

	return id
}

func (r *MemRepository) GetURLBy(id int) *url.URL {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	longURL, ok := r.store[id]
	if !ok {
		return nil
	}
	return &longURL
}
