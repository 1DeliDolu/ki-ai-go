//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ModelInfo struct {
	Name        string
	URL         string
	Size        string
	Description string
}

var models = []ModelInfo{
	{
		Name:        "llama-2-7b-chat.Q4_K_M.gguf",
		URL:         "https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf",
		Size:        "4.1GB",
		Description: "Llama 2 7B Chat model, 4-bit quantized",
	},
	{
		Name:        "mistral-7b-instruct-v0.1.Q4_K_M.gguf",
		URL:         "https://huggingface.co/TheBloke/Mistral-7B-Instruct-v0.1-GGUF/resolve/main/mistral-7b-instruct-v0.1.Q4_K_M.gguf",
		Size:        "4.4GB",
		Description: "Mistral 7B Instruct model, 4-bit quantized",
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Available models:")
		for i, model := range models {
			fmt.Printf("%d. %s (%s) - %s\n", i+1, model.Name, model.Size, model.Description)
		}
		fmt.Println("\nUsage: go run download_models.go <model_number>")
		return
	}

	modelIndex := os.Args[1]
	var selectedModel ModelInfo

	switch modelIndex {
	case "1":
		selectedModel = models[0]
	case "2":
		selectedModel = models[1]
	default:
		fmt.Println("Invalid model number")
		return
	}

	// Create models directory
	homeDir, _ := os.UserHomeDir()
	modelsDir := filepath.Join(homeDir, ".ki-ai-go", "models")
	os.MkdirAll(modelsDir, 0755)

	modelPath := filepath.Join(modelsDir, selectedModel.Name)

	// Check if file already exists
	if _, err := os.Stat(modelPath); err == nil {
		fmt.Printf("Model %s already exists!\n", selectedModel.Name)
		return
	}

	fmt.Printf("Downloading %s (%s)...\n", selectedModel.Name, selectedModel.Size)
	fmt.Printf("URL: %s\n", selectedModel.URL)
	fmt.Printf("Destination: %s\n", modelPath)

	if err := downloadFile(selectedModel.URL, modelPath); err != nil {
		fmt.Printf("Error downloading model: %v\n", err)
		return
	}

	fmt.Printf("Successfully downloaded %s\n", selectedModel.Name)
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
