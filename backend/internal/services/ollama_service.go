package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

// OllamaService handles communication with Ollama API
type OllamaService struct {
	client  *http.Client
	baseURL string
}

func NewOllamaService() *OllamaService {
	return &OllamaService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "http://localhost:11434", // Default Ollama URL
	}
}

func (s *OllamaService) ListModels() ([]*types.Model, error) {
	log.Printf("🔄 Fetching models from Ollama...")

	resp, err := s.client.Get(s.baseURL + "/api/tags")
	if err != nil {
		log.Printf("⚠️ Failed to connect to Ollama, returning fallback models: %v", err)
		return s.getFallbackModels(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ Ollama API error (HTTP %d), returning fallback models", resp.StatusCode)
		return s.getFallbackModels(), nil
	}

	var response struct {
		Models []struct {
			Name       string    `json:"name"`
			Model      string    `json:"model"`
			Size       int64     `json:"size"`
			Digest     string    `json:"digest"`
			ModifiedAt time.Time `json:"modified_at"`
			Details    struct {
				Format            string   `json:"format"`
				Family            string   `json:"family"`
				Families          []string `json:"families"`
				ParameterSize     string   `json:"parameter_size"`
				QuantizationLevel string   `json:"quantization_level"`
			} `json:"details"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("⚠️ Failed to decode Ollama response, returning fallback models: %v", err)
		return s.getFallbackModels(), nil
	}

	var models []*types.Model
	for _, model := range response.Models {
		// Clean model name
		name := model.Name
		if strings.Contains(name, ":") {
			name = strings.Split(name, ":")[0]
		}

		models = append(models, &types.Model{
			ID:          name,
			Name:        name,
			Size:        s.formatBytes(model.Size),
			Type:        "chat",
			Status:      "available",
			Description: fmt.Sprintf("Ollama model: %s (%s)", name, model.Details.Family),
			ModelType:   "ollama",
			URL:         fmt.Sprintf("ollama://%s", model.Name),
		})
	}

	if len(models) == 0 {
		log.Println("⚠️ No models found in Ollama, returning fallback models")
		return s.getFallbackModels(), nil
	}

	log.Printf("✅ Found %d models in Ollama", len(models))
	return models, nil
}

// getFallbackModels returns a list of common models when Ollama is not available
func (s *OllamaService) getFallbackModels() []*types.Model {
	return []*types.Model{
		{
			ID:          "llama2",
			Name:        "llama2",
			Size:        "3.8GB",
			Type:        "chat",
			Status:      "available",
			Description: "Llama 2 7B Chat - Meta's conversational AI model",
			ModelType:   "llama",
			URL:         "ollama://llama2",
		},
		{
			ID:          "phi",
			Name:        "phi",
			Size:        "1.6GB",
			Type:        "chat",
			Status:      "available",
			Description: "Microsoft Phi-2 - Compact but powerful language model",
			ModelType:   "phi",
			URL:         "ollama://phi",
		},
		{
			ID:          "tinyllama",
			Name:        "tinyllama",
			Size:        "600MB",
			Type:        "chat",
			Status:      "available",
			Description: "TinyLlama - Ultra lightweight model for testing",
			ModelType:   "tinyllama",
			URL:         "ollama://tinyllama",
		},
		{
			ID:          "codellama",
			Name:        "codellama",
			Size:        "3.8GB",
			Type:        "code",
			Status:      "available",
			Description: "Code Llama - Specialized for code generation",
			ModelType:   "codellama",
			URL:         "ollama://codellama",
		},
	}
}

func (s *OllamaService) LoadModel(modelName string) error {
	log.Printf("🔄 Testing model availability in Ollama: %s", modelName)

	// Clean model name
	cleanName := strings.Split(modelName, ":")[0]

	// Test if model is available with a simple generation request
	reqBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "test",
		"stream": false,
		"options": map[string]interface{}{
			"num_predict": 1, // Only generate 1 token for testing
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(s.baseURL+"/api/generate", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try without :latest tag
		if strings.HasSuffix(modelName, ":latest") {
			return s.LoadModel(cleanName)
		}

		return fmt.Errorf("model not available in Ollama: %s (HTTP %d)", modelName, resp.StatusCode)
	}

	log.Printf("✅ Model is available and responding: %s", modelName)
	return nil
}

func (s *OllamaService) GenerateText(prompt, modelName string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  modelName,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.7,
			"top_p":       0.9,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(s.baseURL+"/api/generate", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API error: HTTP %d", resp.StatusCode)
	}

	var response struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Response, nil
}

func (s *OllamaService) CreateModel(model *types.Model) error {
	// For now, just return nil as Ollama manages its own models
	return nil
}

func (s *OllamaService) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
