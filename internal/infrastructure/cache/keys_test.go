package cache

import (
	"testing"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestPostKeys(t *testing.T) {
	postID, _ := valueobject.NewPostID("550e8400-e29b-41d4-a716-446655440000")

	t.Run("post key", func(t *testing.T) {
		key := PostKey(postID)
		expected := "post:550e8400-e29b-41d4-a716-446655440000"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("post list key", func(t *testing.T) {
		key := PostListKey(1, 10)
		expected := "post:list:page:1:size:10"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("post list key with different pagination", func(t *testing.T) {
		key := PostListKey(2, 20)
		expected := "post:list:page:2:size:20"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("post by tag key", func(t *testing.T) {
		tagID, _ := valueobject.NewTagID("660e8400-e29b-41d4-a716-446655440001")
		key := PostByTagKey(tagID, 1, 10)
		expected := "post:tag:660e8400-e29b-41d4-a716-446655440001:page:1:size:10"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("post pattern key", func(t *testing.T) {
		pattern := PostPattern()
		expected := "post:*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})

	t.Run("post by ID pattern", func(t *testing.T) {
		pattern := PostByIDPattern(postID)
		expected := "post:550e8400-e29b-41d4-a716-446655440000*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})
}

func TestTagKeys(t *testing.T) {
	tagID, _ := valueobject.NewTagID("660e8400-e29b-41d4-a716-446655440001")
	tagName, _ := valueobject.NewTagName("golang")
	postID, _ := valueobject.NewPostID("550e8400-e29b-41d4-a716-446655440000")

	t.Run("tag key", func(t *testing.T) {
		key := TagKey(tagID)
		expected := "tag:660e8400-e29b-41d4-a716-446655440001"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("tag by name key", func(t *testing.T) {
		key := TagByNameKey(tagName)
		expected := "tag:name:golang"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("tag list key", func(t *testing.T) {
		key := TagListKey(1, 10)
		expected := "tag:list:page:1:size:10"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("tag by post key", func(t *testing.T) {
		key := TagByPostKey(postID)
		expected := "tag:post:550e8400-e29b-41d4-a716-446655440000"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("tag pattern key", func(t *testing.T) {
		pattern := TagPattern()
		expected := "tag:*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})

	t.Run("tag by ID pattern", func(t *testing.T) {
		pattern := TagByIDPattern(tagID)
		expected := "tag:660e8400-e29b-41d4-a716-446655440001*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})
}

func TestCommentKeys(t *testing.T) {
	commentID, _ := valueobject.NewCommentID("770e8400-e29b-41d4-a716-446655440002")
	postID, _ := valueobject.NewPostID("550e8400-e29b-41d4-a716-446655440000")

	t.Run("comment key", func(t *testing.T) {
		key := CommentKey(commentID)
		expected := "comment:770e8400-e29b-41d4-a716-446655440002"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("comment by post key", func(t *testing.T) {
		key := CommentByPostKey(postID, 1, 10)
		expected := "comment:post:550e8400-e29b-41d4-a716-446655440000:page:1:size:10"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("comment count by post key", func(t *testing.T) {
		key := CommentCountByPostKey(postID)
		expected := "comment:count:post:550e8400-e29b-41d4-a716-446655440000"
		if key != expected {
			t.Errorf("Expected %s, got: %s", expected, key)
		}
	})

	t.Run("comment pattern key", func(t *testing.T) {
		pattern := CommentPattern()
		expected := "comment:*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})

	t.Run("comment by ID pattern", func(t *testing.T) {
		pattern := CommentByIDPattern(commentID)
		expected := "comment:770e8400-e29b-41d4-a716-446655440002*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})

	t.Run("comment by post pattern", func(t *testing.T) {
		pattern := CommentByPostPattern(postID)
		expected := "comment:post:550e8400-e29b-41d4-a716-446655440000*"
		if pattern != expected {
			t.Errorf("Expected %s, got: %s", expected, pattern)
		}
	})
}
