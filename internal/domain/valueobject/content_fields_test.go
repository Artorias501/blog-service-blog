package valueobject

import (
	"strings"
	"testing"
)

func TestTitle_Validation(t *testing.T) {
	t.Run("valid title creates Title successfully", func(t *testing.T) {
		title, err := NewTitle("Valid Title")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if title.String() != "Valid Title" {
			t.Errorf("expected 'Valid Title', got '%s'", title.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewTitle("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("title with 201 characters returns validation error", func(t *testing.T) {
		longTitle := strings.Repeat("a", 201)
		_, err := NewTitle(longTitle)
		if err == nil {
			t.Error("expected error for 201 character title, got nil")
		}
	})

	t.Run("title with exactly 200 characters succeeds", func(t *testing.T) {
		exactTitle := strings.Repeat("a", 200)
		title, err := NewTitle(exactTitle)
		if err != nil {
			t.Errorf("expected no error for 200 character title, got: %v", err)
		}
		if title.String() != exactTitle {
			t.Error("title mismatch")
		}
	})

	t.Run("title with exactly 1 character succeeds", func(t *testing.T) {
		title, err := NewTitle("a")
		if err != nil {
			t.Errorf("expected no error for 1 character title, got: %v", err)
		}
		if title.String() != "a" {
			t.Error("title mismatch")
		}
	})
}

func TestContent_Validation(t *testing.T) {
	t.Run("non-empty content creates Content successfully", func(t *testing.T) {
		content, err := NewContent("This is valid content")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if content.String() != "This is valid content" {
			t.Errorf("expected 'This is valid content', got '%s'", content.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewContent("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("content with whitespace only returns validation error", func(t *testing.T) {
		_, err := NewContent("   ")
		if err == nil {
			t.Error("expected error for whitespace-only content, got nil")
		}
	})
}

func TestSummary_Validation(t *testing.T) {
	t.Run("valid summary creates Summary successfully", func(t *testing.T) {
		summary, err := NewSummary("This is a valid summary")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if summary.String() != "This is a valid summary" {
			t.Errorf("expected 'This is a valid summary', got '%s'", summary.String())
		}
	})

	t.Run("empty summary is allowed", func(t *testing.T) {
		summary, err := NewSummary("")
		if err != nil {
			t.Errorf("expected no error for empty summary, got: %v", err)
		}
		if summary.String() != "" {
			t.Errorf("expected empty string, got '%s'", summary.String())
		}
	})

	t.Run("summary with 501 characters returns validation error", func(t *testing.T) {
		longSummary := strings.Repeat("a", 501)
		_, err := NewSummary(longSummary)
		if err == nil {
			t.Error("expected error for 501 character summary, got nil")
		}
	})

	t.Run("summary with exactly 500 characters succeeds", func(t *testing.T) {
		exactSummary := strings.Repeat("a", 500)
		summary, err := NewSummary(exactSummary)
		if err != nil {
			t.Errorf("expected no error for 500 character summary, got: %v", err)
		}
		if summary.String() != exactSummary {
			t.Error("summary mismatch")
		}
	})
}
