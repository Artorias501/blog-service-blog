# Scenario: TTL Configuration

- Given: Cache configuration with TTL values
- When: Setting cache entries
- Then: Entries expire after the configured TTL

## Test Steps

- Case 1 (happy path): Default TTL values are reasonable
- Case 2 (happy path): Custom TTL values can be configured
- Case 3 (happy path): Different entity types have different TTLs

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
