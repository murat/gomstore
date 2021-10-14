package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
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

	backupFile := path.Join(os.TempDir(), "gomstore.json")

	// load initial data from backup file
	if _, err := os.Stat(backupFile); !errors.Is(err, os.ErrNotExist) {
		err := api.store.Load(backupFile)
		if err != nil {
			log.Printf("[error] %v\n", err)
		}
	}

	// start periodic backup
	api.store.PeriodicBackup(backupFile, 1)
	defer func(store Store, filePath string) {
		if err := store.Save(filePath); err != nil {
			log.Printf("[error] %v\n", err)
		}
	}(api.store, backupFile)

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
