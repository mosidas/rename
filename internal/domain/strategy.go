package domain

import (
	"regexp"
	"strings"
)

// RenameStrategy defines the interface for different rename strategies
// Following Strategy Pattern (OCP - Open/Closed Principle)
type RenameStrategy interface {
	Apply(filename string) string
}

// ExactMatchStrategy implements exact string matching
type ExactMatchStrategy struct {
	pattern     string
	replacement string
}

// NewExactMatchStrategy creates a new exact match strategy
func NewExactMatchStrategy(pattern, replacement string) *ExactMatchStrategy {
	return &ExactMatchStrategy{
		pattern:     pattern,
		replacement: replacement,
	}
}

// Apply applies exact string replacement
func (s *ExactMatchStrategy) Apply(filename string) string {
	return strings.ReplaceAll(filename, s.pattern, s.replacement)
}

// RegexMatchStrategy implements regular expression matching
type RegexMatchStrategy struct {
	regex       *regexp.Regexp
	replacement string
}

// NewRegexMatchStrategy creates a new regex match strategy
func NewRegexMatchStrategy(pattern, replacement string) (*RegexMatchStrategy, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RegexMatchStrategy{
		regex:       regex,
		replacement: replacement,
	}, nil
}

// Apply applies regex replacement
func (s *RegexMatchStrategy) Apply(filename string) string {
	return s.regex.ReplaceAllString(filename, s.replacement)
}

// PatternProvider is an interface for strategies that can expose their pattern and replacement
// This allows CaseInsensitiveStrategy to work without type assertions
type PatternProvider interface {
	GetPattern() string
	GetReplacement() string
	NeedsEscaping() bool // Returns true if pattern needs regex escaping (for exact match)
}

// GetPattern returns the pattern for ExactMatchStrategy
func (s *ExactMatchStrategy) GetPattern() string {
	return s.pattern
}

// GetReplacement returns the replacement for ExactMatchStrategy
func (s *ExactMatchStrategy) GetReplacement() string {
	return s.replacement
}

// NeedsEscaping returns true for ExactMatchStrategy (pattern needs regex escaping)
func (s *ExactMatchStrategy) NeedsEscaping() bool {
	return true
}

// GetPattern returns the regex pattern string for RegexMatchStrategy
func (s *RegexMatchStrategy) GetPattern() string {
	return s.regex.String()
}

// GetReplacement returns the replacement for RegexMatchStrategy
func (s *RegexMatchStrategy) GetReplacement() string {
	return s.replacement
}

// NeedsEscaping returns false for RegexMatchStrategy (already a regex pattern)
func (s *RegexMatchStrategy) NeedsEscaping() bool {
	return false
}

// CaseInsensitiveStrategy is a decorator that makes any strategy case-insensitive
// Following Decorator Pattern (OCP - Open/Closed Principle)
type CaseInsensitiveStrategy struct {
	strategy RenameStrategy
}

// NewCaseInsensitiveStrategy creates a case-insensitive decorator
func NewCaseInsensitiveStrategy(strategy RenameStrategy) *CaseInsensitiveStrategy {
	return &CaseInsensitiveStrategy{
		strategy: strategy,
	}
}

// Apply applies the wrapped strategy in a case-insensitive manner
func (s *CaseInsensitiveStrategy) Apply(filename string) string {
	// If the strategy implements PatternProvider, use it
	if provider, ok := s.strategy.(PatternProvider); ok {
		pattern := provider.GetPattern()
		replacement := provider.GetReplacement()

		// Escape pattern if needed (for exact match strategies)
		if provider.NeedsEscaping() {
			pattern = regexp.QuoteMeta(pattern)
		}

		// Add case-insensitive flag if not already present
		if !strings.HasPrefix(pattern, "(?i)") {
			pattern = "(?i)" + pattern
		}

		regex := regexp.MustCompile(pattern)
		return regex.ReplaceAllString(filename, replacement)
	}

	// Fallback to default behavior
	return s.strategy.Apply(filename)
}
