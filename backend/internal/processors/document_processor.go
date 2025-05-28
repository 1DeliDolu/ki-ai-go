package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

// DocumentProcessor interface for different document types
type DocumentProcessor interface {
	Read(path string) (*types.DocumentContent, error)
	GetSupportedTypes() []string
}

// DocumentManager manages different document processors
type DocumentManager struct {
	processors map[string]DocumentProcessor
}

// NewDocumentManager creates a new document manager with all processors
func NewDocumentManager() *DocumentManager {
	dm := &DocumentManager{
		processors: make(map[string]DocumentProcessor),
	}

	// Register working processors first
	dm.RegisterProcessor(&TXTProcessor{})
	dm.RegisterProcessor(&MarkdownProcessor{})
	dm.RegisterProcessor(&HTMLProcessor{})
	// TODO: Add PDF and DOCX when dependencies are resolved
	// dm.RegisterProcessor(&PDFProcessor{})
	// dm.RegisterProcessor(&DOCXProcessor{})

	return dm
}

// RegisterProcessor registers a document processor for specific file types
func (dm *DocumentManager) RegisterProcessor(processor DocumentProcessor) {
	types := processor.GetSupportedTypes()
	for _, t := range types {
		dm.processors[t] = processor
	}
}

// ProcessDocument processes a document based on its file extension
func (dm *DocumentManager) ProcessDocument(path string) (*types.DocumentContent, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:] // Remove the dot
	}

	processor, exists := dm.processors[ext]
	if !exists {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	return processor.Read(path)
}

// GetSupportedTypes returns all supported file extensions
func (dm *DocumentManager) GetSupportedTypes() []string {
	var types []string
	for ext := range dm.processors {
		types = append(types, ext)
	}
	return types
}

// TXTProcessor handles plain text files
type TXTProcessor struct{}

func (p *TXTProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read TXT file: %w", err)
	}

	text := string(content)
	wordCount := len(strings.Fields(text))
	lineCount := len(strings.Split(text, "\n"))

	return &types.DocumentContent{
		Text: text,
		Type: "txt",
		Metadata: map[string]string{
			"word_count": fmt.Sprintf("%d", wordCount),
			"line_count": fmt.Sprintf("%d", lineCount),
			"char_count": fmt.Sprintf("%d", len(text)),
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *TXTProcessor) GetSupportedTypes() []string {
	return []string{"txt", "text"}
}

// MarkdownProcessor handles markdown files (basic implementation)
type MarkdownProcessor struct{}

func (p *MarkdownProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Markdown file: %w", err)
	}

	text := string(content)

	// Count headers (lines starting with #)
	lines := strings.Split(text, "\n")
	headerCount := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			headerCount++
		}
	}

	return &types.DocumentContent{
		Text: text,
		Type: "markdown",
		Metadata: map[string]string{
			"word_count":   fmt.Sprintf("%d", len(strings.Fields(text))),
			"line_count":   fmt.Sprintf("%d", len(lines)),
			"header_count": fmt.Sprintf("%d", headerCount),
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *MarkdownProcessor) GetSupportedTypes() []string {
	return []string{"md", "markdown"}
}

// HTMLProcessor handles HTML files (basic implementation without external libs)
type HTMLProcessor struct{}

func (p *HTMLProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open HTML file: %w", err)
	}

	text := string(content)

	// Basic text extraction - remove HTML tags
	text = p.stripHTMLTags(text)

	// Count basic HTML elements in original content
	originalContent := string(content)
	linkCount := strings.Count(strings.ToLower(originalContent), "<a ")
	imgCount := strings.Count(strings.ToLower(originalContent), "<img ")
	headerCount := 0
	for i := 1; i <= 6; i++ {
		headerCount += strings.Count(strings.ToLower(originalContent), fmt.Sprintf("<h%d", i))
	}

	// Extract title
	title := p.extractTitle(originalContent)

	return &types.DocumentContent{
		Text: text,
		Type: "html",
		Metadata: map[string]string{
			"title":        title,
			"word_count":   fmt.Sprintf("%d", len(strings.Fields(text))),
			"char_count":   fmt.Sprintf("%d", len(text)),
			"link_count":   fmt.Sprintf("%d", linkCount),
			"image_count":  fmt.Sprintf("%d", imgCount),
			"header_count": fmt.Sprintf("%d", headerCount),
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *HTMLProcessor) stripHTMLTags(s string) string {
	// Simple HTML tag removal
	var result strings.Builder
	inTag := false

	for _, char := range s {
		switch char {
		case '<':
			inTag = true
		case '>':
			inTag = false
			result.WriteRune(' ') // Replace tag with space
		default:
			if !inTag {
				result.WriteRune(char)
			}
		}
	}

	// Clean up multiple spaces
	text := result.String()
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Remove multiple consecutive spaces
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return strings.TrimSpace(text)
}

func (p *HTMLProcessor) extractTitle(content string) string {
	lower := strings.ToLower(content)
	start := strings.Index(lower, "<title>")
	if start == -1 {
		return ""
	}
	start += 7 // len("<title>")

	end := strings.Index(lower[start:], "</title>")
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(content[start : start+end])
}

func (p *HTMLProcessor) GetSupportedTypes() []string {
	return []string{"html", "htm"}
}

// FileTypeDetector helps detect file types (basic implementation)
func DetectFileType(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(ext, ".") {
		return ext[1:], nil
	}
	return ext, nil
}
