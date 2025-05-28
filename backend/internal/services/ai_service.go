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

	"github.com/1DeliDolu/ki-ai-go/internal/config"
	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

type AIService struct {
	config        *config.Config
	client        *http.Client
	modelName     string
	currentModel  string // Added missing field
	isModelLoaded bool
	ollamaService *OllamaService // Added missing field
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
		ollamaService: NewOllamaService(), // Initialize ollama service
	}
}

func (s *AIService) LoadModel(modelName string) error {
	log.Printf("🔄 Loading model in AI service: %s", modelName)

	// Clean model name
	cleanModelName := strings.Split(modelName, ":")[0]

	// Try to load with different name variations
	modelVariations := []string{
		cleanModelName,
		modelName,
		cleanModelName + ":latest",
	}

	var lastError error
	for _, variation := range modelVariations {
		log.Printf("🔄 AI Service trying: %s", variation)

		// Test if the model works with a simple generation
		if err := s.testModelGeneration(variation); err != nil {
			log.Printf("⚠️ Model test failed for %s: %v", variation, err)
			lastError = err
			continue
		}

		// Success!
		s.modelName = variation
		s.currentModel = variation
		s.isModelLoaded = true
		log.Printf("✅ AI Service successfully loaded: %s", variation)
		return nil
	}

	return fmt.Errorf("failed to load model in AI service: %w", lastError)
}

// testModelGeneration tests if a model can generate text
func (s *AIService) testModelGeneration(modelName string) error {
	log.Printf("🧪 Testing model generation: %s", modelName)

	response, err := s.generateWithOllama("Hi", modelName)
	if err != nil {
		return fmt.Errorf("generation test failed: %w", err)
	}

	if strings.TrimSpace(response) == "" {
		return fmt.Errorf("model returned empty response")
	}

	log.Printf("✅ Model generation test passed: %s", modelName)
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

func (s *AIService) generateWithOllama(prompt, modelName string) (string, error) {
	log.Printf("🔄 Generating with Ollama: %s", modelName)

	reqBody := OllamaGenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.7,
			"top_p":       0.9,
			"top_k":       40,
			"num_predict": 50, // Limit tokens for faster response
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
		return "", fmt.Errorf("Ollama API error: HTTP %d", resp.StatusCode)
	}

	var response OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Response, nil
}

func (s *AIService) GenerateResponse(query string, documents []types.Document, wikiResults []types.WikiResult) (string, error) {
	log.Printf("🤖 Generating AI response for query: %s", query)

	// Build context from documents with ACTUAL CONTENT
	var context strings.Builder
	context.WriteString("Context from uploaded documents:\n\n")

	for _, doc := range documents {
		// Get actual document content, not just metadata
		if doc.Path != "" {
			// Read file content directly
			if content, err := os.ReadFile(doc.Path); err == nil {
				context.WriteString(fmt.Sprintf("=== Document: %s ===\n", doc.Name))
				context.WriteString(string(content))
				context.WriteString("\n\n")
				log.Printf("📄 Added content from %s (%d bytes)", doc.Name, len(content))
			} else {
				context.WriteString(fmt.Sprintf("=== Document: %s ===\n", doc.Name))
				context.WriteString("(Content could not be read)\n\n")
				log.Printf("❌ Could not read content from %s: %v", doc.Name, err)
			}
		} else {
			context.WriteString(fmt.Sprintf("=== Document: %s ===\n", doc.Name))
			context.WriteString("(No file path available)\n\n")
		}
	}

	// Add wiki context - fix Summary field issue
	if len(wikiResults) > 0 {
		context.WriteString("Additional context from Wikipedia:\n\n")
		for _, wiki := range wikiResults {
			// Use Description instead of Summary if Summary doesn't exist
			summary := wiki.Description
			if summary == "" {
				summary = "No description available"
			}
			context.WriteString(fmt.Sprintf("- %s: %s\n", wiki.Title, summary))
		}
		context.WriteString("\n")
	}

	// Enhanced prompt with document content
	prompt := fmt.Sprintf(`Based on the following documents and context, please answer this question: %s

%s

Please provide a detailed answer based on the content above. If the answer is found in the documents, reference which document contains the information.`,
		query, context.String())

	// Generate response using the current model
	if s.currentModel == "" {
		return "Please load a model first to generate responses.", nil
	}

	// Use generateWithOllama method
	response, err := s.generateWithOllama(prompt, s.currentModel)
	if err != nil {
		log.Printf("❌ Error generating response: %v", err)

		// Fallback: Provide basic response with document content
		if len(documents) > 0 {
			fallback := fmt.Sprintf("I found %d document(s) related to your query:\n\n", len(documents))
			for _, doc := range documents {
				if doc.Path != "" {
					if content, err := os.ReadFile(doc.Path); err == nil {
						fallback += fmt.Sprintf("**%s:**\n%s\n\n", doc.Name, string(content))
					}
				}
			}
			return fallback, nil
		}

		return fmt.Errorf("failed to generate AI response: %w", err).Error(), nil
	}

	log.Printf("✅ Generated AI response (%d characters)", len(response))
	return response, nil
}

func (s *AIService) GetCurrentModel() string {
	if s.currentModel != "" {
		return s.currentModel
	}
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
