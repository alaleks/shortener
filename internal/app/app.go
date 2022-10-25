package app

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/alaleks/shortener/internal/app/shortid"
)

// where key is a url (before shortener)
// value is an uid of short url
type urls map[string]string

// where key is uid of short url
// value is a struct with fields:
// url (from package net/url)
// statistic quantity redirectes (quantity times used this short link)
type shortUrls map[string]struct {
	*url.URL
	*sync.Mutex
	statistic uint
}

// config host for app
type hostConfig struct {
	*url.URL
	port string
}

var host hostConfig
var dataShortUrls shortUrls
var dataUrls urls

func init() {
	host = hostConfig{&url.URL{Scheme: "http", Host: "localhost:8080"}, ":8080"}
	dataShortUrls = make(shortUrls)
	dataUrls = make(urls)

}

// add your host
// default - http://localhost:8080, port :8080
func ConfigureHost(hostName, scheme, port string) {
	host.Host = hostName
	host.Scheme = scheme
	host.port = port
}

func GetPort() string {
	return host.port
}

// getting shortener url for this url
func GetShortUrl(longUrl *url.URL) []byte {

	uid, ok := dataUrls[longUrl.String()]

	if ok {
		host.URL.Path = uid
		return []byte(host.URL.String())
	}

	newUid := shortid.CreateShortId()
	dataUrls[longUrl.String()] = newUid
	host.URL.Path = newUid

	dataShortUrls[newUid] = struct {
		*url.URL
		*sync.Mutex
		statistic uint
	}{longUrl, &sync.Mutex{}, 0}

	return []byte(host.URL.String())
}

// getting long url for this uid
func GetLongUrl(uid string) (url string, found bool) {
	data, ok := dataShortUrls[uid]

	if ok {
		data.Lock()
		data.statistic += 1
		data.Unlock()
		dataShortUrls[uid] = data
		return data.String(), true
	}

	return "", false
}

// getting statistic for shortener
// return url before shortener
// statistic is quantity redirectes (quantity times used this short link)
func GetStatisticShortUrl(uid string) string {
	data, ok := dataShortUrls[uid]

	if ok {
		return data.String() + "\nusage: " + fmt.Sprint(data.statistic)
	}
	return ""
}

// validate received url
func IsUrl(longUrl string) (*url.URL, error) {
	return url.ParseRequestURI(longUrl)
}
