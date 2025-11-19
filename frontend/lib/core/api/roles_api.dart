// core/api/roles_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class RolesApi {
  final Dio _dio;
  RolesApi(DioClient client) : _dio = client.dio;

  /// Get all roles
  Future<List<Map<String, dynamic>>> listRoles() async {
    try {
      final response = await _dio.get('/api/roles');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch roles: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get role by ID
  Future<Map<String, dynamic>> getRole(String roleId) async {
    try {
      final response = await _dio.get('/api/roles/$roleId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch role: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new role
  Future<Map<String, dynamic>> createRole(Map<String, dynamic> roleData) async {
    try {
      final response = await _dio.post('/api/roles', data: roleData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create role: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a role
  Future<Map<String, dynamic>> updateRole(String roleId, Map<String, dynamic> roleData) async {
    try {
      final response = await _dio.patch('/api/roles/$roleId', data: roleData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update role: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a role
  Future<void> deleteRole(String roleId) async {
    try {
      final response = await _dio.delete('/api/roles/$roleId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete role: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

