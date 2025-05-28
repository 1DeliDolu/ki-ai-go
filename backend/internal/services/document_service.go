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
	"github.com/1DeliDolu/ki-ai-go/internal/handlers"
	"github.com/1DeliDolu/ki-ai-go/internal/storage"
	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

type DocumentService struct {
	memDB           *storage.MemoryDB
	config          *config.Config
	documentManager *handlers.DocumentManager
}

func NewDocumentService(db interface{}, cfg *config.Config) *DocumentService {
	// Convert to memory DB
	memDB, ok := db.(*storage.MemoryDB)
	if !ok {
		log.Println("⚠️  Warning: Using memory database fallback")
		memDB = storage.InitMemoryDB()
	}

	// Ensure uploads directory exists
	if err := os.MkdirAll(cfg.UploadsPath, 0755); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}

	return &DocumentService{
		memDB:           memDB,
		config:          cfg,
		documentManager: handlers.NewDocumentManager(),
	}
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

func (s *DocumentService) UploadDocument(fileHeader *multipart.FileHeader) (*types.Document, error) {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(s.config.UploadsPath, 0755); err != nil {
		return nil, err
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	filePath := filepath.Join(s.config.UploadsPath, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, file); err != nil {
		return nil, err
	}

	// Create document with correct fields
	doc := &types.Document{
		Name:       fileHeader.Filename,
		Type:       filepath.Ext(fileHeader.Filename),
		Size:       fileHeader.Size,
		UploadDate: time.Now().Format("2006-01-02 15:04:05"),
		Status:     "ready",
		Path:       filePath,
	}

	// Save to memory database
	if err := s.memDB.CreateDocument(doc); err != nil {
		return nil, err
	}

	log.Printf("Document uploaded successfully: %s", doc.Name)
	return doc, nil
}

func (s *DocumentService) extractTextContent(filePath, originalName string) (string, error) {
	// Use the new document manager for content extraction
	content, err := s.documentManager.ProcessDocument(filePath)
	if err != nil {
		// Fallback to old method if new processor fails
		ext := filepath.Ext(originalName)
		switch ext {
		case ".txt", ".md":
			fileContent, err := os.ReadFile(filePath)
			return string(fileContent), err
		default:
			return "", fmt.Errorf("unsupported file type: %s", ext)
		}
	}

	return content.Text, nil
}

// GetDocumentContent extracts content from a document
func (s *DocumentService) GetDocumentContent(documentID string) (*types.DocumentContent, error) {
	doc, err := s.memDB.GetDocument(documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	if doc.Path == "" {
		return nil, fmt.Errorf("document path not available")
	}

	content, err := s.documentManager.ProcessDocument(doc.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}

	return content, nil
}

// GetSupportedDocumentTypes returns all supported document types
func (s *DocumentService) GetSupportedDocumentTypes() []string {
	return s.documentManager.GetSupportedTypes()
}

func (s *DocumentService) SearchDocuments(query string) ([]types.Document, error) {
	// Simple text search in content
	log.Printf("Searching documents for query: '%s'", query)

	// Get all documents from memory database
	docs, err := s.memDB.ListDocuments()
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// Filter documents based on search query
	var matchedDocs []*types.Document
	for _, doc := range docs {
		// Search in document name (case-insensitive)
		if containsIgnoreCase(doc.Name, query) {
			matchedDocs = append(matchedDocs, doc)
			continue
		}

		// Search in document type
		if containsIgnoreCase(doc.Type, query) {
			matchedDocs = append(matchedDocs, doc)
			continue
		}

		// TODO: Add content search when we implement text extraction
	}

	// Convert pointers to values
	result := make([]types.Document, len(matchedDocs))
	for i, doc := range matchedDocs {
		result[i] = *doc
	}

	log.Printf("Found %d documents matching query '%s'", len(result), query)
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
