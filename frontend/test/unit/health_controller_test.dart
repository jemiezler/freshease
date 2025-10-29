// test/unit/health_controller_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/genai/genai_service.dart';
import 'package:frontend/core/genai/models.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import '../mocks/mock_repositories.dart';

void main() {
  group('HealthController', () {
    late HealthController healthController;
    late MockHealthRepository mockHealthRepository;
    late MockGenAiService mockGenAiService;

    setUp(() {
      mockHealthRepository = MockHealthRepository();
      mockGenAiService = MockGenAiService();
      healthController = HealthController(
        repository: mockHealthRepository,
        genAiService: mockGenAiService,
      );
    });

    tearDown(() {
      healthController.dispose();
    });

    group('User Data Loading', () {
      test('should load user data and calculate age correctly', () async {
        // Arrange
        final mockUser = MockData.mockUserProfile;
        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => mockUser);

        // Act
        await healthController.init();

        // Assert
        expect(healthController.userAge, 34); // 2024 - 1990 = 34
        expect(healthController.userGender, 'Male');
        expect(healthController.userGoal, 'weight_loss');
        expect(healthController.userHeightCm, 175.0);
        expect(healthController.userWeightKg, 70.0);
      });

      test('should map sex to gender correctly', () async {
        // Arrange
        final maleUser = MockData.mockUserProfile.copyWith(sex: 'male');
        final femaleUser = MockData.mockUserProfile.copyWith(sex: 'female');
        final otherUser = MockData.mockUserProfile.copyWith(sex: 'other');

        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => maleUser);

        // Act & Assert for male
        await healthController.init();
        expect(healthController.userGender, 'Male');

        // Reset and test female
        healthController.dispose();
        healthController = HealthController(
          repository: mockHealthRepository,
          genAiService: mockGenAiService,
        );
        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => femaleUser);
        await healthController.init();
        expect(healthController.userGender, 'Female');

        // Reset and test other
        healthController.dispose();
        healthController = HealthController(
          repository: mockHealthRepository,
          genAiService: mockGenAiService,
        );
        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => otherUser);
        await healthController.init();
        expect(healthController.userGender, 'Other');
      });
    });

    group('Meal Plan Generation', () {
      test('should generate daily meal plan successfully', () async {
        // Arrange
        final mockUser = MockData.mockUserProfile;
        final mockResponse = MockData.mockGenAiResponse;

        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => mockUser);
        when(
          mockGenAiService.generateDailyMealPlan(
            gender: anyNamed('gender'),
            age: anyNamed('age'),
            heightCm: anyNamed('heightCm'),
            weightKg: anyNamed('weightKg'),
            stepsToday: anyNamed('stepsToday'),
            activeKcal24h: anyNamed('activeKcal24h'),
            target: anyNamed('target'),
            allergies: anyNamed('allergies'),
            preferences: anyNamed('preferences'),
            userId: anyNamed('userId'),
          ),
        ).thenAnswer((_) async => mockResponse);

        await healthController.init();

        // Act
        await healthController.generateMealPlan();

        // Assert
        expect(healthController.currentMealPlan, isNotNull);
        expect(healthController.currentMealPlan!.plan.length, 1);
        expect(healthController.currentMealPlan!.plan.first.day, 'Monday');
        expect(healthController.mealPlanError, isNull);
        expect(healthController.isGeneratingMealPlan, false);
      });

      test('should generate weekly meal plan successfully', () async {
        // Arrange
        final mockUser = MockData.mockUserProfile;
        final mockResponse = MockData.mockGenAiResponse;

        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => mockUser);
        when(
          mockGenAiService.generateWeeklyMealPlan(
            gender: anyNamed('gender'),
            age: anyNamed('age'),
            heightCm: anyNamed('heightCm'),
            weightKg: anyNamed('weightKg'),
            stepsToday: anyNamed('stepsToday'),
            activeKcal24h: anyNamed('activeKcal24h'),
            target: anyNamed('target'),
            allergies: anyNamed('allergies'),
            preferences: anyNamed('preferences'),
            userId: anyNamed('userId'),
          ),
        ).thenAnswer((_) async => mockResponse);

        await healthController.init();

        // Act
        await healthController.generateWeeklyMealPlan();

        // Assert
        expect(healthController.currentWeeklyPlan, isNotNull);
        expect(healthController.currentWeeklyPlan!.plan.length, 1);
        expect(healthController.weeklyPlanError, isNull);
        expect(healthController.isGeneratingWeeklyPlan, false);
      });

      test('should handle meal plan generation error', () async {
        // Arrange
        final mockUser = MockData.mockUserProfile;

        when(
          mockHealthRepository.getCurrentUser(),
        ).thenAnswer((_) async => mockUser);
        when(
          mockGenAiService.generateDailyMealPlan(
            gender: anyNamed('gender'),
            age: anyNamed('age'),
            heightCm: anyNamed('heightCm'),
            weightKg: anyNamed('weightKg'),
            stepsToday: anyNamed('stepsToday'),
            activeKcal24h: anyNamed('activeKcal24h'),
            target: anyNamed('target'),
            allergies: anyNamed('allergies'),
            preferences: anyNamed('preferences'),
            userId: anyNamed('userId'),
          ),
        ).thenThrow(Exception('API Error'));

        await healthController.init();

        // Act
        await healthController.generateMealPlan();

        // Assert
        expect(healthController.currentMealPlan, isNull);
        expect(healthController.mealPlanError, 'Exception: API Error');
        expect(healthController.isGeneratingMealPlan, false);
      });
    });

    group('Cache Management', () {
      test('should check cache validity correctly', () {
        // Arrange
        healthController._lastGenerationTime = DateTime.now();

        // Act & Assert
        expect(healthController.isCacheValid, true);

        // Test expired cache
        healthController._lastGenerationTime = DateTime.now().subtract(
          const Duration(hours: 25),
        );
        expect(healthController.isCacheValid, false);
      });

      test('should clear meal plans correctly', () {
        // Arrange
        healthController.currentMealPlan = MockData.mockGenAiResponse;
        healthController.currentWeeklyPlan = MockData.mockGenAiResponse;
        healthController.mealPlanError = 'Test error';
        healthController.weeklyPlanError = 'Test error';

        // Act
        healthController.clearAllPlans();

        // Assert
        expect(healthController.currentMealPlan, isNull);
        expect(healthController.currentWeeklyPlan, isNull);
        expect(healthController.mealPlanError, isNull);
        expect(healthController.weeklyPlanError, isNull);
      });
    });

    group('Auto Generation', () {
      test(
        'should trigger auto generation when user has complete profile',
        () async {
          // Arrange
          final mockUser = MockData.mockUserProfile;
          final mockResponse = MockData.mockGenAiResponse;

          when(
            mockHealthRepository.getCurrentUser(),
          ).thenAnswer((_) async => mockUser);
          when(
            mockGenAiService.generateDailyMealPlan(
              gender: anyNamed('gender'),
              age: anyNamed('age'),
              heightCm: anyNamed('heightCm'),
              weightKg: anyNamed('weightKg'),
              stepsToday: anyNamed('stepsToday'),
              activeKcal24h: anyNamed('activeKcal24h'),
              target: anyNamed('target'),
              allergies: anyNamed('allergies'),
              preferences: anyNamed('preferences'),
              userId: anyNamed('userId'),
            ),
          ).thenAnswer((_) async => mockResponse);
          when(
            mockGenAiService.generateWeeklyMealPlan(
              gender: anyNamed('gender'),
              age: anyNamed('age'),
              heightCm: anyNamed('heightCm'),
              weightKg: anyNamed('weightKg'),
              stepsToday: anyNamed('stepsToday'),
              activeKcal24h: anyNamed('activeKcal24h'),
              target: anyNamed('target'),
              allergies: anyNamed('allergies'),
              preferences: anyNamed('preferences'),
              userId: anyNamed('userId'),
            ),
          ).thenAnswer((_) async => mockResponse);

          // Act
          await healthController.triggerAutoGeneration();

          // Assert
          expect(healthController.currentMealPlan, isNotNull);
          expect(healthController.currentWeeklyPlan, isNotNull);
        },
      );

      test(
        'should not auto generate when user profile is incomplete',
        () async {
          // Arrange
          final incompleteUser = MockData.mockUserProfile.copyWith(
            goal: null,
            heightCm: null,
            weightKg: null,
          );

          when(
            mockHealthRepository.getCurrentUser(),
          ).thenAnswer((_) async => incompleteUser);

          // Act
          await healthController.triggerAutoGeneration();

          // Assert
          expect(healthController.currentMealPlan, isNull);
          expect(healthController.currentWeeklyPlan, isNull);
        },
      );
    });
  });
}
