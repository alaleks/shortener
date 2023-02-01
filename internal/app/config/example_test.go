package config_test

import (
	"fmt"

	"github.com/alaleks/shortener/internal/app/config"
)

func Example() {
	// Initialize application settings
	var appConf config.Configurator = config.New(
		config.Options{
			Env:  true, // checking environment variable settings
			Flag: true, // checking flags settings
		})

	// Output host (server addr)
	fmt.Println(appConf.GetServAddr())

	// Output:
	// localhost:8080
}
