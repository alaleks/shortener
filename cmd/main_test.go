package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alaleks/shortener/internal/router"
)

var handler http.Handler

func init() {
	handler = router.Create()
}

func TestUseShortner(t *testing.T) {

	tests := []struct {
		name     string
		code     int
		url      string
		endpoint string
		body     string
	}{
		{
			name:     "Тест url с https",
			code:     201,
			url:      "https://github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/",
			body:     "http://localhost:8080/",
		},
		{
			name:     "Тест url с http",
			code:     201,
			url:      "http://github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/",
			body:     "http://localhost:8080/",
		},
		{
			name:     "Тест url только с www",
			code:     201,
			url:      "www.github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/",
			body:     "http://localhost:8080/",
		},
		{
			name:     "Тест url без http..",
			code:     400,
			url:      "github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/",
			body:     "invalid url",
		},
		{
			name:     "Тест c ошибкой в http",
			code:     400,
			url:      "htts://github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/",
			body:     "invalid url",
		},
		{
			name:     "Тест с пустой строкой после https:// ",
			code:     400,
			url:      "https://",
			endpoint: "http://localhost:8080/",
			body:     "invalid url",
		},
		{
			name:     "Тест c пустым body",
			code:     400,
			url:      "",
			endpoint: "http://localhost:8080/",
			body:     "url is empty",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, v.endpoint, bytes.NewBuffer([]byte(v.url)))
			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, req)
			res := responseRecorder.Result()
			if res != nil {
				defer res.Body.Close()
			}
			resBody, _ := io.ReadAll(res.Body)

			if res.StatusCode != v.code {
				t.Errorf("status code must be %d but is %d", v.code, res.StatusCode)
			}
			if !strings.HasPrefix(string(resBody), v.body) {
				t.Errorf("shortener must was have returned %s#id but no %s", v.body, resBody)
			}
		})
	}
}

func TestParseShortUrl(t *testing.T) {

	shortUrl := func() string {
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewBuffer([]byte("https://github.com/alaleks/shortener")))
		responseRecorder := httptest.NewRecorder()
		handler.ServeHTTP(responseRecorder, request)
		res := responseRecorder.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		return string(resBody)
	}()

	tests := []struct {
		name     string
		code     int
		url      string
		shorturl string
		body     string
	}{
		{
			name:     "Тест с корректным uid",
			code:     307,
			shorturl: shortUrl,
			url:      "https://github.com/alaleks/shortener",
			body:     "",
		},
		{
			name:     "Тест с некорректным uid",
			code:     400,
			shorturl: "http://localhost:8080/1",
			url:      "https://github.com/alaleks/shortener",
			body:     "this short url is invalid",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, v.shorturl, nil)
			resRec := httptest.NewRecorder()
			handler.ServeHTTP(resRec, req)
			res := resRec.Result()
			if res != nil {
				defer res.Body.Close()
			}
			loc, _ := res.Location()
			respBody, _ := io.ReadAll(res.Body)
			if v.code != res.StatusCode {
				t.Errorf("status code must be %d but is %d", v.code, res.StatusCode)
			}
			if loc != nil {
				if v.url != loc.String() {
					t.Errorf("location must be %s but is ", v.url)
				}
			}

			if v.body != strings.TrimSpace(string(respBody)) {
				t.Errorf("body must be %s but is %s", v.body, strings.TrimSpace(string(respBody)))
			}
		})
	}
}
