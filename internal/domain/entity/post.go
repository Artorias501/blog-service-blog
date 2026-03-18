package entity

import (
	"time"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Post represents a blog post aggregate root
type Post struct {
	id          valueobject.PostID
	title       valueobject.Title
	content     valueobject.Content
	summary     *valueobject.Summary
	comments    []*Comment
	tags        []*Tag
	createdAt   valueobject.CreatedAt
	updatedAt   valueobject.UpdatedAt
	publishedAt valueobject.PublishedAt
}

// NewPost creates a new Post with validated data
func NewPost(title valueobject.Title, content valueobject.Content) *Post {
	now := time.Now().UTC()
	return &Post{
		id:          valueobject.GeneratePostID(),
		title:       title,
		content:     content,
		comments:    make([]*Comment, 0),
		tags:        make([]*Tag, 0),
		createdAt:   valueobject.NewCreatedAt(now),
		updatedAt:   valueobject.NewUpdatedAt(now),
		publishedAt: valueobject.NewPublishedAt(nil),
	}
}

// ID returns the post's identifier
func (p *Post) ID() valueobject.PostID {
	return p.id
}

// Title returns the post's title
func (p *Post) Title() valueobject.Title {
	return p.title
}

// Content returns the post's content
func (p *Post) Content() valueobject.Content {
	return p.content
}

// Summary returns the post's summary (may be nil)
func (p *Post) Summary() *valueobject.Summary {
	return p.summary
}

// Comments returns all comments on the post
func (p *Post) Comments() []*Comment {
	return p.comments
}

// Tags returns all tags associated with the post
func (p *Post) Tags() []*Tag {
	return p.tags
}

// CreatedAt returns the post's creation timestamp
func (p *Post) CreatedAt() valueobject.CreatedAt {
	return p.createdAt
}

// UpdatedAt returns the post's last update timestamp
func (p *Post) UpdatedAt() valueobject.UpdatedAt {
	return p.updatedAt
}

// PublishedAt returns the post's publication timestamp (may be nil)
func (p *Post) PublishedAt() *time.Time {
	return p.publishedAt.Time()
}

// SetSummary sets the post's summary
func (p *Post) SetSummary(summary valueobject.Summary) {
	p.summary = &summary
	p.touch()
}

// AddComment adds a new comment to the post
func (p *Post) AddComment(authorName valueobject.AuthorName, authorEmail valueobject.AuthorEmail, content valueobject.Content) {
	comment := NewComment(authorName, authorEmail, content)
	p.comments = append(p.comments, comment)
	p.touch()
}

// AddTag adds a tag to the post
func (p *Post) AddTag(tag *Tag) {
	// Check if tag already exists
	for _, t := range p.tags {
		if t.ID().Equals(tag.ID()) {
			return
		}
	}
	p.tags = append(p.tags, tag)
	p.touch()
}

// RemoveTag removes a tag from the post
func (p *Post) RemoveTag(tagID valueobject.TagID) {
	for i, t := range p.tags {
		if t.ID().Equals(tagID) {
			p.tags = append(p.tags[:i], p.tags[i+1:]...)
			p.touch()
			return
		}
	}
}

// Publish sets the publication timestamp
func (p *Post) Publish() {
	now := time.Now().UTC()
	p.publishedAt = valueobject.NewPublishedAt(&now)
	p.touch()
}

// UpdateContent updates the post's content
func (p *Post) UpdateContent(content valueobject.Content) {
	p.content = content
	p.touch()
}

// UpdateTitle updates the post's title
func (p *Post) UpdateTitle(title valueobject.Title) {
	p.title = title
	p.touch()
}

// touch updates the updatedAt timestamp
func (p *Post) touch() {
	p.updatedAt = valueobject.NewUpdatedAt(time.Now().UTC())
}

// ReconstructPost reconstructs a Post entity from persistence data
// This function is used by the infrastructure layer to rebuild entities from database
func ReconstructPost(post *Post, id valueobject.PostID, title valueobject.Title, content valueobject.Content, createdAt valueobject.CreatedAt, updatedAt valueobject.UpdatedAt, publishedAt valueobject.PublishedAt) {
	post.id = id
	post.title = title
	post.content = content
	post.createdAt = createdAt
	post.updatedAt = updatedAt
	post.publishedAt = publishedAt
	post.comments = make([]*Comment, 0)
	post.tags = make([]*Tag, 0)
}

// AddCommentToPost adds an existing comment to a post (used for reconstruction)
func AddCommentToPost(post *Post, comment *Comment) {
	post.comments = append(post.comments, comment)
}

// SetPostID sets the post ID (used for testing)
func (p *Post) SetPostID(id valueobject.PostID) {
	p.id = id
}

// SetPostTimestamps sets the post timestamps (used for testing)
func (p *Post) SetPostTimestamps(createdAt valueobject.CreatedAt, updatedAt valueobject.UpdatedAt) {
	p.createdAt = createdAt
	p.updatedAt = updatedAt
}

// PostJSON is used for JSON serialization
type PostJSON struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Summary     *string `json:"summary"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	PublishedAt *string `json:"published_at"`
}

// MarshalJSON implements json.Marshaler
func (p Post) MarshalJSON() ([]byte, error) {
	var summary *string
	if p.summary != nil {
		s := p.summary.String()
		summary = &s
	}

	var publishedAt *string
	if p.publishedAt.Time() != nil {
		pa := p.publishedAt.Time().Format(time.RFC3339)
		publishedAt = &pa
	}

	return jsonMarshal(PostJSON{
		ID:          p.id.String(),
		Title:       p.title.String(),
		Content:     p.content.String(),
		Summary:     summary,
		CreatedAt:   p.createdAt.Time().Format(time.RFC3339),
		UpdatedAt:   p.updatedAt.Time().Format(time.RFC3339),
		PublishedAt: publishedAt,
	})
}
