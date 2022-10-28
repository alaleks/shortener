package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alaleks/shortener/internal/service"
	"github.com/alaleks/shortener/internal/storage"
	"github.com/gorilla/mux"
)

var dataStorage storage.Storager

func init() {
	dataStorage = &storage.Urls{
		LongUrls:  make(map[string]*storage.ShortUrl),
		ShortUrls: make(map[string]string),
	}
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uri := strings.TrimSpace(string(body))

	if uri == "" {
		http.Error(w, "url is empty", http.StatusBadRequest)
		return
	}

	err = service.IsUrl(uri)

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		host := "http://" + r.Host + "/"

		if r.TLS != nil {
			host = "https://" + r.Host + "/"
		}

		if uid, ok := dataStorage.Get(uri); ok {
			w.Write([]byte(host + string(uid)))
			return
		}
		newUid := service.GenUid(5)
		dataStorage.Add(uri, newUid)
		w.Write([]byte(host + newUid))
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

	uri, ok := dataStorage.Get(uid)
	if !ok {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	dataStorage.Update(uid)
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GetStat(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]
	uri, ok := dataStorage.Get(uid)
	if !ok {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	host := "http://" + r.Host + "/"

	if r.TLS != nil {
		host = "https://" + r.Host + "/"
	}

	id, counterStat, created := dataStorage.Stat(uri)

	w.Write([]byte(fmt.Sprintf("short link: %s%s \nurl: %s \nusage: %d \ncreated: %s", host, id, uri, counterStat, created)))
}
