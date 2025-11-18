import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/features/shop/data/models/shop_dtos.dart';

void main() {
  group('Product', () {
    test('fromLegacy creates Product with default values', () {
      final product = Product.fromLegacy(
        id: 1,
        name: 'Test Product',
        price: 99.99,
        image: 'image.jpg',
        category: 'Test Category',
      );

      expect(product.id, '1');
      expect(product.name, 'Test Product');
      expect(product.price, 99.99);
      expect(product.image, 'image.jpg');
      expect(product.category, 'Test Category');
      expect(product.description, '');
      expect(product.unitLabel, 'kg');
      expect(product.vendorId, '');
      expect(product.vendorName, '');
      expect(product.categoryId, '');
      expect(product.categoryName, 'Test Category');
      expect(product.stockQuantity, 100);
      expect(product.isInStock, true);
    });

    test('fromShopDTO creates Product from DTO', () {
      final dto = ShopProductDTO(
        id: 'product-1',
        name: 'Test Product',
        price: 99.99,
        description: 'Test description',
        imageUrl: 'image.jpg',
        unitLabel: 'kg',
        isActive: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
        vendorId: 'vendor-1',
        vendorName: 'Test Vendor',
        categoryId: 'cat-1',
        categoryName: 'Test Category',
        stockQuantity: 50,
        isInStock: true,
      );

      final product = Product.fromShopDTO(dto);

      expect(product.id, 'product-1');
      expect(product.name, 'Test Product');
      expect(product.price, 99.99);
      expect(product.image, 'image.jpg');
      expect(product.description, 'Test description');
      expect(product.unitLabel, 'kg');
      expect(product.vendorId, 'vendor-1');
      expect(product.vendorName, 'Test Vendor');
      expect(product.categoryId, 'cat-1');
      expect(product.categoryName, 'Test Category');
      expect(product.stockQuantity, 50);
      expect(product.isInStock, true);
    });

    test('copyWith creates new Product with updated fields', () {
      final original = Product(
        id: 'product-1',
        name: 'Original Product',
        price: 99.99,
        image: 'image.jpg',
        category: 'Category',
        description: 'Description',
        unitLabel: 'kg',
        vendorId: 'vendor-1',
        vendorName: 'Vendor',
        categoryId: 'cat-1',
        categoryName: 'Category',
        stockQuantity: 100,
        isInStock: true,
      );

      final updated = original.copyWith(
        name: 'Updated Product',
        price: 149.99,
        stockQuantity: 50,
      );

      expect(updated.id, 'product-1');
      expect(updated.name, 'Updated Product');
      expect(updated.price, 149.99);
      expect(updated.stockQuantity, 50);
      expect(updated.image, 'image.jpg');
      expect(updated.category, 'Category');
    });

    test('copyWith preserves original values when fields not provided', () {
      final original = Product(
        id: 'product-1',
        name: 'Original Product',
        price: 99.99,
        image: 'image.jpg',
        category: 'Category',
        description: 'Description',
        unitLabel: 'kg',
        vendorId: 'vendor-1',
        vendorName: 'Vendor',
        categoryId: 'cat-1',
        categoryName: 'Category',
        stockQuantity: 100,
        isInStock: true,
      );

      final updated = original.copyWith();

      expect(updated.id, original.id);
      expect(updated.name, original.name);
      expect(updated.price, original.price);
      expect(updated.image, original.image);
      expect(updated.category, original.category);
    });
  });
}

