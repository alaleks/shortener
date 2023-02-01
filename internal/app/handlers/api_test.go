package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/router"
	"github.com/alaleks/shortener/internal/app/serv/middleware"
	"github.com/alaleks/shortener/internal/app/serv/middleware/auth"
	"github.com/alaleks/shortener/internal/app/serv/middleware/compress"
	"github.com/alaleks/shortener/internal/app/storage"
)

const data = `{"url":"https://github.com/alaleks/shortener"}`

func TestShortenURLAPI(t *testing.T) {
	t.Parallel()
	// данные для теста
	appConf := config.New(config.Options{Env: false, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)

	tests := []struct {
		name    string
		data    string
		success bool
	}{
		{name: "url с https", data: data, success: true},
		{name: "url с http", data: `{"url":"http://github.com/alaleks/shortener"}`, success: true},
		{name: "url с www без протокола", data: `{"url":"www.github.com/alaleks/shortener"}`, success: true},
		{name: "url без протокола", data: `{"url":"github.com/alaleks/shortener"}`, success: false},
		{name: "url с ошибкой в протоколе", data: `{"url":"htps://github.com/alaleks/shortener"}`, success: false},
		{name: "невалидный json: имя", data: `{"url1":"htps://github.com/alaleks/shortener"}`, success: false},
		{name: "невалидный json: значение", data: `{"url":false}`, success: false},
	}

	// тестируем
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// создаем запрос, рекордер, хэндлер, запускаем сервер
			w := httptest.NewRecorder()
			h := http.HandlerFunc(testHandler.ShortenURLAPI)
			req := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL(), bytes.NewBuffer([]byte(item.data)))
			h.ServeHTTP(w, req)
			res := w.Result()

			if res != nil {
				defer res.Body.Close()
			}
			resBody, _ := io.ReadAll(res.Body)

			var dataFromRes struct {
				Success bool
			}

			_ = json.Unmarshal(resBody, &dataFromRes)

			// проверка значения success
			if dataFromRes.Success != item.success {
				t.Errorf("api should be return %v but received %v", item.success, dataFromRes.Success)
			}
		})
	}
}

func TestGetStatAPI(t *testing.T) {
	t.Parallel()

	// данные для теста
	appConf := config.New(config.Options{Env: false, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)
	// генерируем uid
	longURL1 := "https://github.com/alaleks/shortener"
	longURL2 := "https://yandex.ru/pogoda/krasnodar"
	// добавляем длинные ссылки в хранилище
	shortURL1, _ := testHandler.Storage.Store.Add(longURL1, "")
	uid1 := strings.SplitAfterN(shortURL1, "/", -1)[3]

	shortURL2, _ := testHandler.Storage.Store.Add(longURL2, "")
	uid2 := strings.SplitAfterN(shortURL2, "/", -1)[3]

	hostStat := appConf.GetBaseURL() + "api/"
	// для uid1 изменяем статистику
	testHandler.Storage.Store.Update(uid1)
	// создаем роутеры
	routers := router.Create(testHandler)

	tests := []struct {
		name    string
		uriStat string
		code    int
		stat    uint
	}{
		{name: "стат uid #1", code: 200, stat: 1, uriStat: hostStat + uid1 + "/statistics"},
		{name: "стат uid #2", code: 200, stat: 0, uriStat: hostStat + uid2 + "/statistics"},
		{name: "стат некорректной короткой ссылки", code: 400, uriStat: hostStat + "badId/statistics"},
	}

	// тестируем
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// создаем запрос, рекордер, хэндлер, запускаем сервер
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, item.uriStat, nil)
			routers.ServeHTTP(w, req)
			res := w.Result()

			if res != nil {
				defer res.Body.Close()
			}

			// проверка возвращаемого кода
			if res.StatusCode != item.code {
				t.Errorf("status code should be %d but received %d", item.code, res.StatusCode)
			}
			resBody, err := io.ReadAll(res.Body)
			if res.StatusCode == 200 && err == nil {
				var stat storage.Statistics
				if json.Unmarshal(resBody, &stat) == nil {
					if stat.Usage != item.stat {
						t.Errorf("mismatch statistics: should be %d but received %d", item.stat, stat.Usage)
					}
				}
			}
		})
	}
}

