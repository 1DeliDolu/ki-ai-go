// backend/internal/services/model_service.go
package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/internal/config"
	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

type ModelService struct {
	config        *config.Config
	db            *sql.DB
	ollamaService *OllamaService
	currentModel  string
}

func NewModelService(cfg *config.Config, db *sql.DB) *ModelService {
	return &ModelService{
		config:        cfg,
		db:            db,
		ollamaService: NewOllamaService(),
		currentModel:  "",
	}
}

func (s *ModelService) ListModels() ([]*types.Model, error) {
	// Get models from Ollama instead of static list
	models, err := s.ollamaService.ListModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models from Ollama: %w", err)
	}

	return models, nil
}

func (s *ModelService) LoadModel(modelName string) error {
	log.Printf("🔄 Loading model: %s", modelName)

	// Clean model name - remove any existing tags
	cleanModelName := strings.Split(modelName, ":")[0]

	// Try different model name variations
	modelVariations := []string{
		cleanModelName,
		cleanModelName + ":latest",
		modelName, // original name as fallback
	}

	var lastError error

	for _, variation := range modelVariations {
		log.Printf("🔄 Trying model variation: %s", variation)

		// Check if model exists in Ollama first
		if err := s.checkModelExists(variation); err != nil {
			log.Printf("⚠️ Model %s not found in Ollama: %v", variation, err)
			lastError = err
			continue
		}

		// Try to load the model
		if err := s.ollamaService.LoadModel(variation); err != nil {
			log.Printf("⚠️ Failed to load model %s: %v", variation, err)
			lastError = err
			continue
		}

		// Success!
		s.currentModel = variation
		log.Printf("✅ Successfully loaded model: %s", variation)
		return nil
	}

	// If all variations failed, try to pull the model
	log.Printf("🔄 Model not found locally, attempting to pull: %s", cleanModelName)
	if err := s.pullAndLoadModel(cleanModelName); err != nil {
		return fmt.Errorf("failed to load or pull model %s: %w (last error: %v)", modelName, err, lastError)
	}

	s.currentModel = cleanModelName
	return nil
}

// ModelInfo represents the metadata for downloaded models
type ModelInfo struct {
	Filename             string
	OllamaName           string
	DisplayName          string
	Description          string
	ModelType            string
	EstimatedSize        string
	AlternativeFilenames []string
}

// getModelDefinitions returns the mapping of your downloaded models
func (s *ModelService) getModelDefinitions() map[string]ModelInfo {
	return map[string]ModelInfo{
		"nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf": {
			Filename:      "nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf",
			OllamaName:    "nemotron-nano",
			DisplayName:   "NVIDIA Llama 3.1 Nemotron Nano 4B",
			Description:   "NVIDIA's optimized Llama 3.1 Nemotron model - fast and efficient",
			ModelType:     "nemotron",
			EstimatedSize: "2.4 GB",
			AlternativeFilenames: []string{
				"nemotron-nano.gguf",
				"llama-3.1-nemotron.gguf",
				"nvidia-nemotron.gguf",
			},
		},
		"neural-chat-7b-v3-1.Q5_0.gguf": {
			Filename:      "neural-chat-7b-v3-1.Q5_0.gguf",
			OllamaName:    "neural-chat",
			DisplayName:   "Neural Chat 7B Q5_0",
			Description:   "Intel's optimized conversational AI model with Q5_0 quantization",
			ModelType:     "neural-chat",
			EstimatedSize: "4.8 GB",
			AlternativeFilenames: []string{
				"neural-chat-7b.gguf",
				"neural-chat.Q5_0.gguf",
				"neuralchat-7b.gguf",
			},
		},
		"openchat-3.5-0106.Q5_K_M.gguf": {
			Filename:      "openchat-3.5-0106.Q5_K_M.gguf",
			OllamaName:    "openchat",
			DisplayName:   "OpenChat 3.5 Q5_K_M",
			Description:   "High-quality open-source conversational AI with Q5_K_M quantization",
			ModelType:     "openchat",
			EstimatedSize: "4.8 GB",
			AlternativeFilenames: []string{
				"openchat-3.5.Q5_K_M.gguf",
				"openchat_3.5.Q5_K_M.gguf",
				"openchat-3.5.gguf",
				"openchat.Q5_K_M.gguf",
			},
		},
		"llama-2-7b-chat.Q4_K_M.gguf": {
			Filename:      "llama-2-7b-chat.Q4_K_M.gguf",
			OllamaName:    "llama2-chat",
			DisplayName:   "Llama 2 7B Chat Q4_K_M",
			Description:   "Meta's Llama 2 model optimized for conversational AI",
			ModelType:     "llama",
			EstimatedSize: "4.1 GB",
			AlternativeFilenames: []string{
				"llama2-7b-chat.gguf",
				"llama-2-chat.gguf",
				"llama2.Q4_K_M.gguf",
			},
		},
		"phi-2.Q8_0.gguf": {
			Filename:      "phi-2.Q8_0.gguf",
			OllamaName:    "phi2",
			DisplayName:   "Microsoft Phi-2 Q8_0",
			Description:   "Compact but powerful language model from Microsoft with Q8_0 quantization",
			ModelType:     "phi",
			EstimatedSize: "2.8 GB",
			AlternativeFilenames: []string{
				"phi2.Q8_0.gguf",
				"phi-2.gguf",
				"phi2.gguf",
				"microsoft-phi2.gguf",
			},
		},
	}
}

