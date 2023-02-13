// Application run
package main

import (
	"errors"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/alaleks/shortener/internal/app/serv"
)

var buildVersion string
var buildDate string
var buildCommit string

func info() {
	io.WriteString(os.Stdout, "Run application Shortener\n")

	if len(buildVersion) != 0 {
		io.WriteString(os.Stdout, "Build version: "+buildVersion+"\n")
	} else {
		io.WriteString(os.Stdout, "Build version: N/A\n")
	}

	if len(buildDate) != 0 {
		io.WriteString(os.Stdout, "Build date: "+buildDate+"\n")
	} else {
		io.WriteString(os.Stdout, "Build date: N/A\n")
	}

	if len(buildCommit) != 0 {
		io.WriteString(os.Stdout, "Build commit: "+buildCommit+"\n")
	} else {
		io.WriteString(os.Stdout, "Build commit: N/A\n")
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
