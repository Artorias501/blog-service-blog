package entity

import (
	"encoding/json"
	"time"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Comment represents a comment entity (part of Post aggregate)
type Comment struct {
	id          valueobject.CommentID
	authorName  valueobject.AuthorName
	authorEmail valueobject.AuthorEmail
	content     valueobject.Content
	status      valueobject.CommentStatus
	createdAt   valueobject.CreatedAt
}

// NewComment creates a new Comment with validated data
func NewComment(authorName valueobject.AuthorName, authorEmail valueobject.AuthorEmail, content valueobject.Content) *Comment {
	return &Comment{
		id:          valueobject.GenerateCommentID(),
		authorName:  authorName,
		authorEmail: authorEmail,
		content:     content,
		status:      valueobject.DefaultCommentStatus(),
		createdAt:   valueobject.NewCreatedAt(time.Now().UTC()),
	}
}

// ID returns the comment's identifier
func (c *Comment) ID() valueobject.CommentID {
	return c.id
}

// AuthorName returns the comment's author name
func (c *Comment) AuthorName() valueobject.AuthorName {
	return c.authorName
}

// AuthorEmail returns the comment's author email
func (c *Comment) AuthorEmail() valueobject.AuthorEmail {
	return c.authorEmail
}

// Content returns the comment's content
func (c *Comment) Content() valueobject.Content {
	return c.content
}

// Status returns the comment's status
func (c *Comment) Status() valueobject.CommentStatus {
	return c.status
}

// CreatedAt returns the comment's creation timestamp
func (c *Comment) CreatedAt() valueobject.CreatedAt {
	return c.createdAt
}

// Approve changes the comment status to approved
func (c *Comment) Approve() {
	c.status, _ = valueobject.NewCommentStatus(valueobject.CommentStatusApproved)
}

// Reject changes the comment status to rejected
func (c *Comment) Reject() {
	c.status, _ = valueobject.NewCommentStatus(valueobject.CommentStatusRejected)
}

// MarkAsSpam changes the comment status to spam
func (c *Comment) MarkAsSpam() {
	c.status, _ = valueobject.NewCommentStatus(valueobject.CommentStatusSpam)
}

// UpdateContent updates the comment's content
func (c *Comment) UpdateContent(content valueobject.Content) {
	c.content = content
}

// NewCommentFromPersistence reconstructs a Comment entity from persistence data
func NewCommentFromPersistence(id valueobject.CommentID, authorName valueobject.AuthorName, authorEmail valueobject.AuthorEmail, content valueobject.Content, status valueobject.CommentStatus, createdAt valueobject.CreatedAt) *Comment {
	return &Comment{
		id:          id,
		authorName:  authorName,
		authorEmail: authorEmail,
		content:     content,
		status:      status,
		createdAt:   createdAt,
	}
}

// CommentJSON is used for JSON serialization
type CommentJSON struct {
	ID          string `json:"id"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// MarshalJSON implements json.Marshaler
func (c Comment) MarshalJSON() ([]byte, error) {
	return jsonMarshal(CommentJSON{
		ID:          c.id.String(),
		AuthorName:  c.authorName.String(),
		AuthorEmail: c.authorEmail.String(),
		Content:     c.content.String(),
		Status:      c.status.String(),
		CreatedAt:   c.createdAt.Time().Format(time.RFC3339),
	})
}

// jsonMarshal is a helper function to marshal JSON
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
