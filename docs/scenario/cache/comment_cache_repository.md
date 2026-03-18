# Scenario: CommentCacheRepository Implementation

- Given: A Comment entity and Redis client
- When: Performing cache operations (Get, Set, Delete, ListByPostID, CountByPostID)
- Then: Data is correctly stored and retrieved from Redis with proper TTL

## Test Steps

- Case 1 (happy path): Set and Get a comment
- Case 2 (happy path): Delete a comment
- Case 3 (happy path): Set and Get comments by post ID with pagination
- Case 4 (happy path): Delete comments by post ID
- Case 5 (happy path): Set and Get comment count by post ID
- Case 6 (happy path): Invalidate comment (pattern-based deletion)
- Case 7 (happy path): Invalidate by post ID
- Case 8 (edge case): Get non-existent comment returns nil
- Case 9 (edge case): Get count for non-existent post returns -1

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
