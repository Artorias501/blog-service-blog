package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
	"github.com/artorias501/blog-service/internal/infrastructure/database"
	"github.com/artorias501/blog-service/internal/infrastructure/persistence/model"
)

// Test helpers

func setupTestDB(t *testing.T) *gorm.DB {
	// Create temp directory for test database
	tmpDir := filepath.Join(os.TempDir(), "blog-test")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create unique database file for each test
	dbPath := filepath.Join(tmpDir, fmt.Sprintf("test_%d.db", time.Now().UnixNano()))

	cfg := database.Config{
		DSN:             dbPath,
		LogLevel:        4, // Silent
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	db, err := database.NewConnectionWithMigrate(cfg)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Cleanup function
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		os.Remove(dbPath)
	})

	return db
}

func createTestPost(t *testing.T) *entity.Post {
	title, err := valueobject.NewTitle("Test Post Title")
	if err != nil {
		t.Fatalf("failed to create title: %v", err)
	}

	content, err := valueobject.NewContent("This is test content for the post.")
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}

	return entity.NewPost(title, content)
}

func createTestTag(t *testing.T) *entity.Tag {
	name, err := valueobject.NewTagName("test-tag")
	if err != nil {
		t.Fatalf("failed to create tag name: %v", err)
	}
	return entity.NewTag(name)
}

func createTestComment(t *testing.T) *entity.Comment {
	authorName, err := valueobject.NewAuthorName("Test Author")
	if err != nil {
		t.Fatalf("failed to create author name: %v", err)
	}

	authorEmail, err := valueobject.NewAuthorEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create author email: %v", err)
	}

	content, err := valueobject.NewContent("This is a test comment.")
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}

	return entity.NewComment(authorName, authorEmail, content)
}

// PostRepository Tests

func TestPostRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	post := createTestPost(t)

	err := repo.Create(ctx, post)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Verify post was created
	var postModel model.PostModel
	result := db.First(&postModel, "id = ?", post.ID().String())
	if result.Error != nil {
		t.Fatalf("failed to find created post: %v", result.Error)
	}

	if postModel.Title != post.Title().String() {
		t.Errorf("expected title %s, got %s", post.Title().String(), postModel.Title)
	}
}

func TestPostRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	// Create a post first
	post := createTestPost(t)
	if err := repo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Get the post
	retrieved, err := repo.GetByID(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to get post: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected post, got nil")
	}

	if retrieved.ID().String() != post.ID().String() {
		t.Errorf("expected ID %s, got %s", post.ID().String(), retrieved.ID().String())
	}
}

func TestPostRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	id := valueobject.GeneratePostID()
	retrieved, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrieved != nil {
		t.Error("expected nil for non-existent post")
	}
}

func TestPostRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	// Create a post first
	post := createTestPost(t)
	if err := repo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Update the post
	newTitle, _ := valueobject.NewTitle("Updated Title")
	post.UpdateTitle(newTitle)

	err := repo.Update(ctx, post)
	if err != nil {
		t.Fatalf("failed to update post: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to get post: %v", err)
	}

	if retrieved.Title().String() != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", retrieved.Title().String())
	}
}

func TestPostRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	// Create a post first
	post := createTestPost(t)
	if err := repo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Delete the post
	err := repo.Delete(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to delete post: %v", err)
	}

	// Verify deletion
	retrieved, err := repo.GetByID(ctx, post.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrieved != nil {
		t.Error("expected nil after deletion")
	}
}

func TestPostRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	// Create multiple posts
	for i := 0; i < 15; i++ {
		post := createTestPost(t)
		if err := repo.Create(ctx, post); err != nil {
			t.Fatalf("failed to create post: %v", err)
		}
		time.Sleep(time.Millisecond) // Ensure different timestamps
	}

	// List with pagination
	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
		SortBy:   "created_at",
		Order:    "desc",
	}

	result, err := repo.List(ctx, params)
	if err != nil {
		t.Fatalf("failed to list posts: %v", err)
	}

	if result.Total != 15 {
		t.Errorf("expected total 15, got %d", result.Total)
	}

	if len(result.Items) != 10 {
		t.Errorf("expected 10 items, got %d", len(result.Items))
	}

	if result.TotalPage != 2 {
		t.Errorf("expected 2 total pages, got %d", result.TotalPage)
	}
}

func TestPostRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)
	ctx := context.Background()

	// Create posts with specific content
	title1, _ := valueobject.NewTitle("Golang Tutorial")
	content1, _ := valueobject.NewContent("Learn Go programming")
	post1 := entity.NewPost(title1, content1)

	title2, _ := valueobject.NewTitle("Python Guide")
	content2, _ := valueobject.NewContent("Learn Python programming")
	post2 := entity.NewPost(title2, content2)

	if err := repo.Create(ctx, post1); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := repo.Create(ctx, post2); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Search for "Golang"
	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result, err := repo.Search(ctx, "Golang", params)
	if err != nil {
		t.Fatalf("failed to search posts: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("expected 1 result, got %d", result.Total)
	}
}

func TestPostRepository_GetByIDWithTags(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	tagRepo := NewTagRepository(db)
	ctx := context.Background()

	// Create post and tag
	post := createTestPost(t)
	tag := createTestTag(t)

	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := tagRepo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Associate tag with post
	if err := postRepo.AddTag(ctx, post.ID(), tag.ID()); err != nil {
		t.Fatalf("failed to add tag: %v", err)
	}

	// Get post with tags
	retrieved, err := postRepo.GetByIDWithTags(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to get post with tags: %v", err)
	}

	if len(retrieved.Tags()) != 1 {
		t.Errorf("expected 1 tag, got %d", len(retrieved.Tags()))
	}
}

func TestPostRepository_AddTag(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	tagRepo := NewTagRepository(db)
	ctx := context.Background()

	// Create post and tag
	post := createTestPost(t)
	tag := createTestTag(t)

	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := tagRepo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Add tag to post
	err := postRepo.AddTag(ctx, post.ID(), tag.ID())
	if err != nil {
		t.Fatalf("failed to add tag: %v", err)
	}

	// Verify association
	var postTag model.PostTagModel
	result := db.First(&postTag, "post_id = ? AND tag_id = ?", post.ID().String(), tag.ID().String())
	if result.Error != nil {
		t.Fatalf("failed to find post_tag association: %v", result.Error)
	}
}

func TestPostRepository_RemoveTag(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	tagRepo := NewTagRepository(db)
	ctx := context.Background()

	// Create post and tag
	post := createTestPost(t)
	tag := createTestTag(t)

	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := tagRepo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Add tag to post
	if err := postRepo.AddTag(ctx, post.ID(), tag.ID()); err != nil {
		t.Fatalf("failed to add tag: %v", err)
	}

	// Remove tag from post
	err := postRepo.RemoveTag(ctx, post.ID(), tag.ID())
	if err != nil {
		t.Fatalf("failed to remove tag: %v", err)
	}

	// Verify removal
	var count int64
	db.Model(&model.PostTagModel{}).Where("post_id = ? AND tag_id = ?", post.ID().String(), tag.ID().String()).Count(&count)
	if count != 0 {
		t.Error("expected tag to be removed from post")
	}
}

