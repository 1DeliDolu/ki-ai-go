package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/internal/processors"
)

// SearchOptions defines search parameters
type SearchOptions struct {
	CaseSensitive bool `json:"case_sensitive"`
	WholeWords    bool `json:"whole_words"`
	UseRegex      bool `json:"use_regex"` // Added missing field
	MaxMatches    int  `json:"max_matches"`
	ContextLines  int  `json:"context_lines"`
}

// SearchResult represents search results for a document
type SearchResult struct {
	FilePath     string    `json:"file_path"`
	FileName     string    `json:"file_name"`
	Matches      []Match   `json:"matches"`
	TotalMatches int       `json:"total_matches"`
	ProcessedAt  time.Time `json:"processed_at"`
}

// Match represents a single search match
type Match struct {
	LineNumber int    `json:"line_number"`
	Content    string `json:"content"`
	Context    string `json:"context"`
}

// DocumentSearcher provides document search functionality
type DocumentSearcher struct {
	manager *processors.DocumentManager
}

// NewDocumentSearcher creates a new document searcher
func NewDocumentSearcher() *DocumentSearcher {
	return &DocumentSearcher{
		manager: processors.NewDocumentManager(),
	}
}

// SearchInMultipleDocuments searches for a query in multiple documents
func (ds *DocumentSearcher) SearchInMultipleDocuments(paths []string, query string, options SearchOptions) (map[string]*SearchResult, error) {
	log.Printf("🔍 Searching in %d documents for: %s", len(paths), query)

	results := make(map[string]*SearchResult)

	for _, path := range paths {
		result, err := ds.SearchInDocument(path, query, options)
		if err != nil {
			log.Printf("❌ Error searching %s: %v", path, err)
			continue
		}

		if result.TotalMatches > 0 {
			results[path] = result
		}
	}

	log.Printf("✅ Search completed. Found matches in %d out of %d documents", len(results), len(paths))
	return results, nil
}

// SearchInDocument searches for a query within a single document
func (ds *DocumentSearcher) SearchInDocument(path, query string, options SearchOptions) (*SearchResult, error) {
	log.Printf("🔍 Searching in document: %s for query: %s", filepath.Base(path), query)

	// Process the document
	content, err := ds.manager.ProcessDocument(path)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}

	// Perform search
	matches := ds.searchInText(content.Text, query, options)

	result := &SearchResult{
		FilePath:     path,
		FileName:     filepath.Base(path),
		Matches:      matches,
		TotalMatches: len(matches),
		ProcessedAt:  time.Now(),
	}

	log.Printf("✅ Found %d matches in %s", len(matches), filepath.Base(path))
	return result, nil
}

// SearchByFileType searches in documents of specific types
func (ds *DocumentSearcher) SearchByFileType(basePath, fileType, query string, options SearchOptions) (map[string]*SearchResult, error) {
	// This would require a file system walker - simplified implementation
	log.Printf("🔍 Searching by file type: %s in %s", fileType, basePath)

	// For now, return empty results - would need file system traversal
	results := make(map[string]*SearchResult)
	return results, nil
}

// SearchWithMetadata searches in both content and metadata
func (ds *DocumentSearcher) SearchWithMetadata(paths []string, query string, options SearchOptions) (map[string]*SearchResult, error) {
	log.Printf("🔍 Searching with metadata in %d documents", len(paths))

	results := make(map[string]*SearchResult)

	for _, path := range paths {
		// Process the document
		content, err := ds.manager.ProcessDocument(path)
		if err != nil {
			continue
		}

		// Search in content
		contentMatches := ds.searchInText(content.Text, query, options)

		// Search in metadata - create Match objects properly
		var metadataMatches []Match
		for key, value := range content.Metadata {
			if ds.matchesQuery(key+": "+value, query, options) {
				metadataMatches = append(metadataMatches, Match{
					LineNumber: 0, // Metadata doesn't have line numbers
					Content:    fmt.Sprintf("[META] %s: %s", key, value),
					Context:    fmt.Sprintf("Metadata field: %s", key),
				})
			}
		}

		// Combine results
		allMatches := append(contentMatches, metadataMatches...)
		if len(allMatches) > 0 {
			results[path] = &SearchResult{
				FilePath:     path,
				FileName:     filepath.Base(path),
				Matches:      allMatches,
				TotalMatches: len(allMatches),
				ProcessedAt:  time.Now(),
			}
		}
	}

	return results, nil
}

