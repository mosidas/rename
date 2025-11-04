package domain

const MaxHistorySize = 100

// HistoryEntry represents a single history entry
type HistoryEntry struct {
	Pattern       string `json:"pattern"`
	Replacement   string `json:"replacement"`
	IsRegex       bool   `json:"isRegex"`
	CaseInsensitive bool `json:"caseInsensitive"`
}

// History manages rename history
// Following SRP (Single Responsibility Principle) - only manages history entries
type History struct {
	entries []HistoryEntry
}

// NewHistory creates a new History
func NewHistory() *History {
	return &History{
		entries: make([]HistoryEntry, 0, MaxHistorySize),
	}
}

// Add adds a new history entry
// If duplicate exists, it moves to front instead of adding
func (h *History) Add(entry HistoryEntry) {
	// Check for duplicate
	for i, existing := range h.entries {
		if existing.Pattern == entry.Pattern &&
			existing.Replacement == entry.Replacement &&
			existing.IsRegex == entry.IsRegex &&
			existing.CaseInsensitive == entry.CaseInsensitive {
			// Move to front using efficient copy
			if i > 0 {
				// Shift entries [0:i] to [1:i+1] and place entry at front
				copy(h.entries[1:i+1], h.entries[0:i])
				h.entries[0] = entry
			}
			return
		}
	}

	// Add new entry to front
	h.entries = append([]HistoryEntry{entry}, h.entries...)

	// Keep only MaxHistorySize entries
	if len(h.entries) > MaxHistorySize {
		h.entries = h.entries[:MaxHistorySize]
	}
}

// GetAll returns all history entries (most recent first)
func (h *History) GetAll() []HistoryEntry {
	return h.entries
}

// Count returns the number of history entries
func (h *History) Count() int {
	return len(h.entries)
}

// Clear removes all history entries
func (h *History) Clear() {
	h.entries = make([]HistoryEntry, 0, MaxHistorySize)
}

// SetEntries sets the history entries directly (for repository loading)
// This is more efficient than calling Add repeatedly
func (h *History) SetEntries(entries []HistoryEntry) {
	h.entries = entries
}
