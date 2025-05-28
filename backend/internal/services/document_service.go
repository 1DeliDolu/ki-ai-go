package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1DeliDolu/ki-ai-go/internal/config"
	"github.com/1DeliDolu/ki-ai-go/internal/processors"
	"github.com/1DeliDolu/ki-ai-go/internal/storage"
	"github.com/1DeliDolu/ki-ai-go/internal/utils"
	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

type DocumentService struct {
	memDB           *storage.MemoryDB
	config          *config.Config
	documentManager *processors.DocumentManager
}

func NewDocumentService(db interface{}, cfg *config.Config) *DocumentService {
	// Convert to memory DB
	memDB, ok := db.(*storage.MemoryDB)
	if !ok {
		log.Println("⚠️  Warning: Using memory database fallback")
		memDB = storage.InitMemoryDB()
	}

	// Ensure both directories exist
	if err := os.MkdirAll(cfg.UploadsPath, 0755); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}

	// Create test_documents directory for frontend uploads
	if err := os.MkdirAll(cfg.TestDocumentsPath, 0755); err != nil {
		log.Printf("Warning: Failed to create test_documents directory: %v", err)
	}

	return &DocumentService{
		memDB:           memDB,
		config:          cfg,
		documentManager: processors.NewDocumentManager(),
	}
}

