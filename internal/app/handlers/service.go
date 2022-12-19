package handlers

import (
	"strings"

	"github.com/alaleks/shortener/internal/app/service"
)

func (h *Handlers) ProcessingURLBatch(userID string, input []InShortenBatch) ([]OutShortenBatch, error) {
	out := make([]OutShortenBatch, 0, len(input))

	for _, item := range input {
		err := service.IsURL(item.OriginalURL)

		if err == nil {
			shortURL := h.Storage.Store.AddBatch(item.OriginalURL, userID, item.CorID)
			out = append(out, OutShortenBatch{CorID: item.CorID, ShortURL: shortURL})
		} else {
			out = append(out, OutShortenBatch{CorID: item.CorID, Err: err.Error()})
		}
	}

	if len(out) == 0 {
		return out, ErrEmptyBatch
	}

	return out, nil
}

func checkShortUID(shortUID ...string) []string {
	var correctShortUID []string

	for _, v := range shortUID {
		if !strings.Contains(v, "/") && v != "" {
			correctShortUID = append(correctShortUID, v)
		} else {
			if sUID := v[strings.LastIndex(v, "/")+1:]; sUID != "" {
				correctShortUID = append(correctShortUID, sUID)
			}
		}
	}

	return correctShortUID
}
