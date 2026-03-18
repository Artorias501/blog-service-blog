package converter

import (
	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
	"github.com/artorias501/blog-service/internal/infrastructure/persistence/model"
)

// PostToModel converts a Post entity to a PostModel
func PostToModel(post *entity.Post) *model.PostModel {
	m := &model.PostModel{
		ID:          post.ID().String(),
		Title:       post.Title().String(),
		Content:     post.Content().String(),
		Summary:     nil,
		PublishedAt: post.PublishedAt(),
		CreatedAt:   post.CreatedAt().Time(),
		UpdatedAt:   post.UpdatedAt().Time(),
		Tags:        make([]model.TagModel, 0),
		Comments:    make([]model.CommentModel, 0),
	}

	// Convert summary if present
	if s := post.Summary(); s != nil {
		summary := s.String()
		m.Summary = &summary
	}

	// Convert tags
	for _, tag := range post.Tags() {
		m.Tags = append(m.Tags, model.TagModel{
			ID:        tag.ID().String(),
			Name:      tag.Name().String(),
			CreatedAt: tag.CreatedAt().Time(),
		})
	}

	// Convert comments
	for _, comment := range post.Comments() {
		m.Comments = append(m.Comments, CommentToModel(comment))
	}

	return m
}

// PostToEntity converts a PostModel to a Post entity
func PostToEntity(m *model.PostModel) (*entity.Post, error) {
	// Create title value object
	title, err := valueobject.NewTitle(m.Title)
	if err != nil {
		return nil, err
	}

	// Create content value object
	content, err := valueobject.NewContent(m.Content)
	if err != nil {
		return nil, err
	}

	// Create post ID
	postID, err := valueobject.NewPostID(m.ID)
	if err != nil {
		return nil, err
	}

	// Create timestamps
	createdAt := valueobject.NewCreatedAt(m.CreatedAt)
	updatedAt := valueobject.NewUpdatedAt(m.UpdatedAt)
	publishedAt := valueobject.NewPublishedAt(m.PublishedAt)

	// Build post using reconstruction pattern
	post := &entity.Post{}
	setPostFields(post, postID, title, content, createdAt, updatedAt, publishedAt)

	// Set summary if present
	if m.Summary != nil {
		summary, err := valueobject.NewSummary(*m.Summary)
		if err != nil {
			return nil, err
		}
		post.SetSummary(summary)
	}

	// Set tags
	for _, tagModel := range m.Tags {
		tag, err := TagToEntity(&tagModel)
		if err != nil {
			return nil, err
		}
		post.AddTag(tag)
	}

	// Set comments
	for _, commentModel := range m.Comments {
		comment, err := CommentToEntity(&commentModel)
		if err != nil {
			return nil, err
		}
		addCommentToPost(post, comment)
	}

	return post, nil
}

// PostToEntityWithoutRelations converts a PostModel to a Post entity without loading relations
func PostToEntityWithoutRelations(m *model.PostModel) (*entity.Post, error) {
	title, err := valueobject.NewTitle(m.Title)
	if err != nil {
		return nil, err
	}

	content, err := valueobject.NewContent(m.Content)
	if err != nil {
		return nil, err
	}

	postID, err := valueobject.NewPostID(m.ID)
	if err != nil {
		return nil, err
	}

	createdAt := valueobject.NewCreatedAt(m.CreatedAt)
	updatedAt := valueobject.NewUpdatedAt(m.UpdatedAt)
	publishedAt := valueobject.NewPublishedAt(m.PublishedAt)

	post := &entity.Post{}
	setPostFields(post, postID, title, content, createdAt, updatedAt, publishedAt)

	if m.Summary != nil {
		summary, err := valueobject.NewSummary(*m.Summary)
		if err != nil {
			return nil, err
		}
		post.SetSummary(summary)
	}

	return post, nil
}

// TagToModel converts a Tag entity to a TagModel
func TagToModel(tag *entity.Tag) *model.TagModel {
	return &model.TagModel{
		ID:        tag.ID().String(),
		Name:      tag.Name().String(),
		CreatedAt: tag.CreatedAt().Time(),
	}
}

