package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/tsuru/autoscalev2/pkg/web"
)

var httpBindAddress string

func main() {
	flag.StringVar(&httpBindAddress, "http-bind-address", ":8081", "The TCP address that the web API should bind to")
	flag.Parse()

	// Web API stuff
	var server web.Server

	fmt.Printf("Starting web server at %s\n", httpBindAddress)
	err := server.Start(httpBindAddress)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "failed to start web server: %s\n", err)
		os.Exit(1)
	}
}
