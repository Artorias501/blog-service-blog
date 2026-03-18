# Scenario: Content Validation

- Given: A string value representing post content
- When: Creating a new Content value object
- Then: The system validates content is not empty

## Test Steps

- Case 1 (happy path): Non-empty content creates Content successfully
- Case 2 (edge case): Empty string returns validation error
- Case 3 (edge case): Content with whitespace only returns validation error

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
