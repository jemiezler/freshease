import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/cart/presentation/pages/cart_page.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/features/cart/data/cart_repository.dart';
import 'package:frontend/features/cart/data/models/cart_dtos.dart';
import 'package:frontend/features/shop/domain/product.dart';

import 'cart_page_test.mocks.dart';

@GenerateMocks([CartRepository])
void main() {
  group('CartPage Widget Tests', () {
    late CartController cartController;
    late MockCartRepository mockRepository;

    setUp(() {
      mockRepository = MockCartRepository();
      cartController = CartController(mockRepository);
    });

    tearDown(() {
      cartController.dispose();
    });

    Widget createTestWidget(Widget child) {
      return MaterialApp.router(
        routerConfig: GoRouter(
          routes: [
            GoRoute(
              path: '/',
              builder: (context, state) => child,
            ),
            GoRoute(
              path: '/cart',
              builder: (context, state) => child,
            ),
            GoRoute(
              path: '/cart/checkout/address',
              builder: (context, state) => const Scaffold(body: Text('Address Page')),
            ),
          ],
        ),
        builder: (context, child) {
          return CartScope(
            controller: cartController,
            child: child!,
          );
        },
      );
    }

    testWidgets('displays empty cart state when cart is empty', (tester) async {
      when(mockRepository.getCart()).thenAnswer((_) async => CartDTO(
            id: '',
            status: 'pending',
            items: [],
            subtotal: 0.0,
            shipping: 0.0,
            tax: 0.0,
            total: 0.0,
            promoDiscount: 0.0,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ));

      await tester.pumpWidget(createTestWidget(const CartPage()));
      await tester.pumpAndSettle();

      expect(find.text('Your cart is empty'), findsOneWidget);
      expect(find.text('Find fresh picks in the shop and add them here.'), findsOneWidget);
      expect(find.text('Go to Shop'), findsOneWidget);
    });

    testWidgets('displays cart items when cart has products', (tester) async {
      when(mockRepository.getCart()).thenAnswer((_) async => CartDTO(
            id: 'cart-1',
            status: 'active',
            items: [
              CartItemDTO(
                id: 'item-1',
                productId: '1',
                productName: 'Test Product',
                productImage: 'test.jpg',
                productPrice: 10.0,
                quantity: 2,
                lineTotal: 20.0,
              ),
            ],
            subtotal: 20.0,
            shipping: 0.0,
            tax: 0.0,
            total: 20.0,
            promoDiscount: 0.0,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ));

      await tester.pumpWidget(createTestWidget(const CartPage()));
      await tester.pumpAndSettle();

      expect(find.text('Test Product'), findsOneWidget);
      expect(find.text('฿10.00'), findsWidgets);
    });

    testWidgets('displays checkout button when cart has items', (tester) async {
      when(mockRepository.getCart()).thenAnswer((_) async => CartDTO(
            id: 'cart-1',
            status: 'active',
            items: [
              CartItemDTO(
                id: 'item-1',
                productId: '1',
                productName: 'Test Product',
                productImage: 'test.jpg',
                productPrice: 10.0,
                quantity: 1,
                lineTotal: 10.0,
              ),
            ],
            subtotal: 10.0,
            shipping: 0.0,
            tax: 0.0,
            total: 10.0,
            promoDiscount: 0.0,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ));

      await tester.pumpWidget(createTestWidget(const CartPage()));
      await tester.pumpAndSettle();

      expect(find.text('Checkout'), findsOneWidget);
    });

    testWidgets('displays clear cart button when cart has items', (tester) async {
      when(mockRepository.getCart()).thenAnswer((_) async => CartDTO(
            id: 'cart-1',
            status: 'active',
            items: [
              CartItemDTO(
                id: 'item-1',
                productId: '1',
                productName: 'Test Product',
                productImage: 'test.jpg',
                productPrice: 10.0,
                quantity: 1,
                lineTotal: 10.0,
              ),
            ],
            subtotal: 10.0,
            shipping: 0.0,
            tax: 0.0,
            total: 10.0,
            promoDiscount: 0.0,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ));

      await tester.pumpWidget(createTestWidget(const CartPage()));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.delete_sweep_outlined), findsOneWidget);
    });

    testWidgets('does not display clear cart button when cart is empty', (tester) async {
      when(mockRepository.getCart()).thenAnswer((_) async => CartDTO(
            id: '',
            status: 'pending',
            items: [],
            subtotal: 0.0,
            shipping: 0.0,
            tax: 0.0,
            total: 0.0,
            promoDiscount: 0.0,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ));

      await tester.pumpWidget(createTestWidget(const CartPage()));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.delete_sweep_outlined), findsNothing);
    });
  });
}

