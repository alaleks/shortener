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
	"github.com/alaleks/shortener/internal/app/router"
	"github.com/alaleks/shortener/internal/app/serv/middleware"
	"github.com/alaleks/shortener/internal/app/storage"
)

const data = `{"url":"https://github.com/alaleks/shortener"}`

func TestShortenURLAPI(t *testing.T) {
	t.Parallel()
	// данные для теста
	appConf := config.New(config.Options{})
	appConf.DefineOptionsFlags([]string{"TestFlags", "-a", "", "-b", "", "-f", ""})
	testHandler := handlers.New(5, appConf)

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
	appConf := config.New(config.Options{})
	appConf.DefineOptionsFlags([]string{"TestFlags", "-a", "", "-b", "", "-f", ""})
	testHandler := handlers.New(5, appConf)
	// генерируем uid
	longURL1 := "https://github.com/alaleks/shortener"
	longURL2 := "https://yandex.ru/pogoda/krasnodar"
	// добавляем длинные ссылки в хранилище
	uid1 := testHandler.DataStorage.Add(longURL1, testHandler.SizeUID)
	uid2 := testHandler.DataStorage.Add(longURL2, testHandler.SizeUID)
	hostStat := appConf.GetBaseURL() + "api/"
	// для uid1 изменяем статистику
	testHandler.DataStorage.Update(uid1)
	// создаем роутеры
	routers := router.Create(testHandler)

	tests := []struct {
		name    string
		code    int
		stat    uint
		uriStat string
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
				var stat handlers.Statistics
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

	// настройки для теста
	options := config.Options{Env: true, Flag: false}
	appConf := config.New(options)
	appConf.DefineOptionsEnv()
	testHandler := handlers.New(5, appConf)

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
		_ = testHandler.DataStorage.Write(appConf.GetFileStoragePath())
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

	// сбрасываеи карту и читаем файл
	testHandler.DataStorage = storage.New()
	_ = testHandler.DataStorage.Read(appConf.GetFileStoragePath())

	if _, ok := testHandler.DataStorage.GetURL(strings.Split(dataFromRes.Result, "/")[3]); !ok {
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
	testHandler := handlers.New(5, appConf)

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
		_ = testHandler.DataStorage.Write(appConf.GetFileStoragePath())
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
	testHandler.DataStorage = storage.New()
	_ = testHandler.DataStorage.Read(appConf.GetFileStoragePath())

	if _, ok := testHandler.DataStorage.GetURL(strings.Split(dataFromRes.Result, "/")[3]); !ok {
		t.Errorf("failed to get data from file storage: %s", dataFromRes.Result)
	}

	// удаляем созданное файловое хранилище
	_ = os.Remove(appConf.GetFileStoragePath())
}

func TestCompress(t *testing.T) {
	t.Parallel()
	// данные для теста
	appConf := config.New(config.Options{})
	appConf.DefineOptionsFlags([]string{"TestFlags", "-a", "", "-b", "", "-f", ""})
	testHandler := handlers.New(5, appConf)

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
			h := middleware.New(middleware.Compress, middleware.DeCompress).
				Configure(http.HandlerFunc(testHandler.ShortenURLAPI))
			req := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL(), bytes.NewBuffer([]byte(data)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", item.acceptEncoding)
			h.ServeHTTP(testRec, req)
			res := testRec.Result()
			if res != nil {
				res.Body.Close()
			}

			if item.contentEncoding != res.Header.Get("Content-Encoding") {
				t.Errorf("content-Encoding should be '%s' but received '%s'",
					item.contentEncoding, res.Header.Get("Content-Encoding"))
			}
		})
	}
}
