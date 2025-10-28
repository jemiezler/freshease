// data/sources/user_api.dart
import 'package:dio/dio.dart';
import '../../../../core/network/dio_client.dart';

class UserApi {
  final Dio _dio;
  UserApi(DioClient client) : _dio = client.dio;

  /// Get current user profile
  Future<Map<String, dynamic>> getCurrentUser() async {
    final response = await _dio.get('/api/whoami');
    return response.data as Map<String, dynamic>;
  }

  /// Get user by ID
  Future<Map<String, dynamic>> getUserById(String userId) async {
    final response = await _dio.get('/api/users/$userId');
    return response.data as Map<String, dynamic>;
  }

  /// Update user profile
  Future<Map<String, dynamic>> updateUser(
    String userId,
    Map<String, dynamic> userData,
  ) async {
    final response = await _dio.put('/api/users/$userId', data: userData);
    return response.data as Map<String, dynamic>;
  }

  /// Update current user profile
  Future<Map<String, dynamic>> updateCurrentUser(
    Map<String, dynamic> userData,
  ) async {
    // First get current user to get the ID
    final currentUser = await getCurrentUser();
    final userId = currentUser['id'] as String;

    return updateUser(userId, userData);
  }
}
