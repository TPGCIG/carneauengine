package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func NewHandler(pool *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{DB: pool, Redis: rdb}
}
