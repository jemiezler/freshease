import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

/// Helper class for test utilities
class TestHelpers {
  /// Create a MaterialApp wrapper for widget tests
  static Widget createTestApp({
    required Widget child,
    List<NavigatorObserver>? navigatorObservers,
  }) {
    return MaterialApp(
      home: child,
      navigatorObservers: navigatorObservers ?? [],
    );
  }

  /// Wait for async operations to complete
  static Future<void> waitForAsync(WidgetTester tester) async {
    await tester.pump();
    await tester.pump(const Duration(seconds: 1));
  }

  /// Tap a widget and wait for animations
  static Future<void> tapAndWait(
    WidgetTester tester,
    Finder finder,
  ) async {
    await tester.tap(finder);
    await tester.pumpAndSettle();
  }

  /// Enter text in a text field
  static Future<void> enterText(
    WidgetTester tester,
    Finder finder,
    String text,
  ) async {
    await tester.enterText(finder, text);
    await tester.pump();
  }

  /// Scroll until a widget is visible
  static Future<void> scrollUntilVisible(
    WidgetTester tester,
    Finder finder,
    Finder scrollable,
    double delta,
  ) async {
    await tester.scrollUntilVisible(
      finder,
      delta,
      scrollable: scrollable,
    );
    await tester.pumpAndSettle();
  }
}

/// Mock data for testing
class MockData {
  static const String testEmail = 'test@example.com';
  static const String testPassword = 'testpassword123';
  static const String testName = 'Test User';
  static const String testUserId = '123';
}

