# Frontend Test Setup Guide

## Initial Setup

### 1. Generate Mock Files
The test files use Mockito for mocking. You need to generate the mock files first:

```bash
cd frontend
flutter pub run build_runner build --delete-conflicting-outputs
```

This will generate the `.mocks.dart` files referenced in the test files.

### 2. Fix Model Field Mismatches
Some test files may have model field mismatches. You need to:

1. Check the actual model definitions in:
   - `lib/features/account/domain/entities/user_profile.dart`
   - `lib/features/auth/domain/entities/user.dart`
   - `lib/features/shop/data/models/shop_dtos.dart`

2. Update the test files to match the actual model fields.

### 3. Common Issues and Fixes

#### UserProfile Model
If `UserProfile` requires additional fields, update `user_cubit_test.dart`:

```dart
UserProfile(
  id: '1',
  email: 'test@example.com',
  name: 'Test User',
  status: 'active', // Add required fields
  createdAt: DateTime.now(),
  updatedAt: DateTime.now(),
)
```

#### ShopProductDTO Model
If `ShopProductDTO` has different fields, update `product_repository_test.dart`:

```dart
ShopProductDTO(
  id: '1',
  name: 'Test Product',
  price: 99.99,
  // Add all required fields based on actual model
)
```

## Running Tests

### Run All Tests
```bash
flutter test
```

### Run Specific Test File
```bash
flutter test test/unit/user_cubit_test.dart
```

### Run with Coverage
```bash
flutter test --coverage
```

### Using Test Script
```bash
./test/run_tests.sh
```

## Test Structure

- **Unit Tests** (`test/unit/`): Test business logic (Cubits, Repositories, Services)
- **Widget Tests** (`test/widgets/`): Test UI components
- **Integration Tests** (`integration_test/`): Test complete user flows
- **Test Utilities** (`test/utils/`): Helper functions and mock data

## Adding New Tests

### Unit Test Example
```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';

@GenerateMocks([YourRepository])
void main() {
  // Your tests here
}
```

### Widget Test Example
```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:your_app/your_widget.dart';

void main() {
  testWidgets('widget test', (tester) async {
    await tester.pumpWidget(YourWidget());
    // Your assertions
  });
}
```

## Troubleshooting

### Mock Generation Fails
- Ensure `build_runner` is in `dev_dependencies`
- Run `flutter pub get` first
- Delete `.dart_tool` and try again

### Model Field Errors
- Check actual model definitions
- Update test files to match model fields
- Use `required` keyword for required fields

### Import Errors
- Ensure all dependencies are in `pubspec.yaml`
- Run `flutter pub get`
- Check import paths are correct

## Next Steps

1. Generate mock files using `build_runner`
2. Fix model field mismatches in test files
3. Add more test coverage for remaining features
4. Add integration tests for key user flows
5. Set up CI/CD to run tests automatically

