package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/router"
)

const host = "http://localhost:8080/"

func TestShortenURL(t *testing.T) {
	t.Parallel()

	// данные для теста
	testHandler := handlers.New(5)
	templateShortURL := "http://localhost:8080/#uids"

	tests := []struct {
		name     string
		code     int
		url      string
		endpoint string
		body     string
	}{
		{name: "url с https", code: 201, url: "https://github.com/alaleks/shortener"},
		{name: "url с http", code: 201, url: "http://github.com/alaleks/shortener"},
		{name: "url только с www без протокола", code: 201, url: "www.github.com/alaleks/shortener"},
		{name: "url без протокола", code: 400, url: "github.com/alaleks/shortener"},
		{name: "url с ошибкой в протоколе", code: 400, url: "htts://github.com/alaleks/shortener"},
		{name: "url с пустым адресом после протокола", code: 400, url: "https://"},
		{name: "пустой body", code: 400, url: ""},
	}

	// тестируем
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// создаем запрос, рекордер, хэндлер, запускаем сервер
			w := httptest.NewRecorder()
			h := http.HandlerFunc(testHandler.ShortenURL)
			req := httptest.NewRequest(http.MethodPost, host, bytes.NewBuffer([]byte(item.url)))
			h.ServeHTTP(w, req)
			res := w.Result()

			if res != nil {
				defer res.Body.Close()
			}
			resBody, _ := io.ReadAll(res.Body)

			// проверка возвращаемого кода ответа
			if res.StatusCode != item.code {
				t.Errorf("status code should be %d but received %d", item.code, res.StatusCode)
			}
			// проверка возвращаемых коротких ссылок на соответствие шаблону
			if res.StatusCode == 201 && len(templateShortURL) != len(string(resBody)) {
				t.Errorf("short url %s does not match pattern", string(resBody))
			}
		})
	}
}

func TestParseShortURL(t *testing.T) {
	t.Parallel()

	// данные для теста
	testHandler := handlers.New(5)
	// генерируем uid
	longURL := "https://github.com/alaleks/shortener"
	// добавляем короткую ссылку
	uid := testHandler.DataStorage.Add(longURL, testHandler.SizeUID)
	// создаем роутеры
	routers := router.Create(testHandler)

	tests := []struct {
		name     string
		code     int
		shortURL string
		longURL  string
	}{
		{name: "парсинг корректной короткой ссылки", code: 307, longURL: longURL, shortURL: host + uid},
		{name: "парсинг некорректной короткой ссылки - 1", code: 405, longURL: "", shortURL: host},
		{name: "парсинг некорректной короткой ссылки - 2", code: 400, longURL: "", shortURL: host + "badId"},
	}

	// тестируем
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// создаем запрос, рекордер, хэндлер, запускаем сервер
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, item.shortURL, nil)
			routers.ServeHTTP(w, req)
			res := w.Result()

			if res != nil {
				defer res.Body.Close()
			}

			// проверка возвращаемого кода ответа
			if res.StatusCode != item.code {
				t.Errorf("status code should be %d but received %d", item.code, res.StatusCode)
			}

			// проверка location при удачном сценарии
			resLoc, err := res.Location()
			if err == nil {
				if resLoc.String() != item.longURL {
					t.Errorf("location should be %s but received %s", item.longURL, resLoc.String())
				}
			}
		})
	}
}
