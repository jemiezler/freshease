// core/api/bundles_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class BundlesApi {
  final Dio _dio;
  BundlesApi(DioClient client) : _dio = client.dio;

  /// Get all bundles
  Future<List<Map<String, dynamic>>> listBundles() async {
    try {
      final response = await _dio.get('/api/bundles');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch bundles: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get bundle by ID
  Future<Map<String, dynamic>> getBundle(String bundleId) async {
    try {
      final response = await _dio.get('/api/bundles/$bundleId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch bundle: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new bundle
  Future<Map<String, dynamic>> createBundle(Map<String, dynamic> bundleData) async {
    try {
      final response = await _dio.post('/api/bundles', data: bundleData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create bundle: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a bundle
  Future<Map<String, dynamic>> updateBundle(String bundleId, Map<String, dynamic> bundleData) async {
    try {
      final response = await _dio.patch('/api/bundles/$bundleId', data: bundleData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update bundle: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a bundle
  Future<void> deleteBundle(String bundleId) async {
    try {
      final response = await _dio.delete('/api/bundles/$bundleId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete bundle: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

