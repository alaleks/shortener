package logger_test

import "github.com/alaleks/shortener/internal/app/logger"

func Example() {
	// Initialize the logger
	logger := logger.NewLogger()

	// Write info log
	logger.LZ.Info("test")
}
