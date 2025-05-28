package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

type OllamaService struct {
	baseURL string
	client  *http.Client
}

type OllamaModel struct {
	Name       string    `json:"name"`
	Model      string    `json:"model"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
	Details    struct {
		Parent            string   `json:"parent"`
		Format            string   `json:"format"`
		Family            string   `json:"family"`
		Families          []string `json:"families"`
		ParameterSize     string   `json:"parameter_size"`
		QuantizationLevel string   `json:"quantization_level"`
	} `json:"details"`
}

type OllamaListResponse struct {
	Models []OllamaModel `json:"models"`
}

func NewOllamaService() *OllamaService {
	return &OllamaService{
		baseURL: "http://localhost:11434",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *OllamaService) ListModels() ([]*types.Model, error) {
	resp, err := s.client.Get(s.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp OllamaListResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	var models []*types.Model
	for _, ollamaModel := range ollamaResp.Models {
		model := s.convertOllamaModelToTypes(ollamaModel)
		models = append(models, model)
	}

	return models, nil
}

func (s *OllamaService) convertOllamaModelToTypes(ollamaModel OllamaModel) *types.Model {
	// Extract base name without :latest suffix
	name := strings.TrimSuffix(ollamaModel.Name, ":latest")

	// Map model names to friendly display names
	displayName := s.getFriendlyName(name)

	// Convert size to string with appropriate unit
	sizeStr := s.formatSize(ollamaModel.Size)

	return &types.Model{
		ID:     name,
		Name:   displayName,
		Size:   sizeStr,
		Status: "available",
		// Removed Path field since it doesn't exist in types.Model
	}
}

func (s *OllamaService) getFriendlyName(modelName string) string {
	nameMap := map[string]string{
		"nemotron-nano":   "NVIDIA Llama 3.1 Nemotron Nano 4B",
		"neural-chat":     "Neural Chat 7B Q5_0",
		"openchat":        "OpenChat 3.5 Q5_K_M",
		"llama2-chat":     "Llama 2 7B Chat Q4_K_M",
		"phi2":            "Microsoft Phi-2 Q8_0",
		"llama3.2":        "Llama 3.2 3B",
		"llama3.2-vision": "Llama 3.2 Vision 11B",
		"mistral":         "Mistral 7B",
	}

	if friendlyName, exists := nameMap[modelName]; exists {
		return friendlyName
	}

	// Default: capitalize first letter and replace hyphens
	return strings.Title(strings.ReplaceAll(modelName, "-", " "))
}

func (s *OllamaService) formatSize(bytes int64) string {
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

func (s *OllamaService) LoadModel(modelName string) error {
	// Prepare the request to load/generate with the model
	reqBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "", // Empty prompt just to load the model
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(
		s.baseURL+"/api/generate",
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to load model %s: status %d", modelName, resp.StatusCode)
	}

	return nil
}

// CreateModel creates a new model entry (placeholder for basic models)
func (o *OllamaService) CreateModel(model *types.Model) error {
	// This is a placeholder implementation
	// In a real scenario, this would register the model with Ollama
	log.Printf("Model registration placeholder for: %s (type: %s)", model.Name, model.Type)
	return nil
}

func (s *OllamaService) IsAvailable() bool {
	resp, err := s.client.Get(s.baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
