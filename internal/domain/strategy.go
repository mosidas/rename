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
	// For ExactMatchStrategy, we need special handling
	if exactStrategy, ok := s.strategy.(*ExactMatchStrategy); ok {
		// Create a case-insensitive regex from the pattern
		pattern := regexp.QuoteMeta(exactStrategy.pattern)
		regex := regexp.MustCompile("(?i)" + pattern)
		return regex.ReplaceAllString(filename, exactStrategy.replacement)
	}

	// For RegexMatchStrategy, modify the pattern to be case-insensitive
	if regexStrategy, ok := s.strategy.(*RegexMatchStrategy); ok {
		// Recreate the regex with case-insensitive flag
		pattern := regexStrategy.regex.String()
		if !strings.HasPrefix(pattern, "(?i)") {
			regex := regexp.MustCompile("(?i)" + pattern)
			return regex.ReplaceAllString(filename, regexStrategy.replacement)
		}
	}

	return s.strategy.Apply(filename)
}
