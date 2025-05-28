package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/1DeliDolu/ki-ai-go/pkg/types"
)

// MemoryDB implements a simple in-memory database using maps and slices
type MemoryDB struct {
	mu           sync.RWMutex
	users        map[int]*User
	prompts      map[int]*Prompt
	documents    map[string]*types.Document
	models       map[string]*types.Model
	chunks       map[string][]*types.DocumentChunk
	nextID       int
	nextUserID   int
	nextPromptID int
}

// User represents a user in the system
type User struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// Prompt represents a prompt and its answer
type Prompt struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	PromptText string `json:"prompt_text"`
	AnswerText string `json:"answer_text"`
	CreatedAt  string `json:"created_at"`
}

// NewMemoryDB creates a new in-memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		users:        make(map[int]*User),
		prompts:      make(map[int]*Prompt),
		documents:    make(map[string]*types.Document),
		models:       make(map[string]*types.Model),
		chunks:       make(map[string][]*types.DocumentChunk),
		nextID:       1,
		nextUserID:   1,
		nextPromptID: 1,
	}
}

// Implement sql.DB interface methods we need
func (db *MemoryDB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Clear all data
	db.documents = make(map[string]*types.Document)
	db.models = make(map[string]*types.Model)
	db.chunks = make(map[string][]*types.DocumentChunk)
	db.users = make(map[int]*User)
	db.prompts = make(map[int]*Prompt)
	db.nextID = 1
	db.nextUserID = 1
	db.nextPromptID = 1

	log.Println("Memory database closed and cleared")
	return nil
}

func (db *MemoryDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	// For memory DB, most exec operations are no-ops or handled internally
	log.Printf("Memory DB Exec (no-op): %s", query)
	return &memoryResult{}, nil
}

func (db *MemoryDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Memory DB doesn't use SQL queries
	return nil, fmt.Errorf("memory DB doesn't support SQL queries")
}

// Document operations
func (db *MemoryDB) CreateDocument(doc *types.Document) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if doc.ID == "" {
		doc.ID = fmt.Sprintf("%d", db.nextID)
		db.nextID++
	}

	if doc.UploadDate == "" {
		doc.UploadDate = time.Now().Format(time.RFC3339)
	}

	db.documents[doc.ID] = doc
	log.Printf("Document created: %s (%s)", doc.Name, doc.ID)
	return nil
}

func (db *MemoryDB) GetDocument(id string) (*types.Document, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	doc, exists := db.documents[id]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	// Return a copy
	docCopy := *doc
	return &docCopy, nil
}

func (db *MemoryDB) ListDocuments() ([]*types.Document, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	docs := make([]*types.Document, 0, len(db.documents))
	for _, doc := range db.documents {
		// Return copies
		docCopy := *doc
		docs = append(docs, &docCopy)
	}

	log.Printf("Listed %d documents", len(docs))
	return docs, nil
}

func (db *MemoryDB) DeleteDocument(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.documents[id]; !exists {
		return fmt.Errorf("document not found: %s", id)
	}

	delete(db.documents, id)
	delete(db.chunks, id) // Also delete associated chunks
	log.Printf("Document deleted: %s", id)
	return nil
}

// Model operations
func (db *MemoryDB) CreateModel(model *types.Model) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.models[model.ID] = model
	log.Printf("Model created: %s", model.ID)
	return nil
}

func (db *MemoryDB) GetModel(id string) (*types.Model, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	model, exists := db.models[id]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", id)
	}

	// Return a copy
	modelCopy := *model
	return &modelCopy, nil
}

func (db *MemoryDB) ListModels() ([]*types.Model, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	models := make([]*types.Model, 0, len(db.models))
	for _, model := range db.models {
		// Return copies
		modelCopy := *model
		models = append(models, &modelCopy)
	}

	return models, nil
}

// Chunk operations
func (db *MemoryDB) CreateChunk(chunk *types.DocumentChunk) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if chunk.ID == "" {
		chunk.ID = fmt.Sprintf("chunk_%d", db.nextID)
		db.nextID++
	}

	db.chunks[chunk.DocumentID] = append(db.chunks[chunk.DocumentID], chunk)
	log.Printf("Chunk created for document: %s", chunk.DocumentID)
	return nil
}

func (db *MemoryDB) GetChunks(documentID string) ([]*types.DocumentChunk, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	chunks := db.chunks[documentID]
	if chunks == nil {
		return []*types.DocumentChunk{}, nil
	}

	// Return copies
	result := make([]*types.DocumentChunk, len(chunks))
	for i, chunk := range chunks {
		chunkCopy := *chunk
		result[i] = &chunkCopy
	}

	return result, nil
}

// User operations
func (db *MemoryDB) CreateUser(username string) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if username already exists
	for _, user := range db.users {
		if user.Username == username {
			return nil, fmt.Errorf("username already exists: %s", username)
		}
	}

	user := &User{
		UserID:    db.nextUserID,
		Username:  username,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	db.users[db.nextUserID] = user
	db.nextUserID++

	log.Printf("User created: %s (ID: %d)", username, user.UserID)
	return user, nil
}

func (db *MemoryDB) GetUser(userID int) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, exists := db.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found: %d", userID)
	}

	userCopy := *user
	return &userCopy, nil
}

// Prompt operations
func (db *MemoryDB) CreatePrompt(userID int, promptText, answerText string) (*Prompt, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if user exists
	if _, exists := db.users[userID]; !exists {
		return nil, fmt.Errorf("user not found: %d", userID)
	}

	prompt := &Prompt{
		ID:         db.nextPromptID,
		UserID:     userID,
		PromptText: promptText,
		AnswerText: answerText,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	db.prompts[db.nextPromptID] = prompt
	db.nextPromptID++

	log.Printf("Prompt created for user %d (ID: %d)", userID, prompt.ID)
	return prompt, nil
}

func (db *MemoryDB) GetUserPrompts(userID int, limit int) ([]*Prompt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var userPrompts []*Prompt
	count := 0

	for _, prompt := range db.prompts {
		if prompt.UserID == userID {
			if limit > 0 && count >= limit {
				break
			}
			promptCopy := *prompt
			userPrompts = append(userPrompts, &promptCopy)
			count++
		}
	}

	return userPrompts, nil
}

// Helper types for sql.Result interface
type memoryResult struct{}

func (r *memoryResult) LastInsertId() (int64, error) { return 0, nil }
func (r *memoryResult) RowsAffected() (int64, error) { return 1, nil }

// Global memory database instance
var memoryDBInstance *MemoryDB

// InitMemoryDB initializes the in-memory database
func InitMemoryDB() *MemoryDB {
	if memoryDBInstance == nil {
		memoryDBInstance = NewMemoryDB()
		log.Println("✅ Memory database initialized")
	}
	return memoryDBInstance
}
