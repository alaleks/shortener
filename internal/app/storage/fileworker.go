package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

// Close perfoms saving data to the file storage and closes its.
func (ds *DefaultStorage) Close() error {
	if ds.conf.GetFileStoragePath() == "" {
		return nil
	}

	file, err := os.Create(correctorFilename(ds.conf.GetFileStoragePath()))
	if err != nil {
		return fmt.Errorf("failed create file storage: %w", err)
	}

	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	if err = encoder.Encode(ds.urls); err != nil {
		return fmt.Errorf("failed encode data for file storage: %w", err)
	}

	writer := bufio.NewWriter(file)

	if _, err = writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed write data to buffer: %w", err)
	}

	err = writer.Flush()

	if err != nil {
		return fmt.Errorf("failed to dump data to file: %w", err)
	}

	return nil
}

// Init performs that initializes the file storage.
func (ds *DefaultStorage) Init() error {
	if ds.conf.GetFileStoragePath() == "" {
		return nil
	}

	file, err := os.Open(correctorFilename(ds.conf.GetFileStoragePath()))
	if err != nil {
		return fmt.Errorf("failed open file storage: %w", err)
	}

	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}

	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&ds.urls)

	if err != nil {
		return fmt.Errorf("failed decode data: %w", err)
	}

	return nil
}

// Ping is a stub method for interface implementation Storage.
func (ds *DefaultStorage) Ping() error {
	return nil
}

func correctorFilename(filepath string) string {
	if strings.HasSuffix(filepath, "/") {
		filepath += "storage"
	}

	if stat, err := os.Stat(filepath); err == nil && stat.IsDir() {
		filepath += "/storage"
	}

	return filepath
}
