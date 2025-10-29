// lib/core/health/health_controller.dart
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:health/health.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/health/util.dart'; // <- where you keep dataTypesAndroid/dataTypesIOS
import 'package:frontend/core/genai/genai_service.dart';
import 'package:frontend/core/genai/models.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';
import 'package:frontend/app/di.dart';

enum HealthState { idle, fetching, ready, noData, authDenied, error }

class HealthController extends ChangeNotifier {
  HealthController({HealthRepository? repository, GenAiService? genAiService})
    : _repo = repository ?? NoopHealthRepository(),
      _genAiService = genAiService;

  final Health _health = Health();
  final HealthRepository _repo;
  final GenAiService? _genAiService;

  HealthState state = HealthState.idle;
  HealthConnectSdkStatus? hcStatus;

  // ====== What the page wants ======
  int stepsToday = 0;
  double kcalToday = 0;

  // ====== User Data ======
  UserProfile? currentUser;
  int? userAge;
  String? userGender;
  String? userGoal;
  double? userHeightCm;
  double? userWeightKg;
  bool isLoadingUserData = false;
  String? userDataError;

  // ====== GenAI Meal Plan Data ======
  GenAiResponse? currentMealPlan;
  GenAiResponse? currentWeeklyPlan;
  bool isGeneratingMealPlan = false;
  bool isGeneratingWeeklyPlan = false;
  String? mealPlanError;
  String? weeklyPlanError;

  // ====== Auto-generation and Caching ======
  DateTime? _lastGenerationTime;
  static const Duration _cacheValidityDuration = Duration(hours: 24);

  // ====== What you want to save ======
  List<HealthDataPoint> allPoints = [];

  // Types: collect EVERYTHING supported for the platform
  List<HealthDataType> get _allTypes => Platform.isAndroid
      ? dataTypesAndroid // from your util.dart (same as in your sample)
      : (Platform.isIOS ? dataTypesIOS : const []);

  // READ/WRITE where allowed; READ-only for iOS-restricted types is handled below
  List<HealthDataAccess> get _permissions => _allTypes.map((type) {
    const iosReadOnly = <HealthDataType>{
      HealthDataType.GENDER,
      HealthDataType.BLOOD_TYPE,
      HealthDataType.BIRTH_DATE,
      HealthDataType.APPLE_MOVE_TIME,
      HealthDataType.APPLE_STAND_HOUR,
      HealthDataType.APPLE_STAND_TIME,
      HealthDataType.WALKING_HEART_RATE,
      HealthDataType.ELECTROCARDIOGRAM,
      HealthDataType.HIGH_HEART_RATE_EVENT,
      HealthDataType.LOW_HEART_RATE_EVENT,
      HealthDataType.IRREGULAR_HEART_RATE_EVENT,
      HealthDataType.EXERCISE_TIME,
    };
    if (Platform.isIOS && iosReadOnly.contains(type)) {
      return HealthDataAccess.READ;
    }
    return HealthDataAccess.READ_WRITE;
  }).toList();

  Future<void> init() async {
    _health.configure();
    if (Platform.isAndroid) {
      hcStatus = await _health.getHealthConnectSdkStatus();
    }
    await _loadUserData();
    notifyListeners();
  }

  /// Load current user data and calculate age
  Future<void> _loadUserData() async {
    isLoadingUserData = true;
    userDataError = null;
    notifyListeners();

    try {
      final userRepository = getIt<UserRepository>();
      currentUser = await userRepository.getCurrentUser();

      if (currentUser != null) {
        // Calculate age from date of birth
        if (currentUser!.dateOfBirth != null) {
          final now = DateTime.now();
          final birthDate = currentUser!.dateOfBirth!;
          userAge = now.year - birthDate.year;

          // Adjust if birthday hasn't occurred this year
          if (now.month < birthDate.month ||
              (now.month == birthDate.month && now.day < birthDate.day)) {
            userAge = userAge! - 1;
          }
        }

        // Map sex to gender for GenAI
        if (currentUser!.sex != null) {
          switch (currentUser!.sex!.toLowerCase()) {
            case 'male':
              userGender = 'Male';
              break;
            case 'female':
              userGender = 'Female';
              break;
            default:
              userGender = 'Other';
          }
        }

        // Load additional user data
        userGoal = currentUser!.goal;
        userHeightCm = currentUser!.heightCm;
        userWeightKg = currentUser!.weightKg;
      }
    } catch (e) {
      userDataError = 'Failed to load user data: $e';
    } finally {
      isLoadingUserData = false;
      notifyListeners();
    }
  }

