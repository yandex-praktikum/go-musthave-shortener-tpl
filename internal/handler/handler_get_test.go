package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/pashagolub/pgxmock"
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
			name:    "Ok",
			request: "/yandex",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://yandex.ru/search/?text=go&lr=11351&clid=9403",
			},
		},
		{
			name:    "Bad request",
			request: "/qwerty",
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
	}
	//init mock db connection
	mock, err := pgxmock.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer mock.Close(context.Background())

	//init main components
	config := configs.NewConfigForTest()
	r := repository.NewStorage(mock)
	s := service.NewService(r, config)
	h := NewHandler(s)

	//init server
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/:id", h.HandlerURLRelocation)

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//set mock
			urlRows := mock.NewRows([]string{"long_url"}).
				AddRow("https://yandex.ru/search/?text=go&lr=11351&clid=9403")

			mock.ExpectQuery("SELECT long_url FROM shortens").
				WithArgs("yandex").
				WillReturnRows(urlRows)

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
