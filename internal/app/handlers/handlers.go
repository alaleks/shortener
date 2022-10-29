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

	longUrl := strings.TrimSpace(string(body))

	if longUrl == "" {
		http.Error(w, "url is empty", http.StatusBadRequest)
		return
	}

	err = service.IsUrl(longUrl)

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		host := "http://" + r.Host + "/"

		if r.TLS != nil {
			host = "https://" + r.Host + "/"
		}

		uid := storage.DataStorage.Add(longUrl)
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

	longUrl, ok := storage.DataStorage.GetURL(uid)
	if !ok {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	storage.DataStorage.Update(uid)
	w.Header().Set("Location", longUrl)
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

	longUrl, counterStat, created := storage.DataStorage.Stat(uid)

	w.Write([]byte(fmt.Sprintf("short link: %s%s \nurl: %s \nusage: %d \ncreated: %s", host, uid, longUrl, counterStat, created)))
}
