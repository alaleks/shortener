package storage

import (
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/service"
)

type Storager interface {
	Add(longURL string) (uid string)
	GetURL(uid string) (string, bool)
	Stat(uid string) (string, uint, string)
	Update(uid string) bool
}

type urlEl struct {
	longURL   string
	created   time.Time
	statistic uint // short URL usage statistics (actually this is the number of redirects)
}

type Urls struct {
	data map[string]*urlEl // where key uid short url
	mu   sync.RWMutex
}

func (u *Urls) Add(longURL string) (uid string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !strings.HasPrefix(longURL, "http") {
		longURL = "http://" + longURL
	}
	uid = service.GenUID(5)
	u.data[uid] = &urlEl{longURL, time.Now(), 0}

	return uid
}

func (u *Urls) GetURL(uid string) (string, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	uri, ok := u.data[uid]
	if ok {
		return uri.longURL, ok
	}
	return "", ok
}

func (u *Urls) Update(uid string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	uri, ok := u.data[uid]
	if ok {
		uri.statistic++
		u.data[uid] = uri
	}
	return ok
}

func (u *Urls) Stat(uid string) (string, uint, string) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	uri, ok := u.data[uid]
	if !ok {
		return "", 0, ""
	}
	return uri.longURL, uri.statistic, uri.created.Format("02.01.2006 15:04:05")
}

func New() Storager {
	return &Urls{
		data: make(map[string]*urlEl),
		mu:   sync.RWMutex{},
	}
}
