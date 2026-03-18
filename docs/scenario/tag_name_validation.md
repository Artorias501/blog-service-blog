# Scenario: TagName Validation

- Given: A string value representing a tag name
- When: Creating a new TagName value object
- Then: The system validates length is between 1-50 characters

## Test Steps

- Case 1 (happy path): Valid tag name (1-50 chars) creates TagName successfully
- Case 2 (edge case): Empty string returns validation error
- Case 3 (edge case): Tag name with 51 characters returns validation error
- Case 4 (edge case): Tag name with exactly 50 characters succeeds
- Case 5 (edge case): Tag name with exactly 1 character succeeds

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
