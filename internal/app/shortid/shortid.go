package shortid

import (
	"crypto/rand"
)

var sizeUid = 5

// generate a default ID of 5 character (letters English Alphabet)
// return uid string
func CreateShortId() string {
	abc := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0")
	b := make([]byte, sizeUid)
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

// get size uid
func GetSizeUid() int {
	return sizeUid
}

// cnange size uid
func ChangeSizeUid(size int) {
	sizeUid = size
}
