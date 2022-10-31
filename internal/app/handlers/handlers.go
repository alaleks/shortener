package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage"

	"github.com/gorilla/mux"
)

func ShortenURL(dataStorage storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		longURL := strings.TrimSpace(string(body))

		if longURL == "" {
			http.Error(w, "url is empty", http.StatusBadRequest)
			return
		}

		err = service.IsURL(longURL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusCreated)
		host := "http://" + r.Host + "/"

		if r.TLS != nil {
			host = "https://" + r.Host + "/"
		}

		uid := dataStorage.Add(longURL)
		w.Write([]byte(host + uid))

	}
}

func ParseShortURL(dataStorage storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := mux.Vars(r)["uid"]

		if uid == "" {
			http.Error(w, "uid is empty", http.StatusBadRequest)
			return
		}

		longURL, ok := dataStorage.GetURL(uid)
		if !ok {
			http.Error(w, "this short url is invalid", http.StatusBadRequest)
			return
		}
		dataStorage.Update(uid)
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func GetStat(dataStorage storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := mux.Vars(r)["uid"]

		if uid == "" {
			http.Error(w, "uid is empty", http.StatusBadRequest)
			return
		}

		host := "http://" + r.Host + "/"

		if r.TLS != nil {
			host = "https://" + r.Host + "/"
		}

		longURL, counterStat, created := dataStorage.Stat(uid)

		if longURL == "" {
			http.Error(w, "uid is invalid", http.StatusBadRequest)
		}

		w.Write([]byte(fmt.Sprintf("short link: %s%s \nurl: %s \nusage: %d \ncreated: %s", host, uid, longURL, counterStat, created)))
	}
}
