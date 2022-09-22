package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	fs := http.FileServer(http.Dir("./tests/static"))
	http.Handle("/", fs)

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ready"))
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Stopping mock server"))
		cancel()
	})

	go func() {
		if err := http.ListenAndServe(":3000", nil); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
