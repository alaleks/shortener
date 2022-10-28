package service

import (
	"testing"
)

func TestIsUrl(t *testing.T) {
	f := IsUrl
	tests := []struct {
		uri   string
		check bool
	}{
		{"https://github.com/alaleks/shortener", true},
		{"http://github.com/alaleks/shortener", true},
		{"www.github.com/alaleks/shortener", true},
		{"github.com/alaleks/shortener", false},
		{"htts://github.com/alaleks/shortener", false},
		{"https://", false},
	}
	for _, v := range tests {
		err := f(v.uri)
		if v.check {
			if err != nil {
				t.Errorf("should error be nil: %v", err)
			}
		} else {
			if err == nil {
				t.Errorf("should error be not nil: %v", err)
			}
		}

	}

}

func TestCreateShortId(t *testing.T) {
	size := 5
	id1 := GenUid(uint(size))
	id2 := GenUid(uint(size))
	if len(id1) != size || len(id2) != size {
		t.Errorf("uid should be сonsist %d characters", size)
	}
	if id1 == id2 {
		t.Errorf("uids should not be egual each other")
	}
}

func TestRemovePrefix(t *testing.T) {
	tests := []struct {
		uriRaw  string
		uriWant string
	}{
		{"https://github.com/alaleks/shortener", "github.com/alaleks/shortener"},
		{"http://github.com/alaleks/shortener", "github.com/alaleks/shortener"},
		{"www.github.com/alaleks/shortener", "github.com/alaleks/shortener"},
	}

	for _, v := range tests {
		urlClean := RemovePrefix(v.uriRaw, "https://", "http://", "www.")
		if urlClean != v.uriWant {
			t.Errorf("url clean should be %s", v.uriWant)
		}
	}

}
