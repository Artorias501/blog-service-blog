# Scenario: AuthorEmail Validation

- Given: A string value representing comment author email
- When: Creating a new AuthorEmail value object
- Then: The system validates the email format

## Test Steps

- Case 1 (happy path): Valid email format creates AuthorEmail successfully
- Case 2 (edge case): Empty string returns validation error
- Case 3 (edge case): Invalid email format (no @) returns validation error
- Case 4 (edge case): Invalid email format (no domain) returns validation error
- Case 5 (edge case): Email with special characters succeeds if valid

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
