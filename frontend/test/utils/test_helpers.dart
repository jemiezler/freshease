// test/utils/test_helpers.dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:get_it/get_it.dart';
import 'package:frontend/app/di.dart';

/// Test helper utilities for setting up test environment
class TestHelpers {
  static GetIt? _testGetIt;

  /// Set up test dependencies
  static void setUpTestDependencies() {
    _testGetIt = GetIt.instance;
    // Reset GetIt instance for tests
    if (_testGetIt!.isReady()) {
      _testGetIt!.reset();
    }
  }

  /// Clean up test dependencies
  static void tearDownTestDependencies() {
    if (_testGetIt != null && _testGetIt!.isReady()) {
      _testGetIt!.reset();
    }
  }

  /// Create a test widget with MaterialApp wrapper
  static Widget createTestWidget(Widget child) {
    return MaterialApp(home: child);
  }

  /// Pump widget with MaterialApp wrapper
  static Future<void> pumpTestWidget(WidgetTester tester, Widget widget) async {
    await tester.pumpWidget(MaterialApp(home: Scaffold(body: widget)));
  }

  /// Wait for async operations to complete
  static Future<void> waitForAsync(WidgetTester tester) async {
    await tester.pumpAndSettle();
  }

  /// Find widget by type and verify it exists
  static void expectWidget<T extends Widget>(WidgetTester tester) {
    expect(find.byType(T), findsOneWidget);
  }

  /// Find widget by type and verify it doesn't exist
  static void expectNoWidget<T extends Widget>(WidgetTester tester) {
    expect(find.byType(T), findsNothing);
  }

  /// Find text and verify it exists
  static void expectText(WidgetTester tester, String text) {
    expect(find.text(text), findsOneWidget);
  }

  /// Find text and verify it doesn't exist
  static void expectNoText(WidgetTester tester, String text) {
    expect(find.text(text), findsNothing);
  }

  /// Tap on a widget by type
  static Future<void> tapWidget<T extends Widget>(WidgetTester tester) async {
    await tester.tap(find.byType(T));
    await tester.pump();
  }

  /// Tap on text
  static Future<void> tapText(WidgetTester tester, String text) async {
    await tester.tap(find.text(text));
    await tester.pump();
  }

  /// Enter text in a text field
  static Future<void> enterText(
    WidgetTester tester,
    String text,
    String input,
  ) async {
    await tester.enterText(find.text(text), input);
    await tester.pump();
  }
}
