# Scenario: Summary Validation

- Given: A string value representing post summary
- When: Creating a new Summary value object
- Then: The system validates summary does not exceed 500 characters

## Test Steps

- Case 1 (happy path): Valid summary (<=500 chars) creates Summary successfully
- Case 2 (edge case): Empty summary is allowed (optional field)
- Case 3 (edge case): Summary with 501 characters returns validation error
- Case 4 (edge case): Summary with exactly 500 characters succeeds

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
