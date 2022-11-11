package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/service"
	"github.com/gorilla/mux"
)

type Statistics struct {
	ShortURL  string `json:"shorturl"`
	LongURL   string `json:"longurl"`
	Usage     uint   `json:"usage"`
	CreatedAt string `json:"createdAt"`
}

type InputShorten struct {
	URL string `json:"url"`
}

type OutputShorten struct {
	Success bool   `json:"success"`
	Result  string `json:"result,omitempty"`
	Err     error  `json:"-"`
	ErrMsg  string `json:"error,omitempty"`
}

var ErrInvalidJSON = errors.New("json is invalid, please check what you send. Should be: {'url':'https://example.ru'}")

func (h *Handlers) ShortenURLAPI(writer http.ResponseWriter, req *http.Request) {
	var (
		buffer bytes.Buffer
		input  InputShorten
		output OutputShorten
	)

	_, err := io.Copy(&buffer, req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	err = json.NewDecoder(&buffer).Decode(&input)

	switch {
	case err != nil:
		output.Err = err
	case input.URL == "":
		output.Err = ErrInvalidJSON
	default:
		output.Err = service.IsURL(input.URL)
	}

	if output.Err != nil {
		buffer.Reset()

		output.ErrMsg = output.Err.Error()
		if err := json.NewEncoder(&buffer).Encode(output); err != nil {
			http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

			return
		}

		writer.WriteHeader(http.StatusResetContent)

		if _, err := writer.Write(buffer.Bytes()); err != nil {
			http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

			return
		}

		return
	}

	buffer.Reset()

	output.Success = true
	output.Result = h.createShortURL(input.URL, req)

	writer.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(&buffer).Encode(output); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) createShortURL(longURL string, req *http.Request) string {
	var shortURL bytes.Buffer

	shortURL.WriteString(func() string {
		if req.TLS != nil {
			return "https://"
		}

		return "http://"
	}())

	shortURL.WriteString(req.Host + "/")

	uid := h.DataStorage.Add(longURL, h.SizeUID)
	shortURL.WriteString(uid)

	return shortURL.String()
}

func (h *Handlers) GetStatAPI(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]
	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	host := "http://" + req.Host + "/"

	if req.TLS != nil {
		host = "https://" + req.Host + "/"
	}

	longURL, counterStat, createdAt := h.DataStorage.Stat(uid)

	if longURL == "" {
		http.Error(writer, ErrUIDInvalid.Error(), http.StatusBadRequest)

		return
	}

	var buffer bytes.Buffer

	stat := Statistics{
		ShortURL:  host + uid,
		LongURL:   longURL,
		Usage:     counterStat,
		CreatedAt: createdAt,
	}

	if err := json.NewEncoder(&buffer).Encode(stat); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}
