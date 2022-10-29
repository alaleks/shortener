package storage

import (
	"sync"
	"testing"
	"time"
)

var urls = Urls{
	data: make(map[string]*urlEl),
	mtx:  &sync.Mutex{},
}
var uri = "https://github.com/alaleks/shortener"
var uriWWWW = "www.github.com/alaleks/shortener"

func TestAddGetUpdate(t *testing.T) {
	uid1 := urls.Add(uri)
	uid2 := urls.Add(uriWWWW)
	tests := []struct {
		uid string
		urlEl
	}{
		{
			uid1, urlEl{uri, time.Now(), 0},
		}, {
			uid2, urlEl{"http://" + uriWWWW, time.Now(), 0},
		},
	}
	for i, v := range tests {
		longURL, _ := urls.GetURL(v.uid)
		if v.longURL != longURL {
			t.Errorf("not correct expected url: %s", v.longURL)
		}
		if i == 1 {
			urls.Update(v.uid)
		}
		_, stat, _ := urls.Stat(v.uid)
		if i == 1 {
			if stat != 1 {
				t.Errorf("stat counter should be 1, not %d", stat)
			}
		} else {
			if stat != 0 {
				t.Errorf("stat counter should be 0, not %d", stat)
			}
		}

	}
}
