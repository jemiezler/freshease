// test/widgets/meal_plan_generator_test.dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/core/genai/widgets.dart';
import 'package:frontend/core/genai/models.dart';
import '../mocks/mock_repositories.dart';
import '../utils/test_helpers.dart';

void main() {
  group('MealPlanGenerator', () {
    testWidgets('should display user profile data correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      const userGender = 'Male';
      const userAge = 30;
      const userGoal = 'weight_loss';
      const userHeightCm = 175.0;
      const userWeightKg = 70.0;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(
          onGenerate: () {},
          isLoading: false,
          userGender: userGender,
          userAge: userAge,
          userGoal: userGoal,
          userHeightCm: userHeightCm,
          userWeightKg: userWeightKg,
          showUserData: true,
        ),
      );

      // Assert
      expect(find.text('Using your profile data:'), findsOneWidget);
      expect(find.text('Male'), findsOneWidget);
      expect(find.text('30 years'), findsOneWidget);
      expect(find.text('WEIGHT LOSS'), findsOneWidget);
      expect(find.text('175 cm'), findsOneWidget);
      expect(find.text('70.0 kg'), findsOneWidget);
    });

    testWidgets('should not display user data when showUserData is false', (
      WidgetTester tester,
    ) async {
      // Arrange
      const userGender = 'Male';
      const userAge = 30;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(
          onGenerate: () {},
          isLoading: false,
          userGender: userGender,
          userAge: userAge,
          showUserData: false,
        ),
      );

      // Assert
      expect(find.text('Using your profile data:'), findsNothing);
      expect(find.text('Male'), findsNothing);
      expect(find.text('30 years'), findsNothing);
    });

    testWidgets('should display generate button when not loading', (
      WidgetTester tester,
    ) async {
      // Arrange
      bool generateCalled = false;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(
          onGenerate: () => generateCalled = true,
          isLoading: false,
        ),
      );

      // Assert
      expect(find.text('Generate Daily Plan'), findsOneWidget);
      expect(find.byType(CircularProgressIndicator), findsNothing);
    });

    testWidgets('should show loading state when generating', (
      WidgetTester tester,
    ) async {
      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(onGenerate: () {}, isLoading: true),
      );

      // Assert
      expect(find.text('Generating...'), findsOneWidget);
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('should display error message when error occurs', (
      WidgetTester tester,
    ) async {
      // Arrange
      const errorMessage = 'Failed to generate meal plan';

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(
          onGenerate: () {},
          isLoading: false,
          error: errorMessage,
        ),
      );

      // Assert
      expect(find.text(errorMessage), findsOneWidget);
      expect(find.byIcon(Icons.error_outline), findsOneWidget);
    });

    testWidgets('should call onGenerate when generate button is tapped', (
      WidgetTester tester,
    ) async {
      // Arrange
      bool generateCalled = false;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(
          onGenerate: () => generateCalled = true,
          isLoading: false,
        ),
      );

      await tester.tap(find.text('Generate Daily Plan'));
      await tester.pump();

      // Assert
      expect(generateCalled, true);
    });

    testWidgets('should disable generate button when loading', (
      WidgetTester tester,
    ) async {
      // Arrange
      bool generateCalled = false;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanGenerator(
          onGenerate: () => generateCalled = true,
          isLoading: true,
        ),
      );

      await tester.tap(find.text('Generating...'));
      await tester.pump();

      // Assert
      expect(generateCalled, false);
    });
  });

  group('MealPlanCard', () {
    testWidgets('should display meal plan information correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mealPlan = MockData.mockGenAiResponse.plan.first;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanCard(mealPlan: mealPlan),
      );

      // Assert
      expect(find.text('Monday'), findsOneWidget);
      expect(find.text('Breakfast'), findsOneWidget);
      expect(find.text('Healthy breakfast'), findsOneWidget);
      expect(find.text('400 kcal'), findsOneWidget);
      expect(find.text('Lunch'), findsOneWidget);
      expect(find.text('Healthy lunch'), findsOneWidget);
      expect(find.text('600 kcal'), findsOneWidget);
      expect(find.text('Dinner'), findsOneWidget);
      expect(find.text('Healthy dinner'), findsOneWidget);
      expect(find.text('500 kcal'), findsOneWidget);
    });

    testWidgets('should display ingredients correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mealPlan = MockData.mockGenAiResponse.plan.first;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanCard(mealPlan: mealPlan),
      );

      // Assert
      expect(find.text('eggs'), findsOneWidget);
      expect(find.text('toast'), findsOneWidget);
      expect(find.text('fruit'), findsOneWidget);
      expect(find.text('chicken'), findsOneWidget);
      expect(find.text('rice'), findsOneWidget);
      expect(find.text('vegetables'), findsOneWidget);
    });

    testWidgets('should display macronutrients correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mealPlan = MockData.mockGenAiResponse.plan.first;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanCard(mealPlan: mealPlan),
      );

      // Assert
      expect(find.text('120.0g'), findsOneWidget); // Protein
      expect(find.text('150.0g'), findsOneWidget); // Carbs
      expect(find.text('60.0g'), findsOneWidget); // Fat
    });

    testWidgets('should display recommendations correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      final mealPlan = MockData.mockGenAiResponse.plan.first;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        MealPlanCard(mealPlan: mealPlan),
      );

      // Assert
      expect(find.text('Stay hydrated'), findsOneWidget);
      expect(find.text('Exercise regularly'), findsOneWidget);
    });
  });

  group('WeeklyPlanGenerator', () {
    testWidgets('should display user profile data correctly', (
      WidgetTester tester,
    ) async {
      // Arrange
      const userGender = 'Female';
      const userAge = 25;
      const userGoal = 'weight_gain';
      const userHeightCm = 165.0;
      const userWeightKg = 55.0;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        WeeklyPlanGenerator(
          onGenerate: () {},
          isLoading: false,
          userGender: userGender,
          userAge: userAge,
          userGoal: userGoal,
          userHeightCm: userHeightCm,
          userWeightKg: userWeightKg,
          showUserData: true,
        ),
      );

      // Assert
      expect(find.text('Using your profile data:'), findsOneWidget);
      expect(find.text('Female'), findsOneWidget);
      expect(find.text('25 years'), findsOneWidget);
      expect(find.text('WEIGHT GAIN'), findsOneWidget);
      expect(find.text('165 cm'), findsOneWidget);
      expect(find.text('55.0 kg'), findsOneWidget);
    });

    testWidgets('should display generate weekly plan button', (
      WidgetTester tester,
    ) async {
      // Arrange
      bool generateCalled = false;

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        WeeklyPlanGenerator(
          onGenerate: () => generateCalled = true,
          isLoading: false,
        ),
      );

      // Assert
      expect(find.text('Generate Weekly Plan'), findsOneWidget);
      expect(find.byIcon(Icons.calendar_view_week), findsOneWidget);
    });

    testWidgets('should show loading state when generating weekly plan', (
      WidgetTester tester,
    ) async {
      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        WeeklyPlanGenerator(onGenerate: () {}, isLoading: true),
      );

      // Assert
      expect(find.text('Generating...'), findsOneWidget);
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('should display error message for weekly plan', (
      WidgetTester tester,
    ) async {
      // Arrange
      const errorMessage = 'Failed to generate weekly plan';

      // Act
      await TestHelpers.pumpTestWidget(
        tester,
        WeeklyPlanGenerator(
          onGenerate: () {},
          isLoading: false,
          error: errorMessage,
        ),
      );

      // Assert
      expect(find.text(errorMessage), findsOneWidget);
      expect(find.byIcon(Icons.error_outline), findsOneWidget);
    });
  });
}
