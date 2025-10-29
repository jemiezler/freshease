# Freshease Backend Testing Guide

This document provides a comprehensive guide to the testing setup for the Freshease backend application.

## 📚 Complete Test Documentation

For comprehensive test documentation, please refer to:

- **[TEST_README.md](./TEST_README.md)** - Complete test suite documentation with coverage analysis
- **[TEST_CASES.md](./TEST_CASES.md)** - Detailed test case specification with 200+ individual test cases
- **[TEST_QUICK_REFERENCE.md](./TEST_QUICK_REFERENCE.md)** - Quick reference guide for daily test operations

## 🚀 Quick Start

```bash
# Run all tests with coverage
./run_tests.sh

# Run specific module tests
go test ./modules/users -v

# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 📊 Test Coverage Summary

| Module | Coverage | Status |
|--------|----------|--------|
| **Users** | 74.6% | ✅ PASS |
| **Products** | 55.7% | ✅ PASS |
| **Middleware** | 91.7% | ✅ PASS |
| **Integration** | 0.0% | ✅ PASS |
| **All Modules** | 43.5% | ✅ PASS |

**Total Test Files**: 40  
**Total Test Cases**: 200+  
**Modules Covered**: 13  

## Overview

The testing suite includes:
- **Unit Tests**: Individual component testing with mocks
- **Integration Tests**: End-to-end API testing with real database
- **Middleware Tests**: Authentication, validation, and logging middleware
- **Coverage Reports**: Detailed code coverage analysis

## Test Structure

```
backend/
├── modules/
│   ├── users/
│   │   ├── service_test.go      # Service layer unit tests
│   │   ├── controller_test.go    # Controller layer unit tests
│   │   └── repo_test.go         # Repository layer unit tests
│   └── products/
│       ├── service_test.go      # Service layer unit tests
│       └── controller_test.go   # Controller layer unit tests
├── internal/
│   ├── common/
│   │   ├── middleware/
│   │   │   └── middleware_test.go # Middleware unit tests
│   │   └── testutils/
│   │       ├── testutils.go      # Test utilities and helpers
│   │       └── integration_test.go # Integration tests
│   └── testutils/
│       └── testutils.go         # Common test utilities
├── run_tests.sh                 # Test runner script
└── coverage/                   # Coverage reports (generated)
```

## Running Tests

### Quick Start

```bash
# Run all tests with coverage
./run_tests.sh

# Run specific test suite
./run_tests.sh -s users
./run_tests.sh -s products
./run_tests.sh -s middleware
./run_tests.sh -s integration

# Run tests with race detection
./run_tests.sh -r

# Run benchmarks
./run_tests.sh -b

# Generate coverage report only
./run_tests.sh -c
```

### Manual Test Execution

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./modules/users/...
go test ./modules/products/...
go test ./internal/common/middleware/...

# Run integration tests
go test -tags=integration ./internal/common/testutils/...
```

## Test Categories

### 1. Unit Tests

Unit tests focus on testing individual components in isolation using mocks.

#### Service Layer Tests
- Test business logic
- Mock repository dependencies
- Verify error handling
- Test data transformation

#### Controller Layer Tests
- Test HTTP request/response handling
- Mock service dependencies
- Verify status codes
- Test parameter validation

#### Repository Layer Tests
- Test database operations
- Use in-memory SQLite for testing
- Verify CRUD operations
- Test error conditions

### 2. Integration Tests

Integration tests verify the complete API workflow:
- End-to-end API testing
- Real database operations
- Complete request/response cycles
- Error handling scenarios

### 3. Middleware Tests

Middleware tests cover:
- JWT authentication
- Request validation
- Logging functionality
- Error handling

## Test Utilities

### MockRepository
Generic mock repository for testing service layers:
```go
type MockRepository[T any] struct {
    ListFunc     func(ctx context.Context) ([]*T, error)
    FindByIDFunc func(ctx context.Context, id uuid.UUID) (*T, error)
    CreateFunc   func(ctx context.Context, dto interface{}) (*T, error)
    UpdateFunc   func(ctx context.Context, dto interface{}) (*T, error)
    DeleteFunc   func(ctx context.Context, id uuid.UUID) error
}
```

