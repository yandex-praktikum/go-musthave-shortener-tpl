package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/pashagolub/pgxmock"
)

//test for endpoint "/"
func Test_HandlerGet(t *testing.T) {
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

	//init router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/:id", h.HandlerURLRelocation)

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//set mock
			urlRows := mock.NewRows([]string{"long_url", "is_deleted"}).
				AddRow("https://yandex.ru/search/?text=go&lr=11351&clid=9403", false)

			mock.ExpectQuery("SELECT long_url,is_deleted FROM shortens").
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

//================================================================================================

//test for endpoint "/user/urls"
func Test_HandlerGetList(t *testing.T) {
	type want struct {
		statusCode int
		list       []model.URL
	}
	tests := []struct {
		name      string
		sessionID string
		want      want
	}{
		{
			name:      "Ok",
			sessionID: "43786d99935441d19fa91d029bf83878",

			want: want{
				statusCode: http.StatusOK,

				list: []model.URL{
					{ShortURL: "http://localhost:8080/b5a41593cf656026",
						LongURL: "https://yandex.ru/search/?text=go&lr=11351&clid=9403"},
					{ShortURL: "http://localhost:8080/b5a41593cfdf6027",
						LongURL: "https://yandex.ru/search/?text=go&lr=11351&clid=6508"},
				},
			},
		},
		{
			name:      "Bad request",
			sessionID: "sdfgsdfgsdfg",

			want: want{
				statusCode: http.StatusNoContent,
				list:       []model.URL{},
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

	//init router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/user/urls", h.HandlerGetList)

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h.publicKey = "0ca876da1faed3aa87ae9d5ccfa5be17"

			//set mock
			urlRows := mock.NewRows([]string{"short_url", "long_url"}).
				AddRow("http://localhost:8080/b5a41593cf656026", "https://yandex.ru/search/?text=go&lr=11351&clid=9403").
				AddRow("http://localhost:8080/b5a41593cfdf6027", "https://yandex.ru/search/?text=go&lr=11351&clid=6508")

			mock.ExpectQuery("SELECT short_url, long_url FROM shortens").
				WithArgs(tt.sessionID).
				WillReturnRows(urlRows)

			//init http components
			req := httptest.NewRequest(http.MethodGet, "/user/urls", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			//read result
			result := w.Result()
			defer result.Body.Close()
			body, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.Marshal(tt.want.list)
			list := string(data)
			if list == "[]" {
				list = ""
			}
			assert.Equal(t, result.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(body), list)

		})
	}
}