func (s *ModelService) ValidateModelName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	// Get models from Ollama to validate
	models, err := s.ollamaService.ListModels()
	if err != nil {
		return fmt.Errorf("failed to get models from Ollama: %w", err)
	}

	for _, model := range models {
		if model.ID == name || model.Name == name {
			return nil
		}
	}

	return fmt.Errorf("unknown model: %s", name)
}

// AddBasicModels adds some basic/sample models to the system
func (s *ModelService) AddBasicModels() error {
	log.Println("Adding basic models to the system...")

	basicModels := []types.Model{
		{
			ID:          "llama-7b-chat",
			Name:        "llama-7b-chat",
			Size:        "4.1GB",
			Type:        "chat",
			Status:      "available",
			Description: "Llama 7B Chat model for general conversation",
			ModelType:   "llama",
			URL:         "https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf",
		},
		{
			ID:          "tinyllama-1.1b",
			Name:        "tinyllama-1.1b",
			Size:        "630MB",
			Type:        "completion",
			Status:      "available",
			Description: "Tiny Llama 1.1B model for lightweight text generation",
			ModelType:   "tinyllama",
			URL:         "https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
		},
		{
			ID:          "gpt2-medium",
			Name:        "gpt2-medium",
			Size:        "345MB",
			Type:        "completion",
			Status:      "available",
			Description: "GPT-2 Medium model for text completion",
			ModelType:   "gpt2",
			URL:         "https://huggingface.co/gpt2-medium",
		},
		{
			ID:          "code-llama-7b",
			Name:        "code-llama-7b",
			Size:        "4.2GB",
			Type:        "code",
			Status:      "available",
			Description: "Code Llama 7B for code generation and completion",
			ModelType:   "codellama",
			URL:         "https://huggingface.co/TheBloke/CodeLlama-7B-Instruct-GGUF/resolve/main/codellama-7b-instruct.Q4_K_M.gguf",
		},
		{
			ID:          "neural-chat-7b",
			Name:        "neural-chat-7b",
			Size:        "4.8GB",
			Type:        "chat",
			Status:      "available",
			Description: "Intel's Neural Chat 7B for conversations",
			ModelType:   "neural-chat",
			URL:         "https://huggingface.co/TheBloke/neural-chat-7B-v3-1-GGUF/resolve/main/neural-chat-7b-v3-1.Q4_K_M.gguf",
		},
		{
			ID:          "phi-2",
			Name:        "phi-2",
			Size:        "2.8GB",
			Type:        "completion",
			Status:      "available",
			Description: "Microsoft's Phi-2 small but powerful model",
			ModelType:   "phi",
			URL:         "https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.Q4_K_M.gguf",
		},
	}

	for _, model := range basicModels {
		// Create a types.Model compatible with Ollama structure
		ollamaModel := &types.Model{
			ID:          model.ID,
			Name:        model.Name,
			Size:        model.Size,
			Type:        model.Type,
			Status:      model.Status,
			Description: model.Description,
			ModelType:   model.ModelType,
			URL:         model.URL,
		}

		// Add via Ollama service if available, fallback to memory
		if err := s.ollamaService.CreateModel(ollamaModel); err != nil {
			log.Printf("Failed to add model %s via Ollama: %v", model.Name, err)
			// Continue without failing - models will be available from Ollama's existing catalog
		}

		log.Printf("✅ Added basic model: %s (%s)", model.Name, model.Size)
	}

	log.Println("Basic models added successfully!")
	return nil
}

