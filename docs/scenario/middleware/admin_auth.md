# Scenario: AdminAuth Middleware

- Given: A request is made to a protected endpoint
- When: The request passes through AdminAuth middleware
- Then: The request is either authorized or rejected with 401

## Test Steps

- Case 1 (happy path): Request with valid Bearer token in Authorization header
- Case 2 (edge case): Request with missing Authorization header
- Case 3 (edge case): Request with invalid token format (not Bearer)
- Case 4 (edge case): Request with wrong token value
- Case 5 (edge case): Request with Bearer prefix but no token
- Case 6 (edge case): Request with extra whitespace in header

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
