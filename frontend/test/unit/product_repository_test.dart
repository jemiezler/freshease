import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/shop/data/product_repository.dart';
import 'package:frontend/features/shop/data/sources/shop_api.dart';
import 'package:frontend/features/shop/data/models/shop_dtos.dart';

import 'product_repository_test.mocks.dart';

@GenerateMocks([ShopApiService])
void main() {
  late MockShopApiService mockApiService;
  late ApiProductRepository repository;

  setUp(() {
    mockApiService = MockShopApiService();
    repository = ApiProductRepository(mockApiService);
  });

  group('ApiProductRepository', () {
    test('list returns products from API', () async {
      // Arrange
      final mockProducts = [
        ShopProductDTO(
          id: '1',
          name: 'Test Product',
          price: 99.99,
          description: 'Test description',
          imageUrl: 'https://example.com/image.jpg',
          unitLabel: 'kg',
          isActive: 'active',
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
          vendorId: 'vendor-1',
          vendorName: 'Test Vendor',
          categoryId: 'category-1',
          categoryName: 'Test Category',
          stockQuantity: 100,
          isInStock: true,
        ),
      ];

      final mockResponse = ShopSearchResponse(
        products: mockProducts,
        total: 1,
        limit: 20,
        offset: 0,
        hasMore: false,
      );

      when(mockApiService.searchProducts(any)).thenAnswer(
        (_) async => mockResponse,
      );

      // Act
      final result = await repository.list();

      // Assert
      expect(result, isNotEmpty);
      expect(result.length, 1);
      expect(result.first.name, 'Test Product');
      verify(mockApiService.searchProducts(any)).called(1);
    });

    test('list returns empty list on error', () async {
      // Arrange
      when(mockApiService.searchProducts(any)).thenThrow(
        Exception('Network error'),
      );

      // Act
      final result = await repository.list();

      // Assert
      expect(result, isEmpty);
    });

    test('getProduct returns product by ID', () async {
      // Arrange
      final mockProduct = ShopProductDTO(
        id: '1',
        name: 'Test Product',
        price: 99.99,
        description: 'Test description',
        imageUrl: 'https://example.com/image.jpg',
        unitLabel: 'kg',
        isActive: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
        vendorId: 'vendor-1',
        vendorName: 'Test Vendor',
        categoryId: 'category-1',
        categoryName: 'Test Category',
        stockQuantity: 100,
        isInStock: true,
      );

      when(mockApiService.getProduct('1')).thenAnswer(
        (_) async => mockProduct,
      );

      // Act
      final result = await repository.getProduct('1');

      // Assert
      expect(result, isNotNull);
      expect(result?.name, 'Test Product');
      verify(mockApiService.getProduct('1')).called(1);
    });

    test('getProduct returns null on error', () async {
      // Arrange
      when(mockApiService.getProduct('1')).thenThrow(
        Exception('Not found'),
      );

      // Act
      final result = await repository.getProduct('1');

      // Assert
      expect(result, isNull);
    });

    test('getCategories returns categories from API', () async {
      // Arrange
      final mockCategories = [
        ShopCategoryDTO(
          id: '1',
          name: 'Category 1',
          description: 'Description 1',
        ),
      ];

      when(mockApiService.getCategories()).thenAnswer(
        (_) async => mockCategories,
      );

      // Act
      final result = await repository.getCategories();

      // Assert
      expect(result, isNotEmpty);
      expect(result.length, 1);
      expect(result.first.name, 'Category 1');
      verify(mockApiService.getCategories()).called(1);
    });

    test('getCategories returns empty list on error', () async {
      // Arrange
      when(mockApiService.getCategories()).thenThrow(
        Exception('Network error'),
      );

      // Act
      final result = await repository.getCategories();

      // Assert
      expect(result, isEmpty);
    });

    test('getVendors returns vendors from API', () async {
      // Arrange
      final mockVendors = [
        ShopVendorDTO(
          id: '1',
          name: 'Vendor 1',
          email: 'vendor@example.com',
          phone: '+1234567890',
          address: '123 Street',
          city: 'City',
          state: 'State',
          country: 'Country',
          postalCode: '12345',
          website: 'https://example.com',
          logoUrl: 'https://example.com/logo.png',
          description: 'Description',
          isActive: 'active',
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
      ];

      when(mockApiService.getVendors()).thenAnswer(
        (_) async => mockVendors,
      );

      // Act
      final result = await repository.getVendors();

      // Assert
      expect(result, isNotEmpty);
      expect(result.length, 1);
      expect(result.first.name, 'Vendor 1');
      verify(mockApiService.getVendors()).called(1);
    });

    test('getVendors returns empty list on error', () async {
      // Arrange
      when(mockApiService.getVendors()).thenThrow(
        Exception('Network error'),
      );

      // Act
      final result = await repository.getVendors();

      // Assert
      expect(result, isEmpty);
    });
  });
}
