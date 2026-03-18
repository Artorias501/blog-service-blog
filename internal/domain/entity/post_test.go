package entity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestPost_Creation(t *testing.T) {
	t.Run("valid post data creates Post successfully", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		if post == nil {
			t.Fatal("expected post, got nil")
		}
		if post.Title().String() != "Test Post" {
			t.Errorf("expected 'Test Post', got '%s'", post.Title().String())
		}
		if post.Content().String() != "Test content" {
			t.Errorf("expected 'Test content', got '%s'", post.Content().String())
		}
	})

	t.Run("post has valid ID after creation", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		if post.ID().String() == "" {
			t.Error("expected non-empty post ID")
		}
	})
}

func TestPost_Invariants(t *testing.T) {
	t.Run("invalid title (empty) returns validation error", func(t *testing.T) {
		_, err := valueobject.NewTitle("")
		if err == nil {
			t.Error("expected error for empty title, got nil")
		}
	})

	t.Run("invalid title (>200 chars) returns validation error", func(t *testing.T) {
		longTitle := strings.Repeat("a", 201)
		_, err := valueobject.NewTitle(longTitle)
		if err == nil {
			t.Error("expected error for long title, got nil")
		}
	})

	t.Run("invalid content (empty) returns validation error", func(t *testing.T) {
		_, err := valueobject.NewContent("")
		if err == nil {
			t.Error("expected error for empty content, got nil")
		}
	})
}

func TestPost_BusinessMethods(t *testing.T) {
	t.Run("AddComment method adds comment to post", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		authorName, _ := valueobject.NewAuthorName("John Doe")
		authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
		commentContent, _ := valueobject.NewContent("Great post!")

		post.AddComment(authorName, authorEmail, commentContent)

		comments := post.Comments()
		if len(comments) != 1 {
			t.Errorf("expected 1 comment, got %d", len(comments))
		}
	})

	t.Run("AddTag method adds tag to post", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		tagName, _ := valueobject.NewTagName("golang")
		tag := NewTag(tagName)

		post.AddTag(tag)

		tags := post.Tags()
		if len(tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(tags))
		}
	})

	t.Run("RemoveTag method removes tag from post", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		tagName, _ := valueobject.NewTagName("golang")
		tag := NewTag(tagName)

		post.AddTag(tag)
		post.RemoveTag(tag.ID())

		tags := post.Tags()
		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("Publish method sets published_at timestamp", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		if post.PublishedAt() != nil {
			t.Error("expected nil published_at before publish")
		}

		post.Publish()

		if post.PublishedAt() == nil {
			t.Error("expected non-nil published_at after publish")
		}
	})

	t.Run("SetSummary method sets summary", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		summary, _ := valueobject.NewSummary("Test summary")
		post.SetSummary(summary)

		if post.Summary() == nil || post.Summary().String() != "Test summary" {
			t.Error("summary mismatch")
		}
	})
}

func TestPost_JSONSerialization(t *testing.T) {
	t.Run("JSON serialization includes all fields with snake_case tags", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		data, err := json.Marshal(post)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		// Check snake_case field names
		expectedFields := []string{"id", "title", "content", "summary", "created_at", "updated_at", "published_at"}
		for _, field := range expectedFields {
			if _, ok := result[field]; !ok {
				t.Errorf("expected field '%s' in JSON output", field)
			}
		}
	})
}

func TestPost_UpdateContent(t *testing.T) {
	t.Run("UpdateContent updates content and sets updated_at", func(t *testing.T) {
		title, _ := valueobject.NewTitle("Test Post")
		content, _ := valueobject.NewContent("Test content")
		post := NewPost(title, content)

		originalUpdatedAt := post.UpdatedAt()
		time.Sleep(time.Millisecond * 10) // ensure time difference

		newContent, _ := valueobject.NewContent("Updated content")
		post.UpdateContent(newContent)

		if post.Content().String() != "Updated content" {
			t.Error("content not updated")
		}
		if !post.UpdatedAt().After(originalUpdatedAt) {
			t.Error("updated_at not updated")
		}
	})
}
