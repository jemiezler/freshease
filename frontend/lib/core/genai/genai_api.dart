import 'package:dio/dio.dart';
import '../network/dio_client.dart';

class GenAiApi {
  final Dio _dio;
  GenAiApi(DioClient client) : _dio = client.dio;

  /// Generate daily meal plan based on health data
  Future<Map<String, dynamic>> generateDailyMeals({
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
    final response = await _dio.post(
      '/api/genai/daily',
      data: {
        'user_id': userId,
        'gender': gender,
        'age': age,
        'height_cm': heightCm,
        'weight_kg': weightKg,
        'steps_today': stepsToday,
        'active_kcal_24h': activeKcal24h,
        'target': target,
        'allergies': allergies ?? [],
        'preferences': preferences ?? [],
      },
    );
    return response.data as Map<String, dynamic>;
  }

  /// Generate weekly meal plan based on health data
  Future<Map<String, dynamic>> generateWeeklyMeals({
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
    final response = await _dio.post(
      '/api/genai/weekly',
      data: {
        'user_id': userId,
        'gender': gender,
        'age': age,
        'height_cm': heightCm,
        'weight_kg': weightKg,
        'steps_today': stepsToday,
        'active_kcal_24h': activeKcal24h,
        'target': target,
        'allergies': allergies ?? [],
        'preferences': preferences ?? [],
      },
      options: Options(
        sendTimeout: const Duration(minutes: 2),
        receiveTimeout: const Duration(minutes: 2),
      ),
    );
    return response.data as Map<String, dynamic>;
  }
}
