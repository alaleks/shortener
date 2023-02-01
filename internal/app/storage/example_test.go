package storage_test

import (
	"fmt"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/storage"
)

func Example() {
	// Initialize application settings
	var appConf config.Configurator = config.New(
		config.Options{})

	// Init Storage
	st := storage.InitStore(appConf, logger.NewLogger())

	// Check Ping
	fmt.Println(st.Store.Ping())

	// Output:
	// <nil>
}