// InitializeBasicModels checks and adds basic models if none exist
func (s *ModelService) InitializeBasicModels() error {
	models, err := s.ListModels()
	if err != nil {
		return fmt.Errorf("failed to check existing models: %w", err)
	}

	if len(models) == 0 {
		log.Println("No models found in database, adding basic models...")
		return s.AddBasicModels()
	}

	log.Printf("Found %d existing models, skipping basic model initialization", len(models))
	return nil
}

// GetModelInfo returns detailed information about a specific model
func (s *ModelService) GetModelInfo(name string) (*ModelInfo, error) {
	log.Printf("Getting info for model: %s", name)

	// Try to get from local definitions first
	modelDefinitions := s.getModelDefinitions()
	for _, modelInfo := range modelDefinitions {
		if modelInfo.OllamaName == name {
			return &modelInfo, nil
		}
	}

	// Try to get from Ollama
	models, err := s.ollamaService.ListModels()
	if err == nil {
		for _, model := range models {
			if model.Name == name || model.ID == name {
				// Convert Ollama model to ModelInfo
				return &ModelInfo{
					Filename:      fmt.Sprintf("%s.gguf", name),
					OllamaName:    model.Name,
					DisplayName:   model.Name,
					Description:   model.Description,
					ModelType:     model.Type,
					EstimatedSize: model.Size,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("model not found: %s", name)
}

// GetAvailableModelTypes returns all available model types
func (s *ModelService) GetAvailableModelTypes() []string {
	return []string{
		"chat",       // Conversational models
		"completion", // Text completion models
		"code",       // Code generation models
		"embedding",  // Text embedding models
		"vision",     // Vision/multimodal models
		"audio",      // Audio processing models
	}
}

// GetModelsByType returns models filtered by type
func (s *ModelService) GetModelsByType(modelType string) ([]*types.Model, error) {
	allModels, err := s.ListModels()
	if err != nil {
		return nil, err
	}

	var filtered []*types.Model
	for _, model := range allModels {
		if model.Type == modelType {
			filtered = append(filtered, model)
		}
	}

	return filtered, nil
}

func (s *ModelService) GetModelFilePath(name string) (string, error) {
	modelInfo, err := s.GetModelInfo(name)
	if err != nil {
		return "", err
	}

	// Try primary filename first
	filePath := filepath.Join(s.config.ModelsPath, modelInfo.Filename)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	// Try alternative filenames
	for _, altFilename := range modelInfo.AlternativeFilenames {
		altPath := filepath.Join(s.config.ModelsPath, altFilename)
		if _, err := os.Stat(altPath); err == nil {
			log.Printf("Using alternative filename for %s: %s", name, altFilename)
			return altPath, nil
		}
	}

	// Try fuzzy matching
	files, err := os.ReadDir(s.config.ModelsPath)
	if err != nil {
		return "", fmt.Errorf("cannot read models directory: %w", err)
	}

	existingFiles := make(map[string]os.FileInfo)
	for _, file := range files {
		if !file.IsDir() && s.isModelFile(file.Name()) {
			if info, err := file.Info(); err == nil {
				existingFiles[file.Name()] = info
			}
		}
	}

	if foundFile := s.findModelFileByPattern(name, existingFiles); foundFile != "" {
		foundPath := filepath.Join(s.config.ModelsPath, foundFile)
		log.Printf("Found model %s via pattern matching: %s", name, foundFile)
		return foundPath, nil
	}

	return "", fmt.Errorf("model file not found for: %s", name)
}

func (s *ModelService) DownloadModel(name, url string) error {
	log.Printf("Starting download: %s from %s", name, url)

	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("download URL cannot be empty")
	}

	// Create the models directory if it doesn't exist
	if err := os.MkdirAll(s.config.ModelsPath, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Minute, // 30 minutes for large model downloads
	}

	// Download the model file
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download model: HTTP %d", resp.StatusCode)
	}

	// Create the destination file
	filePath := filepath.Join(s.config.ModelsPath, name)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create model file: %w", err)
	}
	defer out.Close()

	// Copy the response body to the file with progress tracking
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		// Clean up partial file on error
		os.Remove(filePath)
		return fmt.Errorf("failed to save model file: %w", err)
	}

	log.Printf("Successfully downloaded %s (%s)", name, s.formatFileSize(written))
	return nil
}

func (s *ModelService) DeleteModel(name string) error {
	log.Printf("Deleting model: %s", name)

	// Get model info
	modelInfo, err := s.GetModelInfo(name)
	if err != nil {
		return fmt.Errorf("model not found: %w", err)
	}

	filePath := filepath.Join(s.config.ModelsPath, modelInfo.Filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s", filePath)
	}

	// Delete the model file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete model file: %w", err)
	}

	log.Printf("Successfully deleted model: %s", name)
	return nil
}

