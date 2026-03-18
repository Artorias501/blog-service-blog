# Scenario: RequestLogger Middleware

- Given: A request is made to the API
- When: The request passes through RequestLogger middleware
- Then: Request details are logged with method, path, status, and latency

## Test Steps

- Case 1 (happy path): Successful GET request
- Case 2 (edge case): Request that returns error status
- Case 3 (edge case): Request with request ID header
- Case 4 (edge case): POST request with body
- Case 5 (edge case): Request with query parameters

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
