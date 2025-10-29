// test/test_config.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:get_it/get_it.dart';
import 'package:frontend/app/di.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';

/// Test configuration for setting up test environment
class TestConfig {
  static bool _isInitialized = false;

  /// Initialize test configuration
  static void setUp() {
    if (_isInitialized) return;

    // Set up test environment
    TestWidgetsFlutterBinding.ensureInitialized();

    // Initialize dependency injection for tests
    configureDependencies();

    _isInitialized = true;
  }

  /// Clean up test configuration
  static void tearDown() {
    if (!_isInitialized) return;

    // Reset GetIt instance
    if (GetIt.instance.isReady()) {
      GetIt.instance.resetLazySingleton<UserRepository>();
    }

    _isInitialized = false;
  }

  /// Check if test configuration is initialized
  static bool get isInitialized => _isInitialized;
}
