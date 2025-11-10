// core/api/inventories_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class InventoriesApi {
  final Dio _dio;
  InventoriesApi(DioClient client) : _dio = client.dio;

  /// Get all inventories
  Future<List<Map<String, dynamic>>> listInventories() async {
    try {
      final response = await _dio.get('/api/inventories');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch inventories: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get inventory by ID
  Future<Map<String, dynamic>> getInventory(String inventoryId) async {
    try {
      final response = await _dio.get('/api/inventories/$inventoryId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch inventory: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new inventory
  Future<Map<String, dynamic>> createInventory(Map<String, dynamic> inventoryData) async {
    try {
      final response = await _dio.post('/api/inventories', data: inventoryData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create inventory: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update an inventory
  Future<Map<String, dynamic>> updateInventory(String inventoryId, Map<String, dynamic> inventoryData) async {
    try {
      final response = await _dio.patch('/api/inventories/$inventoryId', data: inventoryData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update inventory: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete an inventory
  Future<void> deleteInventory(String inventoryId) async {
    try {
      final response = await _dio.delete('/api/inventories/$inventoryId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete inventory: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

