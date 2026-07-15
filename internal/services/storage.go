package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"quickshare/internal/models"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("snippet not found or expired")
)

// StorageService manages loading and saving snippets.
type StorageService struct {
	dataDir string
	mu      sync.RWMutex
}

// NewStorageService creates a new StorageService and ensures the directory exists.
func NewStorageService(dataDir string) (*StorageService, error) {
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	service := &StorageService{
		dataDir: dataDir,
	}

	// Start a background worker to periodically clean up expired snippets.
	go service.startCleanupWorker(5 * time.Minute)

	return service, nil
}

// GenerateID generates a secure, URL-safe random ID.
func (s *StorageService) GenerateID() (string, error) {
	bytes := make([]byte, 6) // 12 characters in hex
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Save stores the snippet as a JSON file.
func (s *StorageService) Save(snippet *models.Snippet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.dataDir, snippet.ID+".json")
	data, err := json.MarshalIndent(snippet, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snippet: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write snippet file: %w", err)
	}

	return nil
}

// Get loads a snippet and deletes it if expired.
func (s *StorageService) Get(id string) (*models.Snippet, error) {
	s.mu.RLock()
	filePath := filepath.Join(s.dataDir, id+".json")
	data, err := os.ReadFile(filePath)
	s.mu.RUnlock()

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to read snippet file: %w", err)
	}

	var snippet models.Snippet
	if err := json.Unmarshal(data, &snippet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snippet: %w", err)
	}

	// Check for expiration
	if snippet.IsExpired() {
		// Acquire write lock to delete the expired file
		s.mu.Lock()
		_ = os.Remove(filePath)
		s.mu.Unlock()
		return nil, ErrNotFound
	}

	return &snippet, nil
}

// Delete removes a snippet file.
func (s *StorageService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.dataDir, id+".json")
	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete snippet: %w", err)
	}
	return nil
}

// CleanExpired scans the directory and removes all expired snippets.
func (s *StorageService) CleanExpired() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := os.ReadDir(s.dataDir)
	if err != nil {
		log.Printf("Error reading data directory for cleanup: %v", err)
		return 0
	}

	deletedCount := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(s.dataDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var snippet models.Snippet
		if err := json.Unmarshal(data, &snippet); err != nil {
			// Corrupt JSON or invalid format: remove it
			_ = os.Remove(filePath)
			deletedCount++
			continue
		}

		if snippet.IsExpired() {
			_ = os.Remove(filePath)
			deletedCount++
		}
	}

	return deletedCount
}

// startCleanupWorker runs a background loop to clean up expired snippets.
func (s *StorageService) startCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		count := s.CleanExpired()
		if count > 0 {
			log.Printf("Background cleanup: deleted %d expired snippets", count)
		}
	}
}
