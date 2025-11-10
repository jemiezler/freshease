// ignore_for_file: unused_element

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/core/network/dio_client.dart';
import 'package:frontend/features/auth/domain/repositories/auth_repository.dart';
import 'package:frontend/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:frontend/features/auth/data/sources/auth_api.dart';
import 'package:frontend/features/auth/presentation/state/login_cubit.dart';
import 'package:go_router/go_router.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _phone = TextEditingController();

  late final DioClient _dioClient;
  late final AuthRepository _authRepo;
  late final LoginCubit _cubit;
  bool _hasCheckedToken = false;
  bool _isAutoLogin = false;
  bool _hasNavigated = false; // Prevent multiple navigation attempts

  @override
  void initState() {
    super.initState();
    _dioClient = DioClient();
    _authRepo = AuthRepositoryImpl(AuthApi(_dioClient));
    _cubit = LoginCubit(_authRepo);
    // Check for existing token and verify it
    _verifyExistingToken();
  }

  Future<void> _verifyExistingToken() async {
    // Verify existing token before showing login UI
    try {
      await _cubit.verifyExistingToken();

      // Wait for state to update after verification
      await Future.delayed(const Duration(milliseconds: 200));

      if (!mounted) return;

      final currentState = _cubit.state;

      // If user is set and success, it was auto-login
      if (currentState.user != null && currentState.isSuccess) {
        _isAutoLogin = true;
        // Mark as checked so BlocConsumer renders and listener/builder can navigate
        if (mounted) {
          setState(() {
            _hasCheckedToken = true;
          });
        }
        // Navigation will be handled by listener or builder after widget builds
        return;
      }
    } catch (e) {
      // Token verification failed - this is expected if no token exists
      // Continue to show login screen - don't log as this is normal
    }

    // Token verification failed or no token - show login screen
    if (mounted) {
      setState(() {
        _hasCheckedToken = true;
        _isAutoLogin = false;
      });
    }
  }

  @override
  void dispose() {
    _phone.dispose();
    // Don't close the cubit as it's managed by BlocProvider
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return BlocProvider.value(
      value: _cubit,
      child: BlocConsumer<LoginCubit, LoginState>(
        listener: (context, state) {
          // Handle successful login (either from token verification or manual login)
          // Check current state directly - don't rely only on state changes
          if (state.user != null &&
              state.isSuccess &&
              mounted &&
              !_hasNavigated) {
            _hasNavigated = true;

            if (_isAutoLogin) {
              // Auto-login: Navigate immediately
              Future.delayed(const Duration(milliseconds: 100), () {
                if (mounted) {
                  context.go('/');
                }
              });
            } else {
              // Manual login: Show success message and navigate
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Login successful! Welcome!'),
                  backgroundColor: Colors.green,
                  duration: Duration(seconds: 2),
                ),
              );

              Future.delayed(const Duration(milliseconds: 500), () {
                if (mounted) {
                  context.go('/');
                }
              });
            }
          }

          // Handle errors (only show errors for manual login attempts, not token verification)
          if (state.error != null && state.isFailure && mounted) {
            // Handle different types of errors
            String message;
            Color backgroundColor;

            if (state.error!.contains('canceled')) {
              message = 'Login was canceled. You can try again anytime.';
              backgroundColor = Colors.orange;
            } else if (state.error!.contains('Network error')) {
              message = 'Network error. Please check your internet connection.';
              backgroundColor = Colors.red;
            } else if (state.error!.contains('Authentication failed')) {
              message = 'Authentication failed. Please try again.';
              backgroundColor = Colors.red;
            } else {
              message = 'Login failed: ${state.error!}';
              backgroundColor = Colors.red;
            }

            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(message),
                backgroundColor: backgroundColor,
                duration: const Duration(seconds: 3),
                action: backgroundColor == Colors.orange
                    ? SnackBarAction(
                        label: 'Try Again',
                        textColor: Colors.white,
                        onPressed: () {
                          _cubit.clearError();
                          _cubit.googleLogin();
                        },
                      )
                    : null,
              ),
            );
          }
        },
        builder: (context, state) {
          // Show loading screen while verifying token
          if (!_hasCheckedToken) {
            return Scaffold(
              body: Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [const CircularProgressIndicator()],
                ),
              ),
            );
          }

          // If auto-login was successful, show welcome message and navigate
          if (_isAutoLogin && state.user != null && state.isSuccess) {
            // Trigger navigation once - this is the primary navigation trigger
            if (!_hasNavigated) {
              _hasNavigated = true;
              // Navigate after widget is fully built
              Future.microtask(() {
                Future.delayed(const Duration(milliseconds: 300), () {
                  if (mounted && _isAutoLogin) {
                    context.go('/');
                  }
                });
              });
            }

            return Scaffold(
              body: Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [const CircularProgressIndicator()],
                ),
              ),
            );
          }

          // Show login UI after token verification is complete and not auto-login
          return Scaffold(
            body: Column(
              children: [
                // --- Hero section (image + soft gradient) ---
                SizedBox(
                  height: 320,
                  width: double.infinity,
                  child: Stack(
                    fit: StackFit.expand,
                    children: [
                      Image.asset(
                        'lib/assets/login_hero.png',
                        fit: BoxFit
                            .cover, // fills width and height while keeping aspect ratio
                        errorBuilder: (_, _, _) => const SizedBox.shrink(),
                      ),
                      Container(
                        decoration: const BoxDecoration(
                          gradient: LinearGradient(
                            begin: Alignment.topCenter,
                            end: Alignment.bottomCenter,
                            colors: [Colors.transparent, Color(0x1F6AB7FF)],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                SafeArea(
                  child: LayoutBuilder(
                    builder: (context, c) {
                      return SingleChildScrollView(
                        padding: const EdgeInsets.only(bottom: 24),
                        child: Column(
                          children: [
                            // --- Body content ---
                            ConstrainedBox(
                              constraints: const BoxConstraints(maxWidth: 420),
                              child: Padding(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 24,
                                ),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      'Get your groceries\nwith FreshEase',
                                      style: theme.textTheme.headlineSmall
                                          ?.copyWith(
                                            fontWeight: FontWeight.w800,
                                            height: 1.2,
                                          ),
                                    ),
                                    const SizedBox(height: 20),

                                    const _DividerWithText(
                                      text: 'connect with social media',
                                    ),
                                    const SizedBox(height: 16),

                                    // Google
                                    _SocialButton(
                                      label: 'Continue with Google',
                                      leading: Image.asset(
                                        'assets/icons/google.png',
                                        height: 22,
                                        errorBuilder: (_, _, _) => const Icon(
                                          Icons.g_mobiledata,
                                          size: 26,
                                        ),
                                      ),
                                      onTap: state.loading
                                          ? null
                                          : () => context
                                                .read<LoginCubit>()
                                                .googleLogin(),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ],
                        ),
                      );
                    },
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

/// Compact prefix for the phone TextField (flag + code)
class _CountryCodePrefix extends StatelessWidget {
  final String code;
  final String flag;
  const _CountryCodePrefix({required this.code, required this.flag});

  @override
  Widget build(BuildContext context) {
    final textStyle = Theme.of(context).textTheme.bodyMedium;
    return Padding(
      padding: const EdgeInsets.only(left: 12, right: 8),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(flag, style: textStyle),
          const SizedBox(width: 8),
          Text(code, style: textStyle?.copyWith(fontWeight: FontWeight.w600)),
          const SizedBox(width: 8),
          const VerticalDivider(width: 1, thickness: 1),
          const SizedBox(width: 8),
        ],
      ),
    );
  }
}

/// Centered divider with a caption
class _DividerWithText extends StatelessWidget {
  final String text;
  const _DividerWithText({required this.text});

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).dividerColor;
    return Row(
      children: [
        Expanded(child: Divider(color: color)),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10),
          child: Text(text, style: Theme.of(context).textTheme.bodySmall),
        ),
        Expanded(child: Divider(color: color)),
      ],
    );
  }
}

/// Big rounded social button
class _SocialButton extends StatelessWidget {
  final String label;
  final Widget leading;
  final VoidCallback? onTap;

  const _SocialButton({
    required this.label,
    required this.leading,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isDisabled = onTap == null;
    return SizedBox(
      height: 54,
      width: double.infinity,
      child: OutlinedButton(
        onPressed: onTap,
        style: OutlinedButton.styleFrom(
          side: BorderSide(
            color: Theme.of(
              context,
            ).colorScheme.primary.withValues(alpha: 0.25),
          ),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            leading,
            const SizedBox(width: 12),
            Flexible(
              child: Text(
                label,
                overflow: TextOverflow.ellipsis,
                style: Theme.of(context).textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.w600)
                    .apply(
                      color: isDisabled
                          ? Theme.of(context).disabledColor
                          : null,
                    ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
