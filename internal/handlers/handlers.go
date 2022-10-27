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

func UseShortner(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "for use the shortener, you need to send a POST req", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlForShortener := strings.TrimSpace(string(body))

	if urlForShortener == "" {
		http.Error(w, "url is empty", http.StatusBadRequest)
		return
	}

	err = service.IsUrl(urlForShortener)

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		host := "http://" + r.Host + "/"

		if r.TLS != nil {
			host = "https://" + r.Host + "/"
		}

		if uid, ok := storage.UrlsStorage.FindUidShortUrl(storage.LongURL(urlForShortener)); ok {
			w.Write([]byte(host + string(uid)))
			return
		}
		newUid := service.CreateShortId(5)
		storage.UrlsStorage.Add(storage.LongURL(urlForShortener), storage.Uid(newUid))
		w.Write([]byte(host + newUid))
		return
	}

	http.Error(w, err.Error(), http.StatusBadRequest)

}

func ParseShortUrl(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]

	if uid == "" {
		http.Error(w, "uid is empty", http.StatusBadRequest)
		return
	}

	longUrl, ok := storage.UrlsStorage.FindLongUrl(storage.Uid(uid))
	if !ok {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", string(longUrl))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GiveAwayStatistic(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]
	uri, stat := storage.UrlsStorage.GetStatistic(storage.Uid(uid))
	if uri == "" {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	host := "http://" + r.Host + "/"

	if r.TLS != nil {
		host = "https://" + r.Host + "/"
	}

	w.Write([]byte(fmt.Sprintf("short link: %s%s \nurl: %s \nusage: "+"%d", host, uid, uri, stat)))
}
