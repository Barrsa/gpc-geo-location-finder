// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gpc-geo-location-finder/handlers"
	"github.com/gpc-geo-location-finder/models"
)

var (
	port    = flag.String("port", "", "Port to listen on (default: PORT env var or 8080)")
	timeout = flag.Duration("timeout", 30*time.Second, "Timeout for ping requests")
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Ignore error if .env file doesn't exist
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	flag.Parse()

	// Determine port
	listenPort := *port
	if listenPort == "" {
		listenPort = os.Getenv("PORT")
		if listenPort == "" {
			listenPort = "8080"
		}
	}

	log.Printf("Starting checkPing service on port %s", listenPort)

	// Load endpoints from JSON file
	endpoints := models.AllEndpoints()
	log.Printf("Loaded %d endpoints", len(endpoints))

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Main checkPing endpoint
	router.GET("/api/checkPing", handlers.HandleCheckPing(handlers.HandlerOptions{
		Endpoints: endpoints,
		Client:    client,
		Timeout:   *timeout,
	}))

	// Create server
	server := &http.Server{
		Addr:         ":" + listenPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Start server
	log.Printf("Server listening on :%s", listenPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server stopped")
}
