package storage_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/storage"
)

func InitDB() {
}

// db
func BenchmarkUpdate(b *testing.B) {
	b.StopTimer()
	b.Setenv("DATABASE_DSN", "host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432")
	conf := config.New(config.Options{Env: true}, 5)
	db := storage.NewDB(conf)
	err := db.Init()
	if err != nil {
		return
	}

	short, _ := db.Add("http://example.com/1", strconv.Itoa(1))
	b.StartTimer()

	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.UpdateOld(strings.Split(short, "/")[3])
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			db.Update(strings.Split(short, "/")[3])
		}
	})

	_ = db.Delete(strings.Split(short, "/")[3])
}

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	b.Setenv("DATABASE_DSN", "host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432")
	conf := config.New(config.Options{Env: true}, 5)
	db := storage.NewDB(conf)
	err := db.Init()
	if err != nil {
		return
	}

	var (
		count int
		uids  []string
	)
	b.StartTimer()

	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			count++
			shortURL, err := db.AddOld("http://example.com/"+strconv.Itoa(count),
				strconv.Itoa(1))
			if err == nil {
				uids = append(uids, shortURL)
			}
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortURL, err := db.Add("http://example.com/"+strconv.Itoa(count),
				strconv.Itoa(1))
			if err == nil {
				uids = append(uids, shortURL)
			}
		}
	})

	b.StopTimer()

	for _, v := range uids {
		_ = db.Delete(strings.Split(v, "/")[3])
	}
}

func BenchmarkGetUrlsUser(b *testing.B) {
	b.Setenv("DATABASE_DSN", "host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432")
	conf := config.New(config.Options{Env: true}, 5)
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

	b.Run("before optimize", func(b *testing.B) {
		_, _ = db.GetUrlsUserOld(strconv.Itoa(1))
	})

	b.Run("after optimize", func(b *testing.B) {
		_, _ = db.GetUrlsUser(strconv.Itoa(1))
	})

	for _, v := range uids {
		_ = db.Delete(strings.Split(v, "/")[3])
	}
}
