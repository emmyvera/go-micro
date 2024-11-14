package main

import (
	"fmt"
	"log"
	"net/http"
)

const WEB_PORT = "80"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s", WEB_PORT)

	// Set up http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	// start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
