package repository

import (
	"encoding/json"
	"os"
	"path/filepath"

	"rename/internal/domain"
)

// JSONHistoryRepository implements HistoryRepository using JSON file storage
// Following SRP (Single Responsibility Principle) - only handles persistence
type JSONHistoryRepository struct {
	configPath string
}

// historyData is the JSON structure for persistence
type historyData struct {
	Entries []domain.HistoryEntry `json:"entries"`
}

// NewJSONHistoryRepository creates a new JSON-based history repository
func NewJSONHistoryRepository(configPath string) *JSONHistoryRepository {
	return &JSONHistoryRepository{
		configPath: configPath,
	}
}

// Save persists history to JSON file
func (r *JSONHistoryRepository) Save(history *domain.History) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Prepare data
	data := historyData{
		Entries: history.GetAll(),
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(r.configPath, jsonData, 0644)
}

// Load reads history from JSON file
func (r *JSONHistoryRepository) Load() (*domain.History, error) {
	// Check if file exists
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// Return empty history if file doesn't exist
		return domain.NewHistory(), nil
	}

	// Read file
	jsonData, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var data historyData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	// Reconstruct history
	history := domain.NewHistory()
	// Add entries in reverse order to maintain order (most recent first)
	for i := len(data.Entries) - 1; i >= 0; i-- {
		history.Add(data.Entries[i])
	}

	return history, nil
}
