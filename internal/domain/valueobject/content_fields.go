package valueobject

import (
	"errors"
	"strings"
)

const (
	// Title constraints
	MinTitleLength = 1
	MaxTitleLength = 200

	// Content constraints
	MinContentLength = 1

	// Summary constraints
	MaxSummaryLength = 500
)

// Title represents a post title with validation
type Title struct {
	value string
}

// NewTitle creates a new Title with validation
func NewTitle(value string) (Title, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < MinTitleLength {
		return Title{}, errors.New("title cannot be empty")
	}
	if len(trimmed) > MaxTitleLength {
		return Title{}, errors.New("title cannot exceed 200 characters")
	}
	return Title{value: trimmed}, nil
}

// String returns the string representation of Title
func (t Title) String() string {
	return t.value
}

// Content represents post content with validation
type Content struct {
	value string
}

// NewContent creates a new Content with validation
func NewContent(value string) (Content, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < MinContentLength {
		return Content{}, errors.New("content cannot be empty")
	}
	return Content{value: value}, nil
}

// String returns the string representation of Content
func (c Content) String() string {
	return c.value
}

// Summary represents a post summary with validation
type Summary struct {
	value string
}

// NewSummary creates a new Summary with validation
func NewSummary(value string) (Summary, error) {
	if len(value) > MaxSummaryLength {
		return Summary{}, errors.New("summary cannot exceed 500 characters")
	}
	return Summary{value: value}, nil
}

// String returns the string representation of Summary
func (s Summary) String() string {
	return s.value
}
