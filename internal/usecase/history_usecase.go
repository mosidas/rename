package usecase

import (
	"rename/internal/domain"
)

// HistoryRepository defines history persistence operations
// Following ISP (Interface Segregation Principle)
type HistoryRepository interface {
	Save(history *domain.History) error
	Load() (*domain.History, error)
}

// HistoryUseCase handles history operations
// Following SRP and DIP
type HistoryUseCase struct {
	repository HistoryRepository
	history    *domain.History
}

// NewHistoryUseCase creates a new HistoryUseCase
func NewHistoryUseCase(repository HistoryRepository) *HistoryUseCase {
	// Try to load existing history
	history, err := repository.Load()
	if err != nil || history == nil {
		// If no history exists, create new one
		history = domain.NewHistory()
	}

	return &HistoryUseCase{
		repository: repository,
		history:    history,
	}
}

// AddEntry adds a new history entry and persists it
func (uc *HistoryUseCase) AddEntry(entry domain.HistoryEntry) error {
	uc.history.Add(entry)
	return uc.repository.Save(uc.history)
}

// GetHistory returns all history entries
func (uc *HistoryUseCase) GetHistory() ([]domain.HistoryEntry, error) {
	// Reload from repository to get latest
	history, err := uc.repository.Load()
	if err != nil {
		return nil, err
	}

	uc.history = history
	return history.GetAll(), nil
}

// ClearHistory removes all history entries
func (uc *HistoryUseCase) ClearHistory() error {
	uc.history.Clear()
	return uc.repository.Save(uc.history)
}
