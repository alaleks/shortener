package service_test

import (
	"errors"
	"testing"

	"github.com/alaleks/shortener/internal/app/service"
)

func TestIsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		uri  string
		err  error
	}{
		{"url c https", "https://github.com/alaleks/shortener", nil},
		{"url c http", "http://github.com/alaleks/shortener", nil},
		{"url c www", "www.github.com/alaleks/shortener", nil},
		{"url без протокола", "github.com/alaleks/shortener", service.ErrInvalidURL},
		{"url c ошибкой в протоколе", "htts://github.com/alaleks/shortener", service.ErrInvalidURL},
		{"url c протолом и пустым адресом", "https://", service.ErrInvalidURL},
	}
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			err := service.IsURL(item.uri)
			if !errors.Is(err, item.err) {
				t.Errorf("checking this url %s should be return %v but no %v", item.uri, item.err, err)
			}
		})
	}
}

func TestCreateShortId(t *testing.T) {
	t.Parallel()

	size := 5
	tests := []struct {
		name string
		uid  string
	}{
		{name: "uid #1", uid: service.GenUID(size)},
		{name: "uid #2", uid: service.GenUID(size)},
	}

	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			// проверяем, что был сгенерирован id нужного размера len
			if len(item.uid) != size {
				t.Errorf("uid should be сonsist %d characters", size)
			}
			// проверяем корректность работы рандомайзера, должны получиться разные значения
			// если равны, то это ошибка
			if service.GenUID(size) == item.uid {
				t.Errorf("uids should not be egual each other")
			}
		})
	}
}
