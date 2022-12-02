package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/alaleks/shortener/internal/app/service"
	"github.com/gorilla/mux"
)

func (h *Handlers) ShortenURLAPI(writer http.ResponseWriter, req *http.Request) {
	var (
		input  InputShorten
		output OutputShorten
		userID string
	)

	body, err := io.ReadAll(req.Body)

	if req.Body != nil {
		defer req.Body.Close()
	}

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	err = json.Unmarshal(body, &input)

	if err != nil {
		output.Err = err.Error()
	}

	if err := service.IsURL(input.URL); err != nil {
		output.Err = ErrInvalidJSON.Error()
	}

	if req.URL.User != nil {
		userID = req.URL.User.Username()
	}

	writer.Header().Set("Content-Type", "application/json")

	if output.Err == "" {
		output.Success = true
		output.Result = h.AddShortenURL(userID, input.URL)
	}

	res, err := json.Marshal(output)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.WriteHeader(http.StatusCreated)

	if _, err := writer.Write(res); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) GetStatAPI(writer http.ResponseWriter, req *http.Request) {
	var (
		stat   Statistics
		buffer bytes.Buffer
		uid    = mux.Vars(req)["uid"]
	)

	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	stat, err := h.Statistics(uid)
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
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) GetUsersURL(writer http.ResponseWriter, req *http.Request) {
	var buffer bytes.Buffer

	if req.URL.User == nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	userID, err := strconv.Atoi(req.URL.User.Username())
	if err != nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	out, err := h.GetAllUrlsUser(userID)
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
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

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

	if req.Body != nil {
		defer req.Body.Close()
	}

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
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}
