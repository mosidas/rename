package usecase

import (
	"errors"
	"testing"

	"rename/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHistoryRepository is a mock implementation of HistoryRepository
type MockHistoryRepository struct {
	mock.Mock
}

func (m *MockHistoryRepository) Save(history *domain.History) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockHistoryRepository) Load() (*domain.History, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.History), args.Error(1)
}

func TestHistoryUseCase_AddEntry(t *testing.T) {
	mockRepo := new(MockHistoryRepository)

	// Mock Load for initialization
	mockRepo.On("Load").Return(domain.NewHistory(), nil).Once()

	useCase := NewHistoryUseCase(mockRepo)

	entry := domain.HistoryEntry{
		Pattern:     "test",
		Replacement: "TEST",
		IsRegex:     false,
	}

	mockRepo.On("Save", mock.AnythingOfType("*domain.History")).Return(nil)

	err := useCase.AddEntry(entry)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHistoryUseCase_AddEntry_SaveError(t *testing.T) {
	mockRepo := new(MockHistoryRepository)

	// Mock Load for initialization
	mockRepo.On("Load").Return(domain.NewHistory(), nil).Once()

	useCase := NewHistoryUseCase(mockRepo)

	entry := domain.HistoryEntry{
		Pattern:     "test",
		Replacement: "TEST",
		IsRegex:     false,
	}

	mockRepo.On("Save", mock.AnythingOfType("*domain.History")).Return(errors.New("save failed"))

	err := useCase.AddEntry(entry)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save failed")
	mockRepo.AssertExpectations(t)
}

func TestHistoryUseCase_GetHistory(t *testing.T) {
	mockRepo := new(MockHistoryRepository)

	initialHistory := domain.NewHistory()
	initialHistory.Add(domain.HistoryEntry{
		Pattern:     "test",
		Replacement: "TEST",
		IsRegex:     false,
	})
	mockRepo.On("Load").Return(initialHistory, nil).Once()

	useCase := NewHistoryUseCase(mockRepo)

	// GetHistory should return cached entries, not reload
	entries, err := useCase.GetHistory()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "test", entries[0].Pattern)
	mockRepo.AssertExpectations(t)
}

func TestHistoryUseCase_GetHistory_EmptyWhenNoHistory(t *testing.T) {
	mockRepo := new(MockHistoryRepository)

	// Mock Load for initialization - returns empty history
	mockRepo.On("Load").Return(domain.NewHistory(), nil).Once()

	useCase := NewHistoryUseCase(mockRepo)

	// GetHistory should return cached empty entries, not reload
	entries, err := useCase.GetHistory()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(entries))
	mockRepo.AssertExpectations(t)
}
