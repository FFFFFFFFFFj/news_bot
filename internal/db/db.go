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

	//Create tables if they don't exist
	if err := createTables(pool); err != nil {
		return nil, err
	}
	
	return pool, nil
}

//Create tables
func createTables(pool *pgxpool.Pool) error {
	ctx := context.Background()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				telegram_id BIGINT UNIQUE NOT NULL,
				username TEXT,
				role TEXT NOT NULL,
				timezone TEXT,
				notification_time TIME,
				notification_count INT DEFAULT 5
			);`,
		
		`CREATE TABLE IF NOT EXISTS sources (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				url TEXT NOT NULL,
				category TEXT
			);`,

		`CREATE TABLE IF NOT EXISTS news (
				id SERIAL PRIMARY KEY,
				source_id INT REFERENCES sources(id),
				title TEXT,
				link TEXT,
				published TIMESTAMP,
				category TEXT
			);`,

		`CREATE TABLE IF NOT EXISTS subscriptions (
				id SERIAL PRIMARY KEY,
				user_id INT REFERENCES users(id),
				source_id INT REFERENCES sources(id)
			);`,

		`CREATE TABLE IF NOT EXISTS admin_broadcasts (
				id SERIAL PRIMARY KEY,
				message TEXT,
				created_at TIMESTAMP DEFAULT now()
			);`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	fmt.Println("All tables created or already exist")
	return nil
}

// add a user if he doesn't exist yet
func AddUserWithRoleIfNotExists(pool *pgxpool.Pool, telegramID int64, username, role string) error {
	ctx := context.Background()

	//Default role = "user"
	_, err := pool.Exec(ctx, `
			INSERT INTO users (telegram_id, username, role)
			VALUES ($1, $2, $3)
			ON CONFLICT (telegram_id) DO NOTHING
		`, telegramID, username)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
