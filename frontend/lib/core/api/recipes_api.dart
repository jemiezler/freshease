// core/api/recipes_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class RecipesApi {
  final Dio _dio;
  RecipesApi(DioClient client) : _dio = client.dio;

  /// Get all recipes
  Future<List<Map<String, dynamic>>> listRecipes() async {
    try {
      final response = await _dio.get('/api/recipes');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch recipes: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get recipe by ID
  Future<Map<String, dynamic>> getRecipe(String recipeId) async {
    try {
      final response = await _dio.get('/api/recipes/$recipeId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch recipe: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new recipe
  Future<Map<String, dynamic>> createRecipe(Map<String, dynamic> recipeData) async {
    try {
      final response = await _dio.post('/api/recipes', data: recipeData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create recipe: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a recipe
  Future<Map<String, dynamic>> updateRecipe(String recipeId, Map<String, dynamic> recipeData) async {
    try {
      final response = await _dio.patch('/api/recipes/$recipeId', data: recipeData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update recipe: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a recipe
  Future<void> deleteRecipe(String recipeId) async {
    try {
      final response = await _dio.delete('/api/recipes/$recipeId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete recipe: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

