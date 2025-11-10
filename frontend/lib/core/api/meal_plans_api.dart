// core/api/meal_plans_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class MealPlansApi {
  final Dio _dio;
  MealPlansApi(DioClient client) : _dio = client.dio;

  /// Get all meal plans
  Future<List<Map<String, dynamic>>> listMealPlans() async {
    try {
      final response = await _dio.get('/api/meal_plans');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch meal plans: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get meal plan by ID
  Future<Map<String, dynamic>> getMealPlan(String mealPlanId) async {
    try {
      final response = await _dio.get('/api/meal_plans/$mealPlanId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch meal plan: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new meal plan
  Future<Map<String, dynamic>> createMealPlan(Map<String, dynamic> mealPlanData) async {
    try {
      final response = await _dio.post('/api/meal_plans', data: mealPlanData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create meal plan: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a meal plan
  Future<Map<String, dynamic>> updateMealPlan(String mealPlanId, Map<String, dynamic> mealPlanData) async {
    try {
      final response = await _dio.patch('/api/meal_plans/$mealPlanId', data: mealPlanData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update meal plan: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a meal plan
  Future<void> deleteMealPlan(String mealPlanId) async {
    try {
      final response = await _dio.delete('/api/meal_plans/$mealPlanId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete meal plan: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

