package storage

import (
	"strconv"
)

// GetUrlsUser - getting shorts URLs from default storage for current user.
func (ds *DefaultStorage) GetUrlsUser(userID string) ([]struct {
	ShortUID string `json:"short_url"`
	LongURL  string `json:"original_url"`
}, error) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return nil, ErrUserIDNotValid
	}

	ds.mu.RLock()
	uidsShortURL := ds.users[uint(uid)]

	urls := make([]struct {
		ShortUID string `json:"short_url"`
		LongURL  string `json:"original_url"`
	}, 0, len(uidsShortURL))

	if len(uidsShortURL) == 0 {
		ds.mu.RUnlock()

		return urls, ErrUserUrlsEmpty
	}

	for _, shortUID := range uidsShortURL {
		if originalURL, err := ds.GetURL(shortUID); err == nil {
			urls = append(urls, struct {
				ShortUID string `json:"short_url"`
				LongURL  string `json:"original_url"`
			}{ShortUID: ds.conf.GetBaseURL() + shortUID, LongURL: originalURL})
		}
	}

	ds.mu.RUnlock()

	if len(urls) == 0 {
		return urls, ErrUserUrlsEmpty
	}

	return urls, nil
}

// Create - add new user in DefaultStorage.
func (ds *DefaultStorage) Create() uint {
	ds.mu.Lock()
	uid := uint(len(ds.users) + 1)
	ds.users[uid] = make([]string, 0)
	ds.mu.Unlock()

	return uid
}
