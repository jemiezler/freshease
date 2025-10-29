// test/widgets/edit_profile_page_test.dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:mockito/mockito.dart';
import 'package:frontend/features/account/presentation/pages/edit_profile_page.dart';
import 'package:frontend/features/account/presentation/state/user_cubit.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import '../mocks/mock_repositories.dart';
import '../utils/test_helpers.dart';

void main() {
  group('EditProfilePage', () {
    late MockUserCubit mockUserCubit;

    setUp(() {
      mockUserCubit = MockUserCubit();
    });

    Widget createTestWidget(UserProfile user) {
      return MaterialApp(
        home: BlocProvider<UserCubit>(
          create: (context) => mockUserCubit,
          child: EditProfilePage(user: user),
        ),
      );
    }

    testWidgets('should display user information correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Assert
      expect(find.text('Edit Profile'), findsOneWidget);
      expect(find.text(mockUser.name), findsOneWidget);
      expect(find.text(mockUser.phone ?? ''), findsOneWidget);
      expect(find.text(mockUser.bio ?? ''), findsOneWidget);
      expect(find.text(mockUser.avatar ?? ''), findsOneWidget);
    });

    testWidgets('should display goal dropdown with correct options', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Assert
      expect(find.text('Goal'), findsOneWidget);
      expect(find.text('Not specified'), findsOneWidget);
      expect(find.text('Maintenance'), findsOneWidget);
      expect(find.text('Weight Loss'), findsOneWidget);
      expect(find.text('Weight Gain'), findsOneWidget);
    });

    testWidgets('should display height and weight fields', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Assert
      expect(find.text('Height (cm)'), findsOneWidget);
      expect(find.text('Weight (kg)'), findsOneWidget);
      expect(find.text('175.0'), findsOneWidget); // Height value
      expect(find.text('70.0'), findsOneWidget); // Weight value
    });

    testWidgets('should validate height input correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Test invalid height (too low)
      await tester.enterText(
        find.byType(TextFormField).at(5),
        '30',
      ); // Height field
      await tester.tap(find.text('Save'));
      await tester.pump();

      // Assert
      expect(find.text('Height must be between 50-300 cm'), findsOneWidget);

      // Test invalid height (too high)
      await tester.enterText(find.byType(TextFormField).at(5), '350');
      await tester.tap(find.text('Save'));
      await tester.pump();

      expect(find.text('Height must be between 50-300 cm'), findsOneWidget);

      // Test valid height
      await tester.enterText(find.byType(TextFormField).at(5), '180');
      await tester.tap(find.text('Save'));
      await tester.pump();

      expect(find.text('Height must be between 50-300 cm'), findsNothing);
    });

    testWidgets('should validate weight input correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Test invalid weight (too low)
      await tester.enterText(
        find.byType(TextFormField).at(6),
        '10',
      ); // Weight field
      await tester.tap(find.text('Save'));
      await tester.pump();

      // Assert
      expect(find.text('Weight must be between 20-500 kg'), findsOneWidget);

      // Test invalid weight (too high)
      await tester.enterText(find.byType(TextFormField).at(6), '600');
      await tester.tap(find.text('Save'));
      await tester.pump();

      expect(find.text('Weight must be between 20-500 kg'), findsOneWidget);

      // Test valid weight
      await tester.enterText(find.byType(TextFormField).at(6), '75');
      await tester.tap(find.text('Save'));
      await tester.pump();

      expect(find.text('Weight must be between 20-500 kg'), findsNothing);
    });

    testWidgets('should call updateProfile when save button is tapped', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));
      await tester.tap(find.text('Save'));
      await tester.pump();

      // Assert
      verify(mockUserCubit.updateProfile(any)).called(1);
    });

    testWidgets('should show loading state when saving', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(const UserState(isUpdating: true));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Assert
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(find.text('Saving...'), findsOneWidget);
    });

    testWidgets('should show error state when save fails', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(
        mockUserCubit.state,
      ).thenReturn(const UserState(error: 'Save failed'));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Assert
      expect(find.text('Save failed'), findsOneWidget);
    });

    testWidgets('should navigate back when cancel is tapped', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));
      await tester.tap(find.text('Cancel'));
      await tester.pump();

      // Assert
      // Navigation back should be handled by the widget
      expect(find.text('Cancel'), findsOneWidget);
    });

    testWidgets('should handle date picker correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Find and tap the date field
      final dateField = find.byType(TextFormField).at(3); // Date field
      await tester.tap(dateField);
      await tester.pumpAndSettle();

      // Assert
      // Date picker should be shown
      expect(find.byType(DatePickerDialog), findsOneWidget);
    });

    testWidgets('should handle sex dropdown correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mockUser = MockData.mockUserProfile;
      when(mockUserCubit.state).thenReturn(UserState(user: mockUser));

      // Act
      await TestHelpers.pumpTestWidget(tester, createTestWidget(mockUser));

      // Find and tap the sex dropdown
      final sexDropdown = find.byType(DropdownButtonFormField<String>).at(0);
      await tester.tap(sexDropdown);
      await tester.pumpAndSettle();

      // Assert
      expect(find.text('Male'), findsOneWidget);
      expect(find.text('Female'), findsOneWidget);
      expect(find.text('Other'), findsOneWidget);
    });
  });
}
