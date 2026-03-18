package valueobject

import (
	"errors"
	"strings"
)

const (
	// TagName constraints
	MinTagNameLength = 1
	MaxTagNameLength = 50
)

// TagName represents a tag name with validation
type TagName struct {
	value string
}

// NewTagName creates a new TagName with validation
func NewTagName(value string) (TagName, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < MinTagNameLength {
		return TagName{}, errors.New("tag name cannot be empty")
	}
	if len(trimmed) > MaxTagNameLength {
		return TagName{}, errors.New("tag name cannot exceed 50 characters")
	}
	return TagName{value: trimmed}, nil
}

// String returns the string representation of TagName
func (t TagName) String() string {
	return t.value
}
