package handlers

import (
	"github.com/alaleks/shortener/internal/app/service"
)

func (h *Handlers) AddShortenURL(userID, longURL string) (string, error) {
	var (
		shortURL = h.baseURL
		shortUID string
	)

	if err := h.PingDB(); err != nil {
		h.checkDB = false
	}

	if h.checkDB {
		shortUID, err := h.DB.AddURL(userID, service.GenUID(h.SizeUID), longURL)
		shortURL += shortUID

		return shortURL, err
	}

	shortUID = h.DataStorage.Add(longURL, h.SizeUID)

	h.Users.AddShortUID(userID, shortUID)

	shortURL += shortUID

	return shortURL, nil
}

func (h *Handlers) GetOriginalURL(uid string) (string, error) {
	if err := h.PingDB(); err != nil {
		h.checkDB = false
	}

	if h.checkDB {
		longURL := h.DB.GetOriginalURL(uid)

		if longURL == "" {
			return longURL, ErrInvalidUID
		}

		h.DB.UpdateStat(uid)

		return longURL, nil
	}

	longURL, ok := h.DataStorage.GetURL(uid)

	if !ok {
		return longURL, ErrInvalidUID
	}

	return longURL, nil
}

func (h *Handlers) Statistics(uid string) (Statistics, error) {
	if err := h.PingDB(); err != nil {
		h.checkDB = false
	}

	if h.checkDB {
		stat := h.DB.GetStat(uid)

		if stat.LongURL == "" {
			return stat, ErrInvalidUID
		}

		stat.ShortURL = h.baseURL + stat.ShortURL

		return stat, nil
	}

	var stat Statistics
	longURL, counterStat, createdAt := h.DataStorage.Stat(uid)

	if longURL == "" {
		return stat, ErrInvalidUID
	}

	stat.ShortURL = h.baseURL + uid
	stat.LongURL = longURL
	stat.Usage = counterStat
	stat.CreatedAt = createdAt

	return stat, nil
}

func (h *Handlers) GetAllUrlsUser(userID int) ([]struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}, error,
) {
	if err := h.PingDB(); err != nil {
		h.checkDB = false
	}

	if h.checkDB {
		userUrls := h.DB.GetUrlsUserHandler(userID)

		if len(userUrls) == 0 {
			return userUrls, ErrUserDoesNotExist
		}

		for index := range userUrls {
			userUrls[index].ShortURL = h.baseURL + userUrls[index].ShortURL
		}

		return userUrls, nil
	}

	var userUrls []struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	uidsShorlURL, _ := h.Users.Check(uint(userID))

	if len(uidsShorlURL) == 0 {
		return userUrls, ErrUserDoesNotExist
	}

	for _, item := range uidsShorlURL {
		uri, check := h.DataStorage.GetURL(item)

		if check {
			userUrls = append(userUrls, struct {
				ShortURL    string `json:"short_url"`
				OriginalURL string `json:"original_url"`
			}{ShortURL: h.baseURL + item, OriginalURL: uri})
		}
	}

	return userUrls, nil
}

func (h *Handlers) ProcessingURLBatch(userID string, input []InShortenBatch) ([]OutShortenBatch, error) {
	out := make([]OutShortenBatch, 0, len(input))

	if err := h.PingDB(); err != nil {
		h.checkDB = false
	}

	if h.checkDB {
		for _, item := range input {
			err := service.IsURL(item.OriginalURL)

			if err == nil {
				shortUID := h.DB.AddURLBatch(userID, service.GenUID(h.SizeUID), item.CorID, item.OriginalURL)
				out = append(out, OutShortenBatch{CorID: item.CorID, ShortURL: h.baseURL + shortUID})
			} else {
				out = append(out, OutShortenBatch{CorID: item.CorID, Err: err.Error()})
			}
		}

		if len(out) == 0 {
			return out, ErrEmptyBatch
		}

		return out, nil
	}

	for _, item := range input {
		err := service.IsURL(item.OriginalURL)

		if err == nil {
			shortUID := h.DataStorage.AddBatch(h.SizeUID, item.CorID, item.OriginalURL)

			h.Users.AddShortUID(userID, shortUID)

			out = append(out, OutShortenBatch{CorID: item.CorID, ShortURL: h.baseURL + shortUID})
		} else {
			out = append(out, OutShortenBatch{CorID: item.CorID, Err: err.Error()})
		}
	}

	if len(out) == 0 {
		return out, ErrEmptyBatch
	}

	return out, nil
}
