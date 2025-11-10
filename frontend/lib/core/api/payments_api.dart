// core/api/payments_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class PaymentsApi {
  final Dio _dio;
  PaymentsApi(DioClient client) : _dio = client.dio;

  /// Get all payments
  Future<List<Map<String, dynamic>>> listPayments() async {
    try {
      final response = await _dio.get('/api/payments');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch payments: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get payment by ID
  Future<Map<String, dynamic>> getPayment(String paymentId) async {
    try {
      final response = await _dio.get('/api/payments/$paymentId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch payment: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new payment
  Future<Map<String, dynamic>> createPayment(Map<String, dynamic> paymentData) async {
    try {
      final response = await _dio.post('/api/payments', data: paymentData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create payment: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a payment
  Future<Map<String, dynamic>> updatePayment(String paymentId, Map<String, dynamic> paymentData) async {
    try {
      final response = await _dio.patch('/api/payments/$paymentId', data: paymentData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update payment: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a payment
  Future<void> deletePayment(String paymentId) async {
    try {
      final response = await _dio.delete('/api/payments/$paymentId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete payment: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

