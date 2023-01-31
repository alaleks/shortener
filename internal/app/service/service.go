package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidURL = errors.New("invalid url")

// generate uid string (letters English Alphabet).
func GenUID(size int) string {
	abc := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	buf := make([]byte, uint(size))

	randomizer(buf)

	charCnt := byte(len(abc))

	for i := range buf {
		buf[i] = abc[buf[i]%charCnt]
	}

	return string(buf)
}

func randomizer(buf []byte) {
	var number int

	var err error

	for number < len(buf) && err == nil {
		var i int
		i, err = rand.Reader.Read(buf[number:])
		number += i
	}
}

func IsURLOld(uri string) error {
	switch {
	case strings.HasPrefix(uri, "https://") && strings.TrimPrefix(uri, "https://") != "",
		strings.HasPrefix(uri, "http://") && strings.TrimPrefix(uri, "http://") != "",
		strings.HasPrefix(uri, "www.") && strings.TrimPrefix(uri, "www.") != "":
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidURL, uri)
	}
}

func IsURL(uri string) error {
	switch {
	case strings.HasPrefix(uri, "https://") && len(strings.TrimPrefix(uri, "https://")) > 0:
		return nil
	case strings.HasPrefix(uri, "http://") && len(strings.TrimPrefix(uri, "http://")) > 0:
		return nil
	case strings.HasPrefix(uri, "www.") && len(strings.TrimPrefix(uri, "www.")) > 0:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidURL, uri)
	}
}
