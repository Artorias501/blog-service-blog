# Blog Service API Reference

Base URL: `http://localhost:8080`

All API routes are prefixed with `/api/v1`.

## Authentication

Admin endpoints require Bearer token authentication:

```
Authorization: Bearer <admin_token>
```

## Response Format

### Success Response

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### Error Response

```json
{
  "code": <error_code>,
  "message": "<error_message>",
  "data": null
}
```

### Paginated Response

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "total_page": 10
  }
}
```

## Response Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| 0 | 200/201 | Success |
| 400 | 400 | Bad Request / Validation Error |
| 401 | 401 | Unauthorized |
| 403 | 403 | Forbidden |
| 404 | 404 | Not Found |
| 500 | 500 | Internal Server Error |

---

## Health Check

### GET /health

Returns service health status.

**Request**

No parameters.

**Response**

```json
{
  "status": "healthy",
  "service": "blog-service",
  "version": "1.0.0",
  "checks": {
    "database": {
      "status": "healthy",
      "message": ""
    },
    "redis": {
      "status": "healthy",
      "message": ""
    }
  }
}
```

---

## Posts

### GET /api/v1/posts

List posts with pagination.

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at`, `updated_at`, `title` |
| order | string | No | - | Sort order: `asc`, `desc` |
| tag_id | string (UUID) | No | - | Filter by tag ID |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "uuid",
        "title": "Post Title",
        "content": "Post content...",
        "summary": "Post summary...",
        "tags": [
          {"id": "uuid", "name": "Tag Name"}
        ],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "published_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "total_page": 10
  }
}
```

---

### GET /api/v1/posts/search

Search posts by keyword.

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| keyword | string | Yes | - | Search keyword |
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at`, `updated_at`, `title` |
| order | string | No | - | Sort order: `asc`, `desc` |

**Response**

Same as GET /api/v1/posts.

---

### GET /api/v1/posts/:id

Get a single post by ID.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "title": "Post Title",
    "content": "Post content...",
    "summary": "Post summary...",
    "tags": [
      {"id": "uuid", "name": "Tag Name"}
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "published_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### POST /api/v1/posts

Create a new post. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |
| Content-Type | application/json | Yes |

**Request Body**

```json
{
  "title": "Post Title",
  "content": "Post content...",
  "tag_ids": ["uuid", "uuid"]
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| title | string | Yes | 1-200 characters |
| content | string | Yes | Min 1 character |
| tag_ids | []string (UUID) | No | Array of tag UUIDs |

**Response**

HTTP 201 Created

```json
{
  "code": 0,
  "message": "created",
  "data": {
    "id": "uuid",
    "title": "Post Title",
    "content": "Post content...",
    "summary": "Post summary...",
    "tags": [],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "published_at": null
  }
}
```

---

### PUT /api/v1/posts/:id

Update a post. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |
| Content-Type | application/json | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Request Body**

```json
{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| title | string | No | 1-200 characters |
| content | string | No | Min 1 character |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "title": "Updated Title",
    "content": "Updated content...",
    "summary": "Updated summary...",
    "tags": [],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z",
    "published_at": null
  }
}
```

---

### DELETE /api/v1/posts/:id

Delete a post.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Response**

HTTP 204 No Content (no body)

---

### POST /api/v1/posts/:id/like

Like a post (increment like count).

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Response**

HTTP 204 No Content (no body)

---

### GET /api/v1/posts/:id/tags

Get all tags for a post.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {"id": "uuid", "name": "Tag Name", "created_at": "2024-01-01T00:00:00Z"}
  ]
}
```

---

### POST /api/v1/posts/:id/tags

Add a tag to a post. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |
| Content-Type | application/json | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Request Body**

```json
{
  "tag_id": "uuid"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| tag_id | string (UUID) | Yes | Tag ID |

**Response**

HTTP 204 No Content (no body)

---

### DELETE /api/v1/posts/:id/tags/:tag_id

Remove a tag from a post. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |
| tag_id | string (UUID) | Yes | Tag ID |

**Response**

HTTP 204 No Content (no body)

---

## Tags

### GET /api/v1/tags

List tags with pagination.

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at`, `name` |
| order | string | No | - | Sort order: `asc`, `desc` |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "uuid",
        "name": "Tag Name",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_page": 5
  }
}
```

---

### GET /api/v1/tags/search

Search tags by name.

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| keyword | string | Yes | - | Search keyword |
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at`, `name` |
| order | string | No | - | Sort order: `asc`, `desc` |

