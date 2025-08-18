package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"time"
)

func GetPool() *pgxpool.Pool {
	dsn := os.Getenv("DB_POSTGRES_URL")
	if dsn == "" {
		log.Fatal("DB_POSTGRES_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	return pool
}
