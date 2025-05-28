// backend/internal/storage/database.go
package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS models (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			path TEXT NOT NULL,
			size BIGINT,
			status TEXT DEFAULT 'downloaded',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS documents (
			id SERIAL PRIMARY KEY,
			filename TEXT NOT NULL,
			original_name TEXT NOT NULL,
			path TEXT NOT NULL,
			size BIGINT,
			type TEXT,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS document_chunks (
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
			return err
		}
	}

	return nil
}
