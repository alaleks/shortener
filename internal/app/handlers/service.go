package handlers

import (
	"github.com/alaleks/shortener/internal/app/database/methods"
	"github.com/alaleks/shortener/internal/app/service"
)

func (h *Handlers) AddShortenURL(userID, longURL string) string {
	var (
		shortURL = h.baseURL
		shortUID string
	)

	switch {
	case h.DSN != "":
		dBase := methods.OpenDB(h.DSN)

		if dBase.DB != nil {
			defer dBase.Close()

			shortUID = service.GenUID(h.SizeUID)
			dBase.AddURL(userID, shortUID, longURL)

			shortURL += shortUID

			return shortURL
		}

		fallthrough
	default:
		shortUID = h.DataStorage.Add(longURL, h.SizeUID)

		h.Users.AddShortUID(userID, shortUID)

		shortURL += shortUID

		return shortURL
	}
}

func (h *Handlers) GetOriginalURL(uid string) (string, error) {
	switch {
	case h.DSN != "":
		dBase := methods.OpenDB(h.DSN)

		if dBase.DB != nil {
			defer dBase.Close()

			longURL := dBase.GetOriginalURL(uid)

			if longURL == "" {
				return longURL, ErrUIDInvalid
			}

			dBase.UpdateStat(uid)

			return longURL, nil
		}

		fallthrough
	default:
		longURL, ok := h.DataStorage.GetURL(uid)

		if !ok {
			return longURL, ErrUIDInvalid
		}

		return longURL, nil
	}
}

func (h *Handlers) Statistics(uid string) (Statistics, error) {
	switch {
	case h.DSN != "":
		dBase := methods.OpenDB(h.DSN)

		if dBase.DB != nil {
			defer dBase.Close()

			stat := dBase.GetStat(uid)

			if stat.LongURL == "" {
				return stat, ErrUIDInvalid
			}

			stat.ShortURL = h.baseURL + stat.ShortURL

			return stat, nil
		}

		fallthrough
	default:
		var stat Statistics
		longURL, counterStat, createdAt := h.DataStorage.Stat(uid)

		if longURL == "" {
			return stat, ErrUIDInvalid
		}

		stat.ShortURL = h.baseURL + uid
		stat.LongURL = longURL
		stat.Usage = counterStat
		stat.CreatedAt = createdAt

		return stat, nil
	}
}

func (h *Handlers) GetAllUrlsUser(userID int) ([]struct {
	ShotrURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}, error,
) {
	switch {
	case h.DSN != "":
		dBase := methods.OpenDB(h.DSN)

		if dBase.DB != nil {
			defer dBase.Close()

			userUrls := dBase.GetUrlsUserHandler(userID)

			if len(userUrls) == 0 {
				return userUrls, ErrGetUrlsUser
			}

			for index := range userUrls {
				userUrls[index].ShotrURL = h.baseURL + userUrls[index].ShotrURL
			}

			return userUrls, nil
		}

		fallthrough
	default:
		var userUrls []struct {
			ShotrURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}

		uidsShorlURL, _ := h.Users.Check(uint(userID))

		if len(uidsShorlURL) == 0 {
			return userUrls, ErrGetUrlsUser
		}

		for _, item := range uidsShorlURL {
			uri, check := h.DataStorage.GetURL(item)

			if check {
				userUrls = append(userUrls, struct {
					ShotrURL    string `json:"short_url"`
					OriginalURL string `json:"original_url"`
				}{ShotrURL: h.baseURL + item, OriginalURL: uri})
			}
		}

		return userUrls, nil
	}
}

func (h *Handlers) ProcessingURLBatch(userID string, input []InShortenBatch) ([]OutShortenBatch, error) {
	out := make([]OutShortenBatch, 0, len(input))

	switch {
	case h.DSN != "":
		dBase := methods.OpenDB(h.DSN)

		if dBase.DB != nil {
			defer dBase.Close()

			for _, item := range input {
				err := service.IsURL(item.OriginalURL)

				if err == nil {
					shortUID := service.GenUID(h.SizeUID)
					dBase.AddURL(userID, shortUID, item.OriginalURL)

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

		fallthrough
	default:
		for _, item := range input {
			err := service.IsURL(item.OriginalURL)

			if err == nil {
				shortUID := h.DataStorage.Add(item.OriginalURL, h.SizeUID)

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
}
