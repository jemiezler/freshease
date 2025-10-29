// ignore_for_file: strict_top_level_inference

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:frontend/app/di.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../state/user_cubit.dart';

class AccountPage extends StatefulWidget {
  const AccountPage({super.key});
  @override
  State<AccountPage> createState() => _AccountPageState();
}

class _AccountPageState extends State<AccountPage> {
  late final UserCubit _userCubit;
  bool _notif = true;
  bool _marketing = false;
  bool _darkMode = false; // demo toggle (wire to real theme later)

  @override
  void initState() {
    super.initState();
    _userCubit = getIt<UserCubit>();
    _checkAuthAndLoadUser();
  }

  Future<void> _checkAuthAndLoadUser() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('access_token');

    if (token == null || token.isEmpty) {
      // No token, redirect to login
      if (mounted) {
        context.go('/login');
      }
      return;
    }

    // Token exists, load user data
    _userCubit.loadCurrentUser();
  }

  @override
  void dispose() {
    // Don't close the UserCubit as it's a singleton
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final addr = CheckoutScope.of(context).shippingAddress;

    return BlocProvider.value(
      value: _userCubit,
      child: Scaffold(
        appBar: AppBar(title: const Text('Account')),
        body: BlocBuilder<UserCubit, UserState>(
          bloc: _userCubit,
          builder: (context, state) {
            if (state.loading) {
              return const Center(child: CircularProgressIndicator());
            }

            if (state.error != null) {
              // Check if it's an authentication error
              if (state.error!.contains('Authentication expired') ||
                  state.error!.contains('401') ||
                  state.error!.contains('missing bearer token')) {
                // Redirect to login for auth errors
                WidgetsBinding.instance.addPostFrameCallback((_) {
                  context.go('/login');
                });
                return const Center(child: CircularProgressIndicator());
              }

              // Show other errors with retry option
              return Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.error_outline, size: 64, color: Colors.red[300]),
                    const SizedBox(height: 16),
                    Text(
                      state.error!,
                      style: const TextStyle(fontSize: 16),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 16),
                    ElevatedButton(
                      onPressed: () => _userCubit.loadCurrentUser(),
                      child: const Text('Retry'),
                    ),
                  ],
                ),
              );
            }

            final user = state.user;
            if (user == null) {
              return const Center(child: Text('No user data available'));
            }

            return LayoutBuilder(
              builder: (context, c) {
                final isWide = c.maxWidth >= 1000;

                final profileCard = _CardX(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        CircleAvatar(
                          radius: 32,
                          backgroundImage: user.avatar != null
                              ? NetworkImage(user.avatar!)
                              : null,
                          child: user.avatar == null
                              ? Text(
                                  user.initials,
                                  style: const TextStyle(
                                    fontSize: 20,
                                    fontWeight: FontWeight.w800,
                                  ),
                                )
                              : null,
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                user.displayName,
                                style: const TextStyle(
                                  fontSize: 18,
                                  fontWeight: FontWeight.w800,
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                user.email,
                                style: TextStyle(color: Colors.grey[700]),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                user.phone ?? 'No phone number',
                                style: TextStyle(color: Colors.grey[700]),
                              ),
                              if (user.bio != null) ...[
                                const SizedBox(height: 4),
                                Text(
                                  user.bio!,
                                  style: TextStyle(
                                    color: Colors.grey[600],
                                    fontSize: 12,
                                  ),
                                  maxLines: 2,
                                  overflow: TextOverflow.ellipsis,
                                ),
                              ],
                              const SizedBox(height: 12),
                              Wrap(
                                spacing: 8,
                                runSpacing: 8,
                                children: [
                                  OutlinedButton.icon(
                                    icon: const Icon(Icons.edit),
                                    label: const Text('Edit Profile'),
                                    onPressed: () => _openEditProfile(user),
                                  ),
                                  OutlinedButton.icon(
                                    icon: const Icon(
                                      Icons.location_on_outlined,
                                    ),
                                    label: const Text('Manage Address'),
                                    onPressed: () =>
                                        context.go('/cart/checkout/address'),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                );

                final shortcutsCard = _CardX(
                  child: Column(
                    children: [
                      const _SectionHeader(title: 'Shortcuts'),
                      ListTile(
                        leading: const Icon(Icons.health_and_safety),
                        title: const Text('My Health'),
                        subtitle: const Text('See health data & suggessions'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () => context.go(
                          '/account/health',
                        ), // Orders tab lives there
                      ),
                      const Divider(height: 1),
                      ListTile(
                        leading: const Icon(Icons.receipt_long_outlined),
                        title: const Text('My Orders'),
                        subtitle: const Text(
                          'See order history & delivery status',
                        ),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () =>
                            context.go('/progress'), // Orders tab lives there
                      ),
                      const Divider(height: 1),
                      ListTile(
                        leading: const Icon(Icons.eco_outlined),
                        title: const Text('My Subscriptions'),
                        subtitle: const Text('Manage plan renewals and pauses'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () =>
                            context.go('/progress'), // Subscriptions tab
                      ),
                      const Divider(height: 1),
                      ListTile(
                        leading: const Icon(Icons.payment_outlined),
                        title: const Text('Payment Methods'),
                        subtitle: const Text('Saved cards (demo)'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Payment methods screen (todo)'),
                            ),
                          );
                        },
                      ),
                    ],
                  ),
                );

                final addressCard = _CardX(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      children: [
                        const _SectionHeader(title: 'Default Shipping Address'),
                        const SizedBox(height: 8),
                        Align(
                          alignment: Alignment.centerLeft,
                          child: addr == null
                              ? Text(
                                  'No address saved',
                                  style: TextStyle(color: Colors.grey[700]),
                                )
                              : Text(
                                  '${addr.fullName}\n'
                                  '${addr.line1}${addr.line2 != null ? '\n${addr.line2}' : ''}\n'
                                  '${addr.subDistrict}, ${addr.district}, ${addr.province} ${addr.postalCode}\n'
                                  '☎ ${addr.phone}',
                                  style: const TextStyle(height: 1.35),
                                ),
                        ),
                        const SizedBox(height: 12),
                        Align(
                          alignment: Alignment.centerLeft,
                          child: OutlinedButton.icon(
                            onPressed: () =>
                                context.go('/cart/checkout/address'),
                            icon: const Icon(Icons.edit_location_alt_outlined),
                            label: Text(
                              addr == null ? 'Add Address' : 'Edit Address',
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                );

                final settingsCard = _CardX(
                  child: Column(
                    children: [
                      const _SectionHeader(title: 'Preferences'),
                      SwitchListTile(
                        value: _notif,
                        onChanged: (v) => setState(() => _notif = v),
                        secondary: const Icon(Icons.notifications_outlined),
                        title: const Text('Push notifications'),
                        subtitle: const Text(
                          'Order status and delivery updates',
                        ),
                      ),
                      const Divider(height: 1),
                      SwitchListTile(
                        value: _marketing,
                        onChanged: (v) => setState(() => _marketing = v),
                        secondary: const Icon(Icons.campaign_outlined),
                        title: const Text('Marketing emails'),
                        subtitle: const Text('Deals and recommendations'),
                      ),
                      const Divider(height: 1),
                      SwitchListTile(
                        value: _darkMode,
                        onChanged: (v) {
                          setState(() => _darkMode = v);
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Theme toggle (demo)'),
                            ),
                          );
                          // TODO: wire to real theme mode via app state
                        },
                        secondary: const Icon(Icons.dark_mode_outlined),
                        title: const Text('Dark mode'),
                        subtitle: const Text('Use device theme or override'),
                      ),
                    ],
                  ),
                );

                final dangerCard = _CardX(
                  child: Column(
                    children: [
                      const _SectionHeader(title: 'Security'),
                      ListTile(
                        leading: const Icon(Icons.lock_reset_outlined),
                        title: const Text('Change Password'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Change password (todo)'),
                            ),
                          );
                        },
                      ),
                      const Divider(height: 1),
                      ListTile(
                        leading: const Icon(Icons.logout),
                        title: const Text('Log out'),
                        onTap: () => _logout(),
                      ),
                    ],
                  ),
                );

                // Responsive layout
                if (!isWide) {
                  return ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      profileCard,
                      const SizedBox(height: 12),
                      shortcutsCard,
                      const SizedBox(height: 12),
                      addressCard,
                      const SizedBox(height: 12),
                      settingsCard,
                      const SizedBox(height: 12),
                      dangerCard,
                    ],
                  );
                }

                // Wide screens → 2-column grid feel
                return Padding(
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // left column
                      Expanded(
                        flex: 7,
                        child: Column(
                          children: [
                            profileCard,
                            const SizedBox(height: 12),
                            shortcutsCard,
                            const SizedBox(height: 12),
                            dangerCard,
                          ],
                        ),
                      ),
                      const SizedBox(width: 16),
                      // right column
                      Expanded(
                        flex: 5,
                        child: Column(
                          children: [
                            addressCard,
                            const SizedBox(height: 12),
                            settingsCard,
                          ],
                        ),
                      ),
                    ],
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }

  void _openEditProfile(user) {
    context.push('/account/edit-profile', extra: user);
  }

  Future<void> _logout() async {
    // Clear stored tokens
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('access_token');
    await prefs.remove('refresh_token');
    await prefs.remove('id_token');

    // Navigate to login page
    if (mounted) {
      context.go('/login');
    }
  }
}

/* ----------------- helpers ----------------- */

class _CardX extends StatelessWidget {
  final Widget child;
  const _CardX({required this.child});

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: child,
    );
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;
  const _SectionHeader({required this.title});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
      child: Text(
        title,
        style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800),
      ),
    );
  }
}
