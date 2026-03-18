# Blog Service

> High Performance AI Blog Platform - Core Blog Microservice

## Overview

Blog Service is the core microservice responsible for article CRUD operations, caching, and search functionality. It follows Domain-Driven Design (DDD) principles with clear separation between domain logic and infrastructure concerns.

## Architecture

### Layer Structure (DDD)

```
blog-service/
├── cmd/
│   └── server/           # Application entry point
│       └── main.go
├── internal/
│   ├── domain/           # Domain Layer - Core business logic
│   │   ├── entity/       # Entities & Aggregate Roots
│   │   ├── valueobject/  # Value Objects
│   │   └── repository/   # Repository Interfaces (contracts)
│   ├── service/          # Application Layer - Use cases
│   ├── handler/          # Interface Layer - HTTP handlers
│   │   ├── dto/          # Data Transfer Objects
│   │   └── middleware/   # HTTP middlewares
│   └── infrastructure/   # Infrastructure Layer
│ ├── persistence/
│ │ ├── sqlite/ # SQLite repository implementations
│ │ └── redis/ # Redis repository implementations
│       └── cache/        # Cache service implementations
└── pkg/                  # Shared utilities
    ├── config/
    ├── logger/
    └── response/
```

### Dependency Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Request                            │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   Handler Layer                             │
│  - Request validation                                       │
│  - DTO → Entity conversion                                  │
│  - HTTP response formatting                                 │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   Service Layer                             │
│  - Business logic orchestration                             │
│  - Transaction management                                   │
│  - Cache strategy                                           │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
│  - Entity business rules                                    │
│  - Value object validation                                  │
│  - Domain events                                            │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│ Repository Interfaces │
│ - PostRepository (SQLite) │
│ - PostCacheRepository (Redis) │
│ - TagRepository (SQLite) │
│ - CommentRepository (SQLite) │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│ Infrastructure Implementations │
│ - SQLite repositories (GORM) │
│ - Redis repositories (go-redis) │
└─────────────────────────────────────────────────────────────┘
```

---

## Domain Model Design

### Aggregate Roots

#### 1. Post (Article) - Primary Aggregate Root

```
Post (Aggregate Root)
├── id: PostID (Value Object)
├── title: Title (Value Object)
├── content: Content (Value Object)
├── summary: Summary (Value Object, AI-generated)
├── tags: []Tag (Associated entities, external aggregate)
├── comments: []Comment (Embedded entities, part of aggregate)
├── likeCount: LikeCount (Value Object)
├── publishedAt: PublishedAt (Value Object)
├── createdAt: CreatedAt (Value Object)
└── updatedAt: UpdatedAt (Value Object)
```

**Invariants:**
- Title must not be empty and max 200 characters
- Content must not be empty
- Comments are part of Post aggregate (cascade operations)
- Tags are referenced by ID (separate aggregate)

#### 2. Tag - Secondary Aggregate Root

```
Tag (Aggregate Root)
├── id: TagID (Value Object)
├── name: TagName (Value Object)
└── createdAt: CreatedAt (Value Object)
```

**Invariants:**
- Tag name must be unique
- Tag name must not be empty, max 50 characters

### Entities

#### Comment (Part of Post Aggregate)

```
Comment (Entity)
├── id: CommentID (Value Object)
├── postId: PostID (Reference to aggregate root)
├── authorName: AuthorName (Value Object)
├── authorEmail: AuthorEmail (Value Object)
├── content: Content (Value Object)
├── createdAt: CreatedAt (Value Object)
└── status: CommentStatus (Value Object: visible/hidden)
```

### Value Objects

| Value Object | Type | Validation Rules |
|-------------|------|------------------|
| PostID | UUID v4 | Valid UUID format |
| Title | string | 1-200 characters, not empty |
| Content | string (Markdown) | Not empty |
| Summary | string | Max 500 characters |
| LikeCount | int | >= 0 |
| PublishedAt | time.Time | Not null |
| TagID | UUID v4 | Valid UUID format |
| TagName | string | 1-50 characters, unique |
| CommentID | UUID v4 | Valid UUID format |
| AuthorName | string | 1-100 characters |
| AuthorEmail | string | Valid email format |
| CommentStatus | enum | visible, hidden |

---

## Database Schema Design

### SQLite Tables

#### posts table

```sql
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    like_count INTEGER NOT NULL DEFAULT 0,
    published_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

