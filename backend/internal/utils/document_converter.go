package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/internal/processors"
	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

type DocumentConverter struct {
	manager *processors.DocumentManager
}

func NewDocumentConverter() *DocumentConverter {
	return &DocumentConverter{
		manager: processors.NewDocumentManager(),
	}
}

// ConvertToMarkdown converts any document to Markdown format
func (dc *DocumentConverter) ConvertToMarkdown(inputPath, outputPath string) error {
	// Process the input document
	content, err := dc.manager.ProcessDocument(inputPath)
	if err != nil {
		return fmt.Errorf("failed to process input document: %w", err)
	}

	// Convert content to Markdown format
	markdown := dc.convertContentToMarkdown(content, inputPath)

	// Write to output file
	if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}

// ConvertToHTML converts any document to HTML format
func (dc *DocumentConverter) ConvertToHTML(inputPath, outputPath string) error {
	// Process the input document
	content, err := dc.manager.ProcessDocument(inputPath)
	if err != nil {
		return fmt.Errorf("failed to process input document: %w", err)
	}

	// Convert content to HTML format
	html := dc.convertContentToHTML(content, inputPath)

	// Write to output file
	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}

// ConvertToPlainText extracts plain text from any document
func (dc *DocumentConverter) ConvertToPlainText(inputPath, outputPath string) error {
	content, err := dc.manager.ProcessDocument(inputPath)
	if err != nil {
		return fmt.Errorf("failed to process input document: %w", err)
	}

	// Write plain text
	if err := os.WriteFile(outputPath, []byte(content.Text), 0644); err != nil {
		return fmt.Errorf("failed to write text file: %w", err)
	}

	return nil
}

// BatchConvert converts multiple documents to specified format
func (dc *DocumentConverter) BatchConvert(inputPaths []string, outputDir, format string) (map[string]string, error) {
	results := make(map[string]string)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, inputPath := range inputPaths {
		// Generate output filename
		basename := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
		var outputPath string

		switch strings.ToLower(format) {
		case "markdown", "md":
			outputPath = filepath.Join(outputDir, basename+".md")
			err := dc.ConvertToMarkdown(inputPath, outputPath)
			if err != nil {
				results[inputPath] = fmt.Sprintf("Error: %v", err)
			} else {
				results[inputPath] = "Success: " + outputPath
			}
		case "html":
			outputPath = filepath.Join(outputDir, basename+".html")
			err := dc.ConvertToHTML(inputPath, outputPath)
			if err != nil {
				results[inputPath] = fmt.Sprintf("Error: %v", err)
			} else {
				results[inputPath] = "Success: " + outputPath
			}
		case "txt", "text":
			outputPath = filepath.Join(outputDir, basename+".txt")
			err := dc.ConvertToPlainText(inputPath, outputPath)
			if err != nil {
				results[inputPath] = fmt.Sprintf("Error: %v", err)
			} else {
				results[inputPath] = "Success: " + outputPath
			}
		default:
			results[inputPath] = "Error: Unsupported format: " + format
		}
	}

	return results, nil
}

// convertContentToMarkdown converts document content to Markdown
func (dc *DocumentConverter) convertContentToMarkdown(content *types.DocumentContent, originalPath string) string {
	var md strings.Builder

	// Add header with metadata
	md.WriteString("# Document: " + filepath.Base(originalPath) + "\n\n")
	md.WriteString("**Type:** " + content.Type + "\n\n")
	md.WriteString("**Processed:** " + content.ProcessedAt.Format(time.RFC3339) + "\n\n")

	// Add metadata as table
	if len(content.Metadata) > 0 {
		md.WriteString("## Metadata\n\n")
		md.WriteString("| Property | Value |\n")
		md.WriteString("|----------|-------|\n")
		for key, value := range content.Metadata {
			md.WriteString(fmt.Sprintf("| %s | %s |\n", key, value))
		}
		md.WriteString("\n")
	}

	// Add content
	md.WriteString("## Content\n\n")

	// Format content based on type
	switch content.Type {
	case "code":
		// Detect language from metadata
		lang := content.Metadata["language"]
		if lang == "" {
			lang = "text"
		}
		md.WriteString("```" + strings.ToLower(lang) + "\n")
		md.WriteString(content.Text)
		md.WriteString("\n```\n")
	case "json":
		md.WriteString("```json\n")
		md.WriteString(content.Text)
		md.WriteString("\n```\n")
	case "xml":
		md.WriteString("```xml\n")
		md.WriteString(content.Text)
		md.WriteString("\n```\n")
	default:
		// Regular text content
		md.WriteString(content.Text)
	}

	return md.String()
}