func TestSetEnv(t *testing.T) {
	// устанавливаем переменные окружения
	t.Setenv("SERVER_ADDRESS", "localhost:9090")
	t.Setenv("BASE_URL", "http://example.ru/")
	t.Setenv("FILE_STORAGE_PATH", "./storage")
	t.Setenv("DATABASE_DSN", "")

	// настройки для теста
	appConf := config.New(config.Options{Env: true, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)

	// создаем запрос, рекордер, хэндлер, запускаем сервер
	testRec := httptest.NewRecorder()
	h := http.HandlerFunc(testHandler.ShortenURLAPI)
	req := httptest.NewRequest(http.MethodPost, appConf.GetServAddr(), bytes.NewBuffer([]byte(data)))
	h.ServeHTTP(testRec, req)
	res := testRec.Result()

	if res != nil {
		res.Body.Close()
	}

	resBody, _ := io.ReadAll(res.Body)

	var dataFromRes struct {
		Result string `json:"result"`
	}

	_ = json.Unmarshal(resBody, &dataFromRes)

	if len(appConf.GetFileStoragePath()) != 0 {
		_ = testHandler.Storage.Store.Close()
	}

	if req.URL.String() != "localhost:9090" {
		t.Errorf("host should be localhost:9090 but no %s", req.URL.String())
	}

	if !strings.HasPrefix(dataFromRes.Result, "http://example.ru/") {
		t.Errorf("short url should be contains http://example.ru/ but no %s", req.URL.String())
	}

	if _, err := os.Stat(appConf.GetFileStoragePath()); err != nil {
		t.Errorf("failed to create file storage %s", err.Error())
	}

	// сбрасываем карту и читаем файл
	testHandler.Storage = storage.InitStore(appConf, logger)
	_ = testHandler.Storage.Store.Close()

	if _, err := testHandler.Storage.Store.GetURL(strings.Split(dataFromRes.Result, "/")[3]); err != nil {
		t.Errorf("failed to get data from file storage: %s", dataFromRes.Result)
	}

	// удаляем созданное файловое хранилище
	_ = os.Remove(appConf.GetFileStoragePath())
}

func TestSetFlag(t *testing.T) {
	t.Parallel()

	// настройки для теста
	options := config.Options{Env: true, Flag: true}
	appConf := config.New(options)
	argsTest := []string{
		"TestFlags", "-a", "localhost:9093", "-b",
		"http://localhost:9093/", "-f", "./storage",
	}
	appConf.DefineOptionsFlags(argsTest)

	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)

	// создаем запрос, рекордер, хэндлер, запускаем сервер
	testRec := httptest.NewRecorder()
	h := http.HandlerFunc(testHandler.ShortenURLAPI)
	req := httptest.NewRequest(http.MethodPost, appConf.GetServAddr(), bytes.NewBuffer([]byte(data)))
	h.ServeHTTP(testRec, req)
	res := testRec.Result()

	if res != nil {
		res.Body.Close()
	}

	resBody, _ := io.ReadAll(res.Body)

	var dataFromRes struct {
		Result string `json:"result"`
	}

	_ = json.Unmarshal(resBody, &dataFromRes)

	if len(appConf.GetFileStoragePath()) != 0 {
		_ = testHandler.Storage.Store.Close()
	}

	if req.URL.String() != "localhost:9093" {
		t.Errorf("host should be localhost:9093 but no %s", req.URL.String())
	}

	if !strings.HasPrefix(dataFromRes.Result, "http://localhost:9093/") {
		t.Errorf("short url should be contains http://localhost:9093/ but no %s", req.URL.String())
	}

	if _, err := os.Stat(appConf.GetFileStoragePath()); err != nil {
		t.Errorf("failed to create file storage %s", err.Error())
	}

	// сбрасываеи карту и читаем файл
	testHandler.Storage = storage.InitStore(appConf, logger)
	_ = testHandler.Storage.Store.Close()

	if _, err := testHandler.Storage.Store.GetURL(strings.Split(dataFromRes.Result, "/")[3]); err != nil {
		t.Errorf("failed to get data from file storage: %s", dataFromRes.Result)
	}

	// удаляем созданное файловое хранилище
	_ = os.Remove(appConf.GetFileStoragePath())
}

