package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipReader struct {
	io.ReadCloser
	reader io.Reader
}

func (r gzipReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

// GzipDecompressor middleware gzip-decodes encoded HTTP requests
func GzipDecompressor(next http.Handler) http.Handler {
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
