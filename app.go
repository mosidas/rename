package main

import (
	"context"
	"os"
	"path/filepath"

	"rename/internal/domain"
	"rename/internal/repository"
	"rename/internal/service"
	"rename/internal/usecase"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct - Presentation layer (thin adapter)
// Following DIP (Dependency Inversion Principle) - depends on abstractions (use cases)
type App struct {
	ctx                    context.Context
	renameUseCase          *usecase.RenameUseCase
	historyUseCase         *usecase.HistoryUseCase
	fileSystem             *service.FileSystemService
	currentFiles           []*domain.File
	currentStrategy        domain.RenameStrategy
	currentPattern         string
	currentReplacement     string
	currentIsRegex         bool
	currentCaseInsensitive bool
	initialFiles           []string // Files passed on startup via command-line
}

// NewApp creates a new App application struct with dependency injection
func NewApp() *App {
	// Initialize services and repositories
	fileSystem := service.NewFileSystemService()

	// Get config path
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "rename", "config.json")

	historyRepo := repository.NewJSONHistoryRepository(configPath)

	// Initialize use cases
	renameUseCase := usecase.NewRenameUseCase(fileSystem)
	historyUseCase := usecase.NewHistoryUseCase(historyRepo)

	return &App{
		renameUseCase:  renameUseCase,
		historyUseCase: historyUseCase,
		fileSystem:     fileSystem,
		currentFiles:   make([]*domain.File, 0),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// If initial files were provided via command-line, load them
	if len(a.initialFiles) > 0 {
		a.currentFiles = make([]*domain.File, len(a.initialFiles))
		for i, path := range a.initialFiles {
			a.currentFiles[i] = domain.NewFile(path)
		}

		// Emit event to notify frontend
		runtime.EventsEmit(ctx, "files:loaded", a.initialFiles)
	}
}

// SelectFiles opens file selection dialog
func (a *App) SelectFiles() ([]string, error) {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "ファイルを選択",
	})
	if err != nil {
		return nil, err
	}

	// Convert to File entities
	a.currentFiles = make([]*domain.File, len(files))
	for i, path := range files {
		a.currentFiles[i] = domain.NewFile(path)
	}

	return files, nil
}

// FilePreview represents a file preview for frontend
type FilePreview struct {
	OriginalPath string `json:"originalPath"`
	OriginalName string `json:"originalName"`
	NewName      string `json:"newName"`
	HasChanged   bool   `json:"hasChanged"`
}

// GeneratePreview generates rename preview
func (a *App) GeneratePreview(pattern, replacement string, isRegex, caseInsensitive bool) ([]FilePreview, error) {
	if len(a.currentFiles) == 0 {
		return []FilePreview{}, nil
	}

	// Save current pattern info for later use
	a.currentPattern = pattern
	a.currentReplacement = replacement
	a.currentIsRegex = isRegex
	a.currentCaseInsensitive = caseInsensitive

	// Create strategy
	var strategy domain.RenameStrategy
	var err error

	if isRegex {
		strategy, err = domain.NewRegexMatchStrategy(pattern, replacement)
		if err != nil {
			return nil, err
		}
	} else {
		strategy = domain.NewExactMatchStrategy(pattern, replacement)
	}

	// Apply case-insensitive if needed
	if caseInsensitive {
		strategy = domain.NewCaseInsensitiveStrategy(strategy)
	}

	a.currentStrategy = strategy

	// Generate preview
	files := a.renameUseCase.GeneratePreview(a.currentFiles, strategy)

	// Convert to preview
	previews := make([]FilePreview, len(files))
	for i, file := range files {
		previews[i] = FilePreview{
			OriginalPath: file.OriginalPath(),
			OriginalName: file.OriginalName(),
			NewName:      file.NewName(),
			HasChanged:   file.HasChanged(),
		}
	}

	return previews, nil
}

// ExecuteRename executes the rename operation
func (a *App) ExecuteRename() (usecase.RenameResult, error) {
	if a.currentStrategy == nil {
		return usecase.RenameResult{}, nil
	}

	result := a.renameUseCase.Execute(a.currentFiles)

	// Update currentFiles with new paths after rename
	if len(result.NewFilePaths) > 0 {
		a.currentFiles = make([]*domain.File, len(result.NewFilePaths))
		for i, path := range result.NewFilePaths {
			a.currentFiles[i] = domain.NewFile(path)
		}
	}

	// If successful, add to history
	if result.SuccessCount > 0 {
		entry := domain.HistoryEntry{
			Pattern:         a.currentPattern,
			Replacement:     a.currentReplacement,
			IsRegex:         a.currentIsRegex,
			CaseInsensitive: a.currentCaseInsensitive,
		}
		// Ignore error from history save - it's not critical
		_ = a.historyUseCase.AddEntry(entry)
	}

	return result, nil
}

// GetHistory returns rename history
func (a *App) GetHistory() ([]domain.HistoryEntry, error) {
	return a.historyUseCase.GetHistory()
}

// AddToHistory adds an entry to history
func (a *App) AddToHistory(entry domain.HistoryEntry) error {
	return a.historyUseCase.AddEntry(entry)
}

// SetInitialFiles sets files passed via command-line on startup
func (a *App) SetInitialFiles(files []string) {
	a.initialFiles = files
}

// GetInitialFiles returns files passed on startup (for frontend)
func (a *App) GetInitialFiles() []string {
	return a.initialFiles
}

// LoadFilesFromSecondInstance loads files when a second instance is launched
func (a *App) LoadFilesFromSecondInstance(files []string) {
	// Convert to File entities
	a.currentFiles = make([]*domain.File, len(files))
	for i, path := range files {
		a.currentFiles[i] = domain.NewFile(path)
	}

	// Show the existing window
	runtime.WindowShow(a.ctx)

	// Notify frontend
	runtime.EventsEmit(a.ctx, "files:loaded", files)
}
