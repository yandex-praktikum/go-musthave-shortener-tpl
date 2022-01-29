package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func TestHandler_HandlerPostText(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "test 1",
			requestBody: "https://yandex.ru/search/?text=go&lr=11351&clid=9403sdfasdfasdfasdf",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:        "test 2",
			requestBody: "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "test 3",
			requestBody: "sdfsdfsdfsdfsdfsdf",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config := configs.NewConfigForTest()

			ctx := context.TODO()
			db, err := repository.NewDBClient(ctx, config)
			if err != nil {
				log.Fatal(err)
			}

			r := repository.NewStorage(db)
			s := service.NewService(r, config)
			h := NewHandler(s)

			gin.SetMode(gin.ReleaseMode)
			router := gin.Default()

			router.POST("/", h.HandlerPostURL)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()
			req.Header.Set("content-type", "text/plain")
			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)

		})
	}
}
func TestHandler_HandlerPostJSON(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string `json:"-"`
		RequestBody string `json:"url"`
		contentType string `json:"-"`
		want        want   `json:"-"`
	}{
		{
			name:        "test 1",
			RequestBody: "https://yandex.ru/search/?text=go&lr=11351&clid=9403sdfasdfasdfasdf",
			contentType: "application/json",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:        "test 2",
			RequestBody: "https://yandex.ru/search/?text=go&lr=11351&clid=9403sdfasdfasdfasdf",
			contentType: "application/json",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config := configs.NewConfigForTest()

			ctx := context.TODO()
			db, err := repository.NewDBClient(ctx, config)
			if err != nil {
				log.Fatal(err)
			}

			r := repository.NewStorage(db)
			s := service.NewService(r, config)
			h := NewHandler(s)

			gin.SetMode(gin.ReleaseMode)
			router := gin.Default()

			router.POST("/api/shorten", h.HandlerPostURL)

			body, err := json.Marshal(tt)
			if err != nil {
				log.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			req.Header.Set("content-type", tt.contentType)
			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)

		})
	}
}