// convertContentToHTML converts document content to HTML
func (dc *DocumentConverter) convertContentToHTML(content *types.DocumentContent, originalPath string) string {
	var html strings.Builder

	// HTML header
	html.WriteString("<!DOCTYPE html>\n")
	html.WriteString("<html lang=\"en\">\n<head>\n")
	html.WriteString("    <meta charset=\"UTF-8\">\n")
	html.WriteString("    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	html.WriteString("    <title>" + filepath.Base(originalPath) + "</title>\n")
	html.WriteString("    <style>\n")
	html.WriteString("        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }\n")
	html.WriteString("        .metadata { background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0; }\n")
	html.WriteString("        .content { background: white; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }\n")
	html.WriteString("        pre { background: #f8f8f8; padding: 15px; border-radius: 5px; overflow-x: auto; }\n")
	html.WriteString("        table { border-collapse: collapse; width: 100%; }\n")
	html.WriteString("        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }\n")
	html.WriteString("        th { background-color: #f2f2f2; }\n")
	html.WriteString("    </style>\n")
	html.WriteString("</head>\n<body>\n")

	// Document header
	html.WriteString("    <h1>Document: " + filepath.Base(originalPath) + "</h1>\n")

	// Metadata section
	if len(content.Metadata) > 0 {
		html.WriteString("    <div class=\"metadata\">\n")
		html.WriteString("        <h2>Document Information</h2>\n")
		html.WriteString("        <table>\n")
		html.WriteString("            <tr><th>Property</th><th>Value</th></tr>\n")
		html.WriteString(fmt.Sprintf("            <tr><td>Type</td><td>%s</td></tr>\n", content.Type))
		html.WriteString(fmt.Sprintf("            <tr><td>Processed</td><td>%s</td></tr>\n", content.ProcessedAt.Format(time.RFC3339)))
		for key, value := range content.Metadata {
			html.WriteString(fmt.Sprintf("            <tr><td>%s</td><td>%s</td></tr>\n", key, value))
		}
		html.WriteString("        </table>\n")
		html.WriteString("    </div>\n")
	}

	// Content section
	html.WriteString("    <div class=\"content\">\n")
	html.WriteString("        <h2>Content</h2>\n")

	// Format content based on type
	switch content.Type {
	case "html":
		// Already HTML, just include it
		html.WriteString("        <div>" + content.Text + "</div>\n")
	case "code", "json", "xml":
		// Code content in pre tags
		html.WriteString("        <pre><code>" + dc.escapeHTML(content.Text) + "</code></pre>\n")
	default:
		// Convert plain text to HTML paragraphs
		paragraphs := strings.Split(content.Text, "\n\n")
		for _, para := range paragraphs {
			if strings.TrimSpace(para) != "" {
				html.WriteString("        <p>" + dc.escapeHTML(para) + "</p>\n")
			}
		}
	}

	html.WriteString("    </div>\n")
	html.WriteString("</body>\n</html>\n")

	return html.String()
}

// escapeHTML escapes HTML special characters
func (dc *DocumentConverter) escapeHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(text)
}

// GetSupportedFormats returns supported output formats
func (dc *DocumentConverter) GetSupportedFormats() []string {
	return []string{"markdown", "md", "html", "txt", "text"}
}
