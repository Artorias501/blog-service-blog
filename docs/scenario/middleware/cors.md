# Scenario: CORS Middleware

- Given: A request is made to the API
- When: The request passes through CORS middleware
- Then: Appropriate CORS headers are set

## Test Steps

- Case 1 (happy path): Request from allowed origin
- Case 2 (edge case): Preflight OPTIONS request
- Case 3 (edge case): Request from non-allowed origin
- Case 4 (edge case): Wildcard origin configuration
- Case 5 (edge case): Multiple allowed origins

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
