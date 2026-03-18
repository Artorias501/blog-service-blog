# Scenario: Timestamp Value Objects

- Given: A time.Time value
- When: Creating CreatedAt, UpdatedAt, or PublishedAt value objects
- Then: The system properly serializes/deserializes to/from JSON

## Test Steps

- Case 1 (happy path): CreatedAt with valid time serializes to JSON correctly
- Case 2 (happy path): UpdatedAt with valid time serializes to JSON correctly
- Case 3 (happy path): PublishedAt with valid time serializes to JSON correctly
- Case 4 (edge case): Nil/zero time handles correctly
- Case 5 (edge case): JSON deserialization produces correct time value

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [ ] Refactor implementation without breaking test
- [ ] Run test and confirm still passing after refactor