  Future<void> authorize() async {
    try {
      if (Platform.isAndroid) {
        await Permission.activityRecognition.request();
        await Permission.location.request(); // some calories/workouts gate this
      }

      bool? has = await _health.hasPermissions(
        _allTypes,
        permissions: _permissions,
      );
      has = has == true;

      if (!has) {
        final ok = await _health.requestAuthorization(
          _allTypes,
          permissions: _permissions,
        );
        if (!ok) {
          state = HealthState.authDenied;
          notifyListeners();
          return;
        }
        // optional (Android HC)
        await _health.requestHealthDataHistoryAuthorization().catchError(
          (_) => false,
        );
        await _health.requestHealthDataInBackgroundAuthorization().catchError(
          (_) => false,
        );
      }
    } catch (_) {
      state = HealthState.authDenied;
      notifyListeners();
    }
  }

  /// Fetches ALL points for last 24h (saved) + computes only two KPIs for UI.
  Future<void> fetchAll24hAndComputeKpis() async {
    state = HealthState.fetching;
    notifyListeners();

    final now = DateTime.now();
    final start = now.subtract(const Duration(hours: 24));

    try {
      // 1) Fetch ALL points
      final fetched = await _health.getHealthDataFromTypes(
        types: _allTypes,
        startTime: start,
        endTime: now,
      );

      // dedup
      final dedup = _health.removeDuplicates(fetched);

      // keep in memory
      allPoints = List.unmodifiable(dedup);

      // 2) Persist (optional)
      await _repo.savePoints(allPoints);

      // 3) Compute KPIs for the page
      // Steps (use aggregate helper to match your logs/output)
      final int? s = await _health.getTotalStepsInInterval(
        DateTime(now.year, now.month, now.day),
        now,
        includeManualEntry: true,
      );
      stepsToday = s ?? 0;

      // Calories: sum TOTAL_CALORIES_BURNED points (kcal)
      double sumKcal = 0;
      for (final p in allPoints) {
        if (p.type == HealthDataType.TOTAL_CALORIES_BURNED) {
          final v = p.value;
          if (v is num) {
            final parsed = double.tryParse(v.toString());
            if (parsed != null) sumKcal += parsed;
          } else {
            final parsed = double.tryParse(v.toString());
            if (parsed != null) sumKcal += parsed;
          }
        }
      }
      kcalToday = sumKcal;

      state = allPoints.isEmpty ? HealthState.noData : HealthState.ready;
    } catch (e) {
      state = HealthState.error;
    } finally {
      notifyListeners();
    }
  }

  /// Generate meal plan using current health data and user profile
  Future<void> generateMealPlan({
    String? gender,
    int? age,
    double? heightCm,
    double? weightKg,
    String? target,
    List<String>? allergies,
    List<String>? preferences,
    String? userId,
  }) async {
    if (_genAiService == null) {
      mealPlanError = 'GenAI service not available';
      notifyListeners();
      return;
    }

    // Use user data if available, otherwise use provided parameters
    final finalGender = gender ?? userGender ?? 'Male';
    final finalAge = age ?? userAge ?? 30;
    final finalHeightCm = heightCm ?? userHeightCm ?? 175.0;
    final finalWeightKg = weightKg ?? userWeightKg ?? 70.0;
    final finalTarget = target ?? userGoal ?? 'maintenance';
    final finalUserId = userId ?? currentUser?.id;

    isGeneratingMealPlan = true;
    mealPlanError = null;
    notifyListeners();

    try {
      currentMealPlan = await _genAiService.generateDailyMealPlan(
        gender: finalGender,
        age: finalAge,
        heightCm: finalHeightCm,
        weightKg: finalWeightKg,
        stepsToday: stepsToday,
        activeKcal24h: kcalToday,
        target: finalTarget,
        allergies: allergies,
        preferences: preferences,
        userId: finalUserId,
      );
    } catch (e) {
      mealPlanError = e.toString();
    } finally {
      isGeneratingMealPlan = false;
      notifyListeners();
    }
  }