#### tags table

```sql
CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
```

#### post_tags table (Many-to-Many)

```sql
CREATE TABLE post_tags (
    post_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
```

#### comments table

```sql
CREATE TABLE comments (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    author_name TEXT NOT NULL,
    author_email TEXT NOT NULL,
    content TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'visible' CHECK(status IN ('visible', 'hidden')),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);
```

### Entity Relationship Diagram

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│     posts    │       │   post_tags  │       │     tags     │
├──────────────┤       ├──────────────┤       ├──────────────┤
│ id (PK)      │◄──────│ post_id (FK) │       │ id (PK)      │
│ title        │       │ tag_id (FK)  │──────►│ name         │
│ content      │       └──────────────┘       │ created_at   │
│ summary      │                              └──────────────┘
│ like_count   │
│ published_at │       ┌──────────────┐
│ created_at   │       │   comments   │
│ updated_at   │       ├──────────────┤
└──────────────┘       │ id (PK)      │
        │              │ post_id (FK) │
        │              │ author_name  │
        └─────────────►│ author_email │
                       │ content      │
                       │ status       │
                       │ created_at   │
                       └──────────────┘
```

---

## Cache Design (Redis)

### Key Naming Convention

| Key Pattern | Description | TTL | Data Structure |
|------------|-------------|-----|----------------|
| `post:{id}` | Single post cache | 1 hour | Hash |
| `post:list:page:{page}` | Post list by page | 5 min | String (JSON) |
| `post:list:tag:{tag_id}:page:{page}` | Posts by tag | 5 min | String (JSON) |
| `post:search:{query_hash}` | Search results | 1 min | String (JSON) |
| `tag:all` | All tags list | 24 hours | String (JSON) |
| `comment:post:{post_id}` | Comments by post | 10 min | String (JSON) |

### Cache Data Structures

#### post:{id} (Hash)

```
{
    "id": "uuid-string",
    "title": "article title",
    "content": "markdown content",
    "summary": "ai generated summary",
    "like_count": "42",
    "published_at": "2024-01-15T10:00:00Z",
    "tags": "[{\"id\":\"...\",\"name\":\"Go\"}]"
}
```

### Cache Strategy

#### Read-Through Pattern

```
1. Check Redis cache
2. If hit: return cached data
3. If miss:
   a. Query SQLite
   b. Write to Redis
   c. Return data
```

#### Write-Through Pattern (for Post creation/update)

```
1. Write to SQLite (primary)
2. Invalidate related Redis keys:
   - post:{id}
   - post:list:*
   - post:search:*
   - post:list:tag:*
```

#### Cache Invalidation Rules

| Operation | Keys to Invalidate |
|-----------|-------------------|
| Create Post | `post:list:*`, `tag:all` |
| Update Post | `post:{id}`, `post:list:*`, `post:search:*` |
| Delete Post | `post:{id}`, `post:list:*`, `post:search:*` |
| Add Comment | `comment:post:{post_id}` |
| Like Post | `post:{id}` |

---

## Repository Interfaces Design

### Separation of SQLite and Redis Repositories

```go
// internal/domain/repository/post_repository.go

// PostRepository - Persistence contract
type PostRepository interface {
    // CRUD Operations
    Create(ctx context.Context, post *entity.Post) error
    GetByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error)
    Update(ctx context.Context, post *entity.Post) error
    Delete(ctx context.Context, id valueobject.PostID) error
    
    // Query Operations
    List(ctx context.Context, params ListParams) (*PostListResult, error)
    SearchByTitle(ctx context.Context, keyword string, params ListParams) (*PostListResult, error)
    SearchByDateRange(ctx context.Context, start, end time.Time, params ListParams) (*PostListResult, error)
    SearchByTag(ctx context.Context, tagID valueobject.TagID, params ListParams) (*PostListResult, error)
    
    // Count Operations
    Count(ctx context.Context) (int64, error)
    CountByTag(ctx context.Context, tagID valueobject.TagID) (int64, error)
}

