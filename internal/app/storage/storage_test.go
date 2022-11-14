package storage

import (
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	dataStorage := Urls{
		data: make(map[string]*URLElement),
		mu:   sync.RWMutex{},
	}
	uri1 := "https://github.com/alaleks/shortener"
	uri2 := "www.github.com/alaleks/shortener"
	uid1 := dataStorage.Add(uri1, 5)
	uid2 := dataStorage.Add(uri2, 5)

	tests := []struct {
		name    string
		uid     string
		wantURL string
	}{
		{
			"url c http", uid1, uri1,
		}, {
			"url c www", uid2, "http://" + uri2,
		}, // в handler есть проверка на url
		// поэтому других вариантов в эту функцию не прилетит
	}
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			if el, ok := dataStorage.data[item.uid]; !ok {
				t.Errorf("uid %s should be return true", item.uid)
			} else if item.wantURL != el.LongURL {
				t.Errorf("uid %s should be return this URL %s but no %s", item.uid, item.wantURL, el.LongURL)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	t.Parallel()

	dataStorage := Urls{
		data: make(map[string]*URLElement),
		mu:   sync.RWMutex{},
	}

	uri1 := "https://github.com/alaleks/shortener"
	uri2 := "www.github.com/alaleks/shortener"

	uid1 := dataStorage.Add(uri1, 5)
	uid2 := dataStorage.Add(uri2, 5)

	tests := []struct {
		name    string
		uid     string
		wantURL string
	}{
		{
			"url c http", uid1, uri1,
		}, {
			"url c www", uid2, "http://" + uri2,
		},
	}
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			longURL, _ := dataStorage.GetURL(item.uid)
			if longURL != item.wantURL {
				t.Errorf("uid %s should be return this URL %s but no %s", item.uid, item.wantURL, longURL)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	dataStorage := Urls{
		data: make(map[string]*URLElement),
		mu:   sync.RWMutex{},
	}

	uid1 := dataStorage.Add("https://github.com/alaleks/shortener", 5)
	uid2 := dataStorage.Add("www.github.com/alaleks/shortener", 5)

	dataStorage.Update(uid1)

	tests := []struct {
		name     string
		uid      string
		wantStat uint
	}{
		{
			"обновление статистики", uid1, 1,
		}, {
			"без обновления статистики", uid2, 0,
		},
	}
	for _, v := range tests {
		item := v
		t.Run(item.name, func(t *testing.T) {
			t.Parallel()

			el := dataStorage.data[item.uid]
			if el.Statistics != item.wantStat {
				t.Errorf("uid %s should be return stat %d but no %d", item.uid, item.wantStat, el.Statistics)
			}
		})
	}
}
