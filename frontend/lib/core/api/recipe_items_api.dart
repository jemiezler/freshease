// core/api/recipe_items_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class RecipeItemsApi {
  final Dio _dio;
  RecipeItemsApi(DioClient client) : _dio = client.dio;

  /// Get all recipe items
  Future<List<Map<String, dynamic>>> listRecipeItems() async {
    try {
      final response = await _dio.get('/api/recipe_items');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch recipe items: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get recipe item by ID
  Future<Map<String, dynamic>> getRecipeItem(String recipeItemId) async {
    try {
      final response = await _dio.get('/api/recipe_items/$recipeItemId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch recipe item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new recipe item
  Future<Map<String, dynamic>> createRecipeItem(Map<String, dynamic> recipeItemData) async {
    try {
      final response = await _dio.post('/api/recipe_items', data: recipeItemData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create recipe item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a recipe item
  Future<Map<String, dynamic>> updateRecipeItem(String recipeItemId, Map<String, dynamic> recipeItemData) async {
    try {
      final response = await _dio.patch('/api/recipe_items/$recipeItemId', data: recipeItemData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update recipe item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a recipe item
  Future<void> deleteRecipeItem(String recipeItemId) async {
    try {
      final response = await _dio.delete('/api/recipe_items/$recipeItemId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete recipe item: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

