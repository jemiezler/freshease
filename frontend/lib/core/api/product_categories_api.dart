// core/api/product_categories_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class ProductCategoriesApi {
  final Dio _dio;
  ProductCategoriesApi(DioClient client) : _dio = client.dio;

  /// Get all product categories
  Future<List<Map<String, dynamic>>> listProductCategories() async {
    try {
      final response = await _dio.get('/api/product_categories');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch product categories: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get product category by ID
  Future<Map<String, dynamic>> getProductCategory(String productCategoryId) async {
    try {
      final response = await _dio.get('/api/product_categories/$productCategoryId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch product category: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new product category
  Future<Map<String, dynamic>> createProductCategory(Map<String, dynamic> productCategoryData) async {
    try {
      final response = await _dio.post('/api/product_categories', data: productCategoryData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create product category: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a product category
  Future<Map<String, dynamic>> updateProductCategory(String productCategoryId, Map<String, dynamic> productCategoryData) async {
    try {
      final response = await _dio.patch('/api/product_categories/$productCategoryId', data: productCategoryData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update product category: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a product category
  Future<void> deleteProductCategory(String productCategoryId) async {
    try {
      final response = await _dio.delete('/api/product_categories/$productCategoryId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete product category: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

