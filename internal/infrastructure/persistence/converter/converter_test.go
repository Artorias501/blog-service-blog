package converter

import (
	"testing"
	"time"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
	"github.com/artorias501/blog-service/internal/infrastructure/persistence/model"
)

func TestPostToModel(t *testing.T) {
	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	postModel := PostToModel(post)

	if postModel.ID != post.ID().String() {
		t.Errorf("expected ID %s, got %s", post.ID().String(), postModel.ID)
	}

	if postModel.Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got '%s'", postModel.Title)
	}

	if postModel.Content != "Test content" {
		t.Errorf("expected content 'Test content', got '%s'", postModel.Content)
	}
}

func TestPostToModelWithSummary(t *testing.T) {
	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	summary, _ := valueobject.NewSummary("Test summary")
	post.SetSummary(summary)

	postModel := PostToModel(post)

	if postModel.Summary == nil {
		t.Error("expected summary to be set")
	}

	if *postModel.Summary != "Test summary" {
		t.Errorf("expected summary 'Test summary', got '%s'", *postModel.Summary)
	}
}

func TestPostToModelWithTags(t *testing.T) {
	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	tagName, _ := valueobject.NewTagName("golang")
	tag := entity.NewTag(tagName)
	post.AddTag(tag)

	postModel := PostToModel(post)

	if len(postModel.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(postModel.Tags))
	}

	if postModel.Tags[0].Name != "golang" {
		t.Errorf("expected tag name 'golang', got '%s'", postModel.Tags[0].Name)
	}
}

func TestPostToEntity(t *testing.T) {
	now := time.Now().UTC()
	summary := "Test summary"
	postModel := &model.PostModel{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Title:     "Test Post",
		Content:   "Test content",
		Summary:   &summary,
		CreatedAt: now,
		UpdatedAt: now,
	}

	post, err := PostToEntityWithoutRelations(postModel)
	if err != nil {
		t.Fatalf("failed to convert post: %v", err)
	}

	if post.ID().String() != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected ID '550e8400-e29b-41d4-a716-446655440000', got '%s'", post.ID().String())
	}

	if post.Title().String() != "Test Post" {
		t.Errorf("expected title 'Test Post', got '%s'", post.Title().String())
	}

	if post.Content().String() != "Test content" {
		t.Errorf("expected content 'Test content', got '%s'", post.Content().String())
	}

	if post.Summary() == nil {
		t.Error("expected summary to be set")
	}
}

func TestTagToModel(t *testing.T) {
	name, _ := valueobject.NewTagName("golang")
	tag := entity.NewTag(name)

	tagModel := TagToModel(tag)

	if tagModel.ID != tag.ID().String() {
		t.Errorf("expected ID %s, got %s", tag.ID().String(), tagModel.ID)
	}

	if tagModel.Name != "golang" {
		t.Errorf("expected name 'golang', got '%s'", tagModel.Name)
	}
}

func TestTagToEntity(t *testing.T) {
	now := time.Now().UTC()
	tagModel := &model.TagModel{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		Name:      "golang",
		CreatedAt: now,
	}

	tag, err := TagToEntity(tagModel)
	if err != nil {
		t.Fatalf("failed to convert tag: %v", err)
	}

	if tag.ID().String() != "550e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("expected ID '550e8400-e29b-41d4-a716-446655440001', got '%s'", tag.ID().String())
	}

	if tag.Name().String() != "golang" {
		t.Errorf("expected name 'golang', got '%s'", tag.Name().String())
	}
}

func TestCommentToModel(t *testing.T) {
	authorName, _ := valueobject.NewAuthorName("John Doe")
	authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
	content, _ := valueobject.NewContent("Test comment")
	comment := entity.NewComment(authorName, authorEmail, content)

	commentModel := CommentToModel(comment)

	if commentModel.ID != comment.ID().String() {
		t.Errorf("expected ID %s, got %s", comment.ID().String(), commentModel.ID)
	}

	if commentModel.AuthorName != "John Doe" {
		t.Errorf("expected author name 'John Doe', got '%s'", commentModel.AuthorName)
	}

	if commentModel.AuthorEmail != "john@example.com" {
		t.Errorf("expected author email 'john@example.com', got '%s'", commentModel.AuthorEmail)
	}

	if commentModel.Status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", commentModel.Status)
	}
}

func TestCommentToModelWithPostID(t *testing.T) {
	authorName, _ := valueobject.NewAuthorName("John Doe")
	authorEmail, _ := valueobject.NewAuthorEmail("john@example.com")
	content, _ := valueobject.NewContent("Test comment")
	comment := entity.NewComment(authorName, authorEmail, content)

	postID, _ := valueobject.NewPostID("550e8400-e29b-41d4-a716-446655440000")
	commentModel := CommentToModelWithPostID(comment, postID)

	if commentModel.PostID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected post ID '550e8400-e29b-41d4-a716-446655440000', got '%s'", commentModel.PostID)
	}
}

func TestCommentToEntity(t *testing.T) {
	now := time.Now().UTC()
	commentModel := &model.CommentModel{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		PostID:      "550e8400-e29b-41d4-a716-446655440000",
		AuthorName:  "John Doe",
		AuthorEmail: "john@example.com",
		Content:     "Test comment",
		Status:      "approved",
		CreatedAt:   now,
	}

	comment, err := CommentToEntity(commentModel)
	if err != nil {
		t.Fatalf("failed to convert comment: %v", err)
	}

	if comment.ID().String() != "550e8400-e29b-41d4-a716-446655440002" {
		t.Errorf("expected ID '550e8400-e29b-41d4-a716-446655440002', got '%s'", comment.ID().String())
	}

	if comment.AuthorName().String() != "John Doe" {
		t.Errorf("expected author name 'John Doe', got '%s'", comment.AuthorName().String())
	}

	if comment.Status().String() != "approved" {
		t.Errorf("expected status 'approved', got '%s'", comment.Status().String())
	}
}

func TestPostToEntityWithRelations(t *testing.T) {
	now := time.Now().UTC()
	postModel := &model.PostModel{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Title:     "Test Post",
		Content:   "Test content",
		CreatedAt: now,
		UpdatedAt: now,
		Tags: []model.TagModel{
			{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				Name:      "golang",
				CreatedAt: now,
			},
		},
		Comments: []model.CommentModel{
			{
				ID:          "550e8400-e29b-41d4-a716-446655440002",
				PostID:      "550e8400-e29b-41d4-a716-446655440000",
				AuthorName:  "John Doe",
				AuthorEmail: "john@example.com",
				Content:     "Test comment",
				Status:      "pending",
				CreatedAt:   now,
			},
		},
	}

	post, err := PostToEntity(postModel)
	if err != nil {
		t.Fatalf("failed to convert post: %v", err)
	}

	if len(post.Tags()) != 1 {
		t.Errorf("expected 1 tag, got %d", len(post.Tags()))
	}

	if len(post.Comments()) != 1 {
		t.Errorf("expected 1 comment, got %d", len(post.Comments()))
	}
}
