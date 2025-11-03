package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedName string
		expectedDir  string
	}{
		{
			name:         "simple file",
			path:         "/path/to/file.txt",
			expectedName: "file.txt",
			expectedDir:  "/path/to",
		},
		{
			name:         "file without extension",
			path:         "/path/to/file",
			expectedName: "file",
			expectedDir:  "/path/to",
		},
		{
			name:         "relative path",
			path:         "file.txt",
			expectedName: "file.txt",
			expectedDir:  ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := NewFile(tt.path)
			assert.Equal(t, tt.path, file.OriginalPath())
			assert.Equal(t, tt.expectedName, file.OriginalName())
			assert.Equal(t, tt.expectedDir, file.Directory())
		})
	}
}

func TestFile_SetNewName(t *testing.T) {
	file := NewFile("/path/to/file.txt")

	file.SetNewName("renamed.txt")
	assert.Equal(t, "renamed.txt", file.NewName())
	assert.Equal(t, "/path/to/renamed.txt", file.NewPath())
}

func TestFile_HasChanged(t *testing.T) {
	file := NewFile("/path/to/file.txt")

	assert.False(t, file.HasChanged(), "should not have changed initially")

	file.SetNewName("renamed.txt")
	assert.True(t, file.HasChanged(), "should have changed after rename")

	file.SetNewName("file.txt")
	assert.False(t, file.HasChanged(), "should not have changed when name is same as original")
}
