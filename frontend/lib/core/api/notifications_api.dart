// core/api/notifications_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class NotificationsApi {
  final Dio _dio;
  NotificationsApi(DioClient client) : _dio = client.dio;

  /// Get all notifications
  Future<List<Map<String, dynamic>>> listNotifications() async {
    try {
      final response = await _dio.get('/api/notifications');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch notifications: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get notification by ID
  Future<Map<String, dynamic>> getNotification(String notificationId) async {
    try {
      final response = await _dio.get('/api/notifications/$notificationId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch notification: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new notification
  Future<Map<String, dynamic>> createNotification(Map<String, dynamic> notificationData) async {
    try {
      final response = await _dio.post('/api/notifications', data: notificationData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create notification: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a notification
  Future<Map<String, dynamic>> updateNotification(String notificationId, Map<String, dynamic> notificationData) async {
    try {
      final response = await _dio.patch('/api/notifications/$notificationId', data: notificationData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update notification: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a notification
  Future<void> deleteNotification(String notificationId) async {
    try {
      final response = await _dio.delete('/api/notifications/$notificationId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete notification: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

