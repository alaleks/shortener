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
	data := 2

	// Initialize the pool
	pool := pool.Init(logger)

	// Run pool
	go pool.Run()

	// Add Task
	go pool.AddTask(func() error {
		fmt.Println(data * 2)
		return nil
	})

	go pool.AddTask(func() error {
		fmt.Println(data * 2)
		return nil
	})

	time.Sleep(time.Second)

	// Stop pool
	pool.Stop()

	// Output:
	// 4
	// 4
}
