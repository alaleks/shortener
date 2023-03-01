package storage_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/storage"
)

// db
func BenchmarkUpdate(b *testing.B) {
	b.Setenv("DATABASE_DSN", "host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432")
	conf := config.New(config.Options{Env: true})
	db := storage.NewDB(conf)
	err := db.Init()
	if err != nil {
		return
	}

	var (
		shortURL = "http://example.com/1"
		userID   = "1"
	)

	short, _ := db.Add(shortURL, userID)
	shortUID := strings.Split(short, "/")[3]

	b.ResetTimer()

	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.UpdateOld(shortUID)
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.Update(shortUID)
		}
	})

	b.StopTimer()

	_ = db.Delete(strings.Split(short, "/")[3])
}

func BenchmarkAdd(b *testing.B) {
	b.Setenv("DATABASE_DSN", "host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432")
	conf := config.New(config.Options{Env: true})
	db := storage.NewDB(conf)
	err := db.Init()
	if err != nil {
		return
	}

	var (
		urls   []string
		userID = "1"
	)

	for i := 0; i < 2000; i++ {
		urls = append(urls, "http://example.com/"+strconv.Itoa(i))
	}

	b.ResetTimer()

	b.Run("before optimize", func(b *testing.B) {
		for _, v := range urls[0:1000] {
			_, _ = db.AddOld(v, userID)
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for _, v := range urls[1000:] {
			_, _ = db.Add(v, userID)
		}
	})

	b.StopTimer()

	for _, v := range urls {
		_ = db.DeleteByLongURL(v)
	}
}

func BenchmarkGetUrlsUser(b *testing.B) {
	b.Setenv("DATABASE_DSN", "host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432")
	conf := config.New(config.Options{Env: true})
	db := storage.NewDB(conf)
	err := db.Init()
	if err != nil {
		return
	}

	var uids []string

	for i := 0; i < 1000; i++ {
		shortURL, err := db.Add("http://example.com/"+strconv.Itoa(i),
			strconv.Itoa(1))
		if err == nil {
			uids = append(uids, shortURL)
		}
	}

	b.ResetTimer()

	b.Run("before optimize", func(b *testing.B) {
		_, _ = db.GetUrlsUserOld(strconv.Itoa(1))
	})

	b.Run("after optimize", func(b *testing.B) {
		_, _ = db.GetUrlsUser(strconv.Itoa(1))
	})

	b.StopTimer()

	for _, v := range uids {
		_ = db.Delete(strings.Split(v, "/")[3])
	}
}
