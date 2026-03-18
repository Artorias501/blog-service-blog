# Scenario: CommentID Validation

- Given: A string value representing a comment identifier
- When: Creating a new CommentID value object
- Then: The system validates if the value is a valid UUID

## Test Steps

- Case 1 (happy path): Valid UUID string creates CommentID successfully
- Case 2 (edge case): Empty string returns validation error
- Case 3 (edge case): Invalid UUID format returns validation error

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
