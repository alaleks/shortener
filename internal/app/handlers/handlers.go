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

func ShortenURL(w http.ResponseWriter, r *http.Request) {

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

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		host := "http://" + r.Host + "/"

		if r.TLS != nil {
			host = "https://" + r.Host + "/"
		}

		uid := storage.DataStorage.Add(longURL)
		w.Write([]byte(host + uid))
		return
	}

	http.Error(w, err.Error(), http.StatusBadRequest)

}

func ParseShortURL(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]

	if uid == "" {
		http.Error(w, "uid is empty", http.StatusBadRequest)
		return
	}

	longURL, ok := storage.DataStorage.GetURL(uid)
	if !ok {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	storage.DataStorage.Update(uid)
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GetStat(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]

	if uid == "" {
		http.Error(w, "uid is empty", http.StatusBadRequest)
		return
	}

	host := "http://" + r.Host + "/"

	if r.TLS != nil {
		host = "https://" + r.Host + "/"
	}

	longURL, counterStat, created := storage.DataStorage.Stat(uid)

	if longURL == "" {
		http.Error(w, "uid is invalid", http.StatusBadRequest)
	}

	w.Write([]byte(fmt.Sprintf("short link: %s%s \nurl: %s \nusage: %d \ncreated: %s", host, uid, longURL, counterStat, created)))
}
