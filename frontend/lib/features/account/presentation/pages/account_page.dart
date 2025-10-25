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

  // --- ธีมสีเขียว ---
  final primaryGreen = Colors.green.shade800;
  final lightGreen = const Color.fromARGB(255, 255, 255, 255);
  final dangerRed = Colors.red.shade700;
  // ---

  @override
  Widget build(BuildContext context) {
    final addr = CheckoutScope.of(context).shippingAddress;
    final textTheme = Theme.of(context).textTheme;

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Account',
          style: TextStyle(
            color: primaryGreen, // ใช้สีเขียวที่ AppBar Title
            fontWeight: FontWeight.w800,
          ),
        ),
        backgroundColor:
            Colors.white, // หรือ Theme.of(context).scaffoldBackgroundColor
        elevation: 0,
        scrolledUnderElevation: 0,
      ),
      body: LayoutBuilder(
        builder: (context, c) {
          final isWide = c.maxWidth >= 1000;

          final profileCard = _CardX(
            child: Padding(
              padding: const EdgeInsets.all(20), // เพิ่ม padding
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  CircleAvatar(
                    radius: 40, // ขยายให้ใหญ่ขึ้น
                    backgroundColor: lightGreen,
                    child: Text(
                      _initials(_name),
                      style: TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.w800,
                        color: primaryGreen,
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
                          style: textTheme.headlineSmall?.copyWith(
                            fontWeight: FontWeight.w700,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(_email, style: TextStyle(color: Colors.grey[700])),
                        const SizedBox(height: 4),
                        Text(_phone, style: TextStyle(color: Colors.grey[700])),
                        const SizedBox(height: 16),
                        Wrap(
                          spacing: 8,
                          runSpacing: 8,
                          children: [
                            OutlinedButton.icon(
                              style: OutlinedButton.styleFrom(
                                foregroundColor: primaryGreen,
                                side: BorderSide(color: primaryGreen),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              icon: const Icon(Icons.edit_outlined, size: 18),
                              label: const Text('Edit Profile'),
                              onPressed: () => _openEditProfile(),
                            ),
                            OutlinedButton.icon(
                              style: OutlinedButton.styleFrom(
                                foregroundColor: primaryGreen,
                                side: BorderSide(color: primaryGreen),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              icon: const Icon(
                                Icons.location_on_outlined,
                                size: 18,
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
                _SectionHeader(title: 'Shortcuts', color: primaryGreen),
                ListTile(
                  leading: _IconWrapper(
                    icon: Icons.receipt_long_outlined,
                    color: primaryGreen,
                  ),
                  title: const Text('My Orders'),
                  subtitle: const Text('See order history & delivery status'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () =>
                      context.go('/progress'), // Orders tab lives there
                ),
                const Divider(height: 1, indent: 20, endIndent: 20),
                ListTile(
                  leading: _IconWrapper(
                    icon: Icons.eco_outlined,
                    color: primaryGreen,
                  ),
                  title: const Text('My Subscriptions'),
                  subtitle: const Text('Manage plan renewals and pauses'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.go('/progress'), // Subscriptions tab
                ),
                const Divider(height: 1, indent: 20, endIndent: 20),
                ListTile(
                  leading: _IconWrapper(
                    icon: Icons.payment_outlined,
                    color: primaryGreen,
                  ),
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
                  _SectionHeader(
                    title: 'Default Shipping Address',
                    color: primaryGreen,
                  ),
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
                            style: const TextStyle(height: 1.4, fontSize: 15),
                          ),
                  ),
                  const SizedBox(height: 12),
                  Align(
                    alignment: Alignment.centerLeft,
                    child: OutlinedButton.icon(
                      style: OutlinedButton.styleFrom(
                        foregroundColor: primaryGreen,
                        side: BorderSide(color: primaryGreen),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      onPressed: () => context.go('/cart/checkout/address'),
                      icon: const Icon(
                        Icons.edit_location_alt_outlined,
                        size: 18,
                      ),
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
                _SectionHeader(title: 'Preferences', color: primaryGreen),
                SwitchListTile(
                  value: _notif,
                  onChanged: (v) => setState(() => _notif = v),
                  activeColor: primaryGreen, // ใช้สีเขียว
                  secondary: _IconWrapper(
                    icon: Icons.notifications_outlined,
                    color: primaryGreen,
                  ),
                  title: const Text('Push notifications'),
                  subtitle: const Text('Order status and delivery updates'),
                ),
                const Divider(height: 1, indent: 20, endIndent: 20),
                SwitchListTile(
                  value: _marketing,
                  onChanged: (v) => setState(() => _marketing = v),
                  activeColor: primaryGreen, // ใช้สีเขียว
                  secondary: _IconWrapper(
                    icon: Icons.campaign_outlined,
                    color: primaryGreen,
                  ),
                  title: const Text('Marketing emails'),
                  subtitle: const Text('Deals and recommendations'),
                ),
                const Divider(height: 1, indent: 20, endIndent: 20),
                SwitchListTile(
                  value: _darkMode,
                  onChanged: (v) {
                    setState(() => _darkMode = v);
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Theme toggle (demo)')),
                    );
                    // TODO: wire to real theme mode via app state
                  },
                  activeColor: primaryGreen, // ใช้สีเขียว
                  secondary: _IconWrapper(
                    icon: Icons.dark_mode_outlined,
                    color: primaryGreen,
                  ),
                  title: const Text('Dark mode'),
                  subtitle: const Text('Use device theme or override'),
                ),
              ],
            ),
          );

          final dangerCard = _CardX(
            child: Column(
              children: [
                _SectionHeader(title: 'Security', color: primaryGreen),
                ListTile(
                  leading: _IconWrapper(
                    icon: Icons.lock_reset_outlined,
                    color: Colors.orange.shade700, // สีส้ม_เตือน
                  ),
                  title: const Text('Change Password'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Change password (todo)')),
                    );
                  },
                ),
                const Divider(height: 1, indent: 20, endIndent: 20),
                ListTile(
                  leading: _IconWrapper(
                    icon: Icons.logout,
                    color: dangerRed, // สีแดง_อันตราย
                  ),
                  title: Text('Log out', style: TextStyle(color: dangerRed)),
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
                const SizedBox(height: 16),
                shortcutsCard,
                const SizedBox(height: 16),
                addressCard,
                const SizedBox(height: 16),
                settingsCard,
                const SizedBox(height: 16),
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
                  child: ListView(
                    // ใช้ ListView เพื่อให้เลื่อนได้ถ้าจอล้น
                    children: [
                      profileCard,
                      const SizedBox(height: 16),
                      shortcutsCard,
                      const SizedBox(height: 16),
                      dangerCard,
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                // right column
                Expanded(
                  flex: 5,
                  child: ListView(
                    // ใช้ ListView เพื่อให้เลื่อนได้ถ้าจอล้น
                    children: [
                      addressCard,
                      const SizedBox(height: 16),
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
      backgroundColor: Colors.white, // <-- 1. กำหนดพื้นหลังเป็นสีขาว
      elevation: 4, // <-- 2. เพิ่มเงา (เหมือนกับการ์ด)
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) {
        return Container(
          padding: EdgeInsets.only(
            left: 20,
            right: 20,
            top: 8,
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 20,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                'Edit Profile',
                style: Theme.of(
                  context,
                ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.w800),
              ),
              const SizedBox(height: 20),
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
              const SizedBox(height: 24),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                      ),
                      onPressed: () => Navigator.pop(ctx),
                      child: const Text('Cancel'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: FilledButton(
                      style: FilledButton.styleFrom(
                        backgroundColor: primaryGreen, // ใช้สีเขียว
                        padding: const EdgeInsets.symmetric(vertical: 12),
                      ),
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
    return Container(
      // เปลี่ยนจาก Card เป็น Container เพื่อควบคุม shadow ได้ง่ายขึ้น
      decoration: BoxDecoration(
        color: Colors.white, // พื้นหลังสีขาว
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.08), // สีเงาจางๆ
            spreadRadius: 0,
            blurRadius: 10,
            offset: const Offset(0, 4), // เลื่อนเงาลงมาเล็กน้อย
          ),
        ],
      ),
      clipBehavior: Clip.antiAlias, // เพื่อให้ child ไม่ล้นขอบโค้ง
      child: child,
    );
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;
  final Widget? action;
  final Color color;
  const _SectionHeader({
    required this.title,
    this.action,
    this.color = Colors.black87,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Row(
        children: [
          Text(
            title.toUpperCase(), // ใช้อักษรตัวใหญ่
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: color, // ใช้สีที่ส่งเข้ามา (สีเขียว)
              letterSpacing: 0.5, // เพิ่มระยะห่างตัวอักษร
            ),
          ),
          const Spacer(),
          if (action != null) action!,
        ],
      ),
    );
  }
}

// Helper ใหม่สำหรับหุ้มไอคอน
class _IconWrapper extends StatelessWidget {
  final IconData icon;
  final Color color;
  const _IconWrapper({required this.icon, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1), // สีพื้นหลังอ่อนๆ
        shape: BoxShape.circle,
      ),
      child: Icon(icon, size: 20, color: color), // ไอคอนสีเข้ม
    );
  }
}
