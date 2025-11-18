import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/shop/data/product_repository.dart';
import 'package:frontend/features/shop/data/sources/shop_api.dart';
import 'package:frontend/features/shop/data/models/shop_dtos.dart';
import 'package:frontend/features/shop/domain/product.dart';

import 'product_repository_test.mocks.dart';

@GenerateMocks([ShopApiService])
void main() {
  group('ApiProductRepository', () {
    late ApiProductRepository repository;
    late MockShopApiService mockApi;

    setUp(() {
      mockApi = MockShopApiService();
      repository = ApiProductRepository(mockApi);
    });

    test('list returns products from API', () async {
      final response = ShopSearchResponse(
        products: [
          ShopProductDTO(
            id: 'product-1',
            name: 'Test Product',
            price: 99.99,
            description: 'Description',
            imageUrl: 'image.jpg',
            unitLabel: 'kg',
            isActive: 'active',
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
            vendorId: 'vendor-1',
            vendorName: 'Vendor',
            categoryId: 'cat-1',
            categoryName: 'Category',
            stockQuantity: 100,
            isInStock: true,
          ),
        ],
        total: 1,
        limit: 20,
        offset: 0,
        hasMore: false,
      );

      when(mockApi.searchProducts(any)).thenAnswer((_) async => response);

      final result = await repository.list();

      expect(result.length, 1);
      expect(result.first.id, 'product-1');
      expect(result.first.name, 'Test Product');
      verify(mockApi.searchProducts(any)).called(1);
    });

    test('list with search term filters correctly', () async {
      final response = ShopSearchResponse(
        products: [],
        total: 0,
        limit: 20,
        offset: 0,
        hasMore: false,
      );

      when(mockApi.searchProducts(any)).thenAnswer((_) async => response);

      final result = await repository.list(q: 'test');

      expect(result, isEmpty);
      verify(mockApi.searchProducts(argThat(
        predicate<ShopSearchFilters>((f) => f.searchTerm == 'test'),
      ))).called(1);
    });

    test('list with price filters', () async {
      final response = ShopSearchResponse(
        products: [],
        total: 0,
        limit: 20,
        offset: 0,
        hasMore: false,
      );

      when(mockApi.searchProducts(any)).thenAnswer((_) async => response);

      await repository.list(min: 10.0, max: 100.0);

      verify(mockApi.searchProducts(argThat(
        predicate<ShopSearchFilters>((f) =>
            f.minPrice == 10.0 && f.maxPrice == 100.0),
      ))).called(1);
    });

    test('list returns empty list on error', () async {
      when(mockApi.searchProducts(any)).thenThrow(Exception('Error'));

      final result = await repository.list();

      expect(result, isEmpty);
    });

    test('getProduct returns Product from API', () async {
      final dto = ShopProductDTO(
        id: 'product-1',
        name: 'Test Product',
        price: 99.99,
        description: 'Description',
        imageUrl: 'image.jpg',
        unitLabel: 'kg',
        isActive: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
        vendorId: 'vendor-1',
        vendorName: 'Vendor',
        categoryId: 'cat-1',
        categoryName: 'Category',
        stockQuantity: 100,
        isInStock: true,
      );

      when(mockApi.getProduct('product-1')).thenAnswer((_) async => dto);

      final result = await repository.getProduct('product-1');

      expect(result, isNotNull);
      expect(result!.id, 'product-1');
      expect(result.name, 'Test Product');
      verify(mockApi.getProduct('product-1')).called(1);
    });

    test('getProduct returns null on error', () async {
      when(mockApi.getProduct(any)).thenThrow(Exception('Error'));

      final result = await repository.getProduct('product-1');

      expect(result, isNull);
    });

    test('getCategories returns categories from API', () async {
      final categories = [
        ShopCategoryDTO(
          id: 'cat-1',
          name: 'Category 1',
          description: 'Description 1',
        ),
        ShopCategoryDTO(
          id: 'cat-2',
          name: 'Category 2',
          description: 'Description 2',
        ),
      ];

      when(mockApi.getCategories()).thenAnswer((_) async => categories);

      final result = await repository.getCategories();

      expect(result.length, 2);
      expect(result.first.id, 'cat-1');
      verify(mockApi.getCategories()).called(1);
    });

    test('getCategories returns empty list on error', () async {
      when(mockApi.getCategories()).thenThrow(Exception('Error'));

      final result = await repository.getCategories();

      expect(result, isEmpty);
    });

    test('getVendors returns vendors from API', () async {
      final vendors = [
        ShopVendorDTO(
          id: 'vendor-1',
          name: 'Vendor 1',
          email: 'vendor1@example.com',
          phone: '+1234567890',
          address: '123 Street',
          city: 'City',
          state: 'State',
          country: 'Country',
          postalCode: '12345',
          website: 'https://example.com',
          logoUrl: 'logo.jpg',
          description: 'Description',
          isActive: 'active',
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
      ];

      when(mockApi.getVendors()).thenAnswer((_) async => vendors);

      final result = await repository.getVendors();

      expect(result.length, 1);
      expect(result.first.id, 'vendor-1');
      verify(mockApi.getVendors()).called(1);
    });

    test('getVendors returns empty list on error', () async {
      when(mockApi.getVendors()).thenThrow(Exception('Error'));

      final result = await repository.getVendors();

      expect(result, isEmpty);
    });
  });

  group('MockProductRepository', () {
    late MockProductRepository repository;

    setUp(() {
      repository = MockProductRepository();
    });

    test('list returns filtered products', () async {
      final result = await repository.list(q: 'Salad');

      expect(result.length, 1);
      expect(result.first.name, contains('Salad'));
    });

    test('list filters by category', () async {
      final result = await repository.list(category: 'Fruits');

      expect(result.length, 2);
      expect(result.every((p) => p.category == 'Fruits'), true);
    });

    test('list filters by price range', () async {
      final result = await repository.list(min: 100.0, max: 150.0);

      // Mock data has: 120, 135, 149 in this range
      expect(result.length, 3);
      expect(result.every((p) => p.price >= 100.0 && p.price <= 150.0), true);
    });

    test('getProduct returns product by ID', () async {
      final result = await repository.getProduct('1');

      expect(result, isNotNull);
      expect(result!.id, '1');
    });

    test('getProduct returns null for non-existent product', () async {
      final result = await repository.getProduct('999');

      expect(result, isNull);
    });

    test('getCategories returns mock categories', () async {
      final result = await repository.getCategories();

      expect(result.length, 5);
      expect(result.first.name, 'All');
    });

    test('getVendors returns mock vendors', () async {
      final result = await repository.getVendors();

      expect(result.length, 1);
      expect(result.first.name, 'Fresh Farm Co.');
    });
  });
}
