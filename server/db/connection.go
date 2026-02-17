package db

import (
	"context"
	"fmt"
	"log" // Added log import
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/go-redis/redis/v8"
)

func Connect() (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_URL")
	log.Printf("Attempting to connect to PostgreSQL with DATABASE_URL: %s", dbURL) // Debug log
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1);
	}
	return conn, nil;
}

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	log.Printf("Attempting to connect to Redis with REDIS_URL: %s", redisURL) // Debug log
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0" // Default Redis URL
		log.Printf("REDIS_URL not set, using default: %s", redisURL) // Log default
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)

	// Ping to check if connection is established
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}
	log.Printf("Successfully connected to Redis.") // Success log
	return rdb, nil
}
