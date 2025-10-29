# FreshEase Frontend Test Suite

## Overview
This document provides a comprehensive overview of the frontend test suite for the FreshEase application. The test suite includes unit tests, widget tests, and integration tests to ensure the application's reliability and functionality.

## Test Structure

```
test/
├── utils/
│   └── test_helpers.dart          # Test utility functions
├── mocks/
│   └── mock_repositories.dart     # Mock classes for testing
├── unit/
│   ├── health_controller_test.dart    # HealthController unit tests
│   ├── user_cubit_test.dart          # UserCubit unit tests
│   └── product_repository_test.dart   # ProductRepository unit tests
├── widgets/
│   ├── edit_profile_page_test.dart    # EditProfilePage widget tests
│   └── meal_plan_generator_test.dart  # MealPlanGenerator widget tests
├── integration_test/
│   └── app_test.dart              # Integration tests
├── test_config.dart               # Test configuration
├── run_tests.sh                   # Test runner script
└── widget_test.dart              # Basic widget tests
```

## Test Categories

### 1. Unit Tests
Unit tests focus on testing individual components in isolation:

#### HealthController Tests
- **User Data Loading**: Tests loading user profile data and age calculation
- **Meal Plan Generation**: Tests daily and weekly meal plan generation
- **Cache Management**: Tests meal plan caching and expiration
- **Auto Generation**: Tests automatic meal plan generation logic
- **Error Handling**: Tests error scenarios and recovery

#### UserCubit Tests
- **State Management**: Tests user state transitions
- **Profile Loading**: Tests user profile loading from repository
- **Profile Updates**: Tests user profile update functionality
- **Logout**: Tests user logout functionality
- **Error Handling**: Tests error states and recovery

#### ProductRepository Tests
- **Product Retrieval**: Tests getting products with pagination and filtering
- **Category Management**: Tests category retrieval
- **Vendor Management**: Tests vendor retrieval
- **Error Handling**: Tests network errors and edge cases

### 2. Widget Tests
Widget tests verify UI components work correctly:

#### EditProfilePage Tests
- **Form Display**: Tests form fields are displayed correctly
- **Validation**: Tests input validation for height, weight, and other fields
- **User Interaction**: Tests form submission and user interactions
- **State Management**: Tests loading and error states
- **Navigation**: Tests navigation and form handling

#### MealPlanGenerator Tests
- **User Data Display**: Tests display of user profile data
- **Generation Flow**: Tests meal plan generation process
- **Loading States**: Tests loading indicators and states
- **Error Handling**: Tests error display and recovery
- **Button Interactions**: Tests user interactions with buttons

### 3. Integration Tests
Integration tests verify the complete application flow:

#### App Flow Tests
- **Navigation**: Tests navigation between different tabs
- **Authentication**: Tests user authentication flow
- **Product Browsing**: Tests product browsing and search
- **Health Data**: Tests health data display and interaction
- **Profile Management**: Tests profile editing and updates

## Test Utilities

### TestHelpers
The `TestHelpers` class provides utility functions for:
- Setting up test dependencies
- Creating test widgets with MaterialApp wrapper
- Common test assertions
- Widget interaction helpers

### Mock Data
The `MockData` class provides:
- Sample user profiles for testing
- Sample products for testing
- Sample meal plans for testing
- Consistent test data across all tests

## Running Tests

### Individual Test Files
```bash
# Run specific test file
flutter test test/widget_test.dart

# Run all unit tests
flutter test test/unit/

# Run all widget tests
flutter test test/widgets/

# Run integration tests
flutter test integration_test/
```

### All Tests
```bash
# Run all tests
flutter test

# Run tests with coverage
flutter test --coverage
```

### Using Test Runner Script
```bash
# Make script executable
chmod +x test/run_tests.sh

# Run all tests with the script
./test/run_tests.sh
```

## Test Dependencies

The test suite uses the following packages:
- `flutter_test`: Core Flutter testing framework
- `mockito`: Mock object generation
- `bloc_test`: Testing BLoC/Cubit state management
- `integration_test`: Integration testing framework
- `fake_async`: Testing async operations
- `build_runner`: Code generation for mocks

## Test Coverage

The test suite aims to achieve:
- **Unit Tests**: 80%+ coverage of business logic
- **Widget Tests**: 70%+ coverage of UI components
- **Integration Tests**: Key user flows covered

## Best Practices

### Test Organization
- Group related tests using `group()` function
- Use descriptive test names
- Follow AAA pattern (Arrange, Act, Assert)

### Mock Usage
- Use mocks for external dependencies
- Verify mock interactions when necessary
- Keep mocks simple and focused

### Test Data
- Use consistent test data across tests
- Create reusable mock data classes
- Avoid hardcoded values in tests

### Assertions
- Use specific assertions (`findsOneWidget` vs `findsWidgets`)
- Test both positive and negative cases
- Verify state changes and side effects

## Current Status

### ✅ Completed
- Basic widget tests (working)
- Test structure and utilities
- Mock data classes
- Test configuration

### 🔄 In Progress
- Fixing compilation errors in complex tests
- Updating mock classes to match actual models
- Completing integration tests

### 📋 TODO
- Fix model field mismatches in mocks
- Complete UserCubit tests
- Add more comprehensive widget tests
- Implement end-to-end integration tests

## Troubleshooting

### Common Issues
1. **Mock Generation**: Run `flutter packages pub run build_runner build` to regenerate mocks
2. **Dependency Conflicts**: Check pubspec.yaml for version conflicts
3. **Import Errors**: Verify import paths match actual file locations
4. **Model Mismatches**: Update mock data to match actual model structures

### Debugging Tests
- Use `flutter test --verbose` for detailed output
- Add `debugPrint()` statements for debugging
- Use `tester.pumpAndSettle()` for async operations
- Check test logs for specific error messages

## Future Enhancements

1. **Performance Tests**: Add performance testing for heavy operations
2. **Accessibility Tests**: Add accessibility testing for UI components
3. **Golden Tests**: Add golden file tests for UI consistency
4. **CI/CD Integration**: Integrate tests with continuous integration
5. **Test Automation**: Automate test execution in development workflow

## Conclusion

The FreshEase frontend test suite provides comprehensive coverage of the application's functionality. While some tests are still being refined, the foundation is solid and follows Flutter testing best practices. The test suite will continue to evolve as the application grows and new features are added.
