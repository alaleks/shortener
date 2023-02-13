package service_test

import (
	"errors"
	"testing"

	"github.com/alaleks/shortener/internal/app/service"
)

func TestIsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		name string
		uri  string
	}{
		{name: "url c https", uri: "https://github.com/alaleks/shortener", err: nil},
		{name: "url c http", uri: "http://github.com/alaleks/shortener", err: nil},
		{name: "url c www", uri: "www.github.com/alaleks/shortener", err: nil},
		{name: "url без протокола", uri: "github.com/alaleks/shortener", err: service.ErrInvalidURL},
		{name: "url c ошибкой в протоколе", uri: "htts://github.com/alaleks/shortener", err: service.ErrInvalidURL},
		{name: "url c протолом и пустым адресом", uri: "https://", err: service.ErrInvalidURL},
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

func BenchmarkIsURL(b *testing.B) {
	uri := "https://ya.ru/"
	uri2 := "http://ya.ru/"
	uri3 := "www.ya.ru/"
	uri4 := "htts://ya.ru"

	b.Run("before optimize https", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURLOld(uri)
		}
	})

	b.Run("after optimize https", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURL(uri)
		}
	})

	b.Run("before optimize http", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURLOld(uri2)
		}
	})

	b.Run("after optimize http", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURL(uri2)
		}
	})

	b.Run("before optimize www", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURLOld(uri3)
		}
	})

	b.Run("after optimize www", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURL(uri3)
		}
	})

	b.Run("before optimize wrong url", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURLOld(uri4)
		}
	})

	b.Run("after optimize www wrong url", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = service.IsURL(uri4)
		}
	})
}
