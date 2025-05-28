package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

// FileInfo enhanced file information
type FileInfo struct {
	Name         string            `json:"name"`
	Size         int64             `json:"size"`
	Extension    string            `json:"extension"`
	ModifiedTime time.Time         `json:"modified_time"`
	WordCount    int               `json:"word_count"`
	LineCount    int               `json:"line_count"`
	CharCount    int               `json:"char_count"`
	Metadata     map[string]string `json:"metadata"`
}

// GetFileInfo extracts comprehensive file information
func GetFileInfo(filePath string, content *types.DocumentContent) (*FileInfo, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	words := len(strings.Fields(content.Text))
	lines := len(strings.Split(content.Text, "\n"))

	return &FileInfo{
		Name:         filepath.Base(filePath),
		Size:         stat.Size(),
		Extension:    strings.ToLower(filepath.Ext(filePath)),
		ModifiedTime: stat.ModTime(),
		WordCount:    words,
		LineCount:    lines,
		CharCount:    len(content.Text),
		Metadata:     content.Metadata,
	}, nil
}

// FormatFileSize converts bytes to human readable format
func FormatFileSize(bytes int64) string {
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

// AnalyzeContent provides content analysis
func AnalyzeContent(content string) map[string]interface{} {
	lines := strings.Split(content, "\n")
	words := strings.Fields(content)

	// Count empty lines
	emptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			emptyLines++
		}
	}

	// Find longest and shortest lines
	var maxLineLength, minLineLength int
	if len(lines) > 0 {
		maxLineLength = len(lines[0])
		minLineLength = len(lines[0])
		for _, line := range lines {
			if len(line) > maxLineLength {
				maxLineLength = len(line)
			}
			if len(line) < minLineLength && len(line) > 0 {
				minLineLength = len(line)
			}
		}
	}

	// Calculate averages
	avgLineLength := 0.0
	avgWordLength := 0.0
	if len(lines) > 0 {
		totalChars := len(content)
		avgLineLength = float64(totalChars) / float64(len(lines))
	}
	if len(words) > 0 {
		totalWordChars := 0
		for _, word := range words {
			totalWordChars += len(word)
		}
		avgWordLength = float64(totalWordChars) / float64(len(words))
	}

	return map[string]interface{}{
		"total_lines":     len(lines),
		"empty_lines":     emptyLines,
		"content_lines":   len(lines) - emptyLines,
		"total_words":     len(words),
		"total_chars":     len(content),
		"max_line_length": maxLineLength,
		"min_line_length": minLineLength,
		"avg_line_length": fmt.Sprintf("%.1f", avgLineLength),
		"avg_word_length": fmt.Sprintf("%.1f", avgWordLength),
		"has_content":     len(strings.TrimSpace(content)) > 0,
	}
}
