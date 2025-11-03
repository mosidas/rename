package domain

import (
	"path/filepath"
)

// File represents a file to be renamed
// Following SRP (Single Responsibility Principle) - only manages file information
type File struct {
	originalPath string
	originalName string
	directory    string
	newName      string
}

// NewFile creates a new File entity
func NewFile(path string) *File {
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	return &File{
		originalPath: path,
		originalName: name,
		directory:    dir,
		newName:      name, // Initially same as original
	}
}

// OriginalPath returns the original file path
func (f *File) OriginalPath() string {
	return f.originalPath
}

// OriginalName returns the original file name
func (f *File) OriginalName() string {
	return f.originalName
}

// Directory returns the directory containing the file
func (f *File) Directory() string {
	return f.directory
}

// NewName returns the new file name
func (f *File) NewName() string {
	return f.newName
}

// NewPath returns the new full path
func (f *File) NewPath() string {
	return filepath.Join(f.directory, f.newName)
}

// SetNewName sets the new file name
func (f *File) SetNewName(name string) {
	f.newName = name
}

// HasChanged returns true if the file name has changed
func (f *File) HasChanged() bool {
	return f.originalName != f.newName
}
