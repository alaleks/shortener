package storage

import (
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/service"
)

type Storager interface {
	Add(longUrl string) (uid string)
	GetURL(uid string) (string, bool)
	Stat(uid string) (string, uint, string)
	Update(uid string) bool
}

type uri struct {
	longUrl   string
	created   time.Time
	statistic uint // short URL usage statistics (actually this is the number of redirects)
}

type Urls struct {
	data map[string]*uri // where key uid short url
	mtx  *sync.Mutex
}

func (u *Urls) Add(longUrl string) (uid string) {
	if !strings.HasPrefix(longUrl, "http") {
		longUrl = "http://" + longUrl
	}

	uid = service.GenUid(5)
	u.data[uid] = &uri{longUrl, time.Now(), 0}

	return uid
}

func (u *Urls) GetURL(uid string) (string, bool) {
	uri, ok := u.data[uid]
	if ok {
		return uri.longUrl, ok
	}
	return "", ok
}

func (u *Urls) Update(uid string) bool {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	uri, ok := u.data[uid]
	if ok {
		uri.statistic++
		u.data[uid] = uri
	}
	return ok
}

func (u *Urls) Stat(uid string) (string, uint, string) {
	uri, ok := u.data[uid]
	if !ok {
		return "", 0, ""
	}
	return uri.longUrl, uri.statistic, uri.created.Format("02.01.2006 15:04:05")
}

var DataStorage Storager

func init() {
	DataStorage = &Urls{
		data: make(map[string]*uri),
		mtx:  &sync.Mutex{},
	}
}