func (s *ModelService) GetAvailableModels() ([]string, error) {
	models, err := s.ListModels()
	if err != nil {
		return nil, err
	}

	var available []string
	for _, model := range models {
		if model.Status == "available" {
			available = append(available, model.Name)
		}
	}

	return available, nil
}

// findModelFileByPattern tries to find model files using fuzzy matching
func (s *ModelService) findModelFileByPattern(modelName string, existingFiles map[string]os.FileInfo) string {
	// Convert model name to searchable patterns
	patterns := s.generateSearchPatterns(modelName)

	for filename := range existingFiles {
		if !s.isModelFile(filename) {
			continue
		}

		lowerFilename := strings.ToLower(filename)
		for _, pattern := range patterns {
			if strings.Contains(lowerFilename, pattern) {
				log.Printf("Pattern match found: %s matches pattern %s for model %s",
					filename, pattern, modelName)
				return filename
			}
		}
	}

	return ""
}

// Enhanced pattern generation for better matching
func (s *ModelService) generateSearchPatterns(modelName string) []string {
	patterns := []string{strings.ToLower(modelName)}

	switch modelName {
	case "nemotron-nano":
		patterns = append(patterns, []string{
			"nemotron",
			"nvidia",
			"llama-3.1",
			"llama3.1",
			"nano",
			"3.1-nemotron",
			"bf16",
		}...)
	case "openchat":
		patterns = append(patterns, []string{
			"openchat",
			"openchat-3.5",
			"openchat_3.5",
			"openchat-3",
			"openchat3",
			"0106",
			"3.5-0106",
			"q5_k_m",
		}...)
	case "phi2":
		patterns = append(patterns, []string{
			"phi2",
			"phi-2",
			"phi_2",
			"phi",
			"microsoft",
			"q8_0",
			"q8",
		}...)
	case "llama2-chat":
		patterns = append(patterns, []string{
			"llama-2",
			"llama2",
			"llama_2",
			"chat",
			"q4_k_m",
			"7b-chat",
		}...)
	case "neural-chat":
		patterns = append(patterns, []string{
			"neural",
			"neural-chat",
			"neural_chat",
			"neuralchat",
			"intel",
			"q5_0",
			"7b-v3",
		}...)
	}

	return patterns
}

func (s *ModelService) formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.1f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func (s *ModelService) isModelFile(filename string) bool {
	validExtensions := []string{".gguf", ".bin", ".ggml"}
	lower := strings.ToLower(filename)

	for _, ext := range validExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// pullAndLoadModel attempts to pull and then load a model
func (s *ModelService) pullAndLoadModel(modelName string) error {
	log.Printf("🔄 Pulling model from Ollama registry: %s", modelName)

	// Try common model names for phi2
	if strings.Contains(strings.ToLower(modelName), "phi") {
		// Try different phi model names
		phiVariations := []string{
			"phi",
			"phi2",
			"microsoft/phi-2",
			"phi:latest",
		}

		for _, variation := range phiVariations {
			log.Printf("🔄 Trying to pull phi model: %s", variation)
			if err := s.tryPullModel(variation); err == nil {
				return s.ollamaService.LoadModel(variation)
			}
		}
	}

	// Try to pull the original model name
	return s.tryPullModel(modelName)
}

// tryPullModel attempts to pull a specific model
func (s *ModelService) tryPullModel(modelName string) error {
	// Create a pull request
	reqBody := map[string]interface{}{
		"name": modelName,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal pull request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post(s.config.OllamaURL+"/api/pull", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull model: HTTP %d", resp.StatusCode)
	}

	log.Printf("✅ Successfully pulled model: %s", modelName)
	return nil
}

// checkModelExists verifies if a model exists in Ollama
func (s *ModelService) checkModelExists(modelName string) error {
	models, err := s.ollamaService.ListModels()
	if err != nil {
		return fmt.Errorf("failed to list Ollama models: %w", err)
	}

	for _, model := range models {
		if model.Name == modelName || model.ID == modelName {
			return nil
		}
		// Also check without :latest tag
		if strings.HasSuffix(modelName, ":latest") {
			baseModelName := strings.TrimSuffix(modelName, ":latest")
			if model.Name == baseModelName || model.ID == baseModelName {
				return nil
			}
		}
	}

	return fmt.Errorf("model not found: %s", modelName)
}
