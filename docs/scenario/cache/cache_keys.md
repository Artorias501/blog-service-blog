# Scenario: Cache Key Generators

- Given: Entity IDs and pagination parameters
- When: Generating cache keys
- Then: Keys follow the naming convention exactly

## Test Steps

- Case 1 (happy path): Generate post keys (single, list, by tag)
- Case 2 (happy path): Generate tag keys (single, by name, list, by post)
- Case 3 (happy path): Generate comment keys (single, by post, count)
- Case 4 (edge case): Generate pattern keys for invalidation

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
