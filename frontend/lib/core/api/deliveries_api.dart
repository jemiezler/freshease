// core/api/deliveries_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class DeliveriesApi {
  final Dio _dio;
  DeliveriesApi(DioClient client) : _dio = client.dio;

  /// Get all deliveries
  Future<List<Map<String, dynamic>>> listDeliveries() async {
    try {
      final response = await _dio.get('/api/deliveries');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch deliveries: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get delivery by ID
  Future<Map<String, dynamic>> getDelivery(String deliveryId) async {
    try {
      final response = await _dio.get('/api/deliveries/$deliveryId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch delivery: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new delivery
  Future<Map<String, dynamic>> createDelivery(Map<String, dynamic> deliveryData) async {
    try {
      final response = await _dio.post('/api/deliveries', data: deliveryData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create delivery: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a delivery
  Future<Map<String, dynamic>> updateDelivery(String deliveryId, Map<String, dynamic> deliveryData) async {
    try {
      final response = await _dio.patch('/api/deliveries/$deliveryId', data: deliveryData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update delivery: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a delivery
  Future<void> deleteDelivery(String deliveryId) async {
    try {
      final response = await _dio.delete('/api/deliveries/$deliveryId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete delivery: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

