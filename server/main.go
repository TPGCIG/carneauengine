package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v83"
	"github.com/tpgcig/carneauengine/server/db"
	"github.com/tpgcig/carneauengine/server/handlers"
)

func init() {
    err := godotenv.Load() // loads .env automatically
    if err != nil {
        log.Println("No .env file found")
    }
}

func main() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	    }))

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer conn.Close(context.Background());

	h := handlers.NewHandler(conn);

	r.GET("/api/events", h.GetSummarisedEvents)
	r.GET("/api/events/:id", h.GetEvent)
	r.POST("/api/ticketTypes", h.GetTicketTypes)
	r.POST("/create-checkout-session", h.CreateCheckoutSession)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
