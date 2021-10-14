package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	api := &api{
		store: NewStore(),
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(api.Serve),
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	log.Println("[info] Server is starting...")
	go func() {
		<-quit
		log.Println("[info] Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("[error] Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("[info] Server is ready. Listening on :8080")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[error] could not start server, %v", err)
	}

	<-done
	log.Println("[info] Server stopped")
}