// PostCacheRepository - Redis cache contract
type PostCacheRepository interface {
    // Single Post Operations
    Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error)
    Set(ctx context.Context, post *entity.Post) error
    Delete(ctx context.Context, id valueobject.PostID) error
    
    // List Operations
    GetList(ctx context.Context, key string) (*PostListResult, error)
    SetList(ctx context.Context, key string, result *PostListResult) error
    DeleteList(ctx context.Context, pattern string) error
    
    // Search Cache Operations
    GetSearchResult(ctx context.Context, queryHash string) (*PostListResult, error)
    SetSearchResult(ctx context.Context, queryHash string, result *PostListResult) error
}

// ListParams - Pagination parameters
type ListParams struct {
    Page     int
    PageSize int
    SortBy   string // "published_at", "like_count", "created_at"
    Order    string // "asc", "desc"
}

// PostListResult - Paginated result
type PostListResult struct {
    Posts     []*entity.Post
    Total     int64
    Page      int
    PageSize  int
    TotalPage int
}
```

```go
// internal/domain/repository/tag_repository.go

type TagRepository interface {
    Create(ctx context.Context, tag *entity.Tag) error
    GetByID(ctx context.Context, id valueobject.TagID) (*entity.Tag, error)
    GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error)
    GetAll(ctx context.Context) ([]*entity.Tag, error)
    Delete(ctx context.Context, id valueobject.TagID) error
}

type TagCacheRepository interface {
    GetAll(ctx context.Context) ([]*entity.Tag, error)
    SetAll(ctx context.Context, tags []*entity.Tag) error
    Invalidate(ctx context.Context) error
}
```

```go
// internal/domain/repository/comment_repository.go

type CommentRepository interface {
    Create(ctx context.Context, comment *entity.Comment) error
    GetByPostID(ctx context.Context, postID valueobject.PostID, params ListParams) (*CommentListResult, error)
    Delete(ctx context.Context, id valueobject.CommentID) error
    UpdateStatus(ctx context.Context, id valueobject.CommentID, status valueobject.CommentStatus) error
}

type CommentCacheRepository interface {
    GetByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Comment, error)
    SetByPostID(ctx context.Context, postID valueobject.PostID, comments []*entity.Comment) error
    DeleteByPostID(ctx context.Context, postID valueobject.PostID) error
}
```

---

## Service Layer Design

### PostService

```go
// internal/service/post_service.go

type PostService struct {
    postRepo     repository.PostRepository
    cacheRepo    repository.PostCacheRepository
    tagRepo      repository.TagRepository
    commentRepo  repository.CommentRepository
}

// Business Methods
func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*entity.Post, error)
func (s *PostService) GetPostByID(ctx context.Context, id string) (*entity.Post, error)
func (s *PostService) UpdatePost(ctx context.Context, id string, req *UpdatePostRequest) (*entity.Post, error)
func (s *PostService) DeletePost(ctx context.Context, id string) error
func (s *PostService) ListPosts(ctx context.Context, params *ListParams) (*PostListResult, error)
func (s *PostService) SearchPosts(ctx context.Context, req *SearchRequest) (*PostListResult, error)
func (s *PostService) LikePost(ctx context.Context, id string) error
```

### Search Implementation

```go
// SearchRequest combines multiple search criteria
type SearchRequest struct {
    Keyword   string     // Search in title
    TagID     string     // Filter by tag
    StartDate *time.Time // Date range filter
    EndDate   *time.Time // Date range filter
    Page      int
    PageSize  int
    SortBy    string
    Order     string
}

