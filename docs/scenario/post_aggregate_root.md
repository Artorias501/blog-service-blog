# Scenario: Post Aggregate Root

- Given: Post entity data (title, content, summary, etc.)
- When: Creating or modifying a Post aggregate root
- Then: The system enforces invariants and provides business methods

## Test Steps

- Case 1 (happy path): Valid post data creates Post successfully
- Case 2 (edge case): Invalid title (empty) returns validation error
- Case 3 (edge case): Invalid title (>200 chars) returns validation error
- Case 4 (edge case): Invalid content (empty) returns validation error
- Case 5 (happy path): AddComment method adds comment to post
- Case 6 (happy path): AddTag method adds tag to post
- Case 7 (happy path): RemoveTag method removes tag from post
- Case 8 (happy path): Publish method sets published_at timestamp
- Case 9 (happy path): JSON serialization includes all fields with snake_case tags

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
