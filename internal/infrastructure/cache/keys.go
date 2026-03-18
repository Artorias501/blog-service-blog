package cache

import (
	"fmt"

	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Post key generators

// PostKey generates a cache key for a single post
// Format: post:{id}
func PostKey(id valueobject.PostID) string {
	return fmt.Sprintf("post:%s", id.String())
}

// PostListKey generates a cache key for a paginated list of posts
// Format: post:list:page:{page}:size:{size}
func PostListKey(page, pageSize int) string {
	return fmt.Sprintf("post:list:page:%d:size:%d", page, pageSize)
}

// PostByTagKey generates a cache key for posts filtered by tag
// Format: post:tag:{tagID}:page:{page}:size:{size}
func PostByTagKey(tagID valueobject.TagID, page, pageSize int) string {
	return fmt.Sprintf("post:tag:%s:page:%d:size:%d", tagID.String(), page, pageSize)
}

// PostPattern generates a pattern for all post-related keys
// Format: post:*
func PostPattern() string {
	return "post:*"
}

// PostByIDPattern generates a pattern for all keys related to a specific post
// Format: post:{id}*
func PostByIDPattern(id valueobject.PostID) string {
	return fmt.Sprintf("post:%s*", id.String())
}

// Tag key generators

// TagKey generates a cache key for a single tag
// Format: tag:{id}
func TagKey(id valueobject.TagID) string {
	return fmt.Sprintf("tag:%s", id.String())
}

// TagByNameKey generates a cache key for a tag lookup by name
// Format: tag:name:{name}
func TagByNameKey(name valueobject.TagName) string {
	return fmt.Sprintf("tag:name:%s", name.String())
}

// TagListKey generates a cache key for a paginated list of tags
// Format: tag:list:page:{page}:size:{size}
func TagListKey(page, pageSize int) string {
	return fmt.Sprintf("tag:list:page:%d:size:%d", page, pageSize)
}

// TagByPostKey generates a cache key for tags associated with a post
// Format: tag:post:{postID}
func TagByPostKey(postID valueobject.PostID) string {
	return fmt.Sprintf("tag:post:%s", postID.String())
}

// TagPattern generates a pattern for all tag-related keys
// Format: tag:*
func TagPattern() string {
	return "tag:*"
}

// TagByIDPattern generates a pattern for all keys related to a specific tag
// Format: tag:{id}*
func TagByIDPattern(id valueobject.TagID) string {
	return fmt.Sprintf("tag:%s*", id.String())
}

// Comment key generators

// CommentKey generates a cache key for a single comment
// Format: comment:{id}
func CommentKey(id valueobject.CommentID) string {
	return fmt.Sprintf("comment:%s", id.String())
}

// CommentByPostKey generates a cache key for comments of a specific post
// Format: comment:post:{postID}:page:{page}:size:{size}
func CommentByPostKey(postID valueobject.PostID, page, pageSize int) string {
	return fmt.Sprintf("comment:post:%s:page:%d:size:%d", postID.String(), page, pageSize)
}

// CommentCountByPostKey generates a cache key for comment count of a post
// Format: comment:count:post:{postID}
func CommentCountByPostKey(postID valueobject.PostID) string {
	return fmt.Sprintf("comment:count:post:%s", postID.String())
}

// CommentPattern generates a pattern for all comment-related keys
// Format: comment:*
func CommentPattern() string {
	return "comment:*"
}

// CommentByIDPattern generates a pattern for all keys related to a specific comment
// Format: comment:{id}*
func CommentByIDPattern(id valueobject.CommentID) string {
	return fmt.Sprintf("comment:%s*", id.String())
}

// CommentByPostPattern generates a pattern for all comment keys for a specific post
// Format: comment:post:{postID}*
func CommentByPostPattern(postID valueobject.PostID) string {
	return fmt.Sprintf("comment:post:%s*", postID.String())
}