### TestContext
Helper for HTTP testing:
```go
type TestContext struct {
    t      *testing.T
    app    *fiber.App
    client *http.Client
}
```

### Helper Functions
- `CreateTestUserDTO()` - Generate valid user test data
- `CreateTestProductDTO()` - Generate valid product test data
- `AssertJSONResponse()` - Assert HTTP response format
- `AssertErrorResponse()` - Assert error response format

## Coverage Reports

### HTML Coverage Report
After running tests, view the coverage report:
```bash
open coverage/coverage.html
```

### Coverage Summary
```bash
go tool cover -func=coverage/coverage.out
```

### Coverage Thresholds
- **Minimum**: 80% overall coverage
- **Critical paths**: 95% coverage (auth, validation, core business logic)
- **New code**: 90% coverage requirement

## Test Data Management

### Test Database
- Uses in-memory SQLite for unit tests
- Isolated test data per test case
- Automatic cleanup after each test

### Test Data Creation
```go
// Create test user
userDTO := CreateTestUserDTO()
userDTO["email"] = "test@example.com"
userDTO["name"] = "Test User"

// Create test product
productDTO := CreateTestProductDTO()
productDTO["name"] = "Test Product"
productDTO["price"] = 99.99
```

## Best Practices

### 1. Test Naming
- Use descriptive test names
- Follow pattern: `TestFunction_Scenario_ExpectedResult`
- Example: `TestService_Create_WithValidData_ReturnsUser`

### 2. Test Structure
- Use table-driven tests for multiple scenarios
- Arrange-Act-Assert pattern
- Clear test setup and teardown

### 3. Mocking
- Mock external dependencies
- Use testify/mock for complex mocks
- Verify mock expectations

### 4. Assertions
- Use testify/assert for assertions
- Provide meaningful error messages
- Test both success and error cases

### 5. Test Data
- Use realistic test data
- Avoid hardcoded values
- Generate unique identifiers

## Continuous Integration

### GitHub Actions
Tests run automatically on:
- Pull requests
- Push to main branch
- Scheduled nightly runs

### Test Requirements
- All tests must pass
- Coverage must meet thresholds
- No race conditions
- Performance benchmarks within limits

## Troubleshooting

### Common Issues

#### Test Database Connection
```bash
# Ensure test database is properly configured
export TEST_DB_URL="sqlite3://:memory:"
```

#### Mock Expectations
```go
// Always call AssertExpectations
mockRepo.AssertExpectations(t)
```

#### Race Conditions
```bash
# Run with race detection
go test -race ./...
```

#### Coverage Issues
```bash
# Check coverage file exists
ls -la coverage/coverage.out

# Regenerate coverage
./run_tests.sh -c
```

### Debug Mode
```bash
# Run tests with debug output
go test -v -race -cover ./...
```

## Performance Testing

### Benchmarks
```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkUserService_Create ./modules/users/
```

### Load Testing
```bash
# Run load tests (if implemented)
go test -tags=load ./tests/load/
```

## Contributing

### Adding New Tests
1. Create test file: `*_test.go`
2. Follow naming conventions
3. Include both success and error cases
4. Add to appropriate test suite
5. Update coverage if needed

### Test Review Checklist
- [ ] Tests cover all code paths
- [ ] Error cases are tested
- [ ] Mocks are properly configured
- [ ] Test data is realistic
- [ ] Assertions are meaningful
- [ ] No race conditions
- [ ] Performance is acceptable

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Fiber Testing Guide](https://docs.gofiber.io/guide/testing)
- [Ent Testing Guide](https://entgo.io/docs/testing/)

## Support

For testing-related questions or issues:
1. Check this documentation
2. Review existing test examples
3. Consult the team
4. Create an issue if needed
