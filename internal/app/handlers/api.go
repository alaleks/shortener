package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

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

	res, err := json.Marshal(output)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if _, err := writer.Write(res); err != nil {
		http.Error(writer, ErrInternalError.Error(), http.StatusBadRequest)

		return
	}
}

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

	output, err := h.ProcessingURLBatch(userID, input)
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

	h.Storage.Store.DelUrls(userID, checkShortUID(shortUIDForDel...)...)

	writer.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) ShortenDeletePool(writer http.ResponseWriter, req *http.Request) {
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

	h.Storage.Pool.Add(func() {
		h.Storage.Store.DelUrls(userID, checkShortUID(shortUIDForDel...)...)
	})

	writer.WriteHeader(http.StatusAccepted)
}
