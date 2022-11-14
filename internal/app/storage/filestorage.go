package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
)

type FileStorage interface {
	Write(filepath string) error
	Read(filepath string) error
}

func (u *Urls) Write(filepath string) error {
	file, err := os.Create(filepath + "filestorage.gob")

	if file != nil {
		defer file.Close()
	}

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	encoder.Encode(u.data)

	writer := bufio.NewWriter(file)

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return err
	}

	return writer.Flush()
}

func (u *Urls) Read(filepath string) error {
	file, err := os.Open(filepath + "filestorage.gob")

	if file != nil {
		defer file.Close()
	}

	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&u.data)
}
