package valueobject

import (
	"testing"
)

func TestPostID_Validation(t *testing.T) {
	t.Run("valid UUID creates PostID successfully", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		postID, err := NewPostID(validUUID)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if postID.String() != validUUID {
			t.Errorf("expected %s, got %s", validUUID, postID.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewPostID("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("invalid UUID format returns validation error", func(t *testing.T) {
		_, err := NewPostID("invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID, got nil")
		}
	})

	t.Run("UUID with wrong format returns validation error", func(t *testing.T) {
		_, err := NewPostID("550e8400-e29b-41d4-a716")
		if err == nil {
			t.Error("expected error for malformed UUID, got nil")
		}
	})

	t.Run("Generate creates valid PostID", func(t *testing.T) {
		postID := GeneratePostID()
		if postID.String() == "" {
			t.Error("expected non-empty PostID")
		}
	})
}

func TestTagID_Validation(t *testing.T) {
	t.Run("valid UUID creates TagID successfully", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440001"
		tagID, err := NewTagID(validUUID)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if tagID.String() != validUUID {
			t.Errorf("expected %s, got %s", validUUID, tagID.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewTagID("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("invalid UUID format returns validation error", func(t *testing.T) {
		_, err := NewTagID("invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID, got nil")
		}
	})

	t.Run("Generate creates valid TagID", func(t *testing.T) {
		tagID := GenerateTagID()
		if tagID.String() == "" {
			t.Error("expected non-empty TagID")
		}
	})
}

func TestCommentID_Validation(t *testing.T) {
	t.Run("valid UUID creates CommentID successfully", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440002"
		commentID, err := NewCommentID(validUUID)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if commentID.String() != validUUID {
			t.Errorf("expected %s, got %s", validUUID, commentID.String())
		}
	})

	t.Run("empty string returns validation error", func(t *testing.T) {
		_, err := NewCommentID("")
		if err == nil {
			t.Error("expected error for empty string, got nil")
		}
	})

	t.Run("invalid UUID format returns validation error", func(t *testing.T) {
		_, err := NewCommentID("invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID, got nil")
		}
	})

	t.Run("Generate creates valid CommentID", func(t *testing.T) {
		commentID := GenerateCommentID()
		if commentID.String() == "" {
			t.Error("expected non-empty CommentID")
		}
	})
}