func TestCompress(t *testing.T) {
	t.Parallel()
	// данные для теста
	appConf := config.New(config.Options{Env: false, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)

	tests := []struct {
		name            string
		acceptEncoding  string
		contentEncoding string
	}{
		{name: "проверка сжатия: gzip support", acceptEncoding: "gzip, deflate, br", contentEncoding: "gzip"},
		{name: "проверка сжатия: gzip noupport", acceptEncoding: "", contentEncoding: ""},
	}

	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// создаем запрос, рекордер, хэндлер, запускаем сервер
			testRec := httptest.NewRecorder()
			h := middleware.New(compress.Compression, compress.Unpacking).
				Configure(http.HandlerFunc(testHandler.ShortenURLAPI))
			req := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL(), bytes.NewBuffer([]byte(data)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", item.acceptEncoding)
			h.ServeHTTP(testRec, req)
			res := testRec.Result()
			if res != nil {
				defer res.Body.Close()
			}

			if item.contentEncoding != res.Header.Get("Content-Encoding") {
				t.Errorf("content-Encoding should be '%s' but received '%s'",
					item.contentEncoding, res.Header.Get("Content-Encoding"))
			}
		})
	}
}

func TestGetUsersURL(t *testing.T) {
	t.Parallel()
	// данные для теста
	appConf := config.New(config.Options{Env: false, Flag: false})
	logger := logger.NewLogger()
	testHandler := handlers.New(appConf, logger)
	auth := auth.TurnOn(testHandler.Storage, appConf.GetSecretKey())
	tests := []struct {
		name string
		url  string
		code int
	}{
		{
			name: "проверка когда кука установлена и валидна", code: 200,
			url: "http://github.com/alaleks/shortener",
		},
		{
			name: "проверка, когда кука пуста", code: 204,
			url: "",
		},
		{
			name: "проверка, когда куку пытались поменять", code: 204,
			url: "wrong",
		},
	}

	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			testRec := httptest.NewRecorder()
			handler := middleware.New(auth.Authorization).
				Configure(http.HandlerFunc(testHandler.GetUsersURL))
			req := httptest.NewRequest(http.MethodGet, appConf.GetBaseURL(), nil)

			// тестируем сценарий добавления куки пр сокращении URL
			if item.url != "" {
				testRec2 := httptest.NewRecorder()
				handler2 := middleware.New(auth.Authorization).
					Configure(http.HandlerFunc(testHandler.ShortenURL))
				req2 := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL(), bytes.NewBuffer([]byte(item.url)))
				handler2.ServeHTTP(testRec2, req2)
				res2 := testRec2.Result()
				if res2.Body != nil {
					defer res2.Body.Close()
				}
				http.SetCookie(testRec, res2.Cookies()[0])
				req.Header = http.Header{"Cookie": testRec2.Header()["Set-Cookie"]}
			}

			// здесь меняем куку авторизации
			if item.url == "wrong" {
				testRec2 := httptest.NewRecorder()
				handler2 := middleware.New(auth.Authorization).
					Configure(http.HandlerFunc(testHandler.ShortenURL))
				req2 := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL(), bytes.NewBuffer([]byte(item.url)))
				handler2.ServeHTTP(testRec2, req2)
				res2 := testRec2.Result()
				if res2.Body != nil {
					defer res2.Body.Close()
				}
				cookie := res2.Cookies()[0]
				cookie.Value += "wrong"
				http.SetCookie(testRec, cookie)
				req.Header = http.Header{"Cookie": testRec2.Header()["Set-Cookie"]}
			}

			handler.ServeHTTP(testRec, req)
			res := testRec.Result()
			if res.Body != nil {
				defer res.Body.Close()
			}
			if res.StatusCode != item.code {
				t.Errorf("status code should be %d but received %d", item.code, res.StatusCode)
			}
		})
	}
}