// SearchPosts implements combined search
func (s *PostService) SearchPosts(ctx context.Context, req *SearchRequest) (*PostListResult, error) {
    // 1. Generate cache key from request
    cacheKey := generateSearchCacheKey(req)
    
    // 2. Try cache first
    if result, err := s.cacheRepo.GetSearchResult(ctx, cacheKey); err == nil {
        return result, nil
    }
    
    // 3. Build query based on criteria
    var result *PostListResult
    var err error
    
switch {
    case req.Keyword != "" && req.TagID != "":
        // Combined search
        result, err = s.searchByKeywordAndTag(ctx, req)
    case req.Keyword != "":
        result, err = s.postRepo.SearchByTitle(ctx, req.Keyword, req.toListParams())
    case req.TagID != "":
        tagID := valueobject.NewTagID(req.TagID)
        result, err = s.postRepo.SearchByTag(ctx, tagID, req.toListParams())
    case req.StartDate != nil || req.EndDate != nil:
        result, err = s.postRepo.SearchByDateRange(ctx, req.StartDate, req.EndDate, req.toListParams())
    default:
        result, err = s.postRepo.List(ctx, req.toListParams())
}
    
    // 4. Cache result
    s.cacheRepo.SetSearchResult(ctx, cacheKey, result)
    
    return result, err
}
```

---

## API Interface Design

### RESTful Endpoints

#### Post APIs

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/posts` | List posts (paginated) | Public |
| GET | `/api/posts/:id` | Get single post | Public |
| POST | `/api/posts` | Create post | Admin |
| PUT | `/api/posts/:id` | Update post | Admin |
| DELETE | `/api/posts/:id` | Delete post | Admin |
| POST | `/api/posts/:id/like` | Like a post | Public |
| GET | `/api/posts/search` | Search posts | Public |

#### Tag APIs

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/tags` | List all tags | Public |
| POST | `/api/tags` | Create tag | Admin |
| DELETE | `/api/tags/:id` | Delete tag | Admin |

#### Comment APIs

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/posts/:id/comments` | Get post comments | Public |
| POST | `/api/posts/:id/comments` | Create comment | Public |
| DELETE | `/api/comments/:id` | Delete comment | Admin |
| PUT | `/api/comments/:id/status` | Update comment status | Admin |

### Request/Response DTOs

#### Create Post Request

```json
{
    "title": "Understanding DDD in Go",
    "content": "# Introduction\n\nDomain-Driven Design...",
    "tag_ids": ["uuid-1", "uuid-2"],
    "published_at": "2024-01-15T10:00:00Z"
}
```

#### Post Response

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Understanding DDD in Go",
        "content": "# Introduction\n\nDomain-Driven Design...",
        "summary": "This article explores DDD principles...",
        "like_count": 42,
        "published_at": "2024-01-15T10:00:00Z",
        "created_at": "2024-01-15T09:30:00Z",
        "updated_at": "2024-01-15T10:00:00Z",
        "tags": [
            {"id": "uuid-1", "name": "Go"},
            {"id": "uuid-2", "name": "Architecture"}
        ],
        "comment_count": 5
    }
}
```

#### List Posts Response

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "items": [
            {
                "id": "550e8400-e29b-41d4-a716-446655440000",
                "title": "Understanding DDD in Go",
                "summary": "This article explores...",
                "like_count": 42,
                "published_at": "2024-01-15T10:00:00Z",
                "tags": [{"id": "uuid-1", "name": "Go"}]
            }
        ],
        "total": 100,
        "page": 1,
        "page_size": 10,
        "total_page": 10
    }
}
```

#### Search Request Query Parameters

```
GET /api/posts/search?keyword=DDD&tag_id=uuid-1&start_date=2024-01-01&end_date=2024-12-31&page=1&page_size=10&sort_by=published_at&order=desc
```

#### Create Comment Request

```json
{
    "author_name": "John Doe",
    "author_email": "john@example.com",
    "content": "Great article! Thanks for sharing."
}
```

---

## Handler Layer Design

### Handler Structure

```go
// internal/handler/post_handler.go

type PostHandler struct {
    postService *service.PostService
}

// Handler Methods
func (h *PostHandler) Create(c *gin.Context)      // POST /api/posts
func (h *PostHandler) GetByID(c *gin.Context)     // GET /api/posts/:id
func (h *PostHandler) List(c *gin.Context)        // GET /api/posts
func (h *PostHandler) Update(c *gin.Context)      // PUT /api/posts/:id
func (h *PostHandler) Delete(c *gin.Context)      // DELETE /api/posts/:id
func (h *PostHandler) Like(c *gin.Context)        // POST /api/posts/:id/like
func (h *PostHandler) Search(c *gin.Context)      // GET /api/posts/search
```

### DTO Definitions

