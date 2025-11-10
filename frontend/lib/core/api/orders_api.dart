// core/api/orders_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class OrdersApi {
  final Dio _dio;
  OrdersApi(DioClient client) : _dio = client.dio;

  /// Get all orders
  Future<List<Map<String, dynamic>>> listOrders() async {
    try {
      final response = await _dio.get('/api/orders');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch orders: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get order by ID
  Future<Map<String, dynamic>> getOrder(String orderId) async {
    try {
      final response = await _dio.get('/api/orders/$orderId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch order: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new order
  Future<Map<String, dynamic>> createOrder(Map<String, dynamic> orderData) async {
    try {
      final response = await _dio.post('/api/orders', data: orderData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create order: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update an order
  Future<Map<String, dynamic>> updateOrder(String orderId, Map<String, dynamic> orderData) async {
    try {
      final response = await _dio.patch('/api/orders/$orderId', data: orderData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update order: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete an order
  Future<void> deleteOrder(String orderId) async {
    try {
      final response = await _dio.delete('/api/orders/$orderId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete order: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