// TagToEntity converts a TagModel to a Tag entity
func TagToEntity(m *model.TagModel) (*entity.Tag, error) {
	tagID, err := valueobject.NewTagID(m.ID)
	if err != nil {
		return nil, err
	}

	tagName, err := valueobject.NewTagName(m.Name)
	if err != nil {
		return nil, err
	}

	createdAt := valueobject.NewCreatedAt(m.CreatedAt)

	return newTagFromPersistence(tagID, tagName, createdAt), nil
}

// CommentToModel converts a Comment entity to a CommentModel
func CommentToModel(comment *entity.Comment) model.CommentModel {
	return model.CommentModel{
		ID:          comment.ID().String(),
		PostID:      "", // PostID should be set separately when needed
		AuthorName:  comment.AuthorName().String(),
		AuthorEmail: comment.AuthorEmail().String(),
		Content:     comment.Content().String(),
		Status:      comment.Status().String(),
		CreatedAt:   comment.CreatedAt().Time(),
	}
}

// CommentToModelWithPostID converts a Comment entity to a CommentModel with PostID
func CommentToModelWithPostID(comment *entity.Comment, postID valueobject.PostID) model.CommentModel {
	return model.CommentModel{
		ID:          comment.ID().String(),
		PostID:      postID.String(),
		AuthorName:  comment.AuthorName().String(),
		AuthorEmail: comment.AuthorEmail().String(),
		Content:     comment.Content().String(),
		Status:      comment.Status().String(),
		CreatedAt:   comment.CreatedAt().Time(),
	}
}

// CommentToEntity converts a CommentModel to a Comment entity
func CommentToEntity(m *model.CommentModel) (*entity.Comment, error) {
	commentID, err := valueobject.NewCommentID(m.ID)
	if err != nil {
		return nil, err
	}

	authorName, err := valueobject.NewAuthorName(m.AuthorName)
	if err != nil {
		return nil, err
	}

	authorEmail, err := valueobject.NewAuthorEmail(m.AuthorEmail)
	if err != nil {
		return nil, err
	}

	content, err := valueobject.NewContent(m.Content)
	if err != nil {
		return nil, err
	}

	status, err := valueobject.NewCommentStatus(m.Status)
	if err != nil {
		return nil, err
	}

	createdAt := valueobject.NewCreatedAt(m.CreatedAt)

	return newCommentFromPersistence(commentID, authorName, authorEmail, content, status, createdAt), nil
}

// Helper functions for entity reconstruction (using reflection-like patterns)
// These are implemented via internal entity package exports

// setPostFields sets the internal fields of a Post entity
func setPostFields(post *entity.Post, id valueobject.PostID, title valueobject.Title, content valueobject.Content, createdAt valueobject.CreatedAt, updatedAt valueobject.UpdatedAt, publishedAt valueobject.PublishedAt) {
	// Use the entity's internal reconstruction method if available
	// Otherwise, we need to use exported setters
	entity.ReconstructPost(post, id, title, content, createdAt, updatedAt, publishedAt)
}

// addCommentToPost adds a comment to a post using internal method
func addCommentToPost(post *entity.Post, comment *entity.Comment) {
	entity.AddCommentToPost(post, comment)
}

// newTagFromPersistence creates a Tag from persistence data
func newTagFromPersistence(id valueobject.TagID, name valueobject.TagName, createdAt valueobject.CreatedAt) *entity.Tag {
	return entity.NewTagFromPersistence(id, name, createdAt)
}

// newCommentFromPersistence creates a Comment from persistence data
func newCommentFromPersistence(id valueobject.CommentID, authorName valueobject.AuthorName, authorEmail valueobject.AuthorEmail, content valueobject.Content, status valueobject.CommentStatus, createdAt valueobject.CreatedAt) *entity.Comment {
	return entity.NewCommentFromPersistence(id, authorName, authorEmail, content, status, createdAt)
}
