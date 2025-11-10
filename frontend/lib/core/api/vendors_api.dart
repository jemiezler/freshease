// core/api/vendors_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class VendorsApi {
  final Dio _dio;
  VendorsApi(DioClient client) : _dio = client.dio;

  /// Get all vendors
  Future<List<Map<String, dynamic>>> listVendors() async {
    try {
      final response = await _dio.get('/api/vendors');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch vendors: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get vendor by ID
  Future<Map<String, dynamic>> getVendor(String vendorId) async {
    try {
      final response = await _dio.get('/api/vendors/$vendorId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch vendor: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new vendor
  Future<Map<String, dynamic>> createVendor(Map<String, dynamic> vendorData) async {
    try {
      final response = await _dio.post('/api/vendors', data: vendorData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create vendor: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a vendor
  Future<Map<String, dynamic>> updateVendor(String vendorId, Map<String, dynamic> vendorData) async {
    try {
      final response = await _dio.patch('/api/vendors/$vendorId', data: vendorData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update vendor: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a vendor
  Future<void> deleteVendor(String vendorId) async {
    try {
      final response = await _dio.delete('/api/vendors/$vendorId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete vendor: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

