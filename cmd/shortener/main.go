// Application run
package main

import (
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/alaleks/shortener/internal/app/serv"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func info() {
	fmt.Printf("Run application Shortener\n")

	if len(buildVersion) != 0 {
		fmt.Printf("Build version: %s\n", buildVersion)
	} else {
		fmt.Printf("Build version: N/A\n")
	}

	if len(buildDate) != 0 {
		fmt.Printf("Build date: %s\n", buildDate)
	} else {
		fmt.Printf("Build date: N/A\n")
	}

	if len(buildCommit) != 0 {
		fmt.Printf("Build commit: %s\n", buildCommit)
	} else {
		fmt.Printf("Build commit: N/A\n")
	}
}

func main() {
	info()
	server := serv.New()

	// Run server for pprof
	go func() {
		server.Logger.LZ.Fatal(http.ListenAndServe(":3031", nil))
	}()

	if err := serv.Run(server); !errors.Is(err, http.ErrServerClosed) {
		server.Logger.LZ.Fatal(err)
	}
}
