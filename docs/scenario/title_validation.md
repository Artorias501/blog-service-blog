# Scenario: Title Validation

- Given: A string value representing a post title
- When: Creating a new Title value object
- Then: The system validates length is between 1-200 characters

## Test Steps

- Case 1 (happy path): Valid title (1-200 chars) creates Title successfully
- Case 2 (edge case): Empty string returns validation error
- Case 3 (edge case): Title with 201 characters returns validation error
- Case 4 (edge case): Title with exactly 200 characters succeeds
- Case 5 (edge case): Title with exactly 1 character succeeds

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
