import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/genai/genai_service.dart';
import 'package:frontend/core/genai/models.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';

import 'health_controller_test.mocks.dart';

@GenerateMocks([HealthRepository, GenAiService, UserRepository])
void main() {
  group('HealthController', () {
    late HealthController controller;
    late MockHealthRepository mockRepo;
    late MockGenAiService mockGenAiService;
    late MockUserRepository mockUserRepo;

    setUp(() {
      mockRepo = MockHealthRepository();
      mockGenAiService = MockGenAiService();
      controller = HealthController(
        repository: mockRepo,
        genAiService: mockGenAiService,
      );
    });

    test('initial state is idle', () {
      expect(controller.state, HealthState.idle);
      expect(controller.stepsToday, 0);
      expect(controller.kcalToday, 0.0);
      expect(controller.currentUser, isNull);
    });

    test('clearMealPlan clears daily plan', () {
      controller.currentMealPlan = GenAiResponse(
        stepsToday: 1000,
        activeKcal24h: 500.0,
        plan: [],
      );
      controller.mealPlanError = 'Some error';

      controller.clearMealPlan();

      expect(controller.currentMealPlan, isNull);
      expect(controller.mealPlanError, isNull);
    });

    test('clearWeeklyPlan clears weekly plan', () {
      controller.currentWeeklyPlan = GenAiResponse(
        stepsToday: 1000,
        activeKcal24h: 500.0,
        plan: [],
      );
      controller.weeklyPlanError = 'Some error';

      controller.clearWeeklyPlan();

      expect(controller.currentWeeklyPlan, isNull);
      expect(controller.weeklyPlanError, isNull);
    });

    test('clearAllPlans clears both plans', () {
      controller.currentMealPlan = GenAiResponse(
        stepsToday: 1000,
        activeKcal24h: 500.0,
        plan: [],
      );
      controller.currentWeeklyPlan = GenAiResponse(
        stepsToday: 2000,
        activeKcal24h: 1000.0,
        plan: [],
      );
      controller.mealPlanError = 'Error 1';
      controller.weeklyPlanError = 'Error 2';

      controller.clearAllPlans();

      expect(controller.currentMealPlan, isNull);
      expect(controller.currentWeeklyPlan, isNull);
      expect(controller.mealPlanError, isNull);
      expect(controller.weeklyPlanError, isNull);
    });

    test('isCacheValid returns false when no cache', () {
      expect(controller.isCacheValid, false);
    });

    test('cacheAge returns null when no cache', () {
      expect(controller.cacheAge, isNull);
    });

    test('generateMealPlan uses user data when available', () async {
      controller.userGender = 'Male';
      controller.userAge = 30;
      controller.userHeightCm = 175.0;
      controller.userWeightKg = 70.0;
      controller.userGoal = 'weight_loss';
      controller.stepsToday = 5000;
      controller.kcalToday = 2500.0;
      controller.currentUser = UserProfile(
        id: 'user-1',
        email: 'test@example.com',
        name: 'Test User',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final mockResponse = GenAiResponse(
        stepsToday: 5000,
        activeKcal24h: 2500.0,
        plan: [
          MealPlan(
            day: 'Monday',
            meals: {'breakfast': 'Oatmeal'},
            calories: {'breakfast': 300},
            totalCalories: 300,
          ),
        ],
      );

      when(mockGenAiService.generateDailyMealPlan(
        gender: 'Male',
        age: 30,
        heightCm: 175.0,
        weightKg: 70.0,
        stepsToday: 5000,
        activeKcal24h: 2500.0,
        target: 'weight_loss',
        allergies: anyNamed('allergies'),
        preferences: anyNamed('preferences'),
        userId: 'user-1',
      )).thenAnswer((_) async => mockResponse);

      await controller.generateMealPlan();

      expect(controller.currentMealPlan, isNotNull);
      expect(controller.currentMealPlan!.stepsToday, 5000);
      expect(controller.isGeneratingMealPlan, false);
      expect(controller.mealPlanError, isNull);
    });

    test('generateMealPlan uses provided parameters when user data missing', () async {
      final mockResponse = GenAiResponse(
        stepsToday: 3000,
        activeKcal24h: 2000.0,
        plan: [],
      );

      when(mockGenAiService.generateDailyMealPlan(
        gender: 'Female',
        age: 25,
        heightCm: 160.0,
        weightKg: 60.0,
        stepsToday: 0,
        activeKcal24h: 0.0,
        target: 'muscle_gain',
        allergies: anyNamed('allergies'),
        preferences: anyNamed('preferences'),
        userId: anyNamed('userId'),
      )).thenAnswer((_) async => mockResponse);

      await controller.generateMealPlan(
        gender: 'Female',
        age: 25,
        heightCm: 160.0,
        weightKg: 60.0,
        target: 'muscle_gain',
      );

      expect(controller.currentMealPlan, isNotNull);
      expect(controller.isGeneratingMealPlan, false);
    });

    test('generateMealPlan handles error when GenAI service unavailable', () async {
      final controllerWithoutGenAI = HealthController(
        repository: mockRepo,
        genAiService: null,
      );

      await controllerWithoutGenAI.generateMealPlan();

      expect(controllerWithoutGenAI.mealPlanError, 'GenAI service not available');
      expect(controllerWithoutGenAI.currentMealPlan, isNull);
    });

    test('generateMealPlan handles GenAI service errors', () async {
      when(mockGenAiService.generateDailyMealPlan(
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
      )).thenThrow(Exception('API Error'));

      await controller.generateMealPlan();

      expect(controller.mealPlanError, isNotNull);
      expect(controller.currentMealPlan, isNull);
      expect(controller.isGeneratingMealPlan, false);
    });

    test('generateWeeklyMealPlan uses user data when available', () async {
      controller.userGender = 'Male';
      controller.userAge = 30;
      controller.userHeightCm = 175.0;
      controller.userWeightKg = 70.0;
      controller.userGoal = 'weight_loss';
      controller.stepsToday = 5000;
      controller.kcalToday = 2500.0;

      final mockResponse = GenAiResponse(
        stepsToday: 5000,
        activeKcal24h: 2500.0,
        plan: [],
      );

      when(mockGenAiService.generateWeeklyMealPlan(
        gender: 'Male',
        age: 30,
        heightCm: 175.0,
        weightKg: 70.0,
        stepsToday: 5000,
        activeKcal24h: 2500.0,
        target: 'weight_loss',
        allergies: anyNamed('allergies'),
        preferences: anyNamed('preferences'),
        userId: anyNamed('userId'),
      )).thenAnswer((_) async => mockResponse);

      await controller.generateWeeklyMealPlan();

      expect(controller.currentWeeklyPlan, isNotNull);
      expect(controller.isGeneratingWeeklyPlan, false);
      expect(controller.weeklyPlanError, isNull);
    });

    test('generateWeeklyMealPlan handles error when GenAI service unavailable', () async {
      final controllerWithoutGenAI = HealthController(
        repository: mockRepo,
        genAiService: null,
      );

      await controllerWithoutGenAI.generateWeeklyMealPlan();

      expect(controllerWithoutGenAI.weeklyPlanError, 'GenAI service not available');
      expect(controllerWithoutGenAI.currentWeeklyPlan, isNull);
    });
  });
}

