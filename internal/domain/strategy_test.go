package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExactMatchStrategy tests exact string matching
func TestExactMatchStrategy(t *testing.T) {
	strategy := NewExactMatchStrategy("test", "TEST")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single match", "test.txt", "TEST.txt"},
		{"multiple matches", "test_test.txt", "TEST_TEST.txt"},
		{"no match", "other.txt", "other.txt"},
		{"partial match also replaced", "testing.txt", "TESTing.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Apply(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRegexMatchStrategy tests regular expression matching
func TestRegexMatchStrategy(t *testing.T) {
	t.Run("valid regex", func(t *testing.T) {
		strategy, err := NewRegexMatchStrategy(`test(\d+)`, "result$1")
		assert.NoError(t, err)

		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"match with capture group", "test123.txt", "result123.txt"},
			{"no match", "other.txt", "other.txt"},
			{"multiple matches", "test1_test2.txt", "result1_result2.txt"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := strategy.Apply(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("invalid regex", func(t *testing.T) {
		_, err := NewRegexMatchStrategy(`[invalid(`, "replacement")
		assert.Error(t, err)
	})
}

// TestCaseInsensitiveStrategy tests case-insensitive matching using Decorator pattern
func TestCaseInsensitiveStrategy(t *testing.T) {
	baseStrategy := NewExactMatchStrategy("test", "RESULT")
	strategy := NewCaseInsensitiveStrategy(baseStrategy)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase match", "test.txt", "RESULT.txt"},
		{"uppercase match", "TEST.txt", "RESULT.txt"},
		{"mixed case match", "TeSt.txt", "RESULT.txt"},
		{"no match", "other.txt", "other.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Apply(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCaseInsensitiveWithRegex tests case-insensitive with regex
func TestCaseInsensitiveWithRegex(t *testing.T) {
	baseStrategy, _ := NewRegexMatchStrategy(`test`, "RESULT")
	strategy := NewCaseInsensitiveStrategy(baseStrategy)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "test.txt", "RESULT.txt"},
		{"uppercase", "TEST.txt", "RESULT.txt"},
		{"mixed case", "TeSt.txt", "RESULT.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Apply(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
