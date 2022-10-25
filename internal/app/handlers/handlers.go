package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/alaleks/shortener/internal/app"
	"github.com/alaleks/shortener/internal/app/shortid"
	"github.com/gorilla/mux"
)

func ShortenURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "for use the shortener, you need to send a POST req", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlFromPost := strings.TrimSpace(string(body))

	if urlFromPost == "" {
		http.Error(w, "url is empty", http.StatusBadRequest)
		return
	}

	longUrl, err := app.IsUrl(urlFromPost)

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(app.GetShortUrl(longUrl)))
		return
	}

	http.Error(w, "invalid url: "+string(body), http.StatusBadRequest)

}

func ProcessingShortUrl(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]

	if uid == "" {
		http.Error(w, "uid is empty", http.StatusBadRequest)
		return
	}

	longUrl, check := app.GetLongUrl(uid)

	if !check {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GetStatistic(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]
	stat := app.GetStatisticShortUrl(uid)
	if stat == "" {
		http.Error(w, "this short url is invalid", http.StatusBadRequest)
		return
	}

	w.Write([]byte(stat))
}

func ShortenURLWithoutMux(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		urlSplit := strings.Split(r.URL.Path[1:], "/")
		switch {
		//if path empty
		case len(urlSplit) == 0:
			http.Error(w, "uid is empty", http.StatusBadRequest)
			return
		//check uid
		case len(urlSplit) == 1:
			uid := urlSplit[0]
			if len(uid) == shortid.GetSizeUid() {
				longUrl, check := app.GetLongUrl(uid)

				if !check {
					http.Error(w, "this short url is invalid", http.StatusBadRequest)
					return
				}
				w.Header().Set("Location", longUrl)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			} else {
				http.Error(w, "uid is invalid size characters", http.StatusBadRequest)
				return
			}
		case len(urlSplit) == 2:
			if urlSplit[1] == "statistic" {
				stat := app.GetStatisticShortUrl(urlSplit[0])
				if stat == "" {
					http.Error(w, "this short url is invalid", http.StatusBadRequest)
					return
				}
				w.Write([]byte(stat))
				return
			}
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		default:
			http.Error(w, "404 - page not found", http.StatusNotFound)
			return
		}
	}
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		urlFromPost := strings.TrimSpace(string(body))

		if urlFromPost == "" {
			http.Error(w, "url is empty", http.StatusBadRequest)
			return
		}

		longUrl, err := app.IsUrl(urlFromPost)

		if err == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(app.GetShortUrl(longUrl)))
			return
		}

		http.Error(w, "invalid url: "+string(body), http.StatusBadRequest)
	}
}
