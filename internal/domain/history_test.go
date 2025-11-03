package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHistory(t *testing.T) {
	history := NewHistory()
	assert.NotNil(t, history)
	assert.Equal(t, 0, history.Count())
}

func TestHistory_Add(t *testing.T) {
	history := NewHistory()

	entry := HistoryEntry{
		Pattern:     "test",
		Replacement: "TEST",
		IsRegex:     false,
	}

	history.Add(entry)
	assert.Equal(t, 1, history.Count())

	entries := history.GetAll()
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, entry, entries[0])
}

func TestHistory_MaxLimit(t *testing.T) {
	history := NewHistory()

	// Add 150 unique entries
	for i := 0; i < 150; i++ {
		history.Add(HistoryEntry{
			Pattern:     "pattern" + string(rune('A'+i)),
			Replacement: "replacement",
			IsRegex:     false,
		})
	}

	// Should only keep the last 100 entries
	assert.Equal(t, 100, history.Count())
}

func TestHistory_GetAll(t *testing.T) {
	history := NewHistory()

	entries := []HistoryEntry{
		{Pattern: "test1", Replacement: "TEST1", IsRegex: false},
		{Pattern: "test2", Replacement: "TEST2", IsRegex: true},
		{Pattern: "test3", Replacement: "TEST3", IsRegex: false},
	}

	for _, entry := range entries {
		history.Add(entry)
	}

	retrieved := history.GetAll()
	assert.Equal(t, len(entries), len(retrieved))

	// Most recent should be first
	for i := range entries {
		// Reverse order (most recent first)
		assert.Equal(t, entries[len(entries)-1-i], retrieved[i])
	}
}

func TestHistory_DuplicateHandling(t *testing.T) {
	history := NewHistory()

	entry := HistoryEntry{
		Pattern:     "test",
		Replacement: "TEST",
		IsRegex:     false,
	}

	history.Add(entry)
	history.Add(entry) // Add duplicate

	// Duplicate should be moved to front, not added twice
	assert.Equal(t, 1, history.Count())
}
