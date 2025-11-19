import 'package:flutter/material.dart';
import 'package:frontend/features/account/presentation/pages/account_page.dart';
import 'package:frontend/features/account/presentation/pages/edit_profile_page.dart';
import 'package:frontend/features/account/presentation/pages/health_page.dart';
import 'package:frontend/features/auth/presentation/pages/forgot_password.dart';
import 'package:frontend/features/auth/presentation/pages/signup_page.dart';
import 'package:frontend/features/cart/presentation/pages/cart_page.dart';
import 'package:frontend/features/checkout/presentation/pages/address_page.dart';
import 'package:frontend/features/checkout/presentation/pages/confirmation_page.dart';
import 'package:frontend/features/checkout/presentation/pages/payment_page.dart';
import 'package:frontend/features/checkout/presentation/pages/review_page.dart';
import 'package:frontend/features/plans/presentation/pages/plan_detail_page.dart';
import 'package:frontend/features/plans/presentation/pages/plans_page.dart';
import 'package:frontend/features/progress/presentation/pages/progress_page.dart';
import 'package:frontend/features/shop/domain/product.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import 'package:frontend/features/shop/presentation/pages/product_detail_page.dart';
import 'package:frontend/core/widgets/design_system/soft_bottom_nav.dart';
import 'package:go_router/go_router.dart';
import '../features/auth/presentation/pages/login_page.dart';
import '../features/auth/presentation/pages/auth_callback_page.dart';
import '../features/shop/presentation/pages/shop_page.dart';
import '../core/state/cart_controller.dart';
import '../features/onboarding/onboarding_page.dart';
import '../features/splash/splash_page.dart'; // Create this file

GoRouter buildRouter() {
  return GoRouter(
    initialLocation: '/splash',
    routes: [
      // Splash Screen
      GoRoute(path: '/splash', builder: (_, _) => const SplashPage()),

      // Onboarding
      GoRoute(path: '/onboarding', builder: (_, _) => const OnboardingPage()),

      // Login
      GoRoute(path: '/login', builder: (_, _) => const LoginPage()),

      // OAuth callback (for web)
      GoRoute(path: '/auth/callback', builder: (_, _) => const AuthCallbackPage()),

      //signup
      GoRoute(path: '/signup', builder: (_, _) => const SignupPage()),

      //forgot-password
      GoRoute(
        path: '/forgot-password',
        builder: (_, _) => const ForgotPasswordPage(),
      ),

      // Main app with bottom navigation
      StatefulShellRoute.indexedStack(
        builder: (context, state, nav) {
          final cart = CartScope.of(context);
          return Scaffold(
            body: nav,
            bottomNavigationBar: SoftBottomNav(
              currentIndex: nav.currentIndex,
              onTap: (i) => nav.goBranch(i, initialLocation: i != nav.currentIndex),
              items: [
                const SoftNavItem(
                  icon: Icons.storefront_outlined,
                  selectedIcon: Icons.storefront,
                  label: 'Shop',
                ),
                SoftNavItem(
                  icon: Icons.shopping_cart_outlined,
                  selectedIcon: Icons.shopping_cart,
                  label: 'Cart',
                  badge: cart.count > 0 ? cart.count : null,
                ),
                const SoftNavItem(
                  icon: Icons.list_alt_outlined,
                  selectedIcon: Icons.list_alt,
                  label: 'Plans',
                ),
                const SoftNavItem(
                  icon: Icons.show_chart_outlined,
                  selectedIcon: Icons.show_chart,
                  label: 'Progress',
                ),
                const SoftNavItem(
                  icon: Icons.person_outline,
                  selectedIcon: Icons.person,
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
              GoRoute(
                path: '/account',
                builder: (_, _) => const AccountPage(),
                routes: [
                  GoRoute(
                    path: 'health', // becomes /account/health
                    builder: (_, _) => const HealthPage(),
                  ),
                  GoRoute(
                    path: 'edit-profile', // becomes /account/edit-profile
                    builder: (context, state) {
                      final user = state.extra as UserProfile?;
                      if (user == null) {
                        context.pop();
                        return const SizedBox.shrink();
                      }
                      return EditProfilePage(user: user);
                    },
                  ),
                ],
              ),
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
