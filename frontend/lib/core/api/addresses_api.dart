// core/api/addresses_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class AddressesApi {
  final Dio _dio;
  AddressesApi(DioClient client) : _dio = client.dio;

  /// Get all addresses
  Future<List<Map<String, dynamic>>> listAddresses() async {
    try {
      final response = await _dio.get('/api/addresses');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch addresses: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get address by ID
  Future<Map<String, dynamic>> getAddress(String addressId) async {
    try {
      final response = await _dio.get('/api/addresses/$addressId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch address: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new address
  Future<Map<String, dynamic>> createAddress(Map<String, dynamic> addressData) async {
    try {
      final response = await _dio.post('/api/addresses', data: addressData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create address: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update an address
  Future<Map<String, dynamic>> updateAddress(String addressId, Map<String, dynamic> addressData) async {
    try {
      final response = await _dio.patch('/api/addresses/$addressId', data: addressData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update address: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete an address
  Future<void> deleteAddress(String addressId) async {
    try {
      final response = await _dio.delete('/api/addresses/$addressId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete address: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

