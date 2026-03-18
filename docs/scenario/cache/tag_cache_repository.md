# Scenario: TagCacheRepository Implementation

- Given: A Tag entity and Redis client
- When: Performing cache operations (Get, Set, Delete, List operations, ByName, ByPostID)
- Then: Data is correctly stored and retrieved from Redis with proper TTL

## Test Steps

- Case 1 (happy path): Set and Get a tag
- Case 2 (happy path): Delete a tag
- Case 3 (happy path): Set and Get tag by name
- Case 4 (happy path): Set and Get tag list with pagination
- Case 5 (happy path): Set and Get tags by post ID
- Case 6 (happy path): Invalidate tag (pattern-based deletion)
- Case 7 (edge case): Get non-existent tag returns nil
- Case 8 (edge case): Get non-existent tag by name returns nil

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
