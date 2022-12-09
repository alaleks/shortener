package storage

import (
	"errors"
	"strconv"
)

var (
	ErrUserIDNotValid = errors.New("invalid user id")
	ErrUserUrlsEmpty  = errors.New("shortened URLs for current user is empty")
)

type URLUser struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (ds *DefaultStorage) GetUrlsUser(userID string) ([]URLUser, error) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return []URLUser{}, ErrUserIDNotValid
	}

	ds.mu.RLock()
	defer ds.mu.RUnlock()

	uidsShortURL := ds.users[uint(uid)]
	urls := make([]URLUser, 0, len(uidsShortURL))

	if len(uidsShortURL) == 0 {
		return urls, ErrUserUrlsEmpty
	}

	for _, shortUID := range uidsShortURL {
		if originalURL, err := ds.GetURL(shortUID); err == nil {
			urls = append(urls, URLUser{ShortURL: ds.conf.GetBaseURL() + shortUID, OriginalURL: originalURL})
		}
	}

	if len(urls) == 0 {
		return urls, ErrUserUrlsEmpty
	}

	return urls, nil
}

func (ds *DefaultStorage) Create() uint {
	ds.mu.Lock()
	uid := uint(len(ds.users) + 1)
	ds.users[uid] = make([]string, 0)
	ds.mu.Unlock()

	return uid
}
