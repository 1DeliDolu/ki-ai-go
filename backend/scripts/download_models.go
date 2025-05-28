//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	{
		Name:        "phi-2.Q4_K_M.gguf",
		URL:         "https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.Q4_K_M.gguf",
		Size:        "1.6GB",
		Description: "Microsoft Phi-2 model, 4-bit quantized",
	},
	{
		Name:        "openchat-3.5-0106.Q4_K_M.gguf",
		URL:         "https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.Q4_K_M.gguf",
		Size:        "4.1GB",
		Description: "OpenChat 3.5 model, 4-bit quantized",
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("🤖 Available models for Local AI Project:")
		fmt.Println("=========================================")
		for i, model := range models {
			fmt.Printf("%d. %s (%s)\n", i+1, model.Name, model.Size)
			fmt.Printf("   📝 %s\n", model.Description)
			fmt.Println()
		}
		fmt.Println("📥 Usage: go run download_models.go <model_number>")
		fmt.Println("📁 Models will be saved to: ./models/ directory")
		return
	}

	modelIndex := os.Args[1]
	var selectedModel ModelInfo

	switch modelIndex {
	case "1":
		selectedModel = models[0]
	case "2":
		selectedModel = models[1]
	case "3":
		selectedModel = models[2]
	case "4":
		selectedModel = models[3]
	default:
		fmt.Println("❌ Invalid model number. Please choose 1-4.")
		return
	}

	// Use local models directory relative to project root
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("❌ Error finding project root: %v\n", err)
		return
	}

	modelsDir := filepath.Join(projectRoot, "models")

	// Create models directory if it doesn't exist
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		fmt.Printf("❌ Error creating models directory: %v\n", err)
		return
	}

	modelPath := filepath.Join(modelsDir, selectedModel.Name)

	// Check if file already exists
	if _, err := os.Stat(modelPath); err == nil {
		fmt.Printf("✅ Model %s already exists!\n", selectedModel.Name)
		fmt.Printf("📍 Location: %s\n", modelPath)
		return
	}

	fmt.Printf("📥 Downloading %s (%s)...\n", selectedModel.Name, selectedModel.Size)
	fmt.Printf("🔗 URL: %s\n", selectedModel.URL)
	fmt.Printf("📍 Destination: %s\n", modelPath)
	fmt.Println()

	startTime := time.Now()
	if err := downloadFileWithProgress(selectedModel.URL, modelPath); err != nil {
		fmt.Printf("❌ Error downloading model: %v\n", err)
		// Clean up partial file
		os.Remove(modelPath)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("\n✅ Successfully downloaded %s in %v\n", selectedModel.Name, duration.Truncate(time.Second))
	fmt.Printf("📍 Saved to: %s\n", modelPath)

	// Show file size
	if stat, err := os.Stat(modelPath); err == nil {
		fmt.Printf("📊 File size: %s\n", formatFileSize(stat.Size()))
	}

	fmt.Println("\n💡 Next steps:")
	fmt.Println("1. Start your backend server: ./start.sh")
	fmt.Println("2. Load this model via API or frontend")
	fmt.Printf("3. Model is available as: %s\n", selectedModel.Name)
}

// getProjectRoot finds the project root directory
func getProjectRoot() (string, error) {
	// Start from current directory and go up to find backend directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// If we're already in backend directory, go up one level
	if filepath.Base(currentDir) == "backend" {
		return filepath.Dir(currentDir), nil
	}

	// If we're in scripts directory, go up two levels
	if filepath.Base(currentDir) == "scripts" {
		return filepath.Dir(filepath.Dir(currentDir)), nil
	}

	// Otherwise assume we're in project root
	return currentDir, nil
}

// downloadFileWithProgress downloads file with progress indication
func downloadFileWithProgress(url, filepath string) error {
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

	// Create a progress reader
	fileSize := resp.ContentLength
	if fileSize > 0 {
		fmt.Printf("📊 File size: %s\n", formatFileSize(fileSize))
	}

	// Copy with progress (simplified - no progress bar for now)
	_, err = io.Copy(out, resp.Body)
	return err
}

// formatFileSize formats bytes into human readable format
func formatFileSize(bytes int64) string {
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