func TestPostRepository_ListByTagID(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	tagRepo := NewTagRepository(db)
	ctx := context.Background()

	// Create posts and tag
	post1 := createTestPost(t)
	post2 := createTestPost(t)
	tag := createTestTag(t)

	if err := postRepo.Create(ctx, post1); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := postRepo.Create(ctx, post2); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := tagRepo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Associate tag with post1 only
	if err := postRepo.AddTag(ctx, post1.ID(), tag.ID()); err != nil {
		t.Fatalf("failed to add tag: %v", err)
	}

	// List posts by tag
	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result, err := postRepo.ListByTagID(ctx, tag.ID(), params)
	if err != nil {
		t.Fatalf("failed to list posts by tag: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("expected 1 post, got %d", result.Total)
	}
}

// TagRepository Tests

func TestTagRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)

	err := repo.Create(ctx, tag)
	if err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Verify tag was created
	var tagModel model.TagModel
	result := db.First(&tagModel, "id = ?", tag.ID().String())
	if result.Error != nil {
		t.Fatalf("failed to find created tag: %v", result.Error)
	}
}

func TestTagRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	if err := repo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, tag.ID())
	if err != nil {
		t.Fatalf("failed to get tag: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected tag, got nil")
	}

	if retrieved.ID().String() != tag.ID().String() {
		t.Errorf("expected ID %s, got %s", tag.ID().String(), retrieved.ID().String())
	}
}

func TestTagRepository_GetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	if err := repo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	retrieved, err := repo.GetByName(ctx, tag.Name())
	if err != nil {
		t.Fatalf("failed to get tag by name: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected tag, got nil")
	}

	if retrieved.Name().String() != tag.Name().String() {
		t.Errorf("expected name %s, got %s", tag.Name().String(), retrieved.Name().String())
	}
}

func TestTagRepository_GetOrCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	name, _ := valueobject.NewTagName("unique-tag")

	// First call should create
	tag1, err := repo.GetOrCreate(ctx, name)
	if err != nil {
		t.Fatalf("failed to get or create tag: %v", err)
	}

	// Second call should return existing
	tag2, err := repo.GetOrCreate(ctx, name)
	if err != nil {
		t.Fatalf("failed to get or create tag: %v", err)
	}

	if tag1.ID().String() != tag2.ID().String() {
		t.Error("expected same tag on second call")
	}
}

func TestTagRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	// Create tags
	name1, _ := valueobject.NewTagName("golang")
	name2, _ := valueobject.NewTagName("python")
	tag1 := entity.NewTag(name1)
	tag2 := entity.NewTag(name2)

	if err := repo.Create(ctx, tag1); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}
	if err := repo.Create(ctx, tag2); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result, err := repo.Search(ctx, "go", params)
	if err != nil {
		t.Fatalf("failed to search tags: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("expected 1 result, got %d", result.Total)
	}
}

func TestTagRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	if err := repo.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	err := repo.Delete(ctx, tag.ID())
	if err != nil {
		t.Fatalf("failed to delete tag: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, tag.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrieved != nil {
		t.Error("expected nil after deletion")
	}
}

// CommentRepository Tests

func TestCommentRepository_CreateWithPostID(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post first
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Create comment
	comment := createTestComment(t)
	err := commentRepo.CreateWithPostID(ctx, comment, post.ID())
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	// Verify comment was created
	var commentModel model.CommentModel
	result := db.First(&commentModel, "id = ?", comment.ID().String())
	if result.Error != nil {
		t.Fatalf("failed to find created comment: %v", result.Error)
	}

	if commentModel.PostID != post.ID().String() {
		t.Errorf("expected post_id %s, got %s", post.ID().String(), commentModel.PostID)
	}
}

func TestCommentRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post and comment
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	comment := createTestComment(t)
	if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	// Get comment
	retrieved, err := commentRepo.GetByID(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to get comment: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected comment, got nil")
	}

	if retrieved.ID().String() != comment.ID().String() {
		t.Errorf("expected ID %s, got %s", comment.ID().String(), retrieved.ID().String())
	}
}

func TestCommentRepository_ListByPostID(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Create multiple comments
	for i := 0; i < 5; i++ {
		comment := createTestComment(t)
		if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
			t.Fatalf("failed to create comment: %v", err)
		}
	}

	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result, err := commentRepo.ListByPostID(ctx, post.ID(), params)
	if err != nil {
		t.Fatalf("failed to list comments: %v", err)
	}

	if result.Total != 5 {
		t.Errorf("expected 5 comments, got %d", result.Total)
	}
}

