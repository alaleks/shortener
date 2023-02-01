package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
)

// ShortenURLAPI implements URL shortening.
// POST /api/shorten, JSON: {"url":"http://github.com/alaleks/shortener"}.
func (h *Handlers) ShortenURLAPI(writer http.ResponseWriter, req *http.Request) {
	var (
		input      InputShorten
		output     OutputShorten
		userID     string
		httpStatus = http.StatusCreated
	)

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		output.Err = err.Error()
	}

	if err := service.IsURL(input.URL); err != nil {
		output.Err = ErrInvalidRequest.Error()
	}

	if req.URL.User != nil {
		userID = req.URL.User.Username()
	}

	writer.Header().Set("Content-Type", "application/json")

	shortURL, err := h.Storage.Store.Add(input.URL, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			httpStatus = http.StatusConflict
		} else {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	writer.WriteHeader(httpStatus)

	if output.Err == "" {
		output.Success = true
		output.Result = shortURL
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(output); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		http.Error(writer, ErrInternalError.Error(), http.StatusBadRequest)

		return
	}
}

// GetStatAPI implements getting statistics on the use of a short URL.
// Example: GET /api/{uid}/statistics
func (h *Handlers) GetStatAPI(writer http.ResponseWriter, req *http.Request) {
	var (
		buffer bytes.Buffer
		uid    = mux.Vars(req)["uid"]
	)

	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	stat, err := h.Storage.Store.Stat(uid)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if err := json.NewEncoder(&buffer).Encode(stat); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		http.Error(writer, ErrInternalError.Error(), http.StatusBadRequest)

		return
	}
}

// GetUsersURL returns all user shortened URLs.
// GET /api/user/urls
func (h *Handlers) GetUsersURL(writer http.ResponseWriter, req *http.Request) {
	var buffer bytes.Buffer

	if req.URL.User == nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	out, err := h.Storage.Store.GetUrlsUser(req.URL.User.Username())
	if err != nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	if err := json.NewEncoder(&buffer).Encode(out); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		http.Error(writer, ErrInternalError.Error(), http.StatusBadRequest)

		return
	}
}

// ShortenURLBatch implements url batch shortening.
// POST /api/shorten/batch
// JSON: [{"original_url":"http://github.com/alaleks/shortener", "correlation_id":1}]
func (h *Handlers) ShortenURLBatch(writer http.ResponseWriter, req *http.Request) {
	var (
		input  []InShortenBatch
		userID string
	)

	if req.URL.User != nil {
		userID = req.URL.User.Username()
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	err = json.Unmarshal(body, &input)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	output, err := h.processingURLBatch(userID, input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	res, err := json.MarshalIndent(output, " ", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if _, err := writer.Write(res); err != nil {
		http.Error(writer, ErrInternalError.Error(), http.StatusBadRequest)

		return
	}
}

// ShortenDelete (without Worker Pool) implements
// the ability to delete user shortened URLs.
// DELETE /api/user/urls
func (h *Handlers) ShortenDelete(writer http.ResponseWriter, req *http.Request) {
	var (
		userID         string
		shortUIDForDel []string
	)

	if req.URL.User != nil {
		userID = req.URL.User.Username()
	}

	if err := json.NewDecoder(req.Body).Decode(&shortUIDForDel); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.Storage.Store.DelUrls(userID, checkShortUID(shortUIDForDel...)...); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

// ShortenDelete (with Worker Pool) implements
// the ability to delete user shortened URLs.
// DELETE /api/user/urls
func (h *Handlers) ShortenDeletePool(writer http.ResponseWriter, req *http.Request) {
	var data struct {
		userID         string
		shortUIDForDel []string
	}

	if req.URL.User != nil {
		data.userID = req.URL.User.Username()
	}

	if err := json.NewDecoder(req.Body).Decode(&data.shortUIDForDel); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	h.Storage.Pool.AddTask(data, func(data any) error {
		dataRemoved, ok := data.(struct {
			userID         string
			shortUIDForDel []string
		})

		if !ok {
			return storage.ErrInvalidData
		}

		return fmt.Errorf("deletion error: %w",
			h.Storage.Store.DelUrls(dataRemoved.userID,
				checkShortUID(dataRemoved.shortUIDForDel...)...))
	})

	writer.WriteHeader(http.StatusAccepted)
}
