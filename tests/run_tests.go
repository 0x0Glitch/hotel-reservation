package main

// This file provides documentation for running tests in the Hotel Reservation Project
// It does not contain executable code.

/*
To run all tests in the project:
go test ./tests/...

To run tests for a specific package:
go test ./tests/types -v  (for types tests)
go test ./tests/db -v     (for database tests)
go test ./tests/api -v    (for API tests)
go test ./tests/middleware -v (for middleware tests)

To run a specific test:
go test ./tests/types -v -run TestCreateUserParamsValidate

To run tests in short mode (skipping integration tests):
go test ./tests/... -short

For test coverage:
go test ./tests/... -cover

To view detailed coverage:
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
*/