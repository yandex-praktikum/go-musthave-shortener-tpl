package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/repository"
)

const testUrl = "https://yandex.ru/"

func TestCreateShortenerHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		payload string
		want    want
	}{
		{
			name:    "#1 post request test good payload",
			payload: testUrl,
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "#2 post request test empty payload",
			payload: "",
			want: want{
				code: http.StatusBadRequest,
				//				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	repo, err := repository.New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.payload))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateShortenerHandler(repo))
			h.ServeHTTP(w, request)
			res := w.Result()
			//status code
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d; got %d", tt.want.code, res.StatusCode)
			}
			//content-type
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected content-type %v; got %v", tt.want.contentType, res.Header.Get("Content-type"))
			}

		})
	}
}

func TestGetShortenerHandler(t *testing.T) {
	const testCode = "testtest"
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "#1 get request test",
			path: testCode,
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "application/text",
			},
		},
		{
			name: "#2 get request test",
			path: "_",
			want: want{
				code:        http.StatusNotFound,
				contentType: "application/text",
			},
		},
	}
	repo, err := repository.New()
	if err != nil {
		t.Fatal(err)
	}
	err = repo.Shortener.SaveShortener(testCode, &model.Shortener{URL: testUrl})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.path), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetShortenerHandler(repo))
			h.ServeHTTP(w, request)
			res := w.Result()
			//status code
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d; got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}
