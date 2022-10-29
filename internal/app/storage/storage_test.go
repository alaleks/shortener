package storage

import (
	"sync"
	"testing"
	"time"
)

var urls = Urls{
	data: make(map[string]*fields),
	mtx:  &sync.Mutex{},
}
var uri = "https://github.com/alaleks/shortener"
var uriWWWW = "www.github.com/alaleks/shortener"

func TestAddGetUpdate(t *testing.T) {
	uid1 := urls.Add(uri)
	uid2 := urls.Add(uriWWWW)
	tests := []struct {
		uid string
		fields
	}{
		{
			uid1, fields{uri, time.Now(), 0},
		}, {
			uid2, fields{"http://" + uriWWWW, time.Now(), 0},
		},
	}
	for i, v := range tests {
		longurl, _ := urls.GetURL(v.uid)
		if v.longUrl != longurl {
			t.Errorf("not correct expected url: %s", v.longUrl)
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
