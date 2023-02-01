// Package compress implements gzip compressing
// the web server response and reading the compressed request.
package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type writerGzip struct {
	io.Writer
	http.ResponseWriter
}

// Write method that implements the interface io.Writer.
func (w writerGzip) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("failed write gzip: %w", err)
	}

	return n, nil
}

type readerCloserGzip struct {
	*gzip.Reader
	io.Closer
}

// Close method that implements the interface io.Closer.
func (r readerCloserGzip) Close() error {
	if err := r.Closer.Close(); err != nil {
		return fmt.Errorf("failed readerCloserGzip: %w", err)
	}

	return nil
}

// Compression compresses the server response if the client supports
// gzip reading and the content type is suitable for efficient compression
func Compression(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(writer, req)

			return
		}

		// Checking the compress type of its content for efficient compression.
		if !CheckBeforeCompression(req.Header.Get("Content-Type")) {
			handler.ServeHTTP(writer, req)

			return
		}

		writer.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(writer)
		defer gz.Close()
		gzw := writerGzip{Writer: gz, ResponseWriter: writer}
		handler.ServeHTTP(gzw, req)
	})
}

// CheckBeforeCompression (Deprecated) checks the type of the content to be compressed.
func CheckBeforeCompressionOld(contentType string) bool {
	correctTypes := [...]string{
		"text/css",
		"text/csv",
		"text/html",
		"application/json",
		"text/javascript",
		"image/svg+xml",
		"font/ttf",
		"text/plain",
		"font/woff",
		"font/woff2",
		"text/xml",
	}

	for _, correctType := range correctTypes {
		if strings.Contains(contentType, correctType) {
			return true
		}
	}

	return false
}

// CheckBeforeCompression checks the type of the content to be compressed.
func CheckBeforeCompression(contentType string) bool {
	correctTypes := [...]string{
		"text/css",
		"text/csv",
		"text/html",
		"application/json",
		"text/javascript",
		"image/svg+xml",
		"font/ttf",
		"text/plain",
		"font/woff",
		"font/woff2",
		"text/xml",
	}

	for _, correctType := range correctTypes {
		if contentType == correctType {
			return true
		}
	}

	return false
}

// Unpacking implements unpacking of a gzip-compressed request.
func Unpacking(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		switch req.Header.Get("Content-Encoding") {
		case "gzip":
			var buffer bytes.Buffer
			if _, err := io.Copy(&buffer, req.Body); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)

				return
			}
			req.Header.Del("Content-Length")
			reader, err := gzip.NewReader(&buffer)
			if err != nil {
				handler.ServeHTTP(writer, req)

				return
			}

			if reader != nil {
				defer reader.Close()
			}

			req.Body = readerCloserGzip{
				Reader: reader,
				Closer: req.Body,
			}

			handler.ServeHTTP(writer, req)

			return

		default:
			handler.ServeHTTP(writer, req)

			return
		}
	})
}