**Response**

Same as GET /api/v1/tags.

---

### GET /api/v1/tags/:id

Get a single tag by ID.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Tag ID |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "Tag Name",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### POST /api/v1/tags

Create a new tag.

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Content-Type | application/json | Yes |

**Request Body**

```json
{
  "name": "Tag Name"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| name | string | Yes | 1-50 characters |

**Response**

HTTP 201 Created

```json
{
  "code": 0,
  "message": "created",
  "data": {
    "id": "uuid",
    "name": "Tag Name",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### PUT /api/v1/tags/:id

Update a tag.

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Content-Type | application/json | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Tag ID |

**Request Body**

```json
{
  "name": "Updated Tag Name"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| name | string | Yes | 1-50 characters |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "Updated Tag Name",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### DELETE /api/v1/tags/:id

Delete a tag.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Tag ID |

**Response**

HTTP 204 No Content (no body)

---

## Comments

### GET /api/v1/posts/:id/comments

List comments for a specific post.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at` |
| order | string | No | - | Sort order: `asc`, `desc` |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "uuid",
        "post_id": "",
        "author_name": "John Doe",
        "author_email": "john@example.com",
        "content": "Great post!",
        "status": "approved",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10,
    "total_page": 2
  }
}
```

---

### GET /api/v1/posts/:id/comments/count

Get comment count for a post.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Post ID |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "post_id": "uuid",
    "count": 42
  }
}
```

---

### POST /api/v1/comments

Create a new comment.

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Content-Type | application/json | Yes |

**Request Body**

```json
{
  "post_id": "uuid",
  "author_name": "John Doe",
  "author_email": "john@example.com",
  "content": "Great post!"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| post_id | string (UUID) | Yes | Post ID |
| author_name | string | Yes | 1-100 characters |
| author_email | string | Yes | Valid email format |
| content | string | Yes | 1-5000 characters |

**Response**

HTTP 201 Created

```json
{
  "code": 0,
  "message": "created",
  "data": {
    "id": "uuid",
    "post_id": "",
    "author_name": "John Doe",
    "author_email": "john@example.com",
    "content": "Great post!",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### GET /api/v1/comments/:id

Get a single comment by ID.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Comment ID |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "post_id": "",
    "author_name": "John Doe",
    "author_email": "john@example.com",
    "content": "Great post!",
    "status": "approved",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### PUT /api/v1/comments/:id

Update a comment.

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Content-Type | application/json | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Comment ID |

**Request Body**

```json
{
  "content": "Updated comment content"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| content | string | Yes | 1-5000 characters |

**Response**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "post_id": "",
    "author_name": "John Doe",
    "author_email": "john@example.com",
    "content": "Updated comment content",
    "status": "approved",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### DELETE /api/v1/comments/:id

Delete a comment.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Comment ID |

**Response**

HTTP 204 No Content (no body)

---

### GET /api/v1/comments

List all comments. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at`, `status` |
| order | string | No | - | Sort order: `asc`, `desc` |

**Response**

Same as GET /api/v1/posts/:id/comments.

---

### GET /api/v1/comments/status/:status

List comments by status. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| status | string | Yes | Comment status: `pending`, `approved`, `rejected`, `spam` |

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | int | No | 1 | Page number |
| size | int | No | 10 | Page size (max 100) |
| sort_by | string | No | - | Sort field: `created_at` |
| order | string | No | - | Sort order: `asc`, `desc` |

**Response**

Same as GET /api/v1/posts/:id/comments.

---

### POST /api/v1/comments/:id/approve

Approve a pending comment. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Comment ID |

**Response**

HTTP 204 No Content (no body)

---

### POST /api/v1/comments/:id/reject

Reject a pending comment. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Comment ID |

**Response**

HTTP 204 No Content (no body)

---

### POST /api/v1/comments/:id/spam

Mark a comment as spam. **[Admin]**

**Headers**

| Header | Value | Required |
|--------|-------|----------|
| Authorization | Bearer \<token\> | Yes |

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string (UUID) | Yes | Comment ID |

**Response**

HTTP 204 No Content (no body)

---

## Comment Status Values

| Status | Description |
|--------|-------------|
| pending | Awaiting moderation |
| approved | Visible on post |
| rejected | Hidden, rejected by admin |
| spam | Marked as spam |
