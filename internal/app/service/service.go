package service

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// generate uid string (letters English Alphabet)
func GenUID(size uint) string {
	abc := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0")
	b := make([]byte, 5)
	randomizer(b)
	charCnt := byte(len(abc))
	for i := range b {
		b[i] = abc[b[i]%charCnt]
	}
	return string(b)
}

func randomizer(buf []byte) {
	var n int
	var err error
	for n < len(buf) && err == nil {
		var i int
		i, err = rand.Reader.Read(buf[n:])
		n += i
	}
}

func IsURL(uri string) error {
	switch {
	case strings.HasPrefix(uri, "https://") && strings.TrimPrefix(uri, "https://") != "",
		strings.HasPrefix(uri, "http://") && strings.TrimPrefix(uri, "http://") != "",
		strings.HasPrefix(uri, "www.") && strings.TrimPrefix(uri, "www.") != "":
		return nil
	default:
		return fmt.Errorf("invalid url: %s", uri)
	}
}

func RemovePrefix(str string, prefixs ...string) string {
	for _, prefix := range prefixs {
		str = strings.TrimPrefix(str, prefix)
	}
	return str
}
