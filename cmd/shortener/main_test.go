package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alaleks/shortener/internal/app/router"
)

var handler http.Handler

func init() {
	handler = router.Create()
}

func TestShortenURL(t *testing.T) {

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
			endpoint: "http://localhost:8080/#uids",
		},
		{
			name:     "Тест url с http",
			code:     201,
			url:      "http://github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/#uids",
		},
		{
			name:     "Тест url только с www",
			code:     201,
			url:      "www.github.com/alaleks/shortener",
			endpoint: "http://localhost:8080/#uids",
		},
		{
			name:     "Тест url без http..",
			code:     400,
			url:      "github.com/alaleks/shortener",
			endpoint: "",
		},
		{
			name:     "Тест c ошибкой в http",
			code:     400,
			url:      "htts://github.com/alaleks/shortener",
			endpoint: "",
		},
		{
			name:     "Тест с пустой строкой после https:// ",
			code:     400,
			url:      "https://",
			endpoint: "",
		},
		{
			name:     "Тест c пустым body",
			code:     400,
			url:      "",
			endpoint: "",
		},
	}

	for _, v := range tests {
		v := v
		t.Run(v.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewBuffer([]byte(v.url)))
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
			if res.StatusCode == 201 {
				if len(v.endpoint) != len(string(resBody)) {
					t.Errorf("invalid short url %s: ", resBody)
				}
			}
		})
	}
}

func TestParseshortURL(t *testing.T) {
	shortURL := func() string {
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
		shortURL string
		body     string
	}{
		{
			name:     "Тест с корректным uid",
			code:     307,
			shortURL: shortURL,
			url:      "https://github.com/alaleks/shortener",
			body:     "",
		},
		{
			name:     "Тест с некорректным uid",
			code:     400,
			shortURL: "http://localhost:8080/1",
			url:      "https://github.com/alaleks/shortener",
			body:     "this short url is invalid",
		},
	}

	for _, v := range tests {
		v := v
		t.Run(v.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, v.shortURL, nil)
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

func TestGetStat(t *testing.T) {
	shortURL := func() string {
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewBuffer([]byte("https://github.com/alaleks/shortener")))
		responseRecorder := httptest.NewRecorder()
		handler.ServeHTTP(responseRecorder, request)
		res := responseRecorder.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		return string(resBody)
	}()

	tests := []struct {
		name string
		code int
		uri  string
	}{
		{
			name: "Тест с корректным uid",
			code: 200,
			uri:  shortURL + "/statistic",
		},
		{
			name: "Тест с некорректным uid",
			code: 400,
			uri:  shortURL + "1/statistic",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, v.uri, nil)
			resRec := httptest.NewRecorder()
			handler.ServeHTTP(resRec, req)
			res := resRec.Result()
			if res != nil {
				defer res.Body.Close()
			}
			if v.code != res.StatusCode {
				t.Errorf("status code must be %d but is %d", v.code, res.StatusCode)
			}

		})
	}
}
