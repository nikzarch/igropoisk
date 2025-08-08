package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
)

func GetConnection() *sql.DB {
	dsn := os.Getenv("DB_POSTGRES_URL")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}
