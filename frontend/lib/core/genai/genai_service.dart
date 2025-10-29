import 'genai_api.dart';
import 'models.dart';
import '../network/dio_client.dart';

class GenAiService {
  final GenAiApi _api;

  GenAiService(DioClient client) : _api = GenAiApi(client);

  /// Generate daily meal plan using real health data
  Future<GenAiResponse> generateDailyMealPlan({
    required String gender,
    required int age,
    required double heightCm,
    required double weightKg,
    required int stepsToday,
    required double activeKcal24h,
    required String target,
    List<String>? allergies,
    List<String>? preferences,
    String? userId,
  }) async {
    try {
      final response = await _api.generateDailyMeals(
        gender: gender,
        age: age,
        heightCm: heightCm,
        weightKg: weightKg,
        stepsToday: stepsToday,
        activeKcal24h: activeKcal24h,
        target: target,
        allergies: allergies,
        preferences: preferences,
        userId: userId,
      );
      return GenAiResponse.fromJson(response);
    } catch (e) {
      throw Exception('Failed to generate meal plan: $e');
    }
  }

  /// Generate weekly meal plan using real health data
  Future<GenAiResponse> generateWeeklyMealPlan({
    required String gender,
    required int age,
    required double heightCm,
    required double weightKg,
    required int stepsToday,
    required double activeKcal24h,
    required String target,
    List<String>? allergies,
    List<String>? preferences,
    String? userId,
  }) async {
    try {
      final response = await _api.generateWeeklyMeals(
        gender: gender,
        age: age,
        heightCm: heightCm,
        weightKg: weightKg,
        stepsToday: stepsToday,
        activeKcal24h: activeKcal24h,
        target: target,
        allergies: allergies,
        preferences: preferences,
        userId: userId,
      );
      return GenAiResponse.fromJson(response);
    } catch (e) {
      throw Exception('Failed to generate meal plan: $e');
    }
  }
}
