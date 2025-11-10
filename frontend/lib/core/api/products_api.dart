// core/api/products_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class ProductsApi {
  final Dio _dio;
  ProductsApi(DioClient client) : _dio = client.dio;

  /// Get all products
  Future<List<Map<String, dynamic>>> listProducts() async {
    try {
      final response = await _dio.get('/api/products');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch products: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get product by ID
  Future<Map<String, dynamic>> getProduct(String productId) async {
    try {
      final response = await _dio.get('/api/products/$productId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch product: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new product
  Future<Map<String, dynamic>> createProduct(Map<String, dynamic> productData) async {
    try {
      final response = await _dio.post('/api/products', data: productData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create product: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a product
  Future<Map<String, dynamic>> updateProduct(String productId, Map<String, dynamic> productData) async {
    try {
      final response = await _dio.patch('/api/products/$productId', data: productData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update product: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a product
  Future<void> deleteProduct(String productId) async {
    try {
      final response = await _dio.delete('/api/products/$productId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete product: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