```go
// internal/handler/dto/post_dto.go

type CreatePostRequest struct {
    Title       string   `json:"title" binding:"required,max=200"`
    Content     string   `json:"content" binding:"required"`
    TagIDs      []string `json:"tag_ids" binding:"required"`
    PublishedAt string   `json:"published_at" binding:"required"`
}

type UpdatePostRequest struct {
    Title       string   `json:"title" binding:"max=200"`
    Content     string   `json:"content"`
    TagIDs      []string `json:"tag_ids"`
    PublishedAt string   `json:"published_at"`
}

type PostResponse struct {
    ID          string        `json:"id"`
    Title       string        `json:"title"`
    Content     string        `json:"content"`
    Summary     string        `json:"summary"`
    LikeCount   int           `json:"like_count"`
    PublishedAt time.Time     `json:"published_at"`
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    Tags        []TagResponse `json:"tags"`
    CommentCount int          `json:"comment_count"`
}

type PostListItem struct {
    ID          string        `json:"id"`
    Title       string        `json:"title"`
    Summary     string        `json:"summary"`
    LikeCount   int           `json:"like_count"`
    PublishedAt time.Time     `json:"published_at"`
    Tags        []TagResponse `json:"tags"`
}

type ListPostsResponse struct {
    Items     []PostListItem `json:"items"`
    Total     int64          `json:"total"`
    Page      int            `json:"page"`
    PageSize  int            `json:"page_size"`
    TotalPage int            `json:"total_page"`
}
```

### Middleware

```go
// internal/handler/middleware/auth.go

func AdminAuth() gin.HandlerFunc {
    // Validate admin token from header
}

// internal/handler/middleware/cors.go

func CORS() gin.HandlerFunc {
    // Handle CORS for frontend
}

// internal/handler/middleware/logger.go

func RequestLogger() gin.HandlerFunc {
    // Log request details
}
```

---

## Project Files Structure (Detailed)

```
blog-service/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── post.go               # Post aggregate root
│   │   │   ├── tag.go                # Tag aggregate root
│   │   │   └── comment.go            # Comment entity
│   │   ├── valueobject/
│   │   │   ├── post_id.go            # PostID value object
│   │   │   ├── title.go              # Title value object
│   │   │   ├── content.go            # Content value object
│   │   │   ├── summary.go            # Summary value object
│   │   │   ├── like_count.go         # LikeCount value object
│   │   │   ├── tag_id.go             # TagID value object
│   │   │   ├── tag_name.go           # TagName value object
│   │   │   ├── comment_id.go         # CommentID value object
│   │   │   ├── author_name.go        # AuthorName value object
│   │   │   ├── author_email.go       # AuthorEmail value object
│   │   │   ├── comment_status.go     # CommentStatus value object
│   │   │   └── timestamps.go         # CreatedAt, UpdatedAt, PublishedAt
│   │   └── repository/
│ │ ├── post_repository.go # PostRepository, PostCacheRepository
│ │ ├── tag_repository.go # TagRepository, TagCacheRepository
│ │ ├── comment_repository.go # CommentRepository, CommentCacheRepository
│   │       └── params.go              # ListParams, PostListResult, CommentListResult
│   ├── service/
│   │   ├── post_service.go           # Post business logic
│   │   ├── tag_service.go            # Tag business logic
│   │   └── comment_service.go        # Comment business logic
│   ├── handler/
│   │   ├── post_handler.go           # Post HTTP handlers
│   │   ├── tag_handler.go            # Tag HTTP handlers
│   │   ├── comment_handler.go        # Comment HTTP handlers
│   │   ├── dto/
│   │   │   ├── post_dto.go           # Post request/response DTOs
│   │   │   ├── tag_dto.go            # Tag request/response DTOs
│   │   │   ├── comment_dto.go        # Comment request/response DTOs
│   │   │   └── common.go             # Common response structures
│   │   └── middleware/
│   │       ├── auth.go               # Admin authentication middleware
│   │       ├── cors.go               # CORS middleware
│   │       └── logger.go             # Request logging middleware
│   └── infrastructure/
│ ├── persistence/
│ │ ├── sqlite/
│ │ │ ├── connection.go # SQLite connection setup
│ │ │ ├── post_repo.go # PostRepository implementation
│ │ │ ├── tag_repo.go # TagRepository implementation
│ │ │ ├── comment_repo.go # CommentRepository implementation
│ │ │ └── models.go # GORM models
│       │   └── redis/
│       │       ├── connection.go     # Redis connection setup
│       │       ├── post_cache.go     # PostCacheRepository implementation
│       │       ├── tag_cache.go      # TagCacheRepository implementation
│       │       └── comment_cache.go  # CommentCacheRepository implementation
│       └── cache/
│           └── keys.go               # Cache key generators
├── pkg/
│   ├── config/
│   │   └── config.go                 # Configuration loader
│   ├── logger/
│   │   └── logger.go                 # Structured logger
│   └── response/
│       └── response.go               # Standard API response format
├── go.mod                            # Go module definition
├── go.sum                            # Go dependencies checksum
├── Makefile                          # Build and run commands
└── README.md                         # This file
```

