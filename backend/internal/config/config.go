// backend/internal/config/config.go
package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type Config struct {
	Port              string
	ModelsPath        string
	UploadsPath       string
	TestDocumentsPath string // Frontend'den yüklenen dokümanlar için
	DatabasePath      string
	OllamaURL         string
	MaxFileSize       int64
	AllowedTypes      []string
	// Llama specific settings
	LlamaModelPath   string
	LlamaContextSize int
	LlamaThreads     int
	LlamaGPULayers   int
}

func Load() *Config {
	// Get port from environment or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Default paths
	homeDir, _ := os.UserHomeDir()
	appDir := filepath.Join(homeDir, ".local-ai-project")

	// Get database path from environment or default
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(appDir, "data", "app.db")
	}

	// Create directories if they don't exist
	os.MkdirAll(filepath.Join(appDir, "models"), 0755)
	os.MkdirAll(filepath.Join(appDir, "uploads"), 0755)
	os.MkdirAll(filepath.Join(appDir, "test_documents"), 0755) // Test dokümanları için
	os.MkdirAll(filepath.Join(appDir, "data"), 0755)

	// Auto-detect number of threads
	threads := runtime.NumCPU()
	if threads > 8 {
		threads = 8 // Limit to 8 threads for stability
	}

	return &Config{
		Port:              port,
		ModelsPath:        filepath.Join(appDir, "models"),
		UploadsPath:       filepath.Join(appDir, "uploads"),
		TestDocumentsPath: filepath.Join(appDir, "test_documents"), // Frontend dokümanları
		DatabasePath:      dbPath,
		OllamaURL:         getEnv("OLLAMA_URL", "http://localhost:11434"),
		MaxFileSize:       50 * 1024 * 1024, // 50MB
		AllowedTypes:      []string{".pdf", ".txt", ".docx", ".md"},
		// Llama settings
		LlamaModelPath:   filepath.Join(appDir, "models"),
		LlamaContextSize: getEnvInt("LLAMA_CONTEXT_SIZE", 2048),
		LlamaThreads:     getEnvInt("LLAMA_THREADS", threads),
		LlamaGPULayers:   getEnvInt("LLAMA_GPU_LAYERS", 0), // 0 = CPU only
	}
}

func NewConfig() *Config {
	return Load()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
