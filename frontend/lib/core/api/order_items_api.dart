// core/api/order_items_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class OrderItemsApi {
  final Dio _dio;
  OrderItemsApi(DioClient client) : _dio = client.dio;

  /// Get all order items
  Future<List<Map<String, dynamic>>> listOrderItems() async {
    try {
      final response = await _dio.get('/api/order_items');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch order items: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get order item by ID
  Future<Map<String, dynamic>> getOrderItem(String orderItemId) async {
    try {
      final response = await _dio.get('/api/order_items/$orderItemId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch order item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new order item
  Future<Map<String, dynamic>> createOrderItem(Map<String, dynamic> orderItemData) async {
    try {
      final response = await _dio.post('/api/order_items', data: orderItemData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create order item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update an order item
  Future<Map<String, dynamic>> updateOrderItem(String orderItemId, Map<String, dynamic> orderItemData) async {
    try {
      final response = await _dio.patch('/api/order_items/$orderItemId', data: orderItemData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update order item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete an order item
  Future<void> deleteOrderItem(String orderItemId) async {
    try {
      final response = await _dio.delete('/api/order_items/$orderItemId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete order item: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

