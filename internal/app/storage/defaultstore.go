package storage

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/service"
)

// Data Structures
type (
	DefaultStorage struct {
		conf  config.Configurator
		urls  map[string]*URLElement // where key uid short url
		users map[uint][]string      // where key uid user, value a UID of short URL
		mu    sync.RWMutex
	}

	URLElement struct {
		CreatedAt     time.Time
		LongURL       string
		CorrelationID string
		Statistics    uint // short URL usage statistics (actually this is the number of redirects)
		Removed       bool
	}
)

// NewDefault creates a pointer of DefaultStorage.
func NewDefault(conf config.Configurator) *DefaultStorage {
	return &DefaultStorage{
		urls:  make(map[string]*URLElement),
		users: make(map[uint][]string),
		mu:    sync.RWMutex{},
		conf:  conf,
	}
}

// Add performs adding URL to the default storage.
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

// AddBatch performs adding data to storage when batch processing.
func (ds *DefaultStorage) AddBatch(longURL, userID, corID string) string {
	if !strings.HasPrefix(longURL, "http") {
		longURL = "http://" + longURL
	}

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

// GetURL returns the original url by its short UID.
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

// Update perfoms changing short link usage statistics.
func (ds *DefaultStorage) Update(uid string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	element, check := ds.urls[uid]

	if check {
		element.Statistics++
	}
}

// Stat returns short link statistics by its id.
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

// DelUrls marks as deleted urls added by a specific user.
func (ds *DefaultStorage) DelUrls(userID string, shortsUID ...string) error {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return ErrUserIDNotValid
	}

	// key -> shortID, value -> found in existing.
	shortIDs := make(map[string]bool, len(shortsUID))

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

// DelUrlsOld marks as deleted urls added by a specific user.
func (ds *DefaultStorage) DelUrlsOld(userID string, shortsUID ...string) error {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return ErrUserIDNotValid
	}

	ds.mu.RLock()
	if _, ok := ds.users[uint(uid)]; !ok {
		ds.mu.RUnlock()

		return ErrUserNotExists
	}

	var uidsToDel []string

	for _, v := range shortsUID {
		for _, shortID := range ds.users[uint(uid)] {
			if shortID == v && !ds.urls[shortID].Removed {
				uidsToDel = append(uidsToDel, v)

				break
			}
		}
	}

	ds.mu.RUnlock()

	ds.mu.Lock()
	for _, v := range uidsToDel {
		if _, ok := ds.urls[v]; ok {
			ds.urls[v].Removed = true
		}
	}
	ds.mu.Unlock()

	return nil
}

// GetStatsInternal returns data about the number of shortened URLs
// and the number of users in the app.
func (ds *DefaultStorage) GetStatsInternal() (StatsInternal, error) {
	ds.mu.RLock()
	urlsLen := len(ds.urls)
	usersLen := len(ds.users)
	ds.mu.RUnlock()

	return StatsInternal{
		UrlsSize: urlsLen,
		Users:    usersLen,
	}, nil
}
