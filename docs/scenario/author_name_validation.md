# Scenario: AuthorName Validation

- Given: A string value representing comment author name
- When: Creating a new AuthorName value object
- Then: The system validates length is between 1-100 characters

## Test Steps

- Case 1 (happy path): Valid author name (1-100 chars) creates AuthorName successfully
- Case 2 (edge case): Empty string returns validation error
- Case 3 (edge case): Author name with 101 characters returns validation error
- Case 4 (edge case): Author name with exactly 100 characters succeeds

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
