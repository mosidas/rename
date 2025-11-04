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
	// GeneratePreview now returns clones, so use the returned value
	previews := useCase.GeneratePreview(files, strategy)

    // File existence checks (no conflicts)
    mockFS.On("FileExists", "/path/to/renamed1.txt").Return(false)
    mockFS.On("FileExists", "/path/to/renamed2.txt").Return(false)

    // Mock file system calls
    mockFS.On("RenameFile", "/path/to/test1.txt", "/path/to/renamed1.txt").Return(nil)
    mockFS.On("RenameFile", "/path/to/test2.txt", "/path/to/renamed2.txt").Return(nil)

	result := useCase.Execute(previews)

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
	previews := useCase.GeneratePreview(files, strategy)

    // File existence checks (no conflicts)
    mockFS.On("FileExists", "/path/to/renamed1.txt").Return(false)
    mockFS.On("FileExists", "/path/to/renamed2.txt").Return(false)
    mockFS.On("FileExists", "/path/to/renamed3.txt").Return(false)

    // Mock file system calls - second file fails
    mockFS.On("RenameFile", "/path/to/test1.txt", "/path/to/renamed1.txt").Return(nil)
    mockFS.On("RenameFile", "/path/to/test2.txt", "/path/to/renamed2.txt").Return(errors.New("permission denied"))
    mockFS.On("RenameFile", "/path/to/test3.txt", "/path/to/renamed3.txt").Return(nil)

	result := useCase.Execute(previews)

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
    previews := useCase.GeneratePreview(files, strategy)

    // File existence check (no conflict)
    mockFS.On("FileExists", "/path/to/renamed.txt").Return(false)

    // Only renamed file should be processed
    mockFS.On("RenameFile", "/path/to/test.txt", "/path/to/renamed.txt").Return(nil)

	result := useCase.Execute(previews)

	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailureCount)
	assert.Equal(t, 2, len(result.NewFilePaths))
	assert.Equal(t, "/path/to/renamed.txt", result.NewFilePaths[0])
	assert.Equal(t, "/path/to/other.txt", result.NewFilePaths[1]) // Unchanged

	mockFS.AssertExpectations(t)
	mockFS.AssertNotCalled(t, "RenameFile", "/path/to/other.txt", mock.Anything)
}

func TestRenameUseCase_Execute_ConflictSuffix(t *testing.T) {
    mockFS := new(MockFileSystemService)
    useCase := NewRenameUseCase(mockFS)

    files := []*domain.File{
        domain.NewFile("/path/to/test.txt"),
    }

    strategy := domain.NewExactMatchStrategy("test", "renamed")
    previews := useCase.GeneratePreview(files, strategy)

    // Simulate conflicts: renamed.txt exists, renamed1.txt exists, renamed2.txt does not
    mockFS.On("FileExists", "/path/to/renamed.txt").Return(true)
    mockFS.On("FileExists", "/path/to/renamed1.txt").Return(true)
    mockFS.On("FileExists", "/path/to/renamed2.txt").Return(false)

    // Expect rename to the resolved name with suffix 2
    mockFS.On("RenameFile", "/path/to/test.txt", "/path/to/renamed2.txt").Return(nil)

    result := useCase.Execute(previews)

    assert.Equal(t, 1, result.SuccessCount)
    assert.Equal(t, 0, result.FailureCount)
    assert.Equal(t, 1, len(result.NewFilePaths))
    assert.Equal(t, "/path/to/renamed2.txt", result.NewFilePaths[0])

    mockFS.AssertExpectations(t)
}

func TestRenameUseCase_Execute_ConflictMaxRetries(t *testing.T) {
    mockFS := new(MockFileSystemService)
    useCase := NewRenameUseCase(mockFS)

    files := []*domain.File{
        domain.NewFile("/path/to/test.txt"),
    }

    strategy := domain.NewExactMatchStrategy("test", "renamed")
    previews := useCase.GeneratePreview(files, strategy)

    // Simulate all possible names being taken (up to maxRetries = 1000)
    // Initial check: renamed.txt exists
    mockFS.On("FileExists", "/path/to/renamed.txt").Return(true)

    // Mock FileExists to always return true for any path (all names taken)
    mockFS.On("FileExists", mock.AnythingOfType("string")).Return(true)

    result := useCase.Execute(previews)

    // Should fail because all names are taken
    assert.Equal(t, 0, result.SuccessCount)
    assert.Equal(t, 1, result.FailureCount)
    assert.Equal(t, 1, len(result.Errors))
    assert.Contains(t, result.Errors[0], "Failed to find available name")
    assert.Contains(t, result.Errors[0], "1000 retries")
    assert.Equal(t, 1, len(result.NewFilePaths))
    assert.Equal(t, "/path/to/test.txt", result.NewFilePaths[0]) // Original path kept
}

func TestRenameUseCase_Execute_SameNameRename(t *testing.T) {
    mockFS := new(MockFileSystemService)
    useCase := NewRenameUseCase(mockFS)

    files := []*domain.File{
        domain.NewFile("/path/to/test.txt"),
    }

    // Strategy that doesn't change the name
    strategy := domain.NewExactMatchStrategy("test", "test")
    previews := useCase.GeneratePreview(files, strategy)

    result := useCase.Execute(previews)

    // Should skip because file hasn't changed
    assert.Equal(t, 0, result.SuccessCount)
    assert.Equal(t, 0, result.FailureCount)
    assert.Equal(t, 1, len(result.NewFilePaths))
    assert.Equal(t, "/path/to/test.txt", result.NewFilePaths[0])

    mockFS.AssertNotCalled(t, "RenameFile", mock.Anything, mock.Anything)
}

func TestRenameUseCase_Execute_EmptyPattern(t *testing.T) {
    mockFS := new(MockFileSystemService)
    useCase := NewRenameUseCase(mockFS)

    files := []*domain.File{
        domain.NewFile("/path/to/test.txt"),
    }

    // Empty pattern matches everything and replaces with nothing
    strategy := domain.NewExactMatchStrategy("", "")
    previews := useCase.GeneratePreview(files, strategy)

    result := useCase.Execute(previews)

    // File name becomes empty after replacement, but it's technically "changed"
    // However, the name is same ("" replaces to ""), so HasChanged should be false
    assert.Equal(t, 0, result.SuccessCount)
    assert.Equal(t, 0, result.FailureCount)
}

func TestRenameUseCase_Execute_FileWithoutExtension(t *testing.T) {
    mockFS := new(MockFileSystemService)
    useCase := NewRenameUseCase(mockFS)

    files := []*domain.File{
        domain.NewFile("/path/to/Makefile"),
    }

    strategy := domain.NewExactMatchStrategy("Make", "BUILD")
    previews := useCase.GeneratePreview(files, strategy)

    mockFS.On("FileExists", "/path/to/BUILDfile").Return(false)
    mockFS.On("RenameFile", "/path/to/Makefile", "/path/to/BUILDfile").Return(nil)

    result := useCase.Execute(previews)

    assert.Equal(t, 1, result.SuccessCount)
    assert.Equal(t, 0, result.FailureCount)
    assert.Equal(t, "/path/to/BUILDfile", result.NewFilePaths[0])

    mockFS.AssertExpectations(t)
}
