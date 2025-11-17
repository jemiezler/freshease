import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/features/cart/data/cart_repository.dart'
    hide MockCartRepository;
import 'package:frontend/features/cart/data/models/cart_dtos.dart';
import 'package:frontend/features/shop/domain/product.dart';

import 'cart_controller_test.mocks.dart';

@GenerateMocks([CartRepository])
void main() {
  late MockCartRepository mockRepository;
  late CartController cartController;

  setUp(() {
    mockRepository = MockCartRepository();
    cartController = CartController(mockRepository);
  });

  tearDown(() {
    cartController.dispose();
  });

  group('CartController', () {
    test('initial state is correct', () {
      expect(cartController.cart, null);
      expect(cartController.isLoading, false);
      expect(cartController.error, null);
      expect(cartController.count, 0);
      expect(cartController.itemKinds, 0);
      expect(cartController.subtotal, 0.0);
      expect(cartController.total, 0.0);
    });

    test('loads cart on initialization', () async {
      final cartDTO = CartDTO(
        id: '1',
        status: 'pending',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockRepository.getCart()).thenAnswer((_) async => cartDTO);

      // Wait for async initialization
      await Future.delayed(const Duration(milliseconds: 100));

      verify(mockRepository.getCart()).called(1);
    });

    test('add product to cart', () async {
      final product = Product.fromLegacy(
        id: 1,
        name: 'Test Product',
        price: 99.99,
        image: 'test.jpg',
        category: 'Test',
      );

      final cartDTO = CartDTO(
        id: '1',
        status: 'pending',
        items: [
          CartItemDTO(
            id: 'item1',
            productId: '1',
            productName: 'Test Product',
            productPrice: 99.99,
            productImage: 'test.jpg',
            quantity: 1,
            lineTotal: 99.99,
          ),
        ],
        subtotal: 99.99,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 99.99,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(
        mockRepository.addToCart(product, quantity: 1),
      ).thenAnswer((_) async => cartDTO);

      await cartController.add(product, qty: 1);

      verify(mockRepository.addToCart(product, quantity: 1)).called(1);
      expect(cartController.cart, isNotNull);
      expect(cartController.count, 1);
      expect(cartController.total, 99.99);
    });

    test('remove product from cart', () async {
      final product = Product.fromLegacy(
        id: 1,
        name: 'Test Product',
        price: 99.99,
        image: 'test.jpg',
        category: 'Test',
      );

      final emptyCart = CartDTO(
        id: '1',
        status: 'pending',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // First add a product
      final cartWithItem = CartDTO(
        id: '1',
        status: 'pending',
        items: [
          CartItemDTO(
            id: 'item1',
            productId: '1',
            productName: 'Test Product',
            productPrice: 99.99,
            productImage: 'test.jpg',
            quantity: 1,
            lineTotal: 99.99,
          ),
        ],
        subtotal: 99.99,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 99.99,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockRepository.getCart()).thenAnswer((_) async => cartWithItem);
      when(
        mockRepository.removeCartItem('item1'),
      ).thenAnswer((_) async => emptyCart);

      // Wait for initialization
      await Future.delayed(const Duration(milliseconds: 100));

      await cartController.remove(product);

      verify(mockRepository.removeCartItem('item1')).called(1);
      expect(cartController.cart?.items.length, 0);
    });

    test('update product quantity', () async {
      final product = Product.fromLegacy(
        id: 1,
        name: 'Test Product',
        price: 99.99,
        image: 'test.jpg',
        category: 'Test',
      );

      final updatedCart = CartDTO(
        id: '1',
        status: 'pending',
        items: [
          CartItemDTO(
            id: 'item1',
            productId: '1',
            productName: 'Test Product',
            productPrice: 99.99,
            productImage: 'test.jpg',
            quantity: 3,
            lineTotal: 299.97,
          ),
        ],
        subtotal: 299.97,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 299.97,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      final initialCart = CartDTO(
        id: '1',
        status: 'pending',
        items: [
          CartItemDTO(
            id: 'item1',
            productId: '1',
            productName: 'Test Product',
            productPrice: 99.99,
            productImage: 'test.jpg',
            quantity: 1,
            lineTotal: 99.99,
          ),
        ],
        subtotal: 99.99,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 99.99,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockRepository.getCart()).thenAnswer((_) async => initialCart);
      when(
        mockRepository.updateCartItem('item1', 3),
      ).thenAnswer((_) async => updatedCart);

      // Wait for initialization
      await Future.delayed(const Duration(milliseconds: 100));

      await cartController.setQty(product, 3);

      verify(mockRepository.updateCartItem('item1', 3)).called(1);
      expect(cartController.count, 3);
      expect(cartController.total, 299.97);
    });

    test(
      'decrement product quantity removes when quantity reaches zero',
      () async {
        final product = Product.fromLegacy(
          id: 1,
          name: 'Test Product',
          price: 99.99,
          image: 'test.jpg',
          category: 'Test',
        );

        final cartWithItem = CartDTO(
          id: '1',
          status: 'pending',
          items: [
            CartItemDTO(
              id: 'item1',
              productId: '1',
              productName: 'Test Product',
              productPrice: 99.99,
              productImage: 'test.jpg',
              quantity: 1,
              lineTotal: 99.99,
            ),
          ],
          subtotal: 99.99,
          shipping: 0.0,
          promoDiscount: 0.0,
          tax: 0.0,
          total: 99.99,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        final emptyCart = CartDTO(
          id: '1',
          status: 'pending',
          items: [],
          subtotal: 0.0,
          shipping: 0.0,
          promoDiscount: 0.0,
          tax: 0.0,
          total: 0.0,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        when(mockRepository.getCart()).thenAnswer((_) async => cartWithItem);
        when(
          mockRepository.removeCartItem('item1'),
        ).thenAnswer((_) async => emptyCart);

        // Wait for initialization
        await Future.delayed(const Duration(milliseconds: 100));

        await cartController.decrement(product, qty: 1);

        verify(mockRepository.removeCartItem('item1')).called(1);
      },
    );

    test('clear cart', () async {
      final emptyCart = CartDTO(
        id: '1',
        status: 'pending',
        items: [],
        subtotal: 0.0,
        shipping: 0.0,
        promoDiscount: 0.0,
        tax: 0.0,
        total: 0.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockRepository.clearCart()).thenAnswer((_) async => emptyCart);

      await cartController.clear();

      verify(mockRepository.clearCart()).called(1);
      expect(cartController.cart?.items.length, 0);
    });

    test('apply promo code', () async {
      final cartWithPromo = CartDTO(
        id: '1',
        status: 'pending',
        items: [],
        subtotal: 100.0,
        shipping: 0.0,
        promoDiscount: 10.0,
        tax: 0.0,
        total: 90.0,
        promoCode: 'SAVE10',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(
        mockRepository.applyPromoCode('SAVE10'),
      ).thenAnswer((_) async => cartWithPromo);

      final result = await cartController.applyPromo('SAVE10');

      verify(mockRepository.applyPromoCode('SAVE10')).called(1);
      expect(result, true);
      expect(cartController.promoCode, 'SAVE10');
      expect(cartController.total, 90.0);
    });

    test('apply promo code fails on error', () async {
      when(
        mockRepository.applyPromoCode('INVALID'),
      ).thenThrow(Exception('Invalid promo code'));

      final result = await cartController.applyPromo('INVALID');

      expect(result, false);
      expect(cartController.error, isNotNull);
    });

    test('calculates line totals correctly', () async {
      final cartDTO = CartDTO(
        id: '1',
        status: 'pending',
        items: [
          CartItemDTO(
            id: 'item1',
            productId: '1',
            productName: 'Product 1',
            productPrice: 10.0,
            productImage: 'img1.jpg',
            quantity: 2,
            lineTotal: 20.0,
          ),
          CartItemDTO(
            id: 'item2',
            productId: '2',
            productName: 'Product 2',
            productPrice: 20.0,
            productImage: 'img2.jpg',
            quantity: 3,
            lineTotal: 60.0,
          ),
        ],
        subtotal: 80.0,
        shipping: 5.0,
        promoDiscount: 0.0,
        tax: 8.0,
        total: 93.0,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      when(mockRepository.getCart()).thenAnswer((_) async => cartDTO);

      // Wait for initialization
      await Future.delayed(const Duration(milliseconds: 100));

      expect(cartController.lines.length, 2);
      expect(cartController.lines[0].lineTotal, 20.0);
      expect(cartController.lines[1].lineTotal, 60.0);
      expect(cartController.subtotal, 80.0);
      expect(cartController.shipping, 5.0);
      expect(cartController.vat, 8.0);
      expect(cartController.total, 93.0);
    });
  });
}
