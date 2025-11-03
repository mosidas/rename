package service

import (
	"os"
)

// FileSystemService provides file system operations
// Following SRP (Single Responsibility Principle)
type FileSystemService struct{}

// NewFileSystemService creates a new FileSystemService
func NewFileSystemService() *FileSystemService {
	return &FileSystemService{}
}

// RenameFile renames a file from oldPath to newPath
func (fs *FileSystemService) RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// FileExists checks if a file exists at the given path
func (fs *FileSystemService) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
