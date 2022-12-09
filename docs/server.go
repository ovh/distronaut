package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/", fs)

	go func() {
		if err := http.ListenAndServe(":3000", nil); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
