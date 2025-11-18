import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/cart/data/cart_repository.dart';
import 'package:frontend/features/cart/data/sources/cart_api.dart';
import 'package:frontend/features/cart/data/models/cart_dtos.dart';
import 'package:frontend/features/shop/domain/product.dart';

import 'cart_repository_test.mocks.dart';

@GenerateMocks([CartApiService])
void main() {
  group('ApiCartRepository', () {
    late ApiCartRepository repository;
    late MockCartApiService mockApi;

    setUp(() {
      mockApi = MockCartApiService();
      repository = ApiCartRepository(mockApi);
    });

    test('getCart returns cart from API', () async {
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.getCart()).thenAnswer((_) async => cartDTO);

      final result = await repository.getCart();

      expect(result.id, 'cart-1');
      expect(result.status, 'active');
      verify(mockApi.getCart()).called(1);
    });

    test('getCart returns empty cart on error', () async {
      when(mockApi.getCart()).thenThrow(Exception('Network error'));

      final result = await repository.getCart();

      expect(result.id, isEmpty); // Empty cart has empty ID
      expect(result.items, isEmpty);
      expect(result.subtotal, 0.0);
      expect(result.status, 'pending'); // Empty cart has 'pending' status
    });

    test('addToCart calls API with correct request', () async {
      final product = Product(
        id: '1',
        name: 'Test Product',
        price: 99.99,
        image: 'image.jpg',
        category: 'Test',
        description: 'Test description',
        unitLabel: 'kg',
        vendorId: 'vendor-1',
        vendorName: 'Test Vendor',
        categoryId: 'cat-1',
        categoryName: 'Test Category',
        stockQuantity: 100,
        isInStock: true,
      );
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 99.99,
        shipping: 0.0,
        tax: 0.0,
        total: 99.99,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.addToCart(any)).thenAnswer((_) async => cartDTO);

      final result = await repository.addToCart(product, quantity: 2);

      expect(result.total, 99.99);
      verify(mockApi.addToCart(any)).called(1);
    });

    test('addToCart returns current cart on error', () async {
      final product = Product(
        id: '1',
        name: 'Test Product',
        price: 99.99,
        image: 'image.jpg',
        category: 'Test',
        description: 'Test description',
        unitLabel: 'kg',
        vendorId: 'vendor-1',
        vendorName: 'Test Vendor',
        categoryId: 'cat-1',
        categoryName: 'Test Category',
        stockQuantity: 100,
        isInStock: true,
      );
      final emptyCart = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.addToCart(any)).thenThrow(Exception('Error'));
      when(mockApi.getCart()).thenAnswer((_) async => emptyCart);

      final result = await repository.addToCart(product);

      expect(result, isNotNull);
      verify(mockApi.getCart()).called(1);
    });

    test('updateCartItem calls API with correct request', () async {
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 50.0,
        shipping: 0.0,
        tax: 0.0,
        total: 50.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.updateCartItem(any)).thenAnswer((_) async => cartDTO);

      final result = await repository.updateCartItem('item-1', 3);

      expect(result.total, 50.0);
      verify(mockApi.updateCartItem(any)).called(1);
    });

    test('updateCartItem returns current cart on error', () async {
      final emptyCart = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.updateCartItem(any)).thenThrow(Exception('Error'));
      when(mockApi.getCart()).thenAnswer((_) async => emptyCart);

      final result = await repository.updateCartItem('item-1', 2);

      expect(result, isNotNull);
      verify(mockApi.getCart()).called(1);
    });

    test('removeCartItem calls API with correct ID', () async {
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.removeCartItem('item-1')).thenAnswer((_) async => cartDTO);

      final result = await repository.removeCartItem('item-1');

      expect(result, isNotNull);
      verify(mockApi.removeCartItem('item-1')).called(1);
    });

    test('removeCartItem returns current cart on error', () async {
      final emptyCart = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.removeCartItem(any)).thenThrow(Exception('Error'));
      when(mockApi.getCart()).thenAnswer((_) async => emptyCart);

      final result = await repository.removeCartItem('item-1');

      expect(result, isNotNull);
      verify(mockApi.getCart()).called(1);
    });

    test('applyPromoCode calls API with correct request', () async {
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 100.0,
        shipping: 0.0,
        tax: 0.0,
        promoDiscount: 10.0,
        total: 90.0,
        promoCode: 'SAVE10',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.applyPromoCode(any)).thenAnswer((_) async => cartDTO);

      final result = await repository.applyPromoCode('SAVE10');

      expect(result.promoCode, 'SAVE10');
      expect(result.promoDiscount, 10.0);
      verify(mockApi.applyPromoCode(any)).called(1);
    });

    test('applyPromoCode returns current cart on error', () async {
      final emptyCart = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.applyPromoCode(any)).thenThrow(Exception('Error'));
      when(mockApi.getCart()).thenAnswer((_) async => emptyCart);

      final result = await repository.applyPromoCode('INVALID');

      expect(result, isNotNull);
      verify(mockApi.getCart()).called(1);
    });

    test('removePromoCode calls API', () async {
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 100.0,
        shipping: 0.0,
        tax: 0.0,
        total: 100.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.removePromoCode()).thenAnswer((_) async => cartDTO);

      final result = await repository.removePromoCode();

      expect(result.promoCode, isNull);
      verify(mockApi.removePromoCode()).called(1);
    });

    test('removePromoCode returns current cart on error', () async {
      final emptyCart = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.removePromoCode()).thenThrow(Exception('Error'));
      when(mockApi.getCart()).thenAnswer((_) async => emptyCart);

      final result = await repository.removePromoCode();

      expect(result, isNotNull);
      verify(mockApi.getCart()).called(1);
    });

    test('clearCart calls API', () async {
      final cartDTO = CartDTO(
        id: 'cart-1',
        status: 'active',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        tax: 0.0,
        total: 0.0,
        promoDiscount: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockApi.clearCart()).thenAnswer((_) async => cartDTO);

      final result = await repository.clearCart();

      expect(result.items, isEmpty);
      verify(mockApi.clearCart()).called(1);
    });

    test('clearCart returns empty cart on error', () async {
      when(mockApi.clearCart()).thenThrow(Exception('Error'));

      final result = await repository.clearCart();

      expect(result, isNotNull);
      expect(result.id, isEmpty); // Empty cart has empty ID
      expect(result.status, 'pending'); // Empty cart has 'pending' status
      expect(result.items, isEmpty);
    });
  });
}

