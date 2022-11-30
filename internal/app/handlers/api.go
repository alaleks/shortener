package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

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
	Result  string `json:"result,omitempty"`
	Err     string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

type OutputUsersUrls struct {
	ShotrURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var ErrInvalidJSON = errors.New(`json is invalid, please check what you send. Should be: {"url":"https://example.ru"}`)

func (h *Handlers) ShortenURLAPI(writer http.ResponseWriter, req *http.Request) {
	var (
		input  InputShorten
		output OutputShorten
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

	writer.Header().Set("Content-Type", "application/json")

	if output.Err == "" {
		output.Success = true
		uid := h.DataStorage.Add(input.URL, h.SizeUID)
		output.Result = h.baseURL + uid

		if req.URL.User != nil {
			h.Users.AddShortUID(req.URL.User.Username(), uid)
		}

		writer.WriteHeader(http.StatusCreated)
	}

	res, err := json.Marshal(output)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if _, err := writer.Write(res); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) GetStatAPI(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]
	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	longURL, counterStat, createdAt := h.DataStorage.Stat(uid)

	if longURL == "" {
		http.Error(writer, ErrUIDInvalid.Error(), http.StatusBadRequest)

		return
	}

	var buffer bytes.Buffer

	stat := Statistics{
		ShortURL:  h.baseURL + uid,
		LongURL:   longURL,
		Usage:     counterStat,
		CreatedAt: createdAt,
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

	uidsShorlURL, _ := h.Users.Check(uint(userID))

	if len(uidsShorlURL) == 0 {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	out := []OutputUsersUrls{}

	for _, v := range uidsShorlURL {
		uri, check := h.DataStorage.GetURL(v)

		if check {
			out = append(out, OutputUsersUrls{ShotrURL: h.baseURL + v, OriginalURL: uri})
		}
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