---

## Call Flow Example

### Create Post Flow

```
HTTP POST /api/posts
        │
        ▼
┌───────────────────────────────────────────────────────────────┐
│ PostHandler.Create(c *gin.Context)                            │
│ 1. Bind JSON to CreatePostRequest DTO                         │
│ 2. Validate request                                            │
│ 3. Convert DTO to domain entities                              │
│ 4. Call postService.CreatePost(ctx, req)                      │
└───────────────────────────────────────────────────────────────┘
        │
        ▼
┌───────────────────────────────────────────────────────────────┐
│ PostService.CreatePost(ctx, req)                              │
│ 1. Validate tag IDs exist (tagRepo.GetByID)                   │
│ 2. Create Post entity with value objects                       │
│ 3. Begin transaction                                           │
│ 4. Save post (postRepo.Create) │
│ 5. Save post-tag associations                                  │
│ 6. Commit transaction                                          │
│ 7. Invalidate cache (cacheRepo.DeleteList)                    │
│ 8. Return created post                                         │
└───────────────────────────────────────────────────────────────┘
        │
        ▼
┌───────────────────────────────────────────────────────────────┐
│ PostRepository.Create(ctx, post) │
│ 1. Convert entity to GORM model │
│ 2. Execute INSERT query │
│ 3. Return error or nil │
└───────────────────────────────────────────────────────────────┘
```

### Get Post Flow (with Cache)

```
HTTP GET /api/posts/:id
        │
        ▼
┌───────────────────────────────────────────────────────────────┐
│ PostHandler.GetByID(c *gin.Context)                           │
│ 1. Extract ID from URL params                                  │
│ 2. Call postService.GetPostByID(ctx, id)                      │
│ 3. Convert entity to response DTO                              │
│ 4. Return JSON response                                        │
└───────────────────────────────────────────────────────────────┘
        │
        ▼
┌───────────────────────────────────────────────────────────────┐
│ PostService.GetPostByID(ctx, id)                              │
│ 1. Try cache first (cacheRepo.Get)                            │
│ 2. If cache hit: return cached post                           │
│ 3. If cache miss: │
│ a. Query SQLite (postRepo.GetByID) │
│ b. Query tags (tagRepo.GetByPostID) │
│ c. Query comment count │
│ d. Set cache (cacheRepo.Set) │
│ 4. Return post                                                 │
└───────────────────────────────────────────────────────────────┘
```

---

## Configuration

### Environment Variables

```bash
# Server
SERVER_PORT=8080
SERVER_MODE=debug # debug, release

# SQLite
SQLITE_PATH=./data/blog.db

# Redis (optional for initial development)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Cache
CACHE_ENABLED=false # Set to true when Redis is ready
CACHE_TTL_POST=3600 # 1 hour in seconds
CACHE_TTL_LIST=300 # 5 minutes
```

---

## Build & Run

### Prerequisites

- Go 1.21+
- SQLite 3.x
- Redis 7.0+ (optional)

### Commands

```bash
# Install dependencies
go mod download

# Run in development mode
go run cmd/server/main.go

# Build
go build -o bin/blog-service cmd/server/main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

---

## Testing Strategy

### Unit Tests

- Domain layer: Entity methods, value object validation
- Service layer: Mock repositories, business logic

### Integration Tests

- Repository implementations against real SQLite/Redis
- Handler tests with HTTP test server

### Test Coverage Goals

- Domain layer: 90%+
- Service layer: 80%+
- Handler layer: 70%+

---

## Future Enhancements

- [ ] Add AI summary generation (MQ integration)
- [ ] Implement full-text search with Elasticsearch
- [ ] Add rate limiting
- [ ] Add metrics and monitoring
- [ ] Implement graceful shutdown

---

## License

MIT License
