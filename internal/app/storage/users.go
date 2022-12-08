package storage

import (
	"strconv"
	"sync"
)

type Users struct {
	data map[uint][]string
	mu   sync.RWMutex
}

func NewUsers() Users {
	return Users{data: make(map[uint][]string), mu: sync.RWMutex{}}
}

func (u *Users) Check(uid uint) ([]string, bool) {
	u.mu.RLock()
	uidsShortURL, check := u.data[uid]
	u.mu.RUnlock()

	return uidsShortURL, check
}

func (u *Users) Create() uint {
	u.mu.Lock()
	uid := uint(len(u.data) + 1)
	u.data[uid] = make([]string, 0)
	u.mu.Unlock()

	return uid
}

func (u *Users) AddShortUID(uid, uidShortURL string) {
	uidToInt, err := strconv.Atoi(uid)
	if err != nil {
		return
	}

	u.mu.Lock()
	u.data[uint(uidToInt)] = append(u.data[uint(uidToInt)], uidShortURL)
	u.mu.Unlock()
}
