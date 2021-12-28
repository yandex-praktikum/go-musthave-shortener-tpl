package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func TestHandler_HandlerPost(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := repository.NewStorage()
			s := service.NewService(r)
			h := NewHandler(s)

			router := gin.Default()
			router.POST("/", h.HandlerPost)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)

		})
	}
}
