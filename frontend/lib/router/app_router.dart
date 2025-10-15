import 'package:flutter/material.dart';
import 'package:frontend/features/account/presentation/pages/account_page.dart';
import 'package:frontend/features/cart/presentation/pages/cart_page.dart';
import 'package:frontend/features/checkout/presentation/pages/address_page.dart';
import 'package:frontend/features/checkout/presentation/pages/confirmation_page.dart';
import 'package:frontend/features/checkout/presentation/pages/payment_page.dart';
import 'package:frontend/features/checkout/presentation/pages/review_page.dart';
import 'package:frontend/features/plans/presentation/pages/plan_detail_page.dart';
import 'package:frontend/features/plans/presentation/pages/plans_page.dart';
import 'package:frontend/features/progress/presentation/pages/progress_page.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/features/shop/presentation/pages/product_detail_page.dart';
import 'package:go_router/go_router.dart';
import '../features/auth/presentation/pages/login_page.dart';
import '../features/shop/presentation/pages/shop_page.dart';
// import '../features/cart/presentation/pages/cart_page.dart';
// import '../features/plans/presentation/pages/plans_page.dart';
// import '../features/progress/presentation/pages/progress_page.dart';
// import '../features/account/presentation/pages/account_page.dart';
import '../core/state/cart_controller.dart';

GoRouter buildRouter() {
  return GoRouter(
    routes: [
      GoRoute(path: '/login', builder: (_, _) => const LoginPage()),
      StatefulShellRoute.indexedStack(
        builder: (context, state, nav) {
          final cart = CartScope.of(context);
          return Scaffold(
            body: nav,
            bottomNavigationBar: NavigationBar(
              height: 64,
              selectedIndex: nav.currentIndex,
              onDestinationSelected: (i) =>
                  nav.goBranch(i, initialLocation: i != nav.currentIndex),
              destinations: [
                const NavigationDestination(
                  icon: Icon(Icons.storefront),
                  label: 'Shop',
                ),
                NavigationDestination(
                  icon: Stack(
                    children: [
                      const Icon(Icons.shopping_cart),
                      if (cart.count > 0)
                        Positioned(
                          right: 0,
                          top: 0,
                          child: Container(
                            padding: const EdgeInsets.all(3),
                            decoration: BoxDecoration(
                              color: Colors.red,
                              shape: BoxShape.circle,
                            ),
                            child: Text(
                              '${cart.count}',
                              style: const TextStyle(
                                fontSize: 10,
                                color: Colors.white,
                              ),
                            ),
                          ),
                        ),
                    ],
                  ),
                  label: 'Cart',
                ),
                const NavigationDestination(
                  icon: Icon(Icons.list_alt),
                  label: 'Plans',
                ),
                const NavigationDestination(
                  icon: Icon(Icons.show_chart),
                  label: 'Progress',
                ),
                const NavigationDestination(
                  icon: Icon(Icons.person),
                  label: 'Account',
                ),
              ],
            ),
          );
        },
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/',
                builder: (_, _) => const ShopPage(),
                routes: [
                  GoRoute(
                    path: 'shop/product/:id',
                    builder: (context, state) {
                      final product = state.extra as Product;
                      return ProductDetailPage(product: product);
                    },
                  ),
                ],
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/cart',
                builder: (_, _) => const CartPage(),
                routes: [
                  GoRoute(
                    path: 'checkout/address',
                    builder: (_, _) => const AddressPage(),
                  ),
                  GoRoute(
                    path: 'checkout/payment',
                    builder: (_, _) => const PaymentPage(),
                  ),
                  GoRoute(
                    path: 'checkout/review',
                    builder: (_, _) => const ReviewPage(),
                  ),
                ],
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/plans',
                builder: (_, _) => const PlansPage(),
                routes: [
                  GoRoute(
                    path: ':id',
                    builder: (context, state) {
                      final plan = state.extra as Map<String, dynamic>;
                      return PlanDetailPage(plan: plan);
                    },
                  ),
                ],
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/progress',
                builder: (_, _) => const ProgressPage(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(path: '/account', builder: (_, _) => const AccountPage()),
            ],
          ),
        ],
      ),
      GoRoute(
        path: '/checkout/confirmation',
        builder: (_, _) => const ConfirmationPage(),
      ),
    ],
  );
}
