package entity

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestComment_Creation(t *testing.T) {
	t.Run("valid comment data creates Comment successfully", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)

		if comment == nil {
			t.Fatal("expected comment, got nil")
		}
		if comment.AuthorName().String() != "John Doe" {
			t.Errorf("expected 'John Doe', got '%s'", comment.AuthorName().String())
		}
		if comment.AuthorEmail().String() != "john@example.com" {
			t.Errorf("expected 'john@example.com', got '%s'", comment.AuthorEmail().String())
		}
		if comment.Content().String() != "Great post!" {
			t.Errorf("expected 'Great post!', got '%s'", comment.Content().String())
		}
	})

	t.Run("comment has valid ID after creation", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)

		if comment.ID().String() == "" {
			t.Error("expected non-empty comment ID")
		}
	})

	t.Run("comment has default status 'pending'", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)

		if comment.Status().String() != "pending" {
			t.Errorf("expected 'pending', got '%s'", comment.Status().String())
		}
	})
}

func TestComment_Validation(t *testing.T) {
	t.Run("invalid author_name (empty) returns validation error", func(t *testing.T) {
		_, err := valueobject.NewAuthorName("")
		if err == nil {
			t.Error("expected error for empty author name, got nil")
		}
	})

	t.Run("invalid author_email format returns validation error", func(t *testing.T) {
		_, err := valueobject.NewAuthorEmail("invalid-email")
		if err == nil {
			t.Error("expected error for invalid email, got nil")
		}
	})

	t.Run("invalid status returns validation error", func(t *testing.T) {
		_, err := valueobject.NewCommentStatus("invalid")
		if err == nil {
			t.Error("expected error for invalid status, got nil")
		}
	})
}

func TestComment_BusinessMethods(t *testing.T) {
	t.Run("Approve method changes status to approved", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)
		comment.Approve()

		if comment.Status().String() != "approved" {
			t.Errorf("expected 'approved', got '%s'", comment.Status().String())
		}
	})

	t.Run("Reject method changes status to rejected", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)
		comment.Reject()

		if comment.Status().String() != "rejected" {
			t.Errorf("expected 'rejected', got '%s'", comment.Status().String())
		}
	})

	t.Run("MarkAsSpam method changes status to spam", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)
		comment.MarkAsSpam()

		if comment.Status().String() != "spam" {
			t.Errorf("expected 'spam', got '%s'", comment.Status().String())
		}
	})
}

func TestComment_JSONSerialization(t *testing.T) {
	t.Run("JSON serialization includes all fields with snake_case tags", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)

		data, err := json.Marshal(comment)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		// Check snake_case field names
		expectedFields := []string{"id", "author_name", "author_email", "content", "status", "created_at"}
		for _, field := range expectedFields {
			if _, ok := result[field]; !ok {
				t.Errorf("expected field '%s' in JSON output", field)
			}
		}
	})
}

func TestComment_UpdateContent(t *testing.T) {
	t.Run("UpdateContent updates comment content", func(t *testing.T) {
		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		content, _ := valueobject.NewContent("Great post!")

		comment := NewComment(authorName, authorEmail, content)

		newContent, _ := valueobject.NewContent("Updated comment")
		comment.UpdateContent(newContent)

		if comment.Content().String() != "Updated comment" {
			t.Errorf("expected 'Updated comment', got '%s'", comment.Content().String())
		}
	})
}

func TestComment_InvalidInputs(t *testing.T) {
	t.Run("author name too long returns validation error", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		_, err := valueobject.NewAuthorName(longName)
		if err == nil {
			t.Error("expected error for long author name, got nil")
		}
	})

	t.Run("author email without domain returns validation error", func(t *testing.T) {
		_, err := valueobject.NewAuthorEmail("test@")
		if err == nil {
			t.Error("expected error for email without domain, got nil")
		}
	})
}
