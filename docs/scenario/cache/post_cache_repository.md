# Scenario: PostCacheRepository Implementation

- Given: A Post entity and Redis client
- When: Performing cache operations (Get, Set, Delete, List operations)
- Then: Data is correctly stored and retrieved from Redis with proper TTL

## Test Steps

- Case 1 (happy path): Set and Get a post
- Case 2 (happy path): Delete a post
- Case 3 (happy path): Set and Get post list with pagination
- Case 4 (happy path): Delete post list
- Case 5 (happy path): Set and Get posts by tag ID
- Case 6 (happy path): Invalidate post (pattern-based deletion)
- Case 7 (edge case): Get non-existent post returns nil
- Case 8 (edge case): Cache miss on list returns nil

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
