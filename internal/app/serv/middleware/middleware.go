package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (c *compressWriter) Write(b []byte) (int, error) {
	n, err := c.Writer.Write(b)

	return n, fmt.Errorf("failed writing to byte slice: %w", err)
}

func CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(reswr http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(reswr, req)

			return
		}

		compress, err := gzip.NewWriterLevel(reswr, gzip.DefaultCompression)
		if err != nil {
			if _, err := io.WriteString(reswr, err.Error()); err != nil {
				http.Error(reswr, err.Error(), http.StatusBadRequest)

				return
			}

			return
		}

		if compress != nil {
			defer compress.Close()
		}

		reswr.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(&compressWriter{ResponseWriter: reswr, Writer: compress}, req)
	})
}
