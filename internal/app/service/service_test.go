package service

import (
	"testing"
)

func TestIsURL(t *testing.T) {
	f := IsURL
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
	id1 := GenUID(uint(size))
	id2 := GenUID(uint(size))
	if len(id1) != size || len(id2) != size {
		t.Errorf("uid should be сonsist %d characters", size)
	}
	if id1 == id2 {
		t.Errorf("uids should not be egual each other")
	}
}
