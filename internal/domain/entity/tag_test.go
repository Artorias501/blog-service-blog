package entity

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestTag_Creation(t *testing.T) {
	t.Run("valid tag name creates Tag successfully", func(t *testing.T) {
		tagName, _ := valueobject.NewTagName("golang")
		tag := NewTag(tagName)

		if tag == nil {
			t.Fatal("expected tag, got nil")
		}
		if tag.Name().String() != "golang" {
			t.Errorf("expected 'golang', got '%s'", tag.Name().String())
		}
	})

	t.Run("tag has valid ID after creation", func(t *testing.T) {
		tagName, _ := valueobject.NewTagName("golang")
		tag := NewTag(tagName)

		if tag.ID().String() == "" {
			t.Error("expected non-empty tag ID")
		}
	})
}

func TestTag_Invariants(t *testing.T) {
	t.Run("invalid tag name (empty) returns validation error", func(t *testing.T) {
		_, err := valueobject.NewTagName("")
		if err == nil {
			t.Error("expected error for empty tag name, got nil")
		}
	})

	t.Run("invalid tag name (>50 chars) returns validation error", func(t *testing.T) {
		longTagName := strings.Repeat("a", 51)
		_, err := valueobject.NewTagName(longTagName)
		if err == nil {
			t.Error("expected error for long tag name, got nil")
		}
	})
}

func TestTag_JSONSerialization(t *testing.T) {
	t.Run("JSON serialization includes all fields with snake_case tags", func(t *testing.T) {
		tagName, _ := valueobject.NewTagName("golang")
		tag := NewTag(tagName)

		data, err := json.Marshal(tag)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		// Check snake_case field names
		expectedFields := []string{"id", "name", "created_at"}
		for _, field := range expectedFields {
			if _, ok := result[field]; !ok {
				t.Errorf("expected field '%s' in JSON output", field)
			}
		}
	})
}

func TestTag_UpdateName(t *testing.T) {
	t.Run("UpdateName updates tag name", func(t *testing.T) {
		tagName, _ := valueobject.NewTagName("golang")
		tag := NewTag(tagName)

		newTagName, _ := valueobject.NewTagName("go")
		tag.UpdateName(newTagName)

		if tag.Name().String() != "go" {
			t.Errorf("expected 'go', got '%s'", tag.Name().String())
		}
	})
}
