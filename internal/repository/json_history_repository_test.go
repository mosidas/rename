package repository

import (
	"os"
	"path/filepath"
	"testing"

	"rename/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestJSONHistoryRepository_SaveAndLoad(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "rename-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	repo := NewJSONHistoryRepository(configPath)

	// Create history with entries
	history := domain.NewHistory()
	history.Add(domain.HistoryEntry{
		Pattern:     "test1",
		Replacement: "TEST1",
		IsRegex:     false,
		CaseInsensitive: false,
	})
	history.Add(domain.HistoryEntry{
		Pattern:     "test2",
		Replacement: "TEST2",
		IsRegex:     true,
		CaseInsensitive: true,
	})

	// Save
	err = repo.Save(history)
	assert.NoError(t, err)

	// Verify file exists
	assert.FileExists(t, configPath)

	// Load
	loaded, err := repo.Load()
	assert.NoError(t, err)
	assert.NotNil(t, loaded)

	// Verify loaded data
	entries := loaded.GetAll()
	assert.Equal(t, 2, len(entries))
	assert.Equal(t, "test2", entries[0].Pattern) // Most recent first
	assert.Equal(t, "TEST2", entries[0].Replacement)
	assert.True(t, entries[0].IsRegex)
	assert.True(t, entries[0].CaseInsensitive)
}

func TestJSONHistoryRepository_Load_FileNotExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rename-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "nonexistent.json")
	repo := NewJSONHistoryRepository(configPath)

	// Load should return empty history when file doesn't exist
	loaded, err := repo.Load()
	assert.NoError(t, err)
	assert.NotNil(t, loaded)
	assert.Equal(t, 0, loaded.Count())
}

func TestJSONHistoryRepository_Save_CreatesDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rename-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Config path in non-existent subdirectory
	configPath := filepath.Join(tmpDir, "subdir", "config.json")
	repo := NewJSONHistoryRepository(configPath)

	history := domain.NewHistory()
	history.Add(domain.HistoryEntry{
		Pattern:     "test",
		Replacement: "TEST",
		IsRegex:     false,
	})

	// Save should create directory
	err = repo.Save(history)
	assert.NoError(t, err)

	// Verify file exists
	assert.FileExists(t, configPath)
}

func TestJSONHistoryRepository_Save_MaxHistoryLimit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rename-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	repo := NewJSONHistoryRepository(configPath)

	history := domain.NewHistory()

	// Add 150 entries
	for i := 0; i < 150; i++ {
		history.Add(domain.HistoryEntry{
			Pattern:     "pattern" + string(rune('A'+i)),
			Replacement: "replacement",
			IsRegex:     false,
		})
	}

	// Save
	err = repo.Save(history)
	assert.NoError(t, err)

	// Load and verify only 100 entries
	loaded, err := repo.Load()
	assert.NoError(t, err)
	assert.Equal(t, 100, loaded.Count())
}
