package compress_test

import (
	"testing"

	"github.com/alaleks/shortener/internal/app/serv/middleware/compress"
)

func BenchmarkCheckBeforeCompression(b *testing.B) {
	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = compress.CheckBeforeCompressionOld("text/xml")
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = compress.CheckBeforeCompression("text/xml")
		}
	})
}
