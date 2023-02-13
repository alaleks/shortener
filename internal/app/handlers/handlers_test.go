package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/router"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	// данные для теста
	appConf := config.New(config.Options{Env: false, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)

	tests := []struct {
		name     string
		url      string
		endpoint string
		body     string
		code     int
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
			req := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL(), bytes.NewBuffer([]byte(item.url)))
			h.ServeHTTP(w, req)
			res := w.Result()

			if res != nil {
				defer func() {
					_ = res.Body.Close()
				}()
			}
			resBody, _ := io.ReadAll(res.Body)
			templateShortURL := req.URL.Scheme + "://" + req.URL.Host + "/#uids"

			// проверка возвращаемого кода ответа
			if res.StatusCode != item.code {
				t.Errorf("status code should be %d but received %d", item.code, res.StatusCode)
			}
			// проверка возвращаемых коротких ссылок на соответствие шаблону
			if res.StatusCode == 201 && len(templateShortURL) != len(string(resBody)) {
				t.Errorf("short url %s does not match pattern %s", string(resBody), templateShortURL)
			}
		})
	}
}

func TestParseShortURL(t *testing.T) {
	t.Parallel()

	// данные для теста
	appConf := config.New(config.Options{Env: false, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)
	longURL := "https://github.com/alaleks/shortener"
	shortURL, _ := testHandler.Storage.Store.Add(longURL, "1")
	routers := router.Create(testHandler)

	tests := []struct {
		name     string
		shortURL string
		longURL  string
		code     int
	}{
		{
			name: "парсинг корректной короткой ссылки", code: 307,
			longURL: longURL, shortURL: shortURL,
		},
		{
			name: "парсинг некорректной короткой ссылки - 1", code: 405,
			longURL: "", shortURL: appConf.GetBaseURL(),
		},
		{
			name: "парсинг некорректной короткой ссылки - 2", code: 400,
			longURL: "", shortURL: appConf.GetBaseURL() + "badId",
		},
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
				defer func() {
					_ = res.Body.Close()
				}()
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
