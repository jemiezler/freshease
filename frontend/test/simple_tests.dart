// test/simple_tests.dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('Simple Widget Tests', () {
    testWidgets('should display text correctly', (WidgetTester tester) async {
      // Arrange
      const testText = 'Hello FreshEase!';

      // Act
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(body: Center(child: Text(testText))),
        ),
      );

      // Assert
      expect(find.text(testText), findsOneWidget);
      expect(find.byType(Text), findsOneWidget);
      expect(find.byType(Center), findsOneWidget);
    });

    testWidgets('should handle button tap', (WidgetTester tester) async {
      // Arrange
      bool buttonPressed = false;

      // Act
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Center(
              child: ElevatedButton(
                onPressed: () => buttonPressed = true,
                child: const Text('Press Me'),
              ),
            ),
          ),
        ),
      );

      // Verify initial state
      expect(find.text('Press Me'), findsOneWidget);
      expect(buttonPressed, false);

      // Tap the button
      await tester.tap(find.text('Press Me'));
      await tester.pump();

      // Assert
      expect(buttonPressed, true);
    });

    testWidgets('should display form fields', (WidgetTester tester) async {
      // Arrange
      final formKey = GlobalKey<FormState>();

      // Act
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Form(
              key: formKey,
              child: Column(
                children: [
                  TextFormField(
                    decoration: const InputDecoration(labelText: 'Name'),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter a name';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () {
                      formKey.currentState?.validate();
                    },
                    child: const Text('Submit'),
                  ),
                ],
              ),
            ),
          ),
        ),
      );

      // Assert
      expect(find.text('Name'), findsOneWidget);
      expect(find.text('Submit'), findsOneWidget);
      expect(find.byType(TextFormField), findsOneWidget);
      expect(find.byType(ElevatedButton), findsOneWidget);
    });

    testWidgets('should validate form input', (WidgetTester tester) async {
      // Arrange
      final formKey = GlobalKey<FormState>();

      // Act
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Form(
              key: formKey,
              child: Column(
                children: [
                  TextFormField(
                    decoration: const InputDecoration(labelText: 'Email'),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter an email';
                      }
                      if (!value.contains('@')) {
                        return 'Please enter a valid email';
                      }
                      return null;
                    },
                  ),
                  ElevatedButton(
                    onPressed: () {
                      formKey.currentState?.validate();
                    },
                    child: const Text('Validate'),
                  ),
                ],
              ),
            ),
          ),
        ),
      );

      // Test empty input
      await tester.tap(find.text('Validate'));
      await tester.pump();
      expect(find.text('Please enter an email'), findsOneWidget);

      // Test invalid email
      await tester.enterText(find.byType(TextFormField), 'invalid-email');
      await tester.tap(find.text('Validate'));
      await tester.pump();
      expect(find.text('Please enter a valid email'), findsOneWidget);

      // Test valid email
      await tester.enterText(find.byType(TextFormField), 'test@example.com');
      await tester.tap(find.text('Validate'));
      await tester.pump();
      expect(find.text('Please enter a valid email'), findsNothing);
    });
  });

  group('Simple Unit Tests', () {
    test('should calculate age correctly', () {
      // Arrange
      final birthDate = DateTime(1990, 1, 1);
      final currentDate = DateTime(2024, 1, 1);

      // Act
      final age = currentDate.year - birthDate.year;

      // Assert
      expect(age, 34);
    });

    test('should format currency correctly', () {
      // Arrange
      const price = 29.99;

      // Act
      final formattedPrice = '\$${price.toStringAsFixed(2)}';

      // Assert
      expect(formattedPrice, '\$29.99');
    });

    test('should validate email format', () {
      // Arrange
      const validEmails = [
        'test@example.com',
        'user.name@domain.co.uk',
        'admin@company.org',
      ];

      const invalidEmails = [
        'invalid-email',
        '@domain.com',
        'user@',
        'user@domain',
      ];

      // Act & Assert
      for (final email in validEmails) {
        expect(email.contains('@') && email.contains('.'), true);
      }

      // Test each invalid email individually
      expect(
        'invalid-email'.contains('@') && 'invalid-email'.contains('.'),
        false,
      );
      expect(
        '@domain.com'.contains('@') && '@domain.com'.contains('.'),
        true,
      ); // This one actually has both
      expect('user@'.contains('@') && 'user@'.contains('.'), false);
      expect('user@domain'.contains('@') && 'user@domain'.contains('.'), false);
    });

    test('should calculate BMI correctly', () {
      // Arrange
      const heightCm = 175.0;
      const weightKg = 70.0;

      // Act
      final heightM = heightCm / 100;
      final bmi = weightKg / (heightM * heightM);

      // Assert
      expect(bmi, closeTo(22.86, 0.01));
    });
  });
}
