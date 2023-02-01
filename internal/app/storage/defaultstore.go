package storage

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/service"
)

type DefaultStorage struct {
	conf  config.Configurator
	urls  map[string]*URLElement // where key uid short url
	users map[uint][]string      // where key uid user, value a UID of short URL
	mu    sync.RWMutex
}

type URLElement struct {
	CreatedAt     time.Time
	LongURL       string
	CorrelationID string
	Statistics    uint // short URL usage statistics (actually this is the number of redirects)
	Removed       bool
}

func NewDefault(conf config.Configurator) *DefaultStorage {
	return &DefaultStorage{
		urls:  make(map[string]*URLElement),
		users: make(map[uint][]string),
		mu:    sync.RWMutex{},
		conf:  conf,
	}
}

func (ds *DefaultStorage) Add(longURL, userID string) (string, error) {
	if !strings.HasPrefix(longURL, "http") {
		longURL = "http://" + longURL
	}

	// генерируем id
	uid := service.GenUID(ds.conf.GetSizeUID())

	element := &URLElement{
		LongURL:    longURL,
		CreatedAt:  time.Now(),
		Statistics: 0,
	}

	uidToInt, err := strconv.Atoi(userID)

	ds.mu.Lock()
	ds.urls[uid] = element

	if err == nil {
		ds.users[uint(uidToInt)] = append(ds.users[uint(uidToInt)], uid)
	}

	ds.mu.Unlock()

	return ds.conf.GetBaseURL() + uid, nil
}

func (ds *DefaultStorage) AddBatch(longURL, userID, corID string) string {
	if !strings.HasPrefix(longURL, "http") {
		longURL = "http://" + longURL
	}

	// генерируем id
	uid := service.GenUID(ds.conf.GetSizeUID())

	element := &URLElement{
		LongURL:       longURL,
		CreatedAt:     time.Now(),
		Statistics:    0,
		CorrelationID: corID,
	}

	uidToInt, err := strconv.Atoi(userID)

	ds.mu.Lock()
	ds.urls[uid] = element

	if err == nil {
		ds.users[uint(uidToInt)] = append(ds.users[uint(uidToInt)], uid)
	}

	ds.mu.Unlock()

	return ds.conf.GetBaseURL() + uid
}

func (ds *DefaultStorage) GetURL(uid string) (string, error) {
	ds.mu.RLock()
	uri, check := ds.urls[uid]
	ds.mu.RUnlock()

	if !check {
		return "", ErrUIDNotValid
	}

	if uri.Removed {
		return uri.LongURL, ErrShortURLRemoved
	}

	return uri.LongURL, nil
}

func (ds *DefaultStorage) Update(uid string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	element, check := ds.urls[uid]

	if check {
		element.Statistics++
	}
}

func (ds *DefaultStorage) Stat(uid string) (Statistics, error) {
	ds.mu.RLock()
	uri, check := ds.urls[uid]
	ds.mu.RUnlock()

	if !check {
		return Statistics{}, ErrUIDNotValid
	}

	stat := Statistics{
		ShortURL:  ds.conf.GetBaseURL() + uid,
		LongURL:   uri.LongURL,
		CreatedAt: uri.CreatedAt.Format("02.01.2006 15:04:05"),
		Usage:     uri.Statistics,
	}

	return stat, nil
}

func (ds *DefaultStorage) DelUrls(userID string, shortsUID ...string) error {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return ErrUserIDNotValid
	}

	// key -> shortID, value -> found in existing.
	shortIDs := make(map[string]bool)

	for _, shortID := range shortsUID {
		shortIDs[shortID] = false
	}

	ds.mu.RLock()
	userShortsUID, ok := ds.users[uint(uid)]

	if !ok {
		ds.mu.RUnlock()

		return ErrInvalidData
	}

	for _, userSID := range userShortsUID {
		if _, ok := shortIDs[userSID]; ok {
			shortIDs[userSID] = true
		}
	}

	ds.mu.RUnlock()

	ds.mu.Lock()
	for shortID := range shortIDs {
		if uri, ok := ds.urls[shortID]; ok && !uri.Removed {
			ds.urls[shortID].Removed = true
		}
	}
	ds.mu.Unlock()

	return nil
}
