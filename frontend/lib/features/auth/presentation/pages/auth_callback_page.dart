import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';

/// Page that handles OAuth callback on web
/// This page receives the OAuth code and state from the URL query parameters
/// and processes them through the login flow
class AuthCallbackPage extends StatefulWidget {
  const AuthCallbackPage({super.key});

  @override
  State<AuthCallbackPage> createState() => _AuthCallbackPageState();
}

class _AuthCallbackPageState extends State<AuthCallbackPage> {
  @override
  void initState() {
    super.initState();
    if (kIsWeb) {
      // Process the OAuth callback on web
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _processCallback();
      });
    }
  }

  void _processCallback() {
    final uri = Uri.base;
    final code = uri.queryParameters['code'];
    final state = uri.queryParameters['state'];

    if (code == null || state == null) {
      // Missing parameters, redirect to login
      if (mounted) {
        context.go('/login');
      }
      return;
    }

    // The callback URL contains the code and state
    // FlutterWebAuth2 should have already handled this, but if we're here,
    // we need to process it manually
    // For now, redirect to login page - the auth flow should handle this
    // This page is mainly a fallback for direct navigation
    if (mounted) {
      context.go('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: CircularProgressIndicator(),
      ),
    );
  }
}

