# Scenario: Configuration Loader

- Given: Application needs configuration
- When: Configuration is loaded
- Then: Values are read from environment with defaults applied

## Test Steps

- Case 1 (happy path): Load with all environment variables set
- Case 2 (edge case): Load with missing optional variables (use defaults)
- Case 3 (edge case): Load with invalid port number
- Case 4 (edge case): Validate required fields
- Case 5 (edge case): Load from environment with prefix

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
