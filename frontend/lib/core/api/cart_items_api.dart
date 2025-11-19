// core/api/cart_items_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class CartItemsApi {
  final Dio _dio;
  CartItemsApi(DioClient client) : _dio = client.dio;

  /// Get all cart items
  Future<List<Map<String, dynamic>>> listCartItems() async {
    try {
      final response = await _dio.get('/api/cart_items');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch cart items: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get cart item by ID
  Future<Map<String, dynamic>> getCartItem(String cartItemId) async {
    try {
      final response = await _dio.get('/api/cart_items/$cartItemId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch cart item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new cart item
  Future<Map<String, dynamic>> createCartItem(Map<String, dynamic> cartItemData) async {
    try {
      final response = await _dio.post('/api/cart_items', data: cartItemData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create cart item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a cart item
  Future<Map<String, dynamic>> updateCartItem(String cartItemId, Map<String, dynamic> cartItemData) async {
    try {
      final response = await _dio.patch('/api/cart_items/$cartItemId', data: cartItemData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update cart item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a cart item
  Future<void> deleteCartItem(String cartItemId) async {
    try {
      final response = await _dio.delete('/api/cart_items/$cartItemId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete cart item: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

