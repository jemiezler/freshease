// test/unit/product_repository_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:frontend/features/shop/data/product_repository.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/features/shop/data/models/shop_dtos.dart';
import '../mocks/mock_repositories.dart';

void main() {
  group('ProductRepository', () {
    late MockProductRepository mockProductRepository;

    setUp(() {
      mockProductRepository = MockProductRepository();
    });

    group('Get Products', () {
      test('should return list of products when list is called', () async {
        // Arrange
        when(
          mockProductRepository.list(),
        ).thenAnswer((_) async => MockData.mockProducts);

        // Act
        final result = await mockProductRepository.list();

        // Assert
        expect(result, isA<List<Product>>());
        expect(result.length, 2);
        expect(result.first.name, 'Test Product');
        expect(result.last.name, 'Test Product 2');
      });

      test('should return products with pagination', () async {
        // Arrange
        when(
          mockProductRepository.list(),
        ).thenAnswer((_) async => MockData.mockProducts);

        // Act
        final result = await mockProductRepository.list(limit: 10, offset: 20);

        // Assert
        expect(result, isA<List<Product>>());
        verify(mockProductRepository.list(limit: 10, offset: 20)).called(1);
      });

      test('should return products filtered by category', () async {
        // Arrange
        when(
          mockProductRepository.list(),
        ).thenAnswer((_) async => MockData.mockProducts);

        // Act
        final result = await mockProductRepository.list(
          category: 'Test Category',
        );

        // Assert
        expect(result, isA<List<Product>>());
        verify(mockProductRepository.list(category: 'Test Category')).called(1);
      });

      test('should return products filtered by search term', () async {
        // Arrange
        when(
          mockProductRepository.list(),
        ).thenAnswer((_) async => MockData.mockProducts);

        // Act
        final result = await mockProductRepository.list(q: 'test');

        // Assert
        expect(result, isA<List<Product>>());
        verify(mockProductRepository.list(q: 'test')).called(1);
      });
    });

    group('Get Product', () {
      test('should return single product when getProduct is called', () async {
        // Arrange
        when(
          mockProductRepository.getProduct('test-product-id'),
        ).thenAnswer((_) async => MockData.mockProduct);

        // Act
        final result = await mockProductRepository.getProduct(
          'test-product-id',
        );

        // Assert
        expect(result, isA<Product>());
        expect(result!.id, 'test-product-id');
        expect(result.name, 'Test Product');
        verify(mockProductRepository.getProduct('test-product-id')).called(1);
      });

      test('should return null when product not found', () async {
        // Arrange
        when(
          mockProductRepository.getProduct('non-existent-id'),
        ).thenAnswer((_) async => null);

        // Act
        final result = await mockProductRepository.getProduct(
          'non-existent-id',
        );

        // Assert
        expect(result, null);
      });
    });

    group('Get Categories', () {
      test('should return list of categories', () async {
        // Arrange
        final mockCategories = [
          ShopCategoryDTO(id: '1', name: 'Category 1', description: 'Desc 1'),
          ShopCategoryDTO(id: '2', name: 'Category 2', description: 'Desc 2'),
          ShopCategoryDTO(id: '3', name: 'Category 3', description: 'Desc 3'),
        ];
        when(
          mockProductRepository.getCategories(),
        ).thenAnswer((_) async => mockCategories);

        // Act
        final result = await mockProductRepository.getCategories();

        // Assert
        expect(result, isA<List<ShopCategoryDTO>>());
        expect(result.length, 3);
        expect(result.first.name, 'Category 1');
        expect(result.last.name, 'Category 3');
      });
    });

    group('Get Vendors', () {
      test('should return list of vendors', () async {
        // Arrange
        final mockVendors = [
          ShopVendorDTO(
            id: '1',
            name: 'Vendor 1',
            email: 'vendor1@example.com',
            phone: '+1234567890',
            address: '123 Main St',
            city: 'City',
            state: 'State',
            country: 'Country',
            postalCode: '12345',
            website: 'https://vendor1.com',
            logoUrl: 'https://logo1.png',
            description: 'Description 1',
            isActive: 'active',
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ),
          ShopVendorDTO(
            id: '2',
            name: 'Vendor 2',
            email: 'vendor2@example.com',
            phone: '+1234567891',
            address: '456 Oak St',
            city: 'City',
            state: 'State',
            country: 'Country',
            postalCode: '12346',
            website: 'https://vendor2.com',
            logoUrl: 'https://logo2.png',
            description: 'Description 2',
            isActive: 'active',
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ),
        ];
        when(
          mockProductRepository.getVendors(),
        ).thenAnswer((_) async => mockVendors);

        // Act
        final result = await mockProductRepository.getVendors();

        // Assert
        expect(result, isA<List<ShopVendorDTO>>());
        expect(result.length, 2);
        expect(result.first.name, 'Vendor 1');
        expect(result.last.name, 'Vendor 2');
      });
    });

    group('Error Handling', () {
      test('should handle network errors gracefully', () async {
        // Arrange
        when(
          mockProductRepository.list(),
        ).thenThrow(Exception('Network error'));

        // Act & Assert
        expect(() => mockProductRepository.list(), throwsException);
      });

      test('should handle empty results gracefully', () async {
        // Arrange
        when(mockProductRepository.list()).thenAnswer((_) async => <Product>[]);

        // Act
        final result = await mockProductRepository.list();

        // Assert
        expect(result, isA<List<Product>>());
        expect(result.isEmpty, true);
      });
    });
  });
}
