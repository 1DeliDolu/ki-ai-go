package handlers

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

	// Register all processors
	dm.RegisterProcessor(&TXTProcessor{})
	dm.RegisterProcessor(&MarkdownProcessor{})
	dm.RegisterProcessor(&PDFProcessor{})
	dm.RegisterProcessor(&DOCXProcessor{})
	dm.RegisterProcessor(&HTMLProcessor{})

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

// MarkdownProcessor handles markdown files
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

// PDFProcessor handles PDF files (placeholder - will be implemented later)
type PDFProcessor struct{}

func (p *PDFProcessor) Read(path string) (*types.DocumentContent, error) {
	// Placeholder implementation
	return &types.DocumentContent{
		Text: "PDF processing not yet implemented",
		Type: "pdf",
		Metadata: map[string]string{
			"status": "placeholder",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *PDFProcessor) GetSupportedTypes() []string {
	return []string{"pdf"}
}

// DOCXProcessor handles Word documents (placeholder - will be implemented later)
type DOCXProcessor struct{}

func (p *DOCXProcessor) Read(path string) (*types.DocumentContent, error) {
	// Placeholder implementation
	return &types.DocumentContent{
		Text: "DOCX processing not yet implemented",
		Type: "docx",
		Metadata: map[string]string{
			"status": "placeholder",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *DOCXProcessor) GetSupportedTypes() []string {
	return []string{"docx", "doc"}
}

// HTMLProcessor handles HTML files (placeholder - will be implemented later)
type HTMLProcessor struct{}

func (p *HTMLProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML file: %w", err)
	}

	text := string(content)
	// Simple HTML tag removal for text extraction
	text = strings.ReplaceAll(text, "<", " <")
	text = strings.ReplaceAll(text, ">", "> ")

	return &types.DocumentContent{
		Text: text,
		Type: "html",
		Metadata: map[string]string{
			"char_count": fmt.Sprintf("%d", len(text)),
			"status":     "basic_extraction",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *HTMLProcessor) GetSupportedTypes() []string {
	return []string{"html", "htm"}
}
