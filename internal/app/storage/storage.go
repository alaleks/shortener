package storage

import (
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/service"
)

type Storage interface {
	Procucer
	Consumer
	FileStorage
}

type Procucer interface {
	Add(longURL string, sizeUID int) string
	Update(uid string) bool
}

type Consumer interface {
	GetURL(uid string) (string, bool)
	Stat(uid string) (string, uint, string)
}

type URLElement struct {
	LongURL    string
	CreatedAt  time.Time
	Statistics uint // short URL usage statistics (actually this is the number of redirects)
}

type Urls struct {
	data map[string]*URLElement // where key uid short url
	mu   sync.RWMutex
}

func New() *Urls {
	urls := Urls{
		data: make(map[string]*URLElement),
		mu:   sync.RWMutex{},
	}
	return &urls
}

func (u *Urls) Add(longURL string, sizeUID int) string {
	if !strings.HasPrefix(longURL, "http") {
		longURL = "http://" + longURL
	}

	// генерируем id
	uid := service.GenUID(sizeUID)

	element := &URLElement{
		LongURL:    longURL,
		CreatedAt:  time.Now(),
		Statistics: 0,
	}

	u.mu.Lock()
	u.data[uid] = element
	u.mu.Unlock()

	return uid
}

func (u *Urls) GetURL(uid string) (string, bool) {
	u.mu.RLock()
	uri, check := u.data[uid]
	u.mu.RUnlock()

	if check {
		return uri.LongURL, check
	}

	return "", check
}

func (u *Urls) Update(uid string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	element, check := u.data[uid]

	if check {
		element.Statistics++
	}

	return check
}

func (u *Urls) Stat(uid string) (string, uint, string) {
	u.mu.RLock()
	uri, check := u.data[uid]
	u.mu.RUnlock()

	if !check {
		return "", 0, ""
	}

	return uri.LongURL, uri.Statistics, uri.CreatedAt.Format("02.01.2006 15:04:05")
}
