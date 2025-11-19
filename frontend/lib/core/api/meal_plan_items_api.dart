// core/api/meal_plan_items_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class MealPlanItemsApi {
  final Dio _dio;
  MealPlanItemsApi(DioClient client) : _dio = client.dio;

  /// Get all meal plan items
  Future<List<Map<String, dynamic>>> listMealPlanItems() async {
    try {
      final response = await _dio.get('/api/meal_plan_items');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch meal plan items: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get meal plan item by ID
  Future<Map<String, dynamic>> getMealPlanItem(String mealPlanItemId) async {
    try {
      final response = await _dio.get('/api/meal_plan_items/$mealPlanItemId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch meal plan item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new meal plan item
  Future<Map<String, dynamic>> createMealPlanItem(Map<String, dynamic> mealPlanItemData) async {
    try {
      final response = await _dio.post('/api/meal_plan_items', data: mealPlanItemData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create meal plan item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a meal plan item
  Future<Map<String, dynamic>> updateMealPlanItem(String mealPlanItemId, Map<String, dynamic> mealPlanItemData) async {
    try {
      final response = await _dio.patch('/api/meal_plan_items/$mealPlanItemId', data: mealPlanItemData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update meal plan item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a meal plan item
  Future<void> deleteMealPlanItem(String mealPlanItemId) async {
    try {
      final response = await _dio.delete('/api/meal_plan_items/$mealPlanItemId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete meal plan item: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

