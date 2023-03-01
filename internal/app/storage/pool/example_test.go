package pool_test

import (
	"fmt"
	"time"

	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/storage/pool"
)

func Example() {
	// Initialize the logger
	logger := logger.NewLogger()

	// Initialize the pool
	pool := pool.Init(logger)

	// Run pool
	go pool.Run()

	// Add Task
	go pool.AddTask(2, func(data any) error {
		fmt.Println(data.(int) * 2)
		return nil
	})

	go pool.AddTask(4, func(data any) error {
		fmt.Println(data.(int) * 2)
		return nil
	})

	time.Sleep(time.Second)

	// Stop pool
	pool.Stop()

	// Output:
	// 4
	// 8
}
