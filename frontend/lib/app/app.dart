// ignore_for_file: constant_identifier_names

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:frontend/core/theme/app_theme.dart';
import 'package:frontend/features/account/presentation/state/user_cubit.dart';
import 'package:frontend/router/app_router.dart';
import 'package:frontend/app/di.dart';

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

class AppState extends State<App> {
  @override
  void initState() {
    // Health initialization is now handled in HealthService/HealthController
    // and is conditional based on platform (mobile only)
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
        controller: getIt<CartController>(),
        child: BlocProvider<UserCubit>.value(
          value: getIt<UserCubit>()..loadCurrentUser(),
          child: MaterialApp.router(
            title: 'FreshEase',
            debugShowCheckedModeBanner: false,
            routerConfig: buildRouter(),
            theme: AppTheme.light(),
            themeMode: ThemeMode.light,
          ),
        ),
      ),
    );
  }
}
