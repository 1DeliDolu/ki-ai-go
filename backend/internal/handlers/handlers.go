// backend/internal/handlers/handlers.go
package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/1DeliDolu/ki-ai-go/internal/services"
	"github.com/1DeliDolu/ki-ai-go/pkg/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	modelService    *services.ModelService
	documentService *services.DocumentService
	wikiService     *services.WikiService
	aiService       *services.AIService
	cleanupService  *services.CleanupService
}

func New(modelService *services.ModelService, documentService *services.DocumentService,
	wikiService *services.WikiService, aiService *services.AIService, cleanupService *services.CleanupService) *Handler {
	return &Handler{
		modelService:    modelService,
		documentService: documentService,
		wikiService:     wikiService,
		aiService:       aiService,
		cleanupService:  cleanupService,
	}
}

// Health check
func (h *Handler) HealthCheck(c *gin.Context) {
	log.Printf("Health check requested from %s", c.ClientIP())
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"message":   "Local AI Project API is running",
	})
}

// Model handlers
func (h *Handler) ListModels(c *gin.Context) {
	log.Printf("ListModels requested from %s", c.ClientIP())

	models, err := h.modelService.ListModels()
	if err != nil {
		log.Printf("Error listing models: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Returning %d models", len(models))
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func (h *Handler) DownloadModel(c *gin.Context) {
	log.Printf("DownloadModel requested from %s", c.ClientIP())

	var req struct {
		Name string `json:"name" binding:"required"`
		URL  string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Downloading model %s from %s", req.Name, req.URL)
	if err := h.modelService.DownloadModel(req.Name, req.URL); err != nil {
		log.Printf("Error downloading model: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model downloaded successfully"})
}

func (h *Handler) LoadModel(c *gin.Context) {
	log.Printf("LoadModel requested from %s", c.ClientIP())

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Loading model %s", req.Name)

	// Load model in both model service and AI service
	if err := h.modelService.LoadModel(req.Name); err != nil {
		log.Printf("Error loading model in model service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load model in AI service for inference
	if err := h.aiService.LoadModel(req.Name); err != nil {
		log.Printf("Error loading model in AI service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Model file loaded but AI service failed to initialize: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model loaded successfully"})
}

func (h *Handler) DeleteModel(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model name is required"})
		return
	}

	if err := h.modelService.DeleteModel(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model deleted successfully"})
}

// InitializeBasicModels adds basic models to the system
func (h *Handler) InitializeBasicModels(c *gin.Context) {
	log.Printf("InitializeBasicModels requested from %s", c.ClientIP())

	if err := h.modelService.AddBasicModels(); err != nil {
		log.Printf("Error initializing basic models: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Basic models initialized successfully",
	})
}

// GetModelInfo returns detailed information about a specific model
func (h *Handler) GetModelInfo(c *gin.Context) {
	modelName := c.Param("name")
	if modelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model name is required"})
		return
	}

	model, err := h.modelService.GetModelInfo(modelName)
	if err != nil {
		log.Printf("Error getting model info: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": model,
	})
}

// GetAvailableModelTypes returns all available model types
func (h *Handler) GetAvailableModelTypes(c *gin.Context) {
	types := h.modelService.GetAvailableModelTypes()
	c.JSON(http.StatusOK, gin.H{
		"model_types": types,
	})
}

// GetModelsByType returns models filtered by type
func (h *Handler) GetModelsByType(c *gin.Context) {
	modelType := c.Param("type")
	if modelType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model type is required"})
		return
	}

	models, err := h.modelService.GetModelsByType(modelType)
	if err != nil {
		log.Printf("Error getting models by type: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
		"type":   modelType,
		"count":  len(models),
	})
}

// Document handlers
func (h *Handler) ListDocuments(c *gin.Context) {
	log.Printf("ListDocuments requested from %s", c.ClientIP())

	documents, err := h.documentService.ListDocuments()
	if err != nil {
		log.Printf("Error listing documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Returning %d documents", len(documents))
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func (h *Handler) UploadDocument(c *gin.Context) {
	log.Printf("UploadDocument requested from %s", c.ClientIP())

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	log.Printf("Uploading file: %s (%d bytes)", file.Filename, file.Size)
	document, err := h.documentService.UploadDocument(file)
	if err != nil {
		log.Printf("Error uploading document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Document uploaded successfully: ID %s", document.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Document uploaded successfully",
		"document": document,
	})
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	if err := h.documentService.DeleteDocument(idStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Document deleted successfully",
		"deletedId": idStr,
	})
}

// GetDocumentContent returns the processed content of a document
func (h *Handler) GetDocumentContent(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	content, err := h.documentService.GetDocumentContent(documentID)
	if err != nil {
		log.Printf("Error getting document content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": content,
	})
}

// GetSupportedDocumentTypes returns all supported document types
func (h *Handler) GetSupportedDocumentTypes(c *gin.Context) {
	types := h.documentService.GetSupportedDocumentTypes()
	c.JSON(http.StatusOK, gin.H{
		"supported_types": types,
	})
}

// GetDocumentProcessingStats returns document processing statistics
func (h *Handler) GetDocumentProcessingStats(c *gin.Context) {
	stats := h.documentService.GetDocumentProcessingStats()
	c.JSON(http.StatusOK, gin.H{
		"processing_stats": stats,
	})
}

// ProcessMultipleDocuments processes multiple documents in batch
func (h *Handler) ProcessMultipleDocuments(c *gin.Context) {
	var req struct {
		DocumentIDs []string `json:"document_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get document paths
	var paths []string
	for _, id := range req.DocumentIDs {
		doc, err := h.documentService.GetDocument(id)
		if err == nil && doc.Path != "" {
			paths = append(paths, doc.Path)
		}
	}

	if len(paths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid documents found"})
		return
	}

	// For now, return a placeholder response with document info
	c.JSON(http.StatusOK, gin.H{
		"message":         "Batch processing completed",
		"processed_paths": paths,
		"total_files":     len(paths),
		"note":            "Enhanced batch processing with DocumentManager coming soon",
	})
}

// Wiki handlers
func (h *Handler) SearchWiki(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	results, err := h.wikiService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// AI Query handler
func (h *Handler) Query(c *gin.Context) {
	log.Printf("Query requested from %s", c.ClientIP())

	var req types.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing query: %s", req.Query)
	startTime := time.Now()

	// Check if AI service has a model loaded
	if !h.aiService.IsModelLoaded() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No model loaded. Please load a model first."})
		return
	}

	// Search documents if requested
	var documents []types.Document
	if req.IncludeDocuments {
		docs, err := h.documentService.SearchDocuments(req.Query)
		if err == nil {
			documents = docs
		}
	}

	// Search wiki if requested
	var wikiResults []types.WikiResult
	if req.IncludeWiki {
		wiki, err := h.wikiService.Search(req.Query)
		if err == nil {
			wikiResults = wiki
		}
	}

	// Generate AI response
	response, err := h.aiService.GenerateResponse(req.Query, documents, wikiResults)
	if err != nil {
		log.Printf("Error generating AI response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate response: " + err.Error()})
		return
	}

	processingTime := time.Since(startTime).Seconds()

	result := types.QueryResponse{
		Response:       response,
		ModelUsed:      h.aiService.GetCurrentModel(),
		ProcessingTime: processingTime,
	}
	result.Sources.Documents = documents
	result.Sources.Wiki = wikiResults

	log.Printf("Query processed successfully in %.2f seconds", processingTime)
	c.JSON(http.StatusOK, result)
}

// Cleanup handlers
func (h *Handler) CleanupAll(c *gin.Context) {
	log.Printf("CleanupAll requested from %s", c.ClientIP())

	if err := h.cleanupService.CleanupAll(); err != nil {
		log.Printf("Error during cleanup: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All files cleaned up successfully"})
}

func (h *Handler) CleanupDocuments(c *gin.Context) {
	log.Printf("CleanupDocuments requested from %s", c.ClientIP())

	if err := h.cleanupService.CleanupDocuments(); err != nil {
		log.Printf("Error during document cleanup: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documents cleaned up successfully"})
}
