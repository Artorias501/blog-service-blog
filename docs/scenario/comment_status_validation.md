# Scenario: CommentStatus Validation

- Given: A string value representing comment status
- When: Creating a new CommentStatus value object
- Then: The system validates status is one of: pending, approved, rejected, spam

## Test Steps

- Case 1 (happy path): Valid status "pending" creates CommentStatus successfully
- Case 2 (happy path): Valid status "approved" creates CommentStatus successfully
- Case 3 (happy path): Valid status "rejected" creates CommentStatus successfully
- Case 4 (happy path): Valid status "spam" creates CommentStatus successfully
- Case 5 (edge case): Invalid status returns validation error
- Case 6 (edge case): Empty string returns validation error

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