func TestCommentRepository_CountByPostID(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Create comments
	for i := 0; i < 3; i++ {
		comment := createTestComment(t)
		if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
			t.Fatalf("failed to create comment: %v", err)
		}
	}

	count, err := commentRepo.CountByPostID(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to count comments: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

func TestCommentRepository_Approve(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post and comment
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	comment := createTestComment(t)
	if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	// Approve comment
	err := commentRepo.Approve(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to approve comment: %v", err)
	}

	// Verify status
	retrieved, err := commentRepo.GetByID(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to get comment: %v", err)
	}

	if retrieved.Status().String() != valueobject.CommentStatusApproved {
		t.Errorf("expected status approved, got %s", retrieved.Status().String())
	}
}

func TestCommentRepository_Reject(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post and comment
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	comment := createTestComment(t)
	if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	// Reject comment
	err := commentRepo.Reject(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to reject comment: %v", err)
	}

	// Verify status
	retrieved, err := commentRepo.GetByID(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to get comment: %v", err)
	}

	if retrieved.Status().String() != valueobject.CommentStatusRejected {
		t.Errorf("expected status rejected, got %s", retrieved.Status().String())
	}
}

func TestCommentRepository_MarkAsSpam(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post and comment
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	comment := createTestComment(t)
	if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	// Mark as spam
	err := commentRepo.MarkAsSpam(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to mark comment as spam: %v", err)
	}

	// Verify status
	retrieved, err := commentRepo.GetByID(ctx, comment.ID())
	if err != nil {
		t.Fatalf("failed to get comment: %v", err)
	}

	if retrieved.Status().String() != valueobject.CommentStatusSpam {
		t.Errorf("expected status spam, got %s", retrieved.Status().String())
	}
}

func TestCommentRepository_DeleteByPostID(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Create comments
	for i := 0; i < 3; i++ {
		comment := createTestComment(t)
		if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
			t.Fatalf("failed to create comment: %v", err)
		}
	}

	// Delete all comments for the post
	err := commentRepo.DeleteByPostID(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to delete comments: %v", err)
	}

	// Verify deletion
	count, err := commentRepo.CountByPostID(ctx, post.ID())
	if err != nil {
		t.Fatalf("failed to count comments: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 comments after deletion, got %d", count)
	}
}

func TestCommentRepository_ListByStatus(t *testing.T) {
	db := setupTestDB(t)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)
	ctx := context.Background()

	// Create post
	post := createTestPost(t)
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Create comments
	for i := 0; i < 3; i++ {
		comment := createTestComment(t)
		if err := commentRepo.CreateWithPostID(ctx, comment, post.ID()); err != nil {
			t.Fatalf("failed to create comment: %v", err)
		}
	}

	// Approve one comment
	comments, err := commentRepo.ListByPostID(ctx, post.ID(), repository.ListParams{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("failed to list comments: %v", err)
	}
	if len(comments.Items) > 0 {
		if err := commentRepo.Approve(ctx, comments.Items[0].ID()); err != nil {
			t.Fatalf("failed to approve comment: %v", err)
		}
	}

	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	// List approved comments
	result, err := commentRepo.ListByStatus(ctx, valueobject.CommentStatusApproved, params)
	if err != nil {
		t.Fatalf("failed to list comments by status: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("expected 1 approved comment, got %d", result.Total)
	}
}

// Interface implementation verification

func TestPostRepository_ImplementsInterface(t *testing.T) {
	var _ repository.PostRepository = NewPostRepository(nil)
}

func TestTagRepository_ImplementsInterface(t *testing.T) {
	var _ repository.TagRepository = NewTagRepository(nil)
}

func TestCommentRepository_ImplementsInterface(t *testing.T) {
	var _ repository.CommentRepository = NewCommentRepository(nil)
}
