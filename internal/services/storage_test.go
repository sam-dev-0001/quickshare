package services

import (
	"os"
	"path/filepath"
	"quickshare/internal/models"
	"testing"
	"time"
)

func TestStorageService(t *testing.T) {
	// 1. Create a isolated temporary directory for test storage
	testDir, err := os.MkdirTemp("", "quickshare_test_")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 2. Initialize the service
	service, err := NewStorageService(testDir)
	if err != nil {
		t.Fatalf("Failed to initialize storage service: %v", err)
	}

	// 3. Test GenerateID
	id, err := service.GenerateID()
	if err != nil {
		t.Fatalf("Failed to generate secure ID: %v", err)
	}
	if len(id) != 12 {
		t.Errorf("Expected ID length of 12, got %d", len(id))
	}

	// 4. Test Save & Get
	snippet := &models.Snippet{
		ID:        id,
		Text:      "fmt.Println(\"Hello World!\")",
		Language:  "go",
		Expiry:    models.ExpiryNever,
		CreatedAt: time.Now(),
	}

	err = service.Save(snippet)
	if err != nil {
		t.Fatalf("Failed to save snippet: %v", err)
	}

	retrieved, err := service.Get(id)
	if err != nil {
		t.Fatalf("Failed to retrieve saved snippet: %v", err)
	}

	if retrieved.ID != snippet.ID {
		t.Errorf("Expected ID %s, got %s", snippet.ID, retrieved.ID)
	}
	if retrieved.Text != snippet.Text {
		t.Errorf("Expected Text %q, got %q", snippet.Text, retrieved.Text)
	}
	if retrieved.Language != snippet.Language {
		t.Errorf("Expected Language %q, got %q", snippet.Language, retrieved.Language)
	}

	// 5. Test retrieving a non-existent ID
	_, err = service.Get("nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	// 6. Test Expiration on Get
	expiredID, err := service.GenerateID()
	if err != nil {
		t.Fatalf("Failed to generate secure ID: %v", err)
	}

	expiredSnippet := &models.Snippet{
		ID:        expiredID,
		Text:      "Expired snippet content",
		Language:  "plaintext",
		Expiry:    models.Expiry10Min,
		CreatedAt: time.Now().Add(-15 * time.Minute), // Already expired
	}

	err = service.Save(expiredSnippet)
	if err != nil {
		t.Fatalf("Failed to save expired snippet: %v", err)
	}

	// Fetching should fail with ErrNotFound, and the physical file should be deleted automatically
	_, err = service.Get(expiredID)
	if err != ErrNotFound {
		t.Errorf("Expected expired snippet to be deleted and return ErrNotFound, got %v", err)
	}

	// Check that file is indeed deleted
	filePath := filepath.Join(testDir, expiredID+".json")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Expected expired file to be removed from disk, but it still exists")
	}
}

func TestStorageServiceCleanExpired(t *testing.T) {
	testDir, err := os.MkdirTemp("", "quickshare_clean_test_")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	service, err := NewStorageService(testDir)
	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	// Save active snippet
	activeID, _ := service.GenerateID()
	_ = service.Save(&models.Snippet{
		ID:        activeID,
		Text:      "Active",
		Expiry:    models.ExpiryNever,
		CreatedAt: time.Now(),
	})

	// Save expired snippet
	expiredID, _ := service.GenerateID()
	_ = service.Save(&models.Snippet{
		ID:        expiredID,
		Text:      "Expired",
		Expiry:    models.Expiry10Min,
		CreatedAt: time.Now().Add(-15 * time.Minute),
	})

	// Run clean
	deleted := service.CleanExpired()
	if deleted != 1 {
		t.Errorf("Expected 1 file to be deleted by CleanExpired, got %d", deleted)
	}

	// Verify active remains
	_, err = service.Get(activeID)
	if err != nil {
		t.Errorf("Expected active snippet to still exist, got err: %v", err)
	}

	// Verify expired is gone
	_, err = service.Get(expiredID)
	if err != ErrNotFound {
		t.Errorf("Expected expired snippet to be gone, got err: %v", err)
	}
}
