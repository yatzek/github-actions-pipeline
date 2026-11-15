package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gist-service/client"
	"gist-service/handler"
)

func main() {
	log.Printf("Starting Gist API server")

	githubClient, err := client.NewGitHubClient()
	if err != nil {
		log.Fatalf("Failed to initialize GitHub client: %v", err)
	}
	server := handler.NewServer(githubClient)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	sig := <-stopChan
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Printf("Server exited gracefully")
	}
}
