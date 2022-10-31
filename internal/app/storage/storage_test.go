package storage

import (
	"testing"
	"time"
)

func TestAddGetUpdate(t *testing.T) {
	data := New()
	var uri = "https://github.com/alaleks/shortener"
	var uriWWWW = "www.github.com/alaleks/shortener"
	uid1 := data.Add(uri)
	uid2 := data.Add(uriWWWW)

	tests := []struct {
		name string
		uid  string
		urlEl
	}{
		{
			"url c http", uid1, urlEl{uri, time.Now(), 0},
		}, {
			"url c www", uid2, urlEl{"http://" + uriWWWW, time.Now(), 0},
		},
	}
	for i, v := range tests {
		v := v
		t.Run(v.name, func(t *testing.T) {
			longURL, _ := data.GetURL(v.uid)
			if v.longURL != longURL {
				t.Errorf("not correct expected url: %s", v.longURL)
			}
			if i == 1 {
				data.Update(v.uid)
			}
			_, stat, _ := data.Stat(v.uid)
			if i == 1 {
				if stat != 1 {
					t.Errorf("stat counter should be 1, not %d", stat)
				}
			} else {
				if stat != 0 {
					t.Errorf("stat counter should be 0, not %d", stat)
				}
			}
		})

	}
}
