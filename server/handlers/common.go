package handlers

import (
	"github.com/jackc/pgx/v5"
	"github.com/go-redis/redis/v8"
)

// Handler holds the database connection and provides methods for various handlers.
type Handler struct {
	DB 	*pgx.Conn
	Redis *redis.Client
}

// NewHandler creates a new Handler instance.
func NewHandler(conn *pgx.Conn, rdb *redis.Client) *Handler {
	return &Handler{DB: conn, Redis: rdb}
}