// ConvertDocument converts a document to specified format
func (s *DocumentService) ConvertDocument(documentID, format, outputPath string) error {
	doc, err := s.memDB.GetDocument(documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	converter := utils.NewDocumentConverter()

	switch strings.ToLower(format) {
	case "markdown", "md":
		return converter.ConvertToMarkdown(doc.Path, outputPath)
	case "html":
		return converter.ConvertToHTML(doc.Path, outputPath)
	case "txt", "text":
		return converter.ConvertToPlainText(doc.Path, outputPath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// SearchInDocumentContent searches within a specific document
func (s *DocumentService) SearchInDocumentContent(documentID, query string) ([]string, error) {
	doc, err := s.memDB.GetDocument(documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	return s.documentManager.SearchInDocument(doc.Path, query)
}

// AdvancedSearch performs advanced search with options
func (s *DocumentService) AdvancedSearch(query string, options utils.SearchOptions) (map[string]*utils.SearchResult, error) {
	// Get all documents
	docs, err := s.memDB.ListDocuments()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	// Collect paths
	var paths []string
	for _, doc := range docs {
		if doc.Path != "" {
			paths = append(paths, doc.Path)
		}
	}

	// Perform search
	searcher := utils.NewDocumentSearcher()
	return searcher.SearchInMultipleDocuments(paths, query, options)
}

// GetDocumentPreview returns a preview of document content
func (s *DocumentService) GetDocumentPreview(documentID string, maxLines int) (string, error) {
	doc, err := s.memDB.GetDocument(documentID)
	if err != nil {
		return "", fmt.Errorf("document not found: %w", err)
	}

	return s.documentManager.GetDocumentPreview(doc.Path, maxLines)
}

func (s *DocumentService) ListDocuments() ([]types.Document, error) {
	log.Println("Listing documents from memory database")

	docs, err := s.memDB.ListDocuments()
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// Convert pointers to values
	result := make([]types.Document, len(docs))
	for i, doc := range docs {
		result[i] = *doc
	}

	log.Printf("Found %d documents", len(result))
	return result, nil
}

// GetDocumentContent extracts content from a document with enhanced error handling
func (s *DocumentService) GetDocumentContent(documentID string) (*types.DocumentContent, error) {
	doc, err := s.memDB.GetDocument(documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	if doc.Path == "" {
		return nil, fmt.Errorf("document path not available")
	}

	// Validate file before processing
	if err := s.documentManager.ValidateFile(doc.Path); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	content, err := s.documentManager.ProcessDocument(doc.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}

	return content, nil
}

// GetDocumentProcessingStats returns processing statistics
func (s *DocumentService) GetDocumentProcessingStats() interface{} {
	return s.documentManager.GetProcessingStats()
}

// ValidateUploadedFile validates a file before upload
func (s *DocumentService) ValidateUploadedFile(fileHeader *multipart.FileHeader) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	supportedTypes := s.documentManager.GetSupportedTypes()
	isSupported := false
	for _, supportedType := range supportedTypes {
		if ext == supportedType {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("unsupported file type: %s. Supported types: %v", ext, supportedTypes)
	}

	// Check file size (50MB limit for uploads)
	const maxUploadSize = 50 * 1024 * 1024
	if fileHeader.Size > maxUploadSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", fileHeader.Size, maxUploadSize)
	}

	return nil
}

// UploadDocument with frontend document support
func (s *DocumentService) UploadDocument(fileHeader *multipart.FileHeader) (*types.Document, error) {
	// Validate file before upload
	if err := s.ValidateUploadedFile(fileHeader); err != nil {
		return nil, err
	}

	// Determine save path - frontend uploads go to test_documents
	var savePath string
	isFromFrontend := true // Frontend'den geldiğini varsayıyoruz

	if isFromFrontend {
		// Frontend dokümanları test_documents'e kaydet
		if err := os.MkdirAll(s.config.TestDocumentsPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create test_documents directory: %w", err)
		}
		savePath = s.config.TestDocumentsPath
		log.Printf("📁 Saving frontend document to test_documents: %s", fileHeader.Filename)
	} else {
		// API dokümanları uploads'e kaydet
		if err := os.MkdirAll(s.config.UploadsPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create uploads directory: %w", err)
		}
		savePath = s.config.UploadsPath
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create unique filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, fileHeader.Filename)
	filePath := filepath.Join(savePath, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create document with enhanced metadata
	doc := &types.Document{
		Name:       fileHeader.Filename,
		Type:       filepath.Ext(fileHeader.Filename),
		Size:       fileHeader.Size,
		UploadDate: time.Now().Format("2006-01-02 15:04:05"),
		Status:     "ready",
		Path:       filePath,
	}

	// Add metadata about storage location
	doc.Metadata = map[string]string{
		"storage_location": func() string {
			if isFromFrontend {
				return "test_documents"
			}
			return "uploads"
		}(),
		"original_filename": fileHeader.Filename,
		"saved_filename":    filename,
		"upload_source":     "frontend",
	}

	// Save to memory database
	if err := s.memDB.CreateDocument(doc); err != nil {
		return nil, fmt.Errorf("failed to save to database: %w", err)
	}

	log.Printf("✅ Document uploaded successfully: %s -> %s", doc.Name, filePath)
	return doc, nil
}

// GetTestDocuments returns documents from test_documents folder
func (s *DocumentService) GetTestDocuments() ([]types.Document, error) {
	docs, err := s.memDB.ListDocuments()
	if err != nil {
		return nil, err
	}

	var testDocs []types.Document
	for _, doc := range docs {
		if doc.Metadata != nil && doc.Metadata["storage_location"] == "test_documents" {
			testDocs = append(testDocs, *doc)
		}
	}

	log.Printf("Found %d documents in test_documents", len(testDocs))
	return testDocs, nil
}

// CleanupTestDocuments cleans up test_documents folder
func (s *DocumentService) CleanupTestDocuments() error {
	log.Println("🧹 Cleaning up test_documents folder...")

	// Get all test documents
	testDocs, err := s.GetTestDocuments()
	if err != nil {
		return fmt.Errorf("failed to get test documents: %w", err)
	}

	// Delete each document
	for _, doc := range testDocs {
		if err := s.DeleteDocument(doc.ID); err != nil {
			log.Printf("Warning: Failed to delete test document %s: %v", doc.Name, err)
		}
	}

	// Clean the directory
	if err := os.RemoveAll(s.config.TestDocumentsPath); err != nil {
		return fmt.Errorf("failed to remove test_documents directory: %w", err)
	}

	// Recreate the directory
	if err := os.MkdirAll(s.config.TestDocumentsPath, 0755); err != nil {
		return fmt.Errorf("failed to recreate test_documents directory: %w", err)
	}

	log.Printf("✅ Cleaned up %d test documents", len(testDocs))
	return nil
}

func (s *DocumentService) SearchDocuments(query string) ([]types.Document, error) {
	log.Printf("🔍 Searching documents for query: '%s'", query)

	// Get all documents from memory database
	docs, err := s.memDB.ListDocuments()
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// Filter documents based on search query
	var matchedDocs []*types.Document
	for _, doc := range docs {
		matched := false

		// Search in document name (case-insensitive)
		if containsIgnoreCase(doc.Name, query) {
			matched = true
		}

		// Search in document type
		if containsIgnoreCase(doc.Type, query) {
			matched = true
		}

		// Search in actual file content if query is specific
		if !matched && doc.Path != "" {
			if content, err := os.ReadFile(doc.Path); err == nil {
				if containsIgnoreCase(string(content), query) {
					matched = true
					log.Printf("📄 Content match found in %s", doc.Name)
				}
			}
		}

		// If no specific query, include all documents (for demo.txt case)
		if query == "" || strings.TrimSpace(query) == "" {
			matched = true
		}

		if matched {
			matchedDocs = append(matchedDocs, doc)
		}
	}

	// Convert pointers to values
	result := make([]types.Document, len(matchedDocs))
	for i, doc := range matchedDocs {
		result[i] = *doc
	}

	log.Printf("✅ Found %d documents matching query '%s'", len(result), query)
	return result, nil
}

// Helper function for case-insensitive string matching
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		len(substr) > 0 &&
		(s == substr ||
			strings.ToLower(s) == strings.ToLower(substr) ||
			strings.Contains(strings.ToLower(s), strings.ToLower(substr)))
}

func (s *DocumentService) DeleteDocument(idStr string) error {
	log.Printf("Deleting document with ID: %s", idStr)

	// Get document info first
	doc, err := s.memDB.GetDocument(idStr)
	if err != nil {
		return fmt.Errorf("document with id %s not found: %w", idStr, err)
	}

	// Delete from memory database
	if err := s.memDB.DeleteDocument(idStr); err != nil {
		return fmt.Errorf("failed to delete document from database: %w", err)
	}

	// Delete file from filesystem if path exists
	if doc.Path != "" {
		if err := os.Remove(doc.Path); err != nil {
			// Log the error but don't fail the operation
			// since the database record is already deleted
			log.Printf("Warning: failed to delete file %s: %v", doc.Path, err)
		} else {
			log.Printf("Successfully deleted file: %s", doc.Path)
		}
	}

	log.Printf("Successfully deleted document: %s", doc.Name)
	return nil
}

// GetSupportedDocumentTypes returns all supported document types
func (s *DocumentService) GetSupportedDocumentTypes() []string {
	return s.documentManager.GetSupportedTypes()
}

// GetDocument returns a document by ID
func (s *DocumentService) GetDocument(documentID string) (*types.Document, error) {
	return s.memDB.GetDocument(documentID)
}

// GetDocumentFileInfo returns comprehensive file information
func (s *DocumentService) GetDocumentFileInfo(documentID string) (*utils.FileInfo, error) {
	doc, err := s.memDB.GetDocument(documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	if doc.Path == "" {
		return nil, fmt.Errorf("document path not available")
	}

	// Get document content
	content, err := s.documentManager.ProcessDocument(doc.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}

	// Get comprehensive file info
	return utils.GetFileInfo(doc.Path, content)
}

// GetDocumentAnalysis provides content analysis
func (s *DocumentService) GetDocumentAnalysis(documentID string) (map[string]interface{}, error) {
	content, err := s.GetDocumentContent(documentID)
	if err != nil {
		return nil, err
	}

	analysis := utils.AnalyzeContent(content.Text)
	analysis["processing_metadata"] = content.Metadata
	analysis["processed_at"] = content.ProcessedAt.Format(time.RFC3339)
	analysis["document_type"] = content.Type

	return analysis, nil
}
