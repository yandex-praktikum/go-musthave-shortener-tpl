package app

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/backup"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/inmem"
)

type URLShortener struct {
	*chi.Mux
	Repo    storage.Storage
	BaseURL url.URL
}

type URLShortenerServer struct {
	http.Server
	Repo        storage.BulkStorage
	StorageFile string
}

func (s *URLShortenerServer) ListenAndServe() error {
	if errRestore := backup.Restore(s.StorageFile, s.Repo); errRestore != nil {
		return fmt.Errorf("cannot restore URLs from storage file: %w", errRestore)
	}
	log.Printf("URL repository restored from [%s].", s.StorageFile)
	return s.Server.ListenAndServe()
}

func (s *URLShortenerServer) Shutdown(ctx context.Context) error {
	if errBackup := backup.Backup(s.StorageFile, s.Repo); errBackup != nil {
		return fmt.Errorf("cannot backup URLs to storage file: %w", errBackup)
	}
	log.Printf("URL repository backed up to [%s].", s.StorageFile)
	if errShutdown := s.Server.Shutdown(ctx); errShutdown != nil {
		return fmt.Errorf("cannot shutdown the server: %w", errShutdown)
	}
	return nil
}

func NewServer(conf Config) *URLShortenerServer {
	repo := inmem.New()
	return &URLShortenerServer{
		Server: http.Server{
			Addr:    conf.ServerAddress,
			Handler: NewURLShortener(repo, conf.BaseURL),
		},
		Repo:        repo,
		StorageFile: conf.StorageFile,
	}
}

func NewURLShortener(repo storage.Storage, baseURL url.URL) http.Handler {
	shortener := &URLShortener{
		Mux:     chi.NewMux(),
		Repo:    repo,
		BaseURL: baseURL,
	}
	shortener.Use(gzipDecompressor)
	shortener.Use(gzipCompressor)
	shortener.Post("/", shortener.handlePostLongURL)
	shortener.Post("/api/shorten", shortener.handlePostAPIShorten)
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

	newStorableURL := model.NewStorableURL(longURL)
	shortURL, errShorten := s.ShortenURL(newStorableURL)
	if errShorten != nil {
		log.Printf("Cannot shorten url: %s", errShorten.Error())
		http.Error(w, "Cannot shorten url", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, errWrite := fmt.Fprint(w, shortURL.String()); errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}

func (s *URLShortener) ShortenURL(newURL model.StorableURL) (*url.URL, error) {
	url := s.Repo.Save(newURL)
	urlPath := fmt.Sprintf("%d", url.ID)
	shortURL, err := s.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, fmt.Errorf("cannot shorten URL for id [%d]", url.ID)
	}
	log.Printf("Shortened: %s - %s", url, shortURL)
	return shortURL, nil
}

type LongURLJson struct {
	URL string `json:"url"`
}

type ShortURLJson struct {
	Result string `json:"result"`
}

func (s *URLShortener) handlePostAPIShorten(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	longURLJson := LongURLJson{}
	if errDec := dec.Decode(&longURLJson); errDec != nil {
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

	storableURL := model.NewStorableURL(longURL)
	shortURL, errShorten := s.ShortenURL(storableURL)
	if errShorten != nil {
		log.Printf("Cannot shorten url: %s", errShorten.Error())
		http.Error(w, "Cannot shorten url", http.StatusInternalServerError)
		return
	}

	shortURLJson := ShortURLJson{Result: shortURL.String()}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if errEnc := enc.Encode(shortURLJson); errEnc != nil {
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

	url := s.Repo.GetByID(id)
	if url == nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("Found: %d - %v", id, url)

	w.Header().Add("Location", url.LongURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Println("Cannot create gzip writer", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{w, gz}, r)
	})
}

type gzipReader struct {
	io.ReadCloser
	Reader io.Reader
}

func (r gzipReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func gzipDecompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println("Cannot create gzip reader", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		r.Body = gzipReader{r.Body, gz}

		next.ServeHTTP(w, r)
	})
}
