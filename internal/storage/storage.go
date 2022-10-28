package storage

import (
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/service"
)

type Storager interface {
	Add(uri string, uid string)
	Get(key string) (string, bool)
	Stat(key string) (string, uint, string)
	Update(key string)
}

type ShortUrl struct {
	uid       string
	created   time.Time
	statistic uint
}

type Urls struct {
	LongUrls  map[string]*ShortUrl // where key long url
	ShortUrls map[string]string    // where key is an uid, value is a url before used shortener
}

func (u *Urls) Add(uri string, uid string) {
	if !strings.HasPrefix(uri, "http") {
		uri = "http://" + uri
	}

	u.ShortUrls[uid] = uri
	u.LongUrls[service.RemovePrefix(uri, "https://", "http://", "www.")] = &ShortUrl{uid, time.Now(), 0}
}

func (u *Urls) Get(key string) (string, bool) {
	longUrl, ok := u.ShortUrls[key]

	if ok {
		return longUrl, ok
	}

	cleanUrl := service.RemovePrefix(key, "https://", "http://", "www.")
	shortUrl, ok := u.LongUrls[cleanUrl]

	if ok {
		return shortUrl.uid, ok
	}

	return "", false
}

func (u *Urls) Update(key string) {
	var mtx sync.Mutex
	shortUrl := u.LongUrls[service.RemovePrefix(u.ShortUrls[key], "https://", "http://", "www.")]

	mtx.Lock()
	defer mtx.Unlock()

	shortUrl.statistic++
	u.LongUrls[key] = shortUrl

}

func (u *Urls) Stat(key string) (string, uint, string) {
	cleanUrl := service.RemovePrefix(key, "https://", "http://", "www.")
	shortUrl, ok := u.LongUrls[cleanUrl]
	if !ok {
		return "", 0, ""
	}
	return shortUrl.uid, shortUrl.statistic, shortUrl.created.Format("02.01.2006 15:04:05")
}
