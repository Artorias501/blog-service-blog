package valueobject

import (
	"strings"
	"testing"
)

func TestAuthorName_Validation(t *testing.T) {
	t.Run("valid author name creates AuthorName successfully", func(t *testing.T) {
		authorName, err := NewAuthorName("John Doe")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if authorName.String() != "John Doe" {
			t.Errorf("expected 'John Doe', got '%s'", authorName.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewAuthorName("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("author name with 101 characters returns validation error", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		_, err := NewAuthorName(longName)
		if err == nil {
			t.Error("expected error for 101 character author name, got nil")
		}
	})

	t.Run("author name with exactly 100 characters succeeds", func(t *testing.T) {
		exactName := strings.Repeat("a", 100)
		authorName, err := NewAuthorName(exactName)
		if err != nil {
			t.Errorf("expected no error for 100 character author name, got: %v", err)
		}
		if authorName.String() != exactName {
			t.Error("author name mismatch")
		}
	})
}

func TestAuthorEmail_Validation(t *testing.T) {
	t.Run("valid email format creates AuthorEmail successfully", func(t *testing.T) {
		email, err := NewAuthorEmail("test@example.com")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if email.String() != "test@example.com" {
			t.Errorf("expected 'test@example.com', got '%s'", email.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewAuthorEmail("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("invalid email format (no @) returns validation error", func(t *testing.T) {
		_, err := NewAuthorEmail("invalidemail")
		if err == nil {
			t.Error("expected error for email without @, got nil")
		}
	})

	t.Run("invalid email format (no domain) returns validation error", func(t *testing.T) {
		_, err := NewAuthorEmail("test@")
		if err == nil {
			t.Error("expected error for email without domain, got nil")
		}
	})

	t.Run("email with special characters succeeds if valid", func(t *testing.T) {
		email, err := NewAuthorEmail("test.user+tag@example.com")
		if err != nil {
			t.Errorf("expected no error for valid email with special chars, got: %v", err)
		}
		if email.String() != "test.user+tag@example.com" {
			t.Error("email mismatch")
		}
	})
}

func TestCommentStatus_Validation(t *testing.T) {
	t.Run("valid status 'pending' creates CommentStatus successfully", func(t *testing.T) {
		status, err := NewCommentStatus("pending")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if status.String() != "pending" {
			t.Errorf("expected 'pending', got '%s'", status.String())
		}
	})

	t.Run("valid status 'approved' creates CommentStatus successfully", func(t *testing.T) {
		status, err := NewCommentStatus("approved")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if status.String() != "approved" {
			t.Errorf("expected 'approved', got '%s'", status.String())
		}
	})

	t.Run("valid status 'rejected' creates CommentStatus successfully", func(t *testing.T) {
		status, err := NewCommentStatus("rejected")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if status.String() != "rejected" {
			t.Errorf("expected 'rejected', got '%s'", status.String())
		}
	})

	t.Run("valid status 'spam' creates CommentStatus successfully", func(t *testing.T) {
		status, err := NewCommentStatus("spam")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if status.String() != "spam" {
			t.Errorf("expected 'spam', got '%s'", status.String())
		}
	})

	t.Run("invalid status returns validation error", func(t *testing.T) {
		_, err := NewCommentStatus("invalid")
		if err == nil {
			t.Error("expected error for invalid status, got nil")
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewCommentStatus("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("default status is pending", func(t *testing.T) {
		status := DefaultCommentStatus()
		if status.String() != "pending" {
			t.Errorf("expected default status 'pending', got '%s'", status.String())
		}
	})
}
