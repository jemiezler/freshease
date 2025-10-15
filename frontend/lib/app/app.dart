import 'package:flutter/material.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:frontend/core/theme/app_theme.dart';
import 'package:frontend/router/app_router.dart';

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return CheckoutScope(
      controller: CheckoutController(),
      child: CartScope(
        // provides CartController down the tree
        controller: CartController(),
        child: MaterialApp.router(
          title: 'FreshEase',
          debugShowCheckedModeBanner: false,
          routerConfig: buildRouter(),
          theme: AppTheme.light(),
          themeMode: ThemeMode.light,
        ),
      ),
    );
  }
}
