# Scenario: Tag Aggregate Root

- Given: Tag entity data (name)
- When: Creating a Tag aggregate root
- Then: The system validates and provides business methods

## Test Steps

- Case 1 (happy path): Valid tag name creates Tag successfully
- Case 2 (edge case): Invalid tag name (empty) returns validation error
- Case 3 (edge case): Invalid tag name (>50 chars) returns validation error
- Case 4 (happy path): JSON serialization includes all fields with snake_case tags

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
