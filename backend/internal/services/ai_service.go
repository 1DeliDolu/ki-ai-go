package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1DeliDolu/go_mustAI/local-ai-project/backend/internal/config"
	"github.com/1DeliDolu/go_mustAI/local-ai-project/backend/pkg/types"
)

type AIService struct {
	config        *config.Config
	client        *http.Client
	modelName     string
	isModelLoaded bool
}

type OllamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type OllamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type OllamaPullRequest struct {
	Name string `json:"name"`
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		config: cfg,
		client: &http.Client{
			Timeout: 120 * time.Second, // 2 minutes timeout for AI responses
		},
	}
}

func (s *AIService) LoadModel(modelName string) error {
	log.Printf("Loading model: %s", modelName)

	// First check if model file exists locally
	modelPath := s.findModelFile(modelName)
	if modelPath == "" {
		// If no local file, try to pull from Ollama registry
		return s.pullModelFromOllama(modelName)
	}

	log.Printf("Found local model file: %s", modelPath)

	// For local GGUF files, we need to create Ollama modelfile
	if err := s.createOllamaModelfile(modelName, modelPath); err != nil {
		return fmt.Errorf("failed to create Ollama modelfile: %w", err)
	}

	// Test if the model works with Ollama
	if err := s.testModelWithOllama(modelName); err != nil {
		return fmt.Errorf("model loaded but not responding: %w", err)
	}

	s.modelName = modelName
	s.isModelLoaded = true
	log.Printf("Successfully loaded model: %s", modelName)
	return nil
}

func (s *AIService) createOllamaModelfile(modelName, modelPath string) error {
	// Create a simple Ollama modelfile for the GGUF model
	modelfile := fmt.Sprintf(`FROM %s

TEMPLATE """{{ if .System }}<|system|>
{{ .System }}<|end|>
{{ end }}{{ if .Prompt }}<|user|>
{{ .Prompt }}<|end|>
{{ end }}<|assistant|>
{{ .Response }}<|end|>
"""

PARAMETER stop "<|end|>"
PARAMETER stop "<|user|>"
PARAMETER stop "<|system|>"
PARAMETER temperature 0.7
PARAMETER top_p 0.9
PARAMETER top_k 40
`, modelPath)

	// Create Ollama model using the API
	reqBody := map[string]interface{}{
		"name":      modelName,
		"modelfile": modelfile,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal create request: %w", err)
	}

	resp, err := s.client.Post(s.config.OllamaURL+"/api/create", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create Ollama model: HTTP %d", resp.StatusCode)
	}

	return nil
}

func (s *AIService) pullModelFromOllama(modelName string) error {
	log.Printf("Pulling model from Ollama: %s", modelName)

	reqBody := OllamaPullRequest{
		Name: modelName,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal pull request: %w", err)
	}

	resp, err := s.client.Post(s.config.OllamaURL+"/api/pull", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama at %s: %w", s.config.OllamaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull model from Ollama: HTTP %d", resp.StatusCode)
	}

	s.modelName = modelName
	s.isModelLoaded = true
	return nil
}

func (s *AIService) testModelWithOllama(modelName string) error {
	// Test with a simple prompt
	_, err := s.generateWithOllama("test", modelName)
	return err
}

func (s *AIService) findModelFile(modelName string) string {
	// Try exact match first
	exactPath := filepath.Join(s.config.ModelsPath, modelName)
	if _, err := os.Stat(exactPath); err == nil {
		return exactPath
	}

	// Common model file extensions
	extensions := []string{".gguf", ".bin", ".ggml"}

	for _, ext := range extensions {
		// Try with extension
		path := filepath.Join(s.config.ModelsPath, modelName+ext)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Search for any file containing the model name (case insensitive)
	files, err := os.ReadDir(s.config.ModelsPath)
	if err != nil {
		return ""
	}

	// Try exact filename match first
	for _, file := range files {
		if !file.IsDir() && strings.EqualFold(file.Name(), modelName) {
			return filepath.Join(s.config.ModelsPath, file.Name())
		}
	}

	// Try partial match
	for _, file := range files {
		if !file.IsDir() && strings.Contains(strings.ToLower(file.Name()), strings.ToLower(modelName)) {
			return filepath.Join(s.config.ModelsPath, file.Name())
		}
	}

	return ""
}

func (s *AIService) GenerateResponse(query string, documents []types.Document, wikiResults []types.WikiResult) (string, error) {
	if !s.isModelLoaded {
		return "", fmt.Errorf("no model loaded. Please load a model first")
	}

	// Build context from documents and wiki results
	var context strings.Builder

	// Add document context
	if len(documents) > 0 {
		context.WriteString("Relevant documents:\n")
		for _, doc := range documents {
			context.WriteString(fmt.Sprintf("- %s (Type: %s)\n", doc.Name, doc.Type))
		}
		context.WriteString("\n")
	}

	// Add wiki context
	if len(wikiResults) > 0 {
		context.WriteString("Wikipedia information:\n")
		for _, wiki := range wikiResults {
			if wiki.Extract != "" {
				context.WriteString(fmt.Sprintf("- %s: %s\n", wiki.Title, wiki.Extract))
			} else if wiki.Description != "" {
				context.WriteString(fmt.Sprintf("- %s: %s\n", wiki.Title, wiki.Description))
			}
		}
		context.WriteString("\n")
	}

	// Create the prompt with context
	var prompt string
	if context.Len() > 0 {
		prompt = fmt.Sprintf(`Based on the following context, please answer the question clearly and concisely.

Context:
%s

Question: %s

Answer:`, context.String(), query)
	} else {
		prompt = fmt.Sprintf("Please answer the following question clearly and concisely:\n\nQuestion: %s\n\nAnswer:", query)
	}

	return s.generateWithOllama(prompt, s.modelName)
}

func (s *AIService) generateWithOllama(prompt, modelName string) (string, error) {
	reqBody := OllamaGenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature":    0.7,
			"top_p":          0.9,
			"top_k":          40,
			"num_predict":    512,
			"repeat_penalty": 1.1,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(s.config.OllamaURL+"/api/generate", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API error: HTTP %d", resp.StatusCode)
	}

	var response OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Clean up the response
	result := strings.TrimSpace(response.Response)
	if result == "" {
		return "I apologize, but I couldn't generate a proper response. Please try rephrasing your question.", nil
	}

	return result, nil
}

func (s *AIService) GetCurrentModel() string {
	return s.modelName
}

func (s *AIService) IsModelLoaded() bool {
	return s.isModelLoaded
}

func (s *AIService) Close() {
	// No resources to clean up with HTTP client approach
	s.isModelLoaded = false
	s.modelName = ""
}
