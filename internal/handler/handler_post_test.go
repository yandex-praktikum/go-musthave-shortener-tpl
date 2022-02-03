package handler

import (
	"bytes"
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
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/pashagolub/pgxmock"
)

//test for endpoint "/" and auth middleware
func Test_HandlerPostText(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		ReqBody string
		want    want
	}{
		{
			name:    "Ok",
			ReqBody: "https://yandex.ru/search/test1",
			want: want{
				statusCode: http.StatusConflict,
				body:       "http://localhost:8080/b5a41593cf656026",
			},
		},
		{
			name:    "Bad request",
			ReqBody: "sdfsfsdfsdf",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "error: Not Allowd request",
			},
		},
	}

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
			router.Use(AuthMiddleware(h))
			router.POST("/", h.HandlerPostURL)

			//init http components
			w := httptest.NewRecorder()

			req := *httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.ReqBody))

			//set auth mock
			sesRows := mock.NewRows([]string{"id"}).AddRow(1)

			mock.ExpectQuery("SELECT id FROM sessions").
				WithArgs("43786d99935441d19fa91d029bf83878").
				WillReturnRows(sesRows)

				//set urlsaving mock
			urlRows := mock.NewRows([]string{"id", "short_url"}).
				AddRow(1, "http://localhost:8080/b5a41593cf656026")

			mock.ExpectQuery("INSERT INTO shortens").
				WillReturnRows(urlRows)

			//set content-type
			req.Header.Set("content-type", "text/plain")

			//set cookie
			cookie := &http.Cookie{Name: "session", Value: "0ca876da1faed3aa87ae9d5ccfa5be17"}
			req.AddCookie(cookie)

			//run server
			router.ServeHTTP(w, &req)

			//read result
			result := w.Result()
			body, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer result.Body.Close()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(body), tt.want.body)
		})
	}
}

//========================================================================================================

//test for endpoint "/api/shorten"
func Test_HandlerPostJSON(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string `json:"-"`
		ReqBody string `json:"url"`
		want    want   `json:"-"`
	}{
		{
			name:    "Ok",
			ReqBody: "https://yandex.ru/search/test1",
			want: want{
				statusCode: http.StatusConflict,
				body:       `{"result":"http://localhost:8080/b5a41593cf656026"}`,
			},
		},
		{
			name:    "Bad request",
			ReqBody: "sdfsfsdfdghdfghfsdf",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "error: Not Allowd request",
			},
		},
	}

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
			router.Use(AuthMiddleware(h))
			router.POST("/api/shorten", h.HandlerPostURL)

			//init http components
			w := httptest.NewRecorder()

			body, err := json.Marshal(tt)
			if err != nil {
				log.Fatal(err)
			}

			req := *httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))

			//set auth mock
			sesRows := mock.NewRows([]string{"id"}).AddRow(1)

			mock.ExpectQuery("SELECT id FROM sessions").
				WithArgs("43786d99935441d19fa91d029bf83878").
				WillReturnRows(sesRows)

				//set urlsaving mock
			urlRows := mock.NewRows([]string{"id", "short_url"}).
				AddRow(1, "http://localhost:8080/b5a41593cf656026")

			mock.ExpectQuery("INSERT INTO shortens").
				WillReturnRows(urlRows)

			//set content-type
			req.Header.Set("content-type", "application/json")

			//set cookie
			cookie := &http.Cookie{Name: "session", Value: "0ca876da1faed3aa87ae9d5ccfa5be17"}
			req.AddCookie(cookie)

			//run server
			router.ServeHTTP(w, &req)

			//read result
			result := w.Result()
			respBody, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer result.Body.Close()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(respBody), tt.want.body)
		})
	}
}
