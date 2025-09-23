package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB - connection to the database
func ConnectDB() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to conect to DB: %w", err)
	}

	// connection check
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Cannot ping DB: %w", err)
	}

	fmt.Println("Connected to PostgreSQL (db.go)")
	return pool, nil
}


