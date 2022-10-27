package storage

import "strings"

type LongURL string
type Uid string
type ShortURL struct {
	uri       LongURL
	statistic uint
}

type Urls struct {
	LongUrls  map[LongURL]Uid
	ShortUrls map[Uid]*ShortURL
}

var UrlsStorage *Urls

func init() {
	UrlsStorage = &Urls{
		LongUrls:  make(map[LongURL]Uid),
		ShortUrls: make(map[Uid]*ShortURL),
	}
}

func (u *Urls) Add(long LongURL, id Uid) {
	if !strings.HasPrefix(string(long), "http") {
		long = "http://" + long
	}
	u.ShortUrls[id] = &ShortURL{long, 0}
	switch {
	case strings.HasPrefix(string(long), "https://"):
		long = LongURL(strings.TrimPrefix(string(long), "https://"))
	case strings.HasPrefix(string(long), "http://"):
		long = LongURL(strings.TrimPrefix(string(long), "http://"))
	case strings.HasPrefix(string(long), "www."):
		long = LongURL(strings.TrimPrefix(string(long), "www."))
	}
	u.LongUrls[long] = id
}

func (u *Urls) FindUidShortUrl(long LongURL) (Uid, bool) {
	v, ok := u.LongUrls[long]
	return v, ok
}

func (u *Urls) FindLongUrl(id Uid) (LongURL, bool) {
	v, ok := u.ShortUrls[id]
	if !ok {
		return "", ok
	}
	v.statistic++
	u.ShortUrls[id] = v
	return v.uri, ok
}

func (u *Urls) GetStatistic(id Uid) (string, uint) {
	v, ok := u.ShortUrls[id]
	if ok {
		return string(v.uri), v.statistic
	}
	return "", 0
}
