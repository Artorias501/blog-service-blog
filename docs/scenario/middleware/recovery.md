# Scenario: Recovery Middleware

- Given: A handler panics during request processing
- When: The panic is caught by Recovery middleware
- Then: Server returns 500 and logs the panic details

## Test Steps

- Case 1 (happy path): Handler panics with error message
- Case 2 (edge case): Handler panics with string
- Case 3 (edge case): Handler panics with nil
- Case 4 (edge case): Normal request without panic

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
