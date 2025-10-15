import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/state/checkout_controller.dart';

class AccountPage extends StatefulWidget {
  const AccountPage({super.key});
  @override
  State<AccountPage> createState() => _AccountPageState();
}

class _AccountPageState extends State<AccountPage> {
  // Mock user data (swap with your auth profile later)
  String _name = 'FreshEase User';
  String _email = 'user@example.com';
  String _phone = '+66 80 123 4567';
  bool _notif = true;
  bool _marketing = false;
  bool _darkMode = false; // demo toggle (wire to real theme later)

  @override
  Widget build(BuildContext context) {
    final addr = CheckoutScope.of(context).shippingAddress;

    return Scaffold(
      appBar: AppBar(title: const Text('Account')),
      body: LayoutBuilder(
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
                    child: Text(
                      _initials(_name),
                      style: const TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          _name,
                          style: const TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.w800,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(_email, style: TextStyle(color: Colors.grey[700])),
                        const SizedBox(height: 4),
                        Text(_phone, style: TextStyle(color: Colors.grey[700])),
                        const SizedBox(height: 12),
                        Wrap(
                          spacing: 8,
                          runSpacing: 8,
                          children: [
                            OutlinedButton.icon(
                              icon: const Icon(Icons.edit),
                              label: const Text('Edit Profile'),
                              onPressed: () => _openEditProfile(),
                            ),
                            OutlinedButton.icon(
                              icon: const Icon(Icons.location_on_outlined),
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
                  leading: const Icon(Icons.receipt_long_outlined),
                  title: const Text('My Orders'),
                  subtitle: const Text('See order history & delivery status'),
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
                  onTap: () => context.go('/progress'), // Subscriptions tab
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
                      onPressed: () => context.go('/cart/checkout/address'),
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
                  subtitle: const Text('Order status and delivery updates'),
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
                      const SnackBar(content: Text('Theme toggle (demo)')),
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
                      const SnackBar(content: Text('Change password (todo)')),
                    );
                  },
                ),
                const Divider(height: 1),
                ListTile(
                  leading: const Icon(Icons.logout),
                  title: const Text('Log out'),
                  onTap: () => context.go('/login'),
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
      ),
    );
  }

  void _openEditProfile() {
    final nameCtrl = TextEditingController(text: _name);
    final emailCtrl = TextEditingController(text: _email);
    final phoneCtrl = TextEditingController(text: _phone);

    showModalBottomSheet(
      context: context,
      useSafeArea: true,
      isScrollControlled: true,
      showDragHandle: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) {
        return Padding(
          padding: EdgeInsets.only(
            left: 16,
            right: 16,
            top: 8,
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 16,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text(
                'Edit Profile',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: nameCtrl,
                decoration: const InputDecoration(labelText: 'Name'),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: emailCtrl,
                decoration: const InputDecoration(labelText: 'Email'),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: phoneCtrl,
                decoration: const InputDecoration(labelText: 'Phone'),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () => Navigator.pop(ctx),
                      child: const Text('Cancel'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: FilledButton(
                      onPressed: () {
                        setState(() {
                          _name = nameCtrl.text.trim().isEmpty
                              ? _name
                              : nameCtrl.text.trim();
                          _email = emailCtrl.text.trim().isEmpty
                              ? _email
                              : emailCtrl.text.trim();
                          _phone = phoneCtrl.text.trim().isEmpty
                              ? _phone
                              : phoneCtrl.text.trim();
                        });
                        Navigator.pop(ctx);
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text('Profile updated')),
                        );
                      },
                      child: const Text('Save'),
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      },
    );
  }

  String _initials(String name) {
    final parts = name.trim().split(RegExp(r'\s+'));
    if (parts.isEmpty) return 'U';
    if (parts.length == 1) return parts.first.substring(0, 1).toUpperCase();
    return (parts[0].substring(0, 1) + parts[1].substring(0, 1)).toUpperCase();
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
  final Widget? action;
  const _SectionHeader({required this.title, this.action});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
      child: Row(
        children: [
          Text(
            title,
            style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800),
          ),
          const Spacer(),
          if (action != null) action!,
        ],
      ),
    );
  }
}
