package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EMus88/GO-Yandex-Study/internal/app/service"
	"github.com/EMus88/GO-Yandex-Study/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func TestHandler_HandlerGet(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "test 1",
			request: "/yandex",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://yandex.ru/search/?text=go&lr=11351&clid=9403",
			},
		},
		{
			name:    "test 2",
			request: "/wiki",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://ru.wikipedia.org/wiki/Go",
			},
		},
		{
			name:    "test 3",
			request: "/qwerty",
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.NewStorage()
			storage.SaveURL("yandex", "https://yandex.ru/search/?text=go&lr=11351&clid=9403")
			storage.SaveURL("wiki", "https://ru.wikipedia.org/wiki/Go")

			s := service.NewService(storage)
			h := NewHandler(s)

			router := gin.Default()
			router.GET("/:id", h.HandlerGet)

			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, result.StatusCode, tt.want.statusCode)
			assert.Equal(t, result.Header.Get("Location"), tt.want.location)

		})
	}
}
