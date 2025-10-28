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

  @override
  void initState() {
    super.initState();
    _dioClient = DioClient();
    _authRepo = AuthRepositoryImpl(AuthApi(_dioClient));
    _cubit = LoginCubit(_authRepo);
  }

  @override
  void dispose() {
    _phone.dispose();
    _cubit.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return BlocProvider.value(
      value: _cubit,
      child: BlocConsumer<LoginCubit, LoginState>(
        listenWhen: (prev, curr) =>
            prev.error != curr.error || prev.user != curr.user,
        listener: (context, state) {
          if (state.error != null) {
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(SnackBar(content: Text(state.error!)));
          }
          if (state.user != null) {
            context.go('/'); // navigate after successful login
          }
        },
        builder: (context, state) {
          final loading = state.loading;

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

                                    // Phone field (+ country code)
                                    TextField(
                                      controller: _phone,
                                      keyboardType: TextInputType.phone,
                                      enabled: !loading,
                                      decoration: const InputDecoration(
                                        labelText: 'Phone number',
                                        hintText: '8x-xxx-xxxx',
                                        border: OutlineInputBorder(
                                          borderRadius: BorderRadius.all(
                                            Radius.circular(14),
                                          ),
                                        ),
                                        prefixIcon: _CountryCodePrefix(
                                          code: '+66',
                                          flag: '🇹🇭',
                                        ),
                                        prefixIconConstraints: BoxConstraints(
                                          minWidth: 0,
                                          minHeight: 0,
                                        ),
                                        contentPadding: EdgeInsets.symmetric(
                                          horizontal: 14,
                                          vertical: 16,
                                        ),
                                      ),
                                    ),
                                    const SizedBox(height: 8),
                                    Text(
                                      'We’ll send an OTP to verify your number.',
                                      style: theme.textTheme.bodySmall
                                          ?.copyWith(color: theme.hintColor),
                                    ),
                                    const SizedBox(height: 16),

                                    // Continue (OTP flow stub)
                                    SizedBox(
                                      width: double.infinity,
                                      height: 52,
                                      child: ElevatedButton(
                                        onPressed: loading
                                            ? null
                                            : () {
                                                // TODO: trigger OTP flow
                                                context.go('/');
                                              },
                                        style: ElevatedButton.styleFrom(
                                          shape: RoundedRectangleBorder(
                                            borderRadius: BorderRadius.circular(
                                              16,
                                            ),
                                          ),
                                        ),
                                        child: loading
                                            ? const SizedBox(
                                                width: 20,
                                                height: 20,
                                                child:
                                                    CircularProgressIndicator(
                                                      strokeWidth: 2,
                                                    ),
                                              )
                                            : const Text('Continue'),
                                      ),
                                    ),

                                    const SizedBox(height: 20),
                                    const _DividerWithText(
                                      text: 'Or connect with social media',
                                    ),
                                    const SizedBox(height: 16),

                                    // Google
                                    _SocialButton(
                                      label: 'Continue with Google',
                                      leading: Image.asset(
                                        'assets/icons/google.png',
                                        height: 22,
                                        errorBuilder: (_, __, ___) =>
                                            const Icon(
                                              Icons.g_mobiledata,
                                              size: 26,
                                            ),
                                      ),
                                      onTap: loading
                                          ? null
                                          : () => context
                                                .read<LoginCubit>()
                                                .googleLogin(),
                                    ),

                                    const SizedBox(height: 12),

                                    // Facebook (stub)
                                    _SocialButton(
                                      label: 'Continue with Facebook',
                                      leading: const Icon(
                                        Icons.facebook,
                                        size: 24,
                                      ),
                                      onTap: loading
                                          ? null
                                          : () {
                                              // TODO: Facebook sign-in
                                              context.go('/');
                                            },
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
