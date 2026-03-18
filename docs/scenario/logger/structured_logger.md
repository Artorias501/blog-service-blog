# Scenario: Structured Logger

- Given: Application needs logging
- When: Logger is initialized
- Then: Logger outputs in JSON or text format based on environment

## Test Steps

- Case 1 (happy path): Create JSON logger for production
- Case 2 (edge case): Create text logger for development
- Case 3 (edge case): Log with structured fields
- Case 4 (edge case): Log at different levels (debug, info, warn, error)
- Case 5 (edge case): Invalid environment defaults to text

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
