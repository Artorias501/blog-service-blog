package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

const (
	// AuthorName constraints
	MinAuthorNameLength = 1
	MaxAuthorNameLength = 100
)

// Email regex pattern for basic validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// AuthorName represents a comment author name with validation
type AuthorName struct {
	value string
}

// NewAuthorName creates a new AuthorName with validation
func NewAuthorName(value string) (AuthorName, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < MinAuthorNameLength {
		return AuthorName{}, errors.New("author name cannot be empty")
	}
	if len(trimmed) > MaxAuthorNameLength {
		return AuthorName{}, errors.New("author name cannot exceed 100 characters")
	}
	return AuthorName{value: trimmed}, nil
}

// String returns the string representation of AuthorName
func (a AuthorName) String() string {
	return a.value
}

// AuthorEmail represents a comment author email with validation
type AuthorEmail struct {
	value string
}

// NewAuthorEmail creates a new AuthorEmail with validation
func NewAuthorEmail(value string) (AuthorEmail, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return AuthorEmail{}, errors.New("author email cannot be empty")
	}
	if !emailRegex.MatchString(trimmed) {
		return AuthorEmail{}, errors.New("invalid email format")
	}
	return AuthorEmail{value: trimmed}, nil
}

// String returns the string representation of AuthorEmail
func (a AuthorEmail) String() string {
	return a.value
}

// CommentStatus represents the status of a comment
type CommentStatus struct {
	value string
}

// Valid comment statuses
const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusRejected = "rejected"
	CommentStatusSpam     = "spam"
)

var validCommentStatuses = map[string]bool{
	CommentStatusPending:  true,
	CommentStatusApproved: true,
	CommentStatusRejected: true,
	CommentStatusSpam:     true,
}

// NewCommentStatus creates a new CommentStatus with validation
func NewCommentStatus(value string) (CommentStatus, error) {
	if value == "" {
		return CommentStatus{}, errors.New("comment status cannot be empty")
	}
	if !validCommentStatuses[value] {
		return CommentStatus{}, errors.New("invalid comment status: must be one of pending, approved, rejected, spam")
	}
	return CommentStatus{value: value}, nil
}

// DefaultCommentStatus returns the default comment status (pending)
func DefaultCommentStatus() CommentStatus {
	return CommentStatus{value: CommentStatusPending}
}

// String returns the string representation of CommentStatus
func (c CommentStatus) String() string {
	return c.value
}

// IsApproved checks if the comment is approved
func (c CommentStatus) IsApproved() bool {
	return c.value == CommentStatusApproved
}

// IsPending checks if the comment is pending
func (c CommentStatus) IsPending() bool {
	return c.value == CommentStatusPending
}

// IsRejected checks if the comment is rejected
func (c CommentStatus) IsRejected() bool {
	return c.value == CommentStatusRejected
}

// IsSpam checks if the comment is marked as spam
func (c CommentStatus) IsSpam() bool {
	return c.value == CommentStatusSpam
}
