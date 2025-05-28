package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DocumentConverter provides document conversion functionality
type DocumentConverter struct{}

// NewDocumentConverter creates a new document converter
func NewDocumentConverter() *DocumentConverter {
	return &DocumentConverter{}
}

// ConvertToMarkdown converts any document to Markdown format
func (dc *DocumentConverter) ConvertToMarkdown(inputPath, outputPath string) error {
	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Basic conversion - in practice you'd have more sophisticated conversion
	markdownContent := fmt.Sprintf("# Document: %s\n\n```\n%s\n```\n",
		filepath.Base(inputPath), string(content))

	// Write output file
	return os.WriteFile(outputPath, []byte(markdownContent), 0644)
}

// ConvertToHTML converts any document to HTML format
func (dc *DocumentConverter) ConvertToHTML(inputPath, outputPath string) error {
	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Basic HTML conversion
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <meta charset="UTF-8">
</head>
<body>
    <h1>Document: %s</h1>
    <pre>%s</pre>
</body>
</html>`,
		filepath.Base(inputPath),
		filepath.Base(inputPath),
		strings.ReplaceAll(string(content), "<", "&lt;"))

	// Write output file
	return os.WriteFile(outputPath, []byte(htmlContent), 0644)
}

// ConvertToPlainText extracts plain text from any document
func (dc *DocumentConverter) ConvertToPlainText(inputPath, outputPath string) error {
	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Basic text extraction - remove some common markup
	textContent := string(content)
	textContent = strings.ReplaceAll(textContent, "\r\n", "\n")
	textContent = strings.ReplaceAll(textContent, "\r", "\n")

	// Write output file
	return os.WriteFile(outputPath, []byte(textContent), 0644)
}
