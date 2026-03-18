package entity

import (
	"time"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Tag represents a tag aggregate root
type Tag struct {
	id        valueobject.TagID
	name      valueobject.TagName
	createdAt valueobject.CreatedAt
}

// NewTag creates a new Tag with validated data
func NewTag(name valueobject.TagName) *Tag {
	return &Tag{
		id:        valueobject.GenerateTagID(),
		name:      name,
		createdAt: valueobject.NewCreatedAt(time.Now().UTC()),
	}
}

// ID returns the tag's identifier
func (t *Tag) ID() valueobject.TagID {
	return t.id
}

// Name returns the tag's name
func (t *Tag) Name() valueobject.TagName {
	return t.name
}

// CreatedAt returns the tag's creation timestamp
func (t *Tag) CreatedAt() valueobject.CreatedAt {
	return t.createdAt
}

// UpdateName updates the tag's name
func (t *Tag) UpdateName(name valueobject.TagName) {
	t.name = name
}

// NewTagFromPersistence reconstructs a Tag entity from persistence data
func NewTagFromPersistence(id valueobject.TagID, name valueobject.TagName, createdAt valueobject.CreatedAt) *Tag {
	return &Tag{
		id:        id,
		name:      name,
		createdAt: createdAt,
	}
}

// TagJSON is used for JSON serialization
type TagJSON struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// MarshalJSON implements json.Marshaler
func (t Tag) MarshalJSON() ([]byte, error) {
	return jsonMarshal(TagJSON{
		ID:        t.id.String(),
		Name:      t.name.String(),
		CreatedAt: t.createdAt.Time().Format(time.RFC3339),
	})
}
