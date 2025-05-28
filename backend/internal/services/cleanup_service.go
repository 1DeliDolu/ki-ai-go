package services

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/1DeliDolu/ki-ai-go/internal/config"
)

type CleanupService struct {
	config *config.Config
	db     *sql.DB
}

func NewCleanupService(cfg *config.Config, db *sql.DB) *CleanupService {
	return &CleanupService{
		config: cfg,
		db:     db,
	}
}

func (s *CleanupService) CleanupOnShutdown() error {
	log.Println("🧹 Starting cleanup process...")

	// Clean up uploaded documents
	if err := s.cleanupUploads(); err != nil {
		log.Printf("⚠️  Warning: Failed to cleanup uploads: %v", err)
	}

	// Clean up database
	if err := s.cleanupDatabase(); err != nil {
		log.Printf("⚠️  Warning: Failed to cleanup database: %v", err)
	}

	// Clean up temporary files
	if err := s.cleanupTempFiles(); err != nil {
		log.Printf("⚠️  Warning: Failed to cleanup temp files: %v", err)
	}

	log.Println("✅ Cleanup completed")
	return nil
}

func (s *CleanupService) cleanupUploads() error {
	log.Println("🗂️  Cleaning up uploaded documents...")

	// Check if uploads directory exists
	if _, err := os.Stat(s.config.UploadsPath); os.IsNotExist(err) {
		log.Println("📁 Uploads directory doesn't exist, skipping...")
		return nil
	}

	// Remove all files in uploads directory
	err := filepath.Walk(s.config.UploadsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == s.config.UploadsPath {
			return nil
		}

		// Remove file or directory
		if err := os.RemoveAll(path); err != nil {
			log.Printf("⚠️  Failed to remove %s: %v", path, err)
			return nil // Continue with other files
		}

		log.Printf("🗑️  Removed: %s", path)
		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("✅ Cleaned uploads directory: %s", s.config.UploadsPath)
	return nil
}

func (s *CleanupService) cleanupDatabase() error {
	log.Println("🗄️  Cleaning up database...")

	// Clear all tables
	tables := []string{
		"document_chunks",
		"documents",
		"models",
	}

	for _, table := range tables {
		if _, err := s.db.Exec("DELETE FROM " + table); err != nil {
			log.Printf("⚠️  Warning: Failed to clear table %s: %v", table, err)
			continue
		}
		log.Printf("🗑️  Cleared table: %s", table)
	}

	// Reset auto-increment counters
	for _, table := range tables {
		if _, err := s.db.Exec("DELETE FROM sqlite_sequence WHERE name = ?", table); err != nil {
			log.Printf("⚠️  Warning: Failed to reset sequence for %s: %v", table, err)
		}
	}

	log.Println("✅ Database cleanup completed")
	return nil
}

func (s *CleanupService) cleanupTempFiles() error {
	log.Println("🧹 Cleaning up temporary files...")

	tempDirs := []string{
		os.TempDir(),
		"/tmp",
		filepath.Join(s.config.UploadsPath, ".tmp"),
	}

	for _, tempDir := range tempDirs {
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			continue
		}

		// Clean up files with our app prefix
		err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue on error
			}

			// Only remove files that match our patterns
			if info.IsDir() {
				return nil
			}

			filename := info.Name()
			// Remove files with our app-specific patterns
			if filepath.Ext(filename) == ".tmp" ||
				filepath.HasPrefix(filename, "local-ai-") ||
				filepath.HasPrefix(filename, "upload-") {
				if err := os.Remove(path); err != nil {
					log.Printf("⚠️  Failed to remove temp file %s: %v", path, err)
				} else {
					log.Printf("🗑️  Removed temp file: %s", path)
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("⚠️  Warning: Error walking temp directory %s: %v", tempDir, err)
		}
	}

	log.Println("✅ Temporary files cleanup completed")
	return nil
}

// Optional: Clean up during runtime (for testing)
func (s *CleanupService) CleanupAll() error {
	return s.CleanupOnShutdown()
}

// Clean up only uploaded documents (partial cleanup)
func (s *CleanupService) CleanupDocuments() error {
	log.Println("🗂️  Cleaning up documents only...")

	if err := s.cleanupUploads(); err != nil {
		return err
	}

	// Also clear document tables
	tables := []string{"document_chunks", "documents"}
	for _, table := range tables {
		if _, err := s.db.Exec("DELETE FROM " + table); err != nil {
			log.Printf("⚠️  Warning: Failed to clear table %s: %v", table, err)
		}
	}

	return nil
}
