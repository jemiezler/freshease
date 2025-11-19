import 'package:dio/dio.dart';
import '../models/cart_dtos.dart';

class CartApiService {
  final Dio _dio;

  CartApiService(this._dio);

  Future<CartDTO> getCart() async {
    try {
      final response = await _dio.get('/api/carts/current');

      if (response.statusCode == 200) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to fetch cart: ${response.statusCode}');
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

  Future<CartDTO> addToCart(AddToCartRequest request) async {
    try {
      // Use PATCH since the endpoint allows PATCH (not POST)
      final response = await _dio.patch(
        '/api/carts/add-item',
        data: request.toJson(),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to add item to cart: ${response.statusCode}');
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

  Future<CartDTO> updateCartItem(UpdateCartItemRequest request) async {
    try {
      final response = await _dio.patch(
        '/api/carts/update-item',
        data: request.toJson(),
      );

      if (response.statusCode == 200) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to update cart item: ${response.statusCode}');
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

  Future<CartDTO> removeCartItem(String cartItemId) async {
    try {
      final response = await _dio.delete('/api/carts/remove-item/$cartItemId');

      if (response.statusCode == 200) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to remove cart item: ${response.statusCode}');
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

  Future<CartDTO> applyPromoCode(ApplyPromoRequest request) async {
    try {
      final response = await _dio.post(
        '/api/carts/apply-promo',
        data: request.toJson(),
      );

      if (response.statusCode == 200) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to apply promo code: ${response.statusCode}');
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

  Future<CartDTO> removePromoCode() async {
    try {
      final response = await _dio.delete('/api/carts/remove-promo');

      if (response.statusCode == 200) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to remove promo code: ${response.statusCode}');
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

  Future<CartDTO> clearCart() async {
    try {
      final response = await _dio.delete('/api/carts/clear');

      if (response.statusCode == 200) {
        return CartDTO.fromJson(response.data['data']);
      } else {
        throw Exception('Failed to clear cart: ${response.statusCode}');
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
