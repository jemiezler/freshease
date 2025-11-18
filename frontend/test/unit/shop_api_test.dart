import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dio/dio.dart';
import 'package:frontend/features/shop/data/sources/shop_api.dart';
import 'package:frontend/features/shop/data/models/shop_dtos.dart';

import 'shop_api_test.mocks.dart';

@GenerateMocks([Dio])
void main() {
  group('ShopApiService', () {
    late ShopApiService apiService;
    late MockDio mockDio;

    setUp(() {
      mockDio = MockDio();
      apiService = ShopApiService(mockDio);
    });

    test('searchProducts returns ShopSearchResponse', () async {
      final responseData = {
        'data': {
          'products': [
            {
              'id': 'product-1',
              'name': 'Test Product',
              'price': 99.99,
              'description': 'Description',
              'image_url': 'image.jpg',
              'unit_label': 'kg',
              'is_active': 'active',
              'created_at': '2024-01-01T00:00:00Z',
              'updated_at': '2024-01-01T00:00:00Z',
              'vendor_id': 'vendor-1',
              'vendor_name': 'Vendor',
              'category_id': 'cat-1',
              'category_name': 'Category',
              'stock_quantity': 100,
              'is_in_stock': true,
            }
          ],
          'total': 1,
          'limit': 20,
          'offset': 0,
          'has_more': false,
        }
      };

      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/'),
                statusCode: 200,
                data: responseData,
              ));

      final filters = ShopSearchFilters(searchTerm: 'test');
      final result = await apiService.searchProducts(filters);

      expect(result.products.length, 1);
      expect(result.total, 1);
      expect(result.products.first.name, 'Test Product');
      verify(mockDio.get(
        '/api/shop/products',
        queryParameters: anyNamed('queryParameters'),
      )).called(1);
    });

    test('searchProducts throws exception on error', () async {
      when(mockDio.get(any, queryParameters: anyNamed('queryParameters')))
          .thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 500,
          data: {'message': 'Server error'},
        ),
      ));

      final filters = ShopSearchFilters();
      expect(
        () => apiService.searchProducts(filters),
        throwsException,
      );
    });

    test('getProduct returns ShopProductDTO', () async {
      final responseData = {
        'data': {
          'id': 'product-1',
          'name': 'Test Product',
          'price': 99.99,
          'description': 'Description',
          'image_url': 'image.jpg',
          'unit_label': 'kg',
          'is_active': 'active',
          'created_at': '2024-01-01T00:00:00Z',
          'updated_at': '2024-01-01T00:00:00Z',
          'vendor_id': 'vendor-1',
          'vendor_name': 'Vendor',
          'category_id': 'cat-1',
          'category_name': 'Category',
          'stock_quantity': 100,
          'is_in_stock': true,
        }
      };

      when(mockDio.get('/api/shop/products/product-1'))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/'),
                statusCode: 200,
                data: responseData,
              ));

      final result = await apiService.getProduct('product-1');

      expect(result.id, 'product-1');
      expect(result.name, 'Test Product');
      verify(mockDio.get('/api/shop/products/product-1')).called(1);
    });

    test('getCategories returns list of ShopCategoryDTO', () async {
      final responseData = {
        'data': [
          {
            'id': 'cat-1',
            'name': 'Category 1',
            'description': 'Description 1',
          },
          {
            'id': 'cat-2',
            'name': 'Category 2',
            'description': 'Description 2',
          },
        ]
      };

      when(mockDio.get('/api/shop/categories'))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/'),
                statusCode: 200,
                data: responseData,
              ));

      final result = await apiService.getCategories();

      expect(result.length, 2);
      expect(result.first.id, 'cat-1');
      expect(result.last.id, 'cat-2');
      verify(mockDio.get('/api/shop/categories')).called(1);
    });

    test('getCategory returns ShopCategoryDTO', () async {
      final responseData = {
        'data': {
          'id': 'cat-1',
          'name': 'Category 1',
          'description': 'Description 1',
        }
      };

      when(mockDio.get('/api/shop/categories/cat-1'))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/'),
                statusCode: 200,
                data: responseData,
              ));

      final result = await apiService.getCategory('cat-1');

      expect(result.id, 'cat-1');
      expect(result.name, 'Category 1');
      verify(mockDio.get('/api/shop/categories/cat-1')).called(1);
    });

    test('getVendors returns list of ShopVendorDTO', () async {
      final responseData = {
        'data': [
          {
            'id': 'vendor-1',
            'name': 'Vendor 1',
            'email': 'vendor1@example.com',
            'phone': '+1234567890',
            'address': '123 Street',
            'city': 'City',
            'state': 'State',
            'country': 'Country',
            'postal_code': '12345',
            'website': 'https://example.com',
            'logo_url': 'logo.jpg',
            'description': 'Description',
            'is_active': 'active',
            'created_at': '2024-01-01T00:00:00Z',
            'updated_at': '2024-01-01T00:00:00Z',
          }
        ]
      };

      when(mockDio.get('/api/shop/vendors'))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/'),
                statusCode: 200,
                data: responseData,
              ));

      final result = await apiService.getVendors();

      expect(result.length, 1);
      expect(result.first.id, 'vendor-1');
      verify(mockDio.get('/api/shop/vendors')).called(1);
    });

    test('getVendor returns ShopVendorDTO', () async {
      final responseData = {
        'data': {
          'id': 'vendor-1',
          'name': 'Vendor 1',
          'email': 'vendor1@example.com',
          'phone': '+1234567890',
          'address': '123 Street',
          'city': 'City',
          'state': 'State',
          'country': 'Country',
          'postal_code': '12345',
          'website': 'https://example.com',
          'logo_url': 'logo.jpg',
          'description': 'Description',
          'is_active': 'active',
          'created_at': '2024-01-01T00:00:00Z',
          'updated_at': '2024-01-01T00:00:00Z',
        }
      };

      when(mockDio.get('/api/shop/vendors/vendor-1'))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/'),
                statusCode: 200,
                data: responseData,
              ));

      final result = await apiService.getVendor('vendor-1');

      expect(result.id, 'vendor-1');
      expect(result.name, 'Vendor 1');
      verify(mockDio.get('/api/shop/vendors/vendor-1')).called(1);
    });

    test('handles network errors gracefully', () async {
      when(mockDio.get(any))
          .thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        type: DioExceptionType.connectionError,
        error: 'Network error',
      ));

      expect(
        () => apiService.getProduct('product-1'),
        throwsException,
      );
    });
  });
}

