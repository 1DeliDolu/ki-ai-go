package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitPostgresDB(connectionString string) (*sql.DB, error) {
	// If no connection string provided, use default local credentials
	if connectionString == "" {
		connectionString = "host=localhost port=5432 dbname=local_ai user=postgres password=D0cker sslmode=disable"
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createPostgresTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createPostgresTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`, `CREATE TABLE IF NOT EXISTS prompts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER,
			prompt_text TEXT NOT NULL,
			answer_text TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (user_id)
		)`, `CREATE TABLE IF NOT EXISTS models (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			path TEXT NOT NULL,
			size BIGINT,
			status TEXT DEFAULT 'downloaded',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`, `CREATE TABLE IF NOT EXISTS documents (
			id SERIAL PRIMARY KEY,
			filename TEXT NOT NULL,
			original_name TEXT NOT NULL,
			path TEXT NOT NULL,
			size BIGINT,
			type TEXT,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`, `CREATE TABLE IF NOT EXISTS document_chunks (
			id SERIAL PRIMARY KEY,
			document_id INTEGER,
			content TEXT NOT NULL,
			embedding BYTEA,
			chunk_index INTEGER,
			FOREIGN KEY (document_id) REFERENCES documents (id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}
