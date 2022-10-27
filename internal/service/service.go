package service

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// generate uid string (letters English Alphabet)
func CreateShortId(size uint) string {
	abc := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0")
	b := make([]byte, 5)
	randb(b)
	charCnt := byte(len(abc))
	for i := range b {
		b[i] = abc[b[i]%charCnt]
	}
	return string(b)
}

func randb(buf []byte) {
	var n int
	var err error
	for n < len(buf) && err == nil {
		var i int
		i, err = rand.Reader.Read(buf[n:])
		n += i
	}
}

func IsUrl(uri string) error {
	switch {
	case strings.HasPrefix(uri, "https://"),
		strings.HasPrefix(uri, "http://"),
		strings.HasPrefix(uri, "www."):
		return nil
	default:
		return fmt.Errorf("invalid url: %s", uri)
	}
}
