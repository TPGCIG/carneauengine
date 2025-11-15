package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/tpgcig/carneauengine/server/db"
	"github.com/tpgcig/carneauengine/server/handlers"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer conn.Close(context.Background());

	h := handlers.NewHandler(conn);

	r.GET("/a", h.GetEvent)
	r.GET("/b", h.AddEvent)
	r.GET("/c", h.GetUser)
	r.GET("/d", h.AddUser)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
