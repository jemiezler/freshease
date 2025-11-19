// core/api/bundle_items_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class BundleItemsApi {
  final Dio _dio;
  BundleItemsApi(DioClient client) : _dio = client.dio;

  /// Get all bundle items
  Future<List<Map<String, dynamic>>> listBundleItems() async {
    try {
      final response = await _dio.get('/api/bundle_items');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch bundle items: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get bundle item by ID
  Future<Map<String, dynamic>> getBundleItem(String bundleItemId) async {
    try {
      final response = await _dio.get('/api/bundle_items/$bundleItemId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch bundle item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new bundle item
  Future<Map<String, dynamic>> createBundleItem(Map<String, dynamic> bundleItemData) async {
    try {
      final response = await _dio.post('/api/bundle_items', data: bundleItemData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create bundle item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a bundle item
  Future<Map<String, dynamic>> updateBundleItem(String bundleItemId, Map<String, dynamic> bundleItemData) async {
    try {
      final response = await _dio.patch('/api/bundle_items/$bundleItemId', data: bundleItemData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update bundle item: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a bundle item
  Future<void> deleteBundleItem(String bundleItemId) async {
    try {
      final response = await _dio.delete('/api/bundle_items/$bundleItemId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete bundle item: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

