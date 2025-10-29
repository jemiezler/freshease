import 'package:dio/dio.dart';
import '../models/shop_dtos.dart';

class ShopApiService {
  final Dio _dio;

  ShopApiService(this._dio);

  Future<ShopSearchResponse> searchProducts(ShopSearchFilters filters) async {
    try {
      final response = await _dio.get(
        '/api/shop/products',
        queryParameters: filters.toQueryParams(),
      );

      if (response.statusCode == 200) {
        return ShopSearchResponse.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to fetch products: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          'API Error: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception('Network Error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  Future<ShopProductDTO> getProduct(String productId) async {
    try {
      final response = await _dio.get('/api/shop/products/$productId');

      if (response.statusCode == 200) {
        return ShopProductDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to fetch product: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          'API Error: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception('Network Error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  Future<List<ShopCategoryDTO>> getCategories() async {
    try {
      final response = await _dio.get('/api/shop/categories');

      if (response.statusCode == 200) {
        final List<dynamic> categoriesJson = response.data['data'];
        return categoriesJson
            .map(
              (json) => ShopCategoryDTO.fromJson(json as Map<String, dynamic>),
            )
            .toList();
      } else {
        throw Exception('Failed to fetch categories: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          'API Error: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception('Network Error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  Future<ShopCategoryDTO> getCategory(String categoryId) async {
    try {
      final response = await _dio.get('/api/shop/categories/$categoryId');

      if (response.statusCode == 200) {
        return ShopCategoryDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to fetch category: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          'API Error: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception('Network Error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  Future<List<ShopVendorDTO>> getVendors() async {
    try {
      final response = await _dio.get('/api/shop/vendors');

      if (response.statusCode == 200) {
        final List<dynamic> vendorsJson = response.data['data'];
        return vendorsJson
            .map((json) => ShopVendorDTO.fromJson(json as Map<String, dynamic>))
            .toList();
      } else {
        throw Exception('Failed to fetch vendors: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          'API Error: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception('Network Error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  Future<ShopVendorDTO> getVendor(String vendorId) async {
    try {
      final response = await _dio.get('/api/shop/vendors/$vendorId');

      if (response.statusCode == 200) {
        return ShopVendorDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to fetch vendor: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          'API Error: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception('Network Error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }
}
