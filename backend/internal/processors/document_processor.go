package processors

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
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
	stats      ProcessingStats
}

// ProcessingStats tracks document processing statistics
type ProcessingStats struct {
	TotalProcessed     int
	SuccessfullyParsed int
	Failed             int
	TypeCounts         map[string]int
	LastProcessed      time.Time
}

// NewDocumentManager creates a new document manager with all processors
func NewDocumentManager() *DocumentManager {
	dm := &DocumentManager{
		processors: make(map[string]DocumentProcessor),
		stats: ProcessingStats{
			TypeCounts: make(map[string]int),
		},
	}

	// Register basic processors
	dm.RegisterProcessor(&TXTProcessor{})
	dm.RegisterProcessor(&MarkdownProcessor{})
	dm.RegisterProcessor(&HTMLProcessor{})

	// Register advanced processors
	dm.RegisterProcessor(&PDFProcessor{})
	dm.RegisterProcessor(&DOCXProcessor{})
	dm.RegisterProcessor(&JSONProcessor{})
	dm.RegisterProcessor(&XMLProcessor{})
	dm.RegisterProcessor(&CSVProcessor{})
	dm.RegisterProcessor(&LogProcessor{})
	dm.RegisterProcessor(&CodeProcessor{})

	log.Printf("📄 DocumentManager initialized with %d processors", len(dm.processors))
	return dm
}

// RegisterProcessor registers a document processor for specific file types
func (dm *DocumentManager) RegisterProcessor(processor DocumentProcessor) {
	types := processor.GetSupportedTypes()
	for _, t := range types {
		dm.processors[t] = processor
	}
}

// ProcessDocument processes a document based on its file extension with enhanced features
func (dm *DocumentManager) ProcessDocument(path string) (*types.DocumentContent, error) {
	log.Printf("🔄 Processing document: %s", filepath.Base(path))

	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:] // Remove the dot
	}

	processor, exists := dm.processors[ext]
	if !exists {
		dm.stats.Failed++
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	// Update processing stats
	dm.stats.TotalProcessed++
	dm.stats.LastProcessed = time.Now()

	content, err := processor.Read(path)
	if err != nil {
		dm.stats.Failed++
		return nil, fmt.Errorf("failed to process %s: %w", filepath.Base(path), err)
	}

	// Update success stats
	dm.stats.SuccessfullyParsed++
	dm.stats.TypeCounts[ext]++

	log.Printf("✅ Successfully processed %s (%s)", filepath.Base(path), ext)
	return content, nil
}

// ProcessMultipleDocuments processes multiple documents and returns results
func (dm *DocumentManager) ProcessMultipleDocuments(paths []string) map[string]*types.DocumentContent {
	results := make(map[string]*types.DocumentContent)

	log.Printf("📦 Processing %d documents...", len(paths))

	for _, path := range paths {
		content, err := dm.ProcessDocument(path)
		if err != nil {
			log.Printf("❌ Error processing %s: %v", filepath.Base(path), err)
			continue
		}
		results[path] = content
	}

	log.Printf("✅ Successfully processed %d out of %d documents", len(results), len(paths))
	return results
}

// GetProcessingStats returns current processing statistics
func (dm *DocumentManager) GetProcessingStats() ProcessingStats {
	return dm.stats
}

// ResetStats resets processing statistics
func (dm *DocumentManager) ResetStats() {
	dm.stats = ProcessingStats{
		TypeCounts: make(map[string]int),
	}
	log.Println("📊 Processing stats reset")
}

// GetProcessorInfo returns information about a specific processor
func (dm *DocumentManager) GetProcessorInfo(fileType string) map[string]interface{} {
	processor, exists := dm.processors[fileType]
	if !exists {
		return map[string]interface{}{
			"supported": false,
			"error":     fmt.Sprintf("No processor available for type: %s", fileType),
		}
	}

	return map[string]interface{}{
		"supported":       true,
		"processor_type":  fmt.Sprintf("%T", processor),
		"supported_types": processor.GetSupportedTypes(),
		"processed_count": dm.stats.TypeCounts[fileType],
	}
}

// ValidateFile checks if a file can be processed
func (dm *DocumentManager) ValidateFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	if _, exists := dm.processors[ext]; !exists {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	// Check file size (optional limit)
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot read file info: %w", err)
	}

	// Set a reasonable file size limit (100MB)
	const maxFileSize = 100 * 1024 * 1024
	if stat.Size() > maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", stat.Size(), maxFileSize)
	}

	return nil
}

