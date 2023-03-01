package storage_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/storage"
)

func BenchmarkDelUrls(b *testing.B) {
	conf := config.New(config.Options{})
	storeDefault := storage.NewDefault(conf)

	var (
		shortUIDs []string
		userID    = "1"
	)

	for i := 0; i < 2000; i++ {
		shortURL, _ := storeDefault.Add("http://example.com/"+strconv.Itoa(i), userID)
		shortUIDs = append(shortUIDs, strings.Split(shortURL, "/")[3])
	}

	b.ResetTimer()

	b.Run("before optimize", func(b *testing.B) {
		for _, v := range shortUIDs[0:1000] {
			_ = storeDefault.DelUrlsOld(userID, v)
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for _, v := range shortUIDs[1000:] {
			_ = storeDefault.DelUrls(userID, v)
		}
	})

}
