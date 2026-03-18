package valueobject

import (
	"errors"

	"github.com/google/uuid"
)

// PostID represents a unique identifier for a Post entity
type PostID struct {
	value uuid.UUID
}

// NewPostID creates a new PostID from a string, validating UUID format
func NewPostID(value string) (PostID, error) {
	if value == "" {
		return PostID{}, errors.New("post ID cannot be empty")
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return PostID{}, errors.New("invalid post ID format: must be a valid UUID")
	}

	return PostID{value: id}, nil
}

// GeneratePostID generates a new random PostID
func GeneratePostID() PostID {
	return PostID{value: uuid.New()}
}

// String returns the string representation of PostID
func (p PostID) String() string {
	return p.value.String()
}

// Equals checks if two PostIDs are equal
func (p PostID) Equals(other PostID) bool {
	return p.value == other.value
}

// TagID represents a unique identifier for a Tag entity
type TagID struct {
	value uuid.UUID
}

// NewTagID creates a new TagID from a string, validating UUID format
func NewTagID(value string) (TagID, error) {
	if value == "" {
		return TagID{}, errors.New("tag ID cannot be empty")
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return TagID{}, errors.New("invalid tag ID format: must be a valid UUID")
	}

	return TagID{value: id}, nil
}

// GenerateTagID generates a new random TagID
func GenerateTagID() TagID {
	return TagID{value: uuid.New()}
}

// String returns the string representation of TagID
func (t TagID) String() string {
	return t.value.String()
}

// Equals checks if two TagIDs are equal
func (t TagID) Equals(other TagID) bool {
	return t.value == other.value
}

// CommentID represents a unique identifier for a Comment entity
type CommentID struct {
	value uuid.UUID
}

// NewCommentID creates a new CommentID from a string, validating UUID format
func NewCommentID(value string) (CommentID, error) {
	if value == "" {
		return CommentID{}, errors.New("comment ID cannot be empty")
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return CommentID{}, errors.New("invalid comment ID format: must be a valid UUID")
	}

	return CommentID{value: id}, nil
}

// GenerateCommentID generates a new random CommentID
func GenerateCommentID() CommentID {
	return CommentID{value: uuid.New()}
}

// String returns the string representation of CommentID
func (c CommentID) String() string {
	return c.value.String()
}

// Equals checks if two CommentIDs are equal
func (c CommentID) Equals(other CommentID) bool {
	return c.value == other.value
}
