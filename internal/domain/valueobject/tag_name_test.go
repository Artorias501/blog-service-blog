package valueobject

import (
	"strings"
	"testing"
)

func TestTagName_Validation(t *testing.T) {
	t.Run("valid tag name creates TagName successfully", func(t *testing.T) {
		tagName, err := NewTagName("golang")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if tagName.String() != "golang" {
			t.Errorf("expected 'golang', got '%s'", tagName.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewTagName("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("tag name with 51 characters returns validation error", func(t *testing.T) {
		longTagName := strings.Repeat("a", 51)
		_, err := NewTagName(longTagName)
		if err == nil {
			t.Error("expected error for 51 character tag name, got nil")
		}
	})

	t.Run("tag name with exactly 50 characters succeeds", func(t *testing.T) {
		exactTagName := strings.Repeat("a", 50)
		tagName, err := NewTagName(exactTagName)
		if err != nil {
			t.Errorf("expected no error for 50 character tag name, got: %v", err)
		}
		if tagName.String() != exactTagName {
			t.Error("tag name mismatch")
		}
	})

	t.Run("tag name with exactly 1 character succeeds", func(t *testing.T) {
		tagName, err := NewTagName("a")
		if err != nil {
			t.Errorf("expected no error for 1 character tag name, got: %v", err)
		}
		if tagName.String() != "a" {
			t.Error("tag name mismatch")
		}
	})
}