// TruncateString helper function for content preview
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// GetSupportedExtensions returns all supported file extensions with their processors
func (dm *DocumentManager) GetSupportedExtensions() map[string]string {
	extensions := make(map[string]string)

	for ext, processor := range dm.processors {
		extensions[ext] = fmt.Sprintf("%T", processor)
	}

	return extensions
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

// PDFProcessor handles PDF files with fallback implementation
type PDFProcessor struct{}

func (p *PDFProcessor) Read(path string) (*types.DocumentContent, error) {
	// Try to use a simple fallback for now
	content, err := p.extractPDFContentBasic(path)
	if err != nil {
		return &types.DocumentContent{
			Text: "PDF content extraction not available - file detected but cannot read content",
			Type: "pdf",
			Metadata: map[string]string{
				"status": "extraction_failed",
				"error":  err.Error(),
			},
			ProcessedAt: time.Now(),
		}, nil
	}

	stat, _ := os.Stat(path)

	return &types.DocumentContent{
		Text: content,
		Type: "pdf",
		Metadata: map[string]string{
			"file_size": fmt.Sprintf("%d", stat.Size()),
			"status":    "basic_extraction",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *PDFProcessor) extractPDFContentBasic(path string) (string, error) {
	// Basic PDF content extraction placeholder
	// In real implementation, this would use a PDF library
	return fmt.Sprintf("PDF file detected: %s\nContent extraction requires PDF processing library.",
		filepath.Base(path)), nil
}

func (p *PDFProcessor) GetSupportedTypes() []string {
	return []string{"pdf"}
}

// DOCXProcessor handles Word documents with fallback implementation
type DOCXProcessor struct{}

func (p *DOCXProcessor) Read(path string) (*types.DocumentContent, error) {
	// Try to extract DOCX content
	content, err := p.extractDOCXContentBasic(path)
	if err != nil {
		return &types.DocumentContent{
			Text: "DOCX content extraction not available - file detected but cannot read content",
			Type: "docx",
			Metadata: map[string]string{
				"status": "extraction_failed",
				"error":  err.Error(),
			},
			ProcessedAt: time.Now(),
		}, nil
	}

	stat, _ := os.Stat(path)
	wordCount := len(strings.Fields(content))

	return &types.DocumentContent{
		Text: content,
		Type: "docx",
		Metadata: map[string]string{
			"word_count": fmt.Sprintf("%d", wordCount),
			"file_size":  fmt.Sprintf("%d", stat.Size()),
			"status":     "basic_extraction",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *DOCXProcessor) extractDOCXContentBasic(path string) (string, error) {
	// Basic DOCX content extraction placeholder
	// In real implementation, this would use a DOCX library
	return fmt.Sprintf("DOCX file detected: %s\nContent extraction requires DOCX processing library.",
		filepath.Base(path)), nil
}

func (p *DOCXProcessor) GetSupportedTypes() []string {
	return []string{"docx", "doc"}
}

// JSONProcessor handles JSON files
type JSONProcessor struct{}

func (p *JSONProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	text := string(content)

	// Basic JSON validation
	var jsonData interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return &types.DocumentContent{
			Text: text,
			Type: "json",
			Metadata: map[string]string{
				"status":     "invalid_json",
				"error":      err.Error(),
				"char_count": fmt.Sprintf("%d", len(text)),
			},
			ProcessedAt: time.Now(),
		}, nil
	}

	// Count JSON elements
	lineCount := len(strings.Split(text, "\n"))

	return &types.DocumentContent{
		Text: text,
		Type: "json",
		Metadata: map[string]string{
			"line_count": fmt.Sprintf("%d", lineCount),
			"char_count": fmt.Sprintf("%d", len(text)),
			"status":     "valid_json",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *JSONProcessor) GetSupportedTypes() []string {
	return []string{"json"}
}

// XMLProcessor handles XML files
type XMLProcessor struct{}

func (p *XMLProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	text := string(content)

	// Basic XML validation
	decoder := xml.NewDecoder(strings.NewReader(text))
	elementCount := 0
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &types.DocumentContent{
				Text: text,
				Type: "xml",
				Metadata: map[string]string{
					"status":     "invalid_xml",
					"error":      err.Error(),
					"char_count": fmt.Sprintf("%d", len(text)),
				},
				ProcessedAt: time.Now(),
			}, nil
		}
		elementCount++
	}

	return &types.DocumentContent{
		Text: text,
		Type: "xml",
		Metadata: map[string]string{
			"element_count": fmt.Sprintf("%d", elementCount),
			"char_count":    fmt.Sprintf("%d", len(text)),
			"status":        "valid_xml",
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *XMLProcessor) GetSupportedTypes() []string {
	return []string{"xml"}
}

// FileTypeDetector helps detect file types (basic implementation)
func DetectFileType(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	return ext, nil
}

// CSVProcessor handles CSV files - ONLY DECLARATION
type CSVProcessor struct{}

func (p *CSVProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Count non-empty lines
	actualLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			actualLines++
		}
	}

	// Estimate columns from first line
	columns := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		columns = len(strings.Split(lines[0], ","))
	}

	return &types.DocumentContent{
		Text: text,
		Type: "csv",
		Metadata: map[string]string{
			"lines":          fmt.Sprintf("%d", actualLines),
			"columns":        fmt.Sprintf("%d", columns),
			"estimated_rows": fmt.Sprintf("%d", actualLines-1), // minus header
			"char_count":     fmt.Sprintf("%d", len(text)),
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *CSVProcessor) GetSupportedTypes() []string {
	return []string{"csv"}
}

// LogProcessor handles log files - ONLY DECLARATION
type LogProcessor struct{}

func (p *LogProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Count different log levels
	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "error") || strings.Contains(lower, "err") {
			errorCount++
		} else if strings.Contains(lower, "warning") || strings.Contains(lower, "warn") {
			warningCount++
		} else if strings.Contains(lower, "info") {
			infoCount++
		}
	}

	return &types.DocumentContent{
		Text: text,
		Type: "log",
		Metadata: map[string]string{
			"total_lines":   fmt.Sprintf("%d", len(lines)),
			"error_lines":   fmt.Sprintf("%d", errorCount),
			"warning_lines": fmt.Sprintf("%d", warningCount),
			"info_lines":    fmt.Sprintf("%d", infoCount),
			"char_count":    fmt.Sprintf("%d", len(text)),
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *LogProcessor) GetSupportedTypes() []string {
	return []string{"log", "logs"}
}

// CodeProcessor handles source code files - ONLY DECLARATION
type CodeProcessor struct{}

func (p *CodeProcessor) Read(path string) (*types.DocumentContent, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read code file: %w", err)
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Count code statistics
	codeLines := 0
	commentLines := 0
	emptyLines := 0

	ext := strings.ToLower(filepath.Ext(path))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			emptyLines++
		} else if p.isCommentLine(trimmed, ext) {
			commentLines++
		} else {
			codeLines++
		}
	}

	return &types.DocumentContent{
		Text: text,
		Type: "code",
		Metadata: map[string]string{
			"total_lines":   fmt.Sprintf("%d", len(lines)),
			"code_lines":    fmt.Sprintf("%d", codeLines),
			"comment_lines": fmt.Sprintf("%d", commentLines),
			"empty_lines":   fmt.Sprintf("%d", emptyLines),
			"language":      p.detectLanguage(ext),
			"char_count":    fmt.Sprintf("%d", len(text)),
		},
		ProcessedAt: time.Now(),
	}, nil
}

func (p *CodeProcessor) isCommentLine(line, ext string) bool {
	switch ext {
	case ".go", ".js", ".java", ".c", ".cpp", ".cs":
		return strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*")
	case ".py", ".sh", ".bash":
		return strings.HasPrefix(line, "#")
	case ".html", ".xml":
		return strings.HasPrefix(line, "<!--")
	default:
		return strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#")
	}
}

func (p *CodeProcessor) detectLanguage(ext string) string {
	languages := map[string]string{
		".go":   "Go",
		".py":   "Python",
		".js":   "JavaScript",
		".java": "Java",
		".c":    "C",
		".cpp":  "C++",
		".cs":   "C#",
		".php":  "PHP",
		".rb":   "Ruby",
		".sh":   "Shell",
		".bash": "Bash",
		".sql":  "SQL",
		".html": "HTML",
		".css":  "CSS",
		".xml":  "XML",
	}

	if lang, exists := languages[ext]; exists {
		return lang
	}
	return "Unknown"
}

func (p *CodeProcessor) GetSupportedTypes() []string {
	return []string{"go", "py", "js", "java", "c", "cpp", "cs", "php", "rb", "sh", "bash", "sql", "css"}
}

// SearchInDocument searches for text within a document
func (dm *DocumentManager) SearchInDocument(path, query string) ([]string, error) {
	log.Printf("🔍 Searching in document: %s for: %s", filepath.Base(path), query)

	content, err := dm.ProcessDocument(path)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}

	var matches []string
	lines := strings.Split(content.Text, "\n")

	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			// Add context: line number and content
			match := fmt.Sprintf("Line %d: %s", i+1, strings.TrimSpace(line))
			matches = append(matches, match)
		}
	}

	log.Printf("✅ Found %d matches in %s", len(matches), filepath.Base(path))
	return matches, nil
}

// SearchInMultipleDocuments searches for text in multiple documents
func (dm *DocumentManager) SearchInMultipleDocuments(paths []string, query string) (map[string][]string, error) {
	log.Printf("🔍 Searching in %d documents for: %s", len(paths), query)

	results := make(map[string][]string)

	for _, path := range paths {
		matches, err := dm.SearchInDocument(path, query)
		if err != nil {
			log.Printf("❌ Error searching %s: %v", filepath.Base(path), err)
			continue
		}

		if len(matches) > 0 {
			results[path] = matches
		}
	}

	log.Printf("✅ Search completed. Found matches in %d out of %d documents", len(results), len(paths))
	return results, nil
}

// GetDocumentPreview returns a preview of document content
func (dm *DocumentManager) GetDocumentPreview(path string, maxLines int) (string, error) {
	content, err := dm.ProcessDocument(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(content.Text, "\n")
	if len(lines) <= maxLines {
		return content.Text, nil
	}

	preview := strings.Join(lines[:maxLines], "\n")
	preview += fmt.Sprintf("\n... (%d more lines)", len(lines)-maxLines)

	return preview, nil
}
