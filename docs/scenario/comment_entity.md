# Scenario: Comment Entity

- Given: Comment entity data (author_name, author_email, content, status)
- When: Creating a Comment entity
- Then: The system validates all fields and enforces email format

## Test Steps

- Case 1 (happy path): Valid comment data creates Comment successfully
- Case 2 (edge case): Invalid author_name (empty) returns validation error
- Case 3 (edge case): Invalid author_email format returns validation error
- Case 4 (edge case): Invalid status returns validation error
- Case 5 (happy path): Approve method changes status to approved
- Case 6 (happy path): Reject method changes status to rejected
- Case 7 (happy path): MarkAsSpam method changes status to spam
- Case 8 (happy path): JSON serialization includes all fields with snake_case tags

## Status

- [x] Write scenario document document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
