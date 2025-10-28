import 'package:flutter/material.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:frontend/core/theme/app_theme.dart';
import 'package:frontend/router/app_router.dart';
import 'package:health/health.dart';

class App extends StatefulWidget {
  const App({super.key});

  @override
  AppState createState() => AppState();
}

enum HealthState {
  DATA_NOT_FETCHED,
  FETCHING_DATA,
  DATA_READY,
  NO_DATA,
  AUTHORIZED,
  AUTH_NOT_GRANTED,
  DATA_ADDED,
  DATA_DELETED,
  DATA_NOT_ADDED,
  DATA_NOT_DELETED,
  STEPS_READY,
  HEALTH_CONNECT_STATUS,
  PERMISSIONS_REVOKING,
  PERMISSIONS_REVOKED,
  PERMISSIONS_NOT_REVOKED,
}

final health = Health();

class AppState extends State<App> {
  @override
  void initState() {
    health.configure();
    health.getHealthConnectSdkStatus();
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

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
