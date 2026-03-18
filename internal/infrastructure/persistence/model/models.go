package model

import (
	"time"
)

// PostModel represents the database model for Post entity
type PostModel struct {
	ID          string  `gorm:"primaryKey;type:varchar(36)"`
	Title       string  `gorm:"type:varchar(200);not null"`
	Content     string  `gorm:"type:text;not null"`
	Summary     *string `gorm:"type:varchar(500)"`
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tags        []TagModel     `gorm:"many2many:post_tags;joinForeignKey:PostID;joinReferences:TagID"`
	Comments    []CommentModel `gorm:"foreignKey:PostID"`
}

// TableName returns the table name for PostModel
func (PostModel) TableName() string {
	return "posts"
}

// TagModel represents the database model for Tag entity
type TagModel struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time
	Posts     []PostModel `gorm:"many2many:post_tags;joinForeignKey:TagID;joinReferences:PostID"`
}

// TableName returns the table name for TagModel
func (TagModel) TableName() string {
	return "tags"
}

// CommentModel represents the database model for Comment entity
type CommentModel struct {
	ID          string `gorm:"primaryKey;type:varchar(36)"`
	PostID      string `gorm:"type:varchar(36);index;not null"`
	AuthorName  string `gorm:"type:varchar(100);not null"`
	AuthorEmail string `gorm:"type:varchar(255);not null"`
	Content     string `gorm:"type:text;not null"`
	Status      string `gorm:"type:varchar(20);default:'pending';not null"`
	CreatedAt   time.Time
}

// TableName returns the table name for CommentModel
func (CommentModel) TableName() string {
	return "comments"
}

// PostTagModel represents the many-to-many relationship between Post and Tag
type PostTagModel struct {
	PostID    string    `gorm:"primaryKey;type:varchar(36)"`
	TagID     string    `gorm:"primaryKey;type:varchar(36)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for PostTagModel
func (PostTagModel) TableName() string {
	return "post_tags"
}