// searchInText performs the actual text search
func (ds *DocumentSearcher) searchInText(text, query string, options SearchOptions) []Match {
	var matches []Match

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if ds.matchesQuery(line, query, options) {
			// Extract context around match
			context := ds.extractContext(lines, i, options.ContextLines)
			matches = append(matches, Match{
				LineNumber: i + 1, // 1-based line numbers
				Content:    line,
				Context:    context,
			})

			// Check max matches limit
			if options.MaxMatches > 0 && len(matches) >= options.MaxMatches {
				break
			}
		}
	}

	return matches
}

// matchesQuery checks if a line matches the search query
func (ds *DocumentSearcher) matchesQuery(line, query string, options SearchOptions) bool {
	searchLine := line
	searchQuery := query

	// Handle case sensitivity
	if !options.CaseSensitive {
		searchLine = strings.ToLower(searchLine)
		searchQuery = strings.ToLower(searchQuery)
	}

	// Handle regex search
	if options.UseRegex {
		regex, err := regexp.Compile(searchQuery)
		if err != nil {
			return false
		}
		return regex.MatchString(searchLine)
	}

	// Handle whole words
	if options.WholeWords {
		regex, err := regexp.Compile(`\b` + regexp.QuoteMeta(searchQuery) + `\b`)
		if err != nil {
			return strings.Contains(searchLine, searchQuery)
		}
		return regex.MatchString(searchLine)
	}

	// Simple substring search
	return strings.Contains(searchLine, searchQuery)
}

// extractContext extracts context lines around a match
func (ds *DocumentSearcher) extractContext(lines []string, matchIndex, contextLines int) string {
	start := matchIndex - contextLines
	if start < 0 {
		start = 0
	}

	end := matchIndex + contextLines + 1
	if end > len(lines) {
		end = len(lines)
	}

	contextSlice := lines[start:end]

	// Mark the actual match line
	for i := range contextSlice {
		if start+i == matchIndex {
			contextSlice[i] = ">>> " + contextSlice[i]
		} else {
			contextSlice[i] = "    " + contextSlice[i]
		}
	}

	return strings.Join(contextSlice, "\n")
}

// GetSearchStatistics returns search statistics
func (ds *DocumentSearcher) GetSearchStatistics(results map[string]*SearchResult) map[string]interface{} {
	totalFiles := len(results)
	totalMatches := 0

	fileTypes := make(map[string]int)

	for path, result := range results {
		totalMatches += result.TotalMatches
		ext := strings.ToLower(filepath.Ext(path))
		if ext != "" {
			ext = ext[1:] // Remove dot
		} else {
			ext = "no_extension"
		}
		fileTypes[ext]++
	}

	return map[string]interface{}{
		"total_files_searched":     totalFiles,
		"total_matches":            totalMatches,
		"file_types":               fileTypes,
		"average_matches_per_file": float64(totalMatches) / float64(totalFiles),
	}
}

// HighlightMatches adds HTML highlighting to search results
func (ds *DocumentSearcher) HighlightMatches(text, query string, options SearchOptions) string {
	if options.UseRegex {
		regex, err := regexp.Compile(query)
		if err != nil {
			return text
		}
		return regex.ReplaceAllStringFunc(text, func(match string) string {
			return fmt.Sprintf("<mark>%s</mark>", match)
		})
	}

	searchQuery := query
	if !options.CaseSensitive {
		// For case-insensitive search, we need to find actual matches
		regex, err := regexp.Compile("(?i)" + regexp.QuoteMeta(query))
		if err != nil {
			return text
		}
		return regex.ReplaceAllStringFunc(text, func(match string) string {
			return fmt.Sprintf("<mark>%s</mark>", match)
		})
	}

	return strings.ReplaceAll(text, searchQuery, fmt.Sprintf("<mark>%s</mark>", searchQuery))
}
