package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/router"
)

func TestShortenURLAPI(t *testing.T) {
	t.Parallel()
	// данные для теста
	appConf := config.New()
	testHandler := handlers.New(5, appConf)

	tests := []struct {
		name    string
		data    string
		success bool
		code    int
	}{
		{name: "url с https", data: `{"url":"https://github.com/alaleks/shortener"}`, success: true, code: 201},
		{name: "url с http", data: `{"url":"http://github.com/alaleks/shortener"}`, success: true, code: 201},
		{name: "url с www без протокола", data: `{"url":"www.github.com/alaleks/shortener"}`, success: true, code: 201},
		{name: "url без протокола", data: `{"url":"github.com/alaleks/shortener"}`, success: false, code: 205},
		{name: "url с ошибкой в протоколе", data: `{"url":"htps://github.com/alaleks/shortener"}`, success: false, code: 205},
		{name: "невалидный json: имя", data: `{"url1":"htps://github.com/alaleks/shortener"}`, success: false, code: 205},
		{name: "невалидный json: значение", data: `{"url":false}`, success: false, code: 205},
	}

	// тестируем
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// создаем запрос, рекордер, хэндлер, запускаем сервер
			w := httptest.NewRecorder()
			h := http.HandlerFunc(testHandler.ShortenURLAPI)
			req := httptest.NewRequest(http.MethodPost, appConf.GetBaseURL().String(), bytes.NewBuffer([]byte(item.data)))
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

			// проверка возвращаемого кода ответа
			if res.StatusCode != item.code {
				t.Errorf("status code should be %d but received %d", item.code, res.StatusCode)
			}
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
	appConf := config.New()
	testHandler := handlers.New(5, appConf)
	// генерируем uid
	longURL1 := "https://github.com/alaleks/shortener"
	longURL2 := "https://yandex.ru/pogoda/krasnodar"
	// добавляем длинные ссылки в хранилище
	uid1 := testHandler.DataStorage.Add(longURL1, testHandler.SizeUID)
	uid2 := testHandler.DataStorage.Add(longURL2, testHandler.SizeUID)
	hostStat := appConf.GetBaseURL().String() + "api/"
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
