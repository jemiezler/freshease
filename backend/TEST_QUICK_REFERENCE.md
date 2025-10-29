# 🚀 Freshease Backend - Test Quick Reference Guide

## Quick Commands

### Run All Tests
```bash
# Complete test suite with coverage
./run_tests.sh

# All tests with verbose output
go test ./... -v

# All tests with coverage
go test ./... -cover
```

### Run Module-Specific Tests
```bash
# Users module
go test ./modules/users -v

# Products module
go test ./modules/products -v

# Authentication module
go test ./modules/auth/authoidc -v

# GenAI module
go test ./modules/genai -v

# Cart items module
go test ./modules/cart_items -v
```

### Run Test Categories
```bash
# Unit tests only
go test ./modules/... -v

# Integration tests only
go test ./internal/common/testutils -v

# Middleware tests only
go test ./internal/common/middleware -v
```

### Coverage Analysis
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open coverage report
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

## Test File Locations

```
backend/
├── TEST_README.md              # Complete test documentation
├── TEST_CASES.md               # Detailed test case specification
├── run_tests.sh                # Test execution script
├── internal/common/
│   ├── middleware/middleware_test.go
│   └── testutils/integration_test.go
└── modules/
    ├── addresses/
    │   ├── controller_test.go
    │   ├── repo_test.go
    │   └── service_test.go
    ├── auth/authoidc/
    │   ├── controller_test.go
    │   └── service_test.go
    ├── cart_items/
    │   ├── controller_test.go
    │   ├── repo_test.go
    │   └── service_test.go
    ├── genai/
    │   ├── controller_test.go
    │   ├── repo_test.go
    │   └── service_test.go
    ├── users/
    │   ├── controller_test.go
    │   ├── repo_test.go
    │   └── service_test.go
    ├── products/
    │   ├── controller_test.go
    │   ├── repo_test.go
    │   └── service_test.go
    └── [8 more modules with same structure]
```

## Test Coverage by Module

| Module | Coverage | Test Files | Status |
|--------|----------|------------|--------|
| **Middleware** | 91.7% | 1 | ✅ |
| **Carts** | 82.0% | 3 | ✅ |
| **Inventories** | 80.7% | 3 | ✅ |
| **Permissions** | 81.2% | 3 | ✅ |
| **Product Categories** | 81.6% | 3 | ✅ |
| **Roles** | 81.2% | 3 | ✅ |
| **Shop** | 77.7% | 3 | ✅ |
| **Vendors** | 75.7% | 3 | ✅ |
| **Users** | 74.6% | 3 | ✅ |
| **Addresses** | 55.0% | 3 | ✅ |
| **Products** | 55.7% | 3 | ✅ |
| **Cart Items** | 41.0% | 3 | ✅ |
| **GenAI** | 41.2% | 3 | ✅ |
| **Auth/OIDC** | 26.8% | 2 | ✅ |

## Common Test Patterns

### Controller Test Structure
```go
func TestController_Method(t *testing.T) {
    tests := []struct {
        name           string
        input          InputType
        mockSetup      func(*MockService)
        expectedStatus int
        expectedError  bool
    }{
        {
            name: "success - normal operation",
            mockSetup: func(mockSvc *MockService) {
                mockSvc.On("Method", mock.Anything).Return(expectedResult, nil)
            },
            expectedStatus: http.StatusOK,
            expectedError:  false,
        },
        {
            name: "error - service failure",
            mockSetup: func(mockSvc *MockService) {
                mockSvc.On("Method", mock.Anything).Return(nil, errors.New("error"))
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Service Test Structure
```go
func TestService_Method(t *testing.T) {
    tests := []struct {
        name           string
        mockSetup      func(*MockRepository)
        expectedResult ResultType
        expectedError  error
    }{
        {
            name: "success - normal operation",
            mockSetup: func(mockRepo *MockRepository) {
                mockRepo.On("Method", mock.Anything).Return(expectedResult, nil)
            },
            expectedResult: expectedResult,
            expectedError:  nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Test Data Helpers

### Common Helper Functions
```go
// String pointer helper
func stringPtr(s string) *string {
    return &s
}

// UUID pointer helper
func uuidPtr(u uuid.UUID) *uuid.UUID {
    return &u
}

// Create test user DTO
func CreateTestUserDTO() CreateUserDTO {
    return CreateUserDTO{
        ID:    uuid.New(),
        Name:  "Test User",
        Email: "test@example.com",
    }
}
```

## Debugging Tests

### Run Single Test
```bash
# Run specific test function
go test ./modules/users -run TestController_GetUser -v

# Run test with specific pattern
go test ./modules/users -run "TestController.*Get" -v
```

### Test with Race Detection
```bash
# Run tests with race detection
go test ./... -race

# Run specific module with race detection
go test ./modules/users -race -v
```

### Test with Coverage and Race Detection
```bash
# Full test suite with coverage and race detection
go test ./... -cover -race -v
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go test ./... -v -cover
      - run: go test ./... -race
```

### Pre-commit Hook
```bash
#!/bin/sh
# Run tests before commit
go test ./... -v
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

## Troubleshooting

### Common Issues

1. **Test Failures**
   ```bash
   # Run with verbose output to see details
   go test ./modules/users -v
   ```

2. **Coverage Issues**
   ```bash
   # Check coverage for specific module
   go test ./modules/users -cover
   ```

3. **Race Conditions**
   ```bash
   # Run with race detection
   go test ./... -race
   ```

4. **Integration Test Failures**
   ```bash
   # Check database connection
   go test ./internal/common/testutils -v
   ```

### Test Environment Setup
```bash
# Set test environment variables
export DB_URL="sqlite://:memory:"
export JWT_SECRET="test-secret"
export GOOGLE_API_KEY="test-key"

# Run tests
go test ./... -v
```

## Performance Testing

### Benchmark Tests
```bash
# Run benchmark tests
go test ./... -bench=.

# Run specific benchmark
go test ./modules/users -bench=BenchmarkController_GetUser
```

### Memory Profiling
```bash
# Run with memory profiling
go test ./... -memprofile=mem.prof

# Analyze memory profile
go tool pprof mem.prof
```

---

**Quick Reference**: Keep this guide handy for daily test operations!  
**Full Documentation**: See `TEST_README.md` for complete details  
**Test Cases**: See `TEST_CASES.md` for detailed test specifications
