// Package service implements helper functions for the application.
package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
)

const (
	prefixHTTPS = "https://"
	prefixHTTP  = "http://"
	prefixWWW   = "www."
)

// ErrInvalidURL is an indicator that the invalid URL.
var ErrInvalidURL = errors.New("invalid url")

// Generate uid uses letters English Alphabet.
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

// IsURLOld (Deprecated) checks if a string is a URL.
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

// IsURL checks if a string is a URL and
// returns true if the string is a valid URL.
func IsURL(uri string) error {
	switch {
	case strings.HasPrefix(uri, prefixHTTPS) && uri != prefixHTTPS:
		return nil
	case strings.HasPrefix(uri, prefixHTTP) && uri != prefixHTTP:
		return nil
	case strings.HasPrefix(uri, prefixWWW) && uri != prefixWWW:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidURL, uri)
	}
}
