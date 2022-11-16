package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

type FileStorage interface {
	Write(filepath string) error
	Read(filepath string) error
}

func (u *Urls) Write(filepath string) error {
	if !strings.HasSuffix(filepath, ".gob") {
		filepath += "filestorage.gob"
	}

	file, err := os.Create(filepath)

	if file != nil {
		defer file.Close()
	}

	if err != nil {
		return fmt.Errorf("failed create file storage: %w", err)
	}

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	if err := encoder.Encode(u.data); err != nil {
		return fmt.Errorf("failed encode data for file storage: %w", err)
	}

	writer := bufio.NewWriter(file)

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed write data to buffer: %w", err)
	}

	err = writer.Flush()

	if err != nil {
		return fmt.Errorf("failed to dump data to file: %w", writer.Flush())
	}

	return nil
}

func (u *Urls) Read(filepath string) error {
	if !strings.HasSuffix(filepath, ".gob") {
		filepath += "filestorage.gob"
	}

	file, err := os.Open(filepath)

	if file != nil {
		defer file.Close()
	}

	if err != nil {
		return fmt.Errorf("failed open file storage: %w", err)
	}

	decoder := gob.NewDecoder(file)

	return fmt.Errorf("failed decode data: %w", decoder.Decode(&u.data))
}