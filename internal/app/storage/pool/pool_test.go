package pool_test

import (
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/storage/pool"
)

func TestPool(t *testing.T) {
	t.Parallel()
	defer os.Remove("log.json")
	count := 100
	var test []int

	logger := logger.NewLogger()
	pool := pool.Init(logger)
	defer pool.Stop()
	go pool.Run()

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			pool.AddTask(i, func(data any) error {
				defer wg.Done()
				test = append(test, i)
				return nil
			})
		}(i)
	}
	wg.Wait()

	if count != len(test) {
		t.Errorf("pool lost %d items", count-len(test))
	}
}

func BenchmarkPool(b *testing.B) {
	defer os.Remove("log.json")
	logger := logger.NewLogger()
	pool := pool.Init(logger)
	go pool.Run()
	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pool.AddTask(i, func(data any) error {
				_ = i * i
				return nil
			})
		}
	})
	pool.SetNumWorker(10 * runtime.NumCPU())
	b.Run("after optimize 10 x NumCPU", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pool.AddTask(i, func(data any) error {
				_ = i * i
				return nil
			})
		}
	})
	pool.SetNumWorker(100 * runtime.NumCPU())
	b.Run("after optimize 100 x NumCPU", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pool.AddTask(i, func(data any) error {
				_ = i * i
				return nil
			})
		}
	})
	pool.SetNumWorker(1000 * runtime.NumCPU())
	b.Run("after optimize 1000 x NumCPU", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pool.AddTask(i, func(data any) error {
				_ = i * i
				return nil
			})
		}
	})
	pool.Stop()
}
