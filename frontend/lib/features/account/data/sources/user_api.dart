// data/sources/user_api.dart
import 'package:dio/dio.dart';
import '../../../../core/network/dio_client.dart';

class UserApi {
  final Dio _dio;
  UserApi(DioClient client) : _dio = client.dio;

  /// Get current user profile from /api/whoami endpoint
  /// Returns user data directly (not wrapped in data field)
  Future<Map<String, dynamic>> getCurrentUser() async {
    try {
      final response = await _dio.get('/api/whoami');

      if (response.statusCode == 200) {
        // Backend returns user directly as JSON object
        final data = response.data;
        if (data is Map<String, dynamic>) {
          return data;
        } else {
          throw Exception('Invalid response format from server');
        }
      } else {
        throw Exception('Failed to fetch user profile: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response?.statusCode;
        final message = e.response?.data is Map<String, dynamic>
            ? e.response?.data['message'] as String?
            : e.message;

        if (statusCode == 401) {
          throw Exception('Authentication required. Please login again.');
        } else if (statusCode == 403) {
          throw Exception('Access forbidden: ${message ?? 'Insufficient permissions'}');
        } else if (statusCode == 404) {
          throw Exception('User not found');
        } else {
          throw Exception('API Error: ${message ?? e.message}');
        }
      } else {
        throw Exception('Network Error: ${e.message ?? 'Unable to connect to server'}');
      }
    } catch (e) {
      if (e is Exception) {
        rethrow;
      }
      throw Exception('Unexpected error: $e');
    }
  }

  /// Get user by ID from /api/users/:id endpoint
  /// Returns user data directly (not wrapped in data field)
  Future<Map<String, dynamic>> getUserById(String userId) async {
    try {
      final response = await _dio.get('/api/users/$userId');

      if (response.statusCode == 200) {
        // Backend returns user directly as JSON object
        final data = response.data;
        if (data is Map<String, dynamic>) {
          return data;
        } else {
          throw Exception('Invalid response format from server');
        }
      } else {
        throw Exception('Failed to fetch user: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response?.statusCode;
        final message = e.response?.data is Map<String, dynamic>
            ? e.response?.data['message'] as String?
            : e.message;

        if (statusCode == 401) {
          throw Exception('Authentication required. Please login again.');
        } else if (statusCode == 403) {
          throw Exception('Access forbidden: ${message ?? 'Insufficient permissions'}');
        } else if (statusCode == 404) {
          throw Exception('User not found');
        } else {
          throw Exception('API Error: ${message ?? e.message}');
        }
      } else {
        throw Exception('Network Error: ${e.message ?? 'Unable to connect to server'}');
      }
    } catch (e) {
      if (e is Exception) {
        rethrow;
      }
      throw Exception('Unexpected error: $e');
    }
  }

  /// Update user profile via PUT /api/users/:id
  /// Accepts partial update data (only fields that need to be updated)
  /// Returns updated user data directly (not wrapped in data field)
  Future<Map<String, dynamic>> updateUser(
    String userId,
    Map<String, dynamic> userData,
  ) async {
    try {
      // Remove null values and empty strings to send only changed fields
      final cleanData = <String, dynamic>{};
      userData.forEach((key, value) {
        if (value != null && value != '') {
          cleanData[key] = value;
        }
      });

      final response = await _dio.put(
        '/api/users/$userId',
        data: cleanData,
        options: Options(
          headers: {'Content-Type': 'application/json'},
        ),
      );

      if (response.statusCode == 200) {
        // Backend returns updated user directly as JSON object
        final data = response.data;
        if (data is Map<String, dynamic>) {
          return data;
        } else {
          throw Exception('Invalid response format from server');
        }
      } else {
        throw Exception('Failed to update user: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response?.statusCode;
        final message = e.response?.data is Map<String, dynamic>
            ? e.response?.data['message'] as String?
            : e.message;

        if (statusCode == 401) {
          throw Exception('Authentication required. Please login again.');
        } else if (statusCode == 403) {
          throw Exception('Access forbidden: ${message ?? 'You can only update your own profile'}');
        } else if (statusCode == 404) {
          throw Exception('User not found');
        } else if (statusCode == 400) {
          throw Exception('Invalid data: ${message ?? 'Please check your input'}');
        } else {
          throw Exception('API Error: ${message ?? e.message}');
        }
      } else {
        throw Exception('Network Error: ${e.message ?? 'Unable to connect to server'}');
      }
    } catch (e) {
      if (e is Exception) {
        rethrow;
      }
      throw Exception('Unexpected error: $e');
    }
  }

  /// Update current user profile
  /// Automatically gets the current user ID and updates the profile
  Future<Map<String, dynamic>> updateCurrentUser(
    Map<String, dynamic> userData,
  ) async {
    try {
      // First get current user to get the ID
      final currentUser = await getCurrentUser();
      final userId = currentUser['id'] as String?;

      if (userId == null || userId.isEmpty) {
        throw Exception('User ID not found in current user data');
      }

      return await updateUser(userId, userData);
    } on Exception {
      rethrow;
    } catch (e) {
      throw Exception('Failed to update current user: $e');
    }
  }
}
