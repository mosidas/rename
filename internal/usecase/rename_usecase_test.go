package usecase

import (
	"errors"
	"testing"

	"rename/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileSystemService is a mock implementation of FileSystemService
type MockFileSystemService struct {
	mock.Mock
}

func (m *MockFileSystemService) RenameFile(oldPath, newPath string) error {
	args := m.Called(oldPath, newPath)
	return args.Error(0)
}

func (m *MockFileSystemService) FileExists(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func TestRenameUseCase_GeneratePreview(t *testing.T) {
	mockFS := new(MockFileSystemService)
	useCase := NewRenameUseCase(mockFS)

	files := []*domain.File{
		domain.NewFile("/path/to/test1.txt"),
		domain.NewFile("/path/to/test2.txt"),
		domain.NewFile("/path/to/other.txt"),
	}

	strategy := domain.NewExactMatchStrategy("test", "renamed")

	previews := useCase.GeneratePreview(files, strategy)

	assert.Equal(t, 3, len(previews))
	assert.Equal(t, "renamed1.txt", previews[0].NewName())
	assert.Equal(t, "renamed2.txt", previews[1].NewName())
	assert.Equal(t, "other.txt", previews[2].NewName())
}

func TestRenameUseCase_Execute_Success(t *testing.T) {
	mockFS := new(MockFileSystemService)
	useCase := NewRenameUseCase(mockFS)

	files := []*domain.File{
		domain.NewFile("/path/to/test1.txt"),
		domain.NewFile("/path/to/test2.txt"),
	}

	strategy := domain.NewExactMatchStrategy("test", "renamed")
	useCase.GeneratePreview(files, strategy)

	// Mock file system calls
	mockFS.On("RenameFile", "/path/to/test1.txt", "/path/to/renamed1.txt").Return(nil)
	mockFS.On("RenameFile", "/path/to/test2.txt", "/path/to/renamed2.txt").Return(nil)

	result := useCase.Execute(files)

	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailureCount)
	assert.Equal(t, 0, len(result.Errors))
	assert.Equal(t, 2, len(result.NewFilePaths))
	assert.Equal(t, "/path/to/renamed1.txt", result.NewFilePaths[0])
	assert.Equal(t, "/path/to/renamed2.txt", result.NewFilePaths[1])

	mockFS.AssertExpectations(t)
}

func TestRenameUseCase_Execute_SkipOnError(t *testing.T) {
	mockFS := new(MockFileSystemService)
	useCase := NewRenameUseCase(mockFS)

	files := []*domain.File{
		domain.NewFile("/path/to/test1.txt"),
		domain.NewFile("/path/to/test2.txt"),
		domain.NewFile("/path/to/test3.txt"),
	}

	strategy := domain.NewExactMatchStrategy("test", "renamed")
	useCase.GeneratePreview(files, strategy)

	// Mock file system calls - second file fails
	mockFS.On("RenameFile", "/path/to/test1.txt", "/path/to/renamed1.txt").Return(nil)
	mockFS.On("RenameFile", "/path/to/test2.txt", "/path/to/renamed2.txt").Return(errors.New("permission denied"))
	mockFS.On("RenameFile", "/path/to/test3.txt", "/path/to/renamed3.txt").Return(nil)

	result := useCase.Execute(files)

	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailureCount)
	assert.Equal(t, 1, len(result.Errors))
	assert.Contains(t, result.Errors[0], "test2.txt")
	assert.Contains(t, result.Errors[0], "permission denied")
	assert.Equal(t, 3, len(result.NewFilePaths))
	assert.Equal(t, "/path/to/renamed1.txt", result.NewFilePaths[0])
	assert.Equal(t, "/path/to/test2.txt", result.NewFilePaths[1]) // Failed, keeps original
	assert.Equal(t, "/path/to/renamed3.txt", result.NewFilePaths[2])

	mockFS.AssertExpectations(t)
}

func TestRenameUseCase_Execute_SkipUnchangedFiles(t *testing.T) {
	mockFS := new(MockFileSystemService)
	useCase := NewRenameUseCase(mockFS)

	files := []*domain.File{
		domain.NewFile("/path/to/test.txt"),
		domain.NewFile("/path/to/other.txt"),
	}

	// Strategy that only renames "test" to "renamed"
	strategy := domain.NewExactMatchStrategy("test", "renamed")
	useCase.GeneratePreview(files, strategy)

	// Only renamed file should be processed
	mockFS.On("RenameFile", "/path/to/test.txt", "/path/to/renamed.txt").Return(nil)

	result := useCase.Execute(files)

	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailureCount)
	assert.Equal(t, 2, len(result.NewFilePaths))
	assert.Equal(t, "/path/to/renamed.txt", result.NewFilePaths[0])
	assert.Equal(t, "/path/to/other.txt", result.NewFilePaths[1]) // Unchanged

	mockFS.AssertExpectations(t)
	mockFS.AssertNotCalled(t, "RenameFile", "/path/to/other.txt", mock.Anything)
}
