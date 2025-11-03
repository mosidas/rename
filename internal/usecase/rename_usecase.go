package usecase

import (
    "fmt"
    "path/filepath"
    "strconv"
    "strings"

    "rename/internal/domain"
)

// FileSystemService defines file system operations
// Following ISP (Interface Segregation Principle) - only the methods we need
type FileSystemService interface {
	RenameFile(oldPath, newPath string) error
	FileExists(path string) bool
}

// RenameResult represents the result of a rename operation
type RenameResult struct {
	SuccessCount  int
	FailureCount  int
	Errors        []string
	NewFilePaths  []string
}

// RenameUseCase handles file renaming operations
// Following SRP (Single Responsibility Principle) and DIP (Dependency Inversion Principle)
type RenameUseCase struct {
	fileSystem FileSystemService
}

// NewRenameUseCase creates a new RenameUseCase
func NewRenameUseCase(fileSystem FileSystemService) *RenameUseCase {
	return &RenameUseCase{
		fileSystem: fileSystem,
	}
}

// GeneratePreview applies the strategy to files and returns preview
func (uc *RenameUseCase) GeneratePreview(files []*domain.File, strategy domain.RenameStrategy) []*domain.File {
	for _, file := range files {
		newName := strategy.Apply(file.OriginalName())
		file.SetNewName(newName)
	}
	return files
}

// Execute performs the actual file renaming
// Skips files on error (as per requirements)
func (uc *RenameUseCase) Execute(files []*domain.File) RenameResult {
    result := RenameResult{
        Errors:       make([]string, 0),
        NewFilePaths: make([]string, 0),
    }

    const maxRetries = 1000

    for _, file := range files {
        // Skip files that haven't changed
        if !file.HasChanged() {
            // Keep original path for unchanged files
            result.NewFilePaths = append(result.NewFilePaths, file.OriginalPath())
            continue
        }

        // If target path already exists, find an available name by adding numeric suffix
        targetPath := file.NewPath()
        if targetPath != file.OriginalPath() && uc.fileSystem.FileExists(targetPath) {
            // Split name and extension, then try base+1, base+2, ...
            ext := filepath.Ext(file.NewName())
            base := strings.TrimSuffix(file.NewName(), ext)

            found := false
            for i := 1; i < maxRetries; i++ {
                candidateName := base + strconv.Itoa(i) + ext
                candidatePath := filepath.Join(file.Directory(), candidateName)
                if !uc.fileSystem.FileExists(candidatePath) {
                    // Update file's new name to resolved unique name
                    file.SetNewName(candidateName)
                    targetPath = candidatePath
                    found = true
                    break
                }
            }

            if !found {
                result.FailureCount++
                errorMsg := fmt.Sprintf("Failed to find available name for %s after %d retries", file.OriginalName(), maxRetries)
                result.Errors = append(result.Errors, errorMsg)
                // Keep original path for failed files
                result.NewFilePaths = append(result.NewFilePaths, file.OriginalPath())
                continue
            }
        }

        err := uc.fileSystem.RenameFile(file.OriginalPath(), targetPath)
        if err != nil {
            result.FailureCount++
            errorMsg := fmt.Sprintf("Failed to rename %s: %v", file.OriginalName(), err)
            result.Errors = append(result.Errors, errorMsg)
            // Keep original path for failed files
            result.NewFilePaths = append(result.NewFilePaths, file.OriginalPath())
            // Continue to next file (skip on error)
            continue
        }

        result.SuccessCount++
        // Add new path for successfully renamed files
        result.NewFilePaths = append(result.NewFilePaths, file.NewPath())
    }

    return result
}
