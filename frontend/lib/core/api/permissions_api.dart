// core/api/permissions_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class PermissionsApi {
  final Dio _dio;
  PermissionsApi(DioClient client) : _dio = client.dio;

  /// Get all permissions
  Future<List<Map<String, dynamic>>> listPermissions() async {
    try {
      final response = await _dio.get('/api/permissions');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch permissions: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get permission by ID
  Future<Map<String, dynamic>> getPermission(String permissionId) async {
    try {
      final response = await _dio.get('/api/permissions/$permissionId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch permission: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new permission
  Future<Map<String, dynamic>> createPermission(Map<String, dynamic> permissionData) async {
    try {
      final response = await _dio.post('/api/permissions', data: permissionData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create permission: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a permission
  Future<Map<String, dynamic>> updatePermission(String permissionId, Map<String, dynamic> permissionData) async {
    try {
      final response = await _dio.patch('/api/permissions/$permissionId', data: permissionData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update permission: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a permission
  Future<void> deletePermission(String permissionId) async {
    try {
      final response = await _dio.delete('/api/permissions/$permissionId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete permission: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

