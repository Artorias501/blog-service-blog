# Scenario: Redis Connection Setup

- Given: Redis configuration parameters (host, port, password, db)
- When: Creating a new Redis client connection
- Then: Client is successfully created and can ping Redis server

## Test Steps

- Case 1 (happy path): Connect to Redis with valid configuration
- Case 2 (edge case): Connection fails when Redis is unavailable - returns error, not panic
- Case 3 (edge case): Connection with custom TTL configuration

## Status

- [x] Write scenario document
- [x] Write solid test according to document
- [x] Run test and watch it failing
- [x] Implement to make test pass
- [x] Run test and confirm it passed
- [x] Refactor implementation without breaking test
- [x] Run test and confirm still passing after refactor