  /// Generate weekly meal plan using current health data and user profile
  Future<void> generateWeeklyMealPlan({
    String? gender,
    int? age,
    double? heightCm,
    double? weightKg,
    String? target,
    List<String>? allergies,
    List<String>? preferences,
    String? userId,
  }) async {
    if (_genAiService == null) {
      weeklyPlanError = 'GenAI service not available';
      notifyListeners();
      return;
    }

    // Use user data if available, otherwise use provided parameters
    final finalGender = gender ?? userGender ?? 'Male';
    final finalAge = age ?? userAge ?? 30;
    final finalHeightCm = heightCm ?? userHeightCm ?? 175.0;
    final finalWeightKg = weightKg ?? userWeightKg ?? 70.0;
    final finalTarget = target ?? userGoal ?? 'maintenance';
    final finalUserId = userId ?? currentUser?.id;

    isGeneratingWeeklyPlan = true;
    weeklyPlanError = null;
    notifyListeners();

    try {
      currentWeeklyPlan = await _genAiService.generateWeeklyMealPlan(
        gender: finalGender,
        age: finalAge,
        heightCm: finalHeightCm,
        weightKg: finalWeightKg,
        stepsToday: stepsToday,
        activeKcal24h: kcalToday,
        target: finalTarget,
        allergies: allergies,
        preferences: preferences,
        userId: finalUserId,
      );
    } catch (e) {
      weeklyPlanError = e.toString();
    } finally {
      isGeneratingWeeklyPlan = false;
      notifyListeners();
    }
  }

  /// Clear current meal plans
  void clearMealPlan() {
    currentMealPlan = null;
    mealPlanError = null;
    notifyListeners();
  }

  /// Clear current weekly meal plan
  void clearWeeklyPlan() {
    currentWeeklyPlan = null;
    weeklyPlanError = null;
    notifyListeners();
  }

  /// Clear all meal plans
  void clearAllPlans() {
    currentMealPlan = null;
    currentWeeklyPlan = null;
    mealPlanError = null;
    weeklyPlanError = null;
    notifyListeners();
  }

  /// Trigger auto-generation of meal plans (called from shop page)
  Future<void> triggerAutoGeneration() async {
    await _autoGenerateMealPlans();
  }

  /// Auto-generate meal plans if user has complete profile data and cache is expired
  Future<void> _autoGenerateMealPlans() async {
    // Check if we have complete user data
    if (userAge == null ||
        userGender == null ||
        userHeightCm == null ||
        userWeightKg == null ||
        userGoal == null) {
      return; // Don't auto-generate if profile is incomplete
    }

    // Check if cache is still valid
    if (_lastGenerationTime != null &&
        DateTime.now().difference(_lastGenerationTime!) <
            _cacheValidityDuration) {
      return; // Cache is still valid
    }

    // Check if we already have plans and they're recent
    if (currentMealPlan != null &&
        currentWeeklyPlan != null &&
        _lastGenerationTime != null &&
        DateTime.now().difference(_lastGenerationTime!) <
            _cacheValidityDuration) {
      return; // We have recent plans
    }

    // Auto-generate both daily and weekly plans
    _lastGenerationTime = DateTime.now();

    // Generate daily plan
    if (currentMealPlan == null) {
      await generateMealPlan();
    }

    // Generate weekly plan
    if (currentWeeklyPlan == null) {
      await generateWeeklyMealPlan();
    }
  }

  /// Check if cache is valid
  bool get isCacheValid {
    if (_lastGenerationTime == null) return false;
    return DateTime.now().difference(_lastGenerationTime!) <
        _cacheValidityDuration;
  }

  /// Get cache age
  Duration? get cacheAge {
    if (_lastGenerationTime == null) return null;
    return DateTime.now().difference(_lastGenerationTime!);
  }

  /// Force refresh all meal plans
  Future<void> refreshAllPlans() async {
    clearAllPlans();
    _lastGenerationTime = null;
    await _autoGenerateMealPlans();
  }
}
