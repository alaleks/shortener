package storage

import (
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/service"
)

type Storager interface {
	Add(longURL string, sizeUID int) string
	GetURL(uid string) (string, bool)
	Stat(uid string) (string, uint, string)
	Update(uid string) bool
}

type URLElement struct {
	longURL    string
	createdAt  time.Time
	statistics uint // short URL usage statistics (actually this is the number of redirects)
}

type Urls struct {
	data map[string]*URLElement // where key uid short url
	mu   sync.RWMutex
}

func New() *Urls {
	return &Urls{
		data: make(map[string]*URLElement),
		mu:   sync.RWMutex{},
	}
}

func (u *Urls) Add(longURL string, sizeUID int) string {
	if !strings.HasPrefix(longURL, "http") {
		longURL = "http://" + longURL
	}

	// генерируем id
	uid := service.GenUID(sizeUID)

	element := &URLElement{
		longURL:    longURL,
		createdAt:  time.Now(),
		statistics: 0,
	}

	u.mu.Lock()
	u.data[uid] = element
	u.mu.Unlock()

	return uid
}

func (u *Urls) GetURL(uid string) (string, bool) {
	uri, check := u.data[uid]

	if check {
		return uri.longURL, check
	}

	return "", check
}

func (u *Urls) Update(uid string) bool {
	element, check := u.data[uid]

	if check {
		u.mu.Lock()
		defer u.mu.Unlock()
		element.statistics++
	}

	return check
}

func (u *Urls) Stat(uid string) (string, uint, string) {
	uri, check := u.data[uid]

	if !check {
		return "", 0, ""
	}

	return uri.longURL, uri.statistics, uri.createdAt.Format("02.01.2006 15:04:05")
}
