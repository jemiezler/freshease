import 'dart:math';

import 'package:flutter/material.dart';
// import 'package:go_router/go_router.dart'; // <-- Changed: ไม่ได้ใช้ GoRouter ในหน้านี้แล้ว

class ProgressPage extends StatefulWidget {
  const ProgressPage({super.key});

  @override
  State<ProgressPage> createState() => _ProgressPageState();
}

// ----------------- MODELS -----------------

class _Kpi {
  final String label;
  final String value;
  final IconData icon;
  const _Kpi({required this.label, required this.value, required this.icon});
}

class _Habit {
  final String label;
  final double value; // 0..1
  const _Habit({required this.label, required this.value});
}

class _Activity {
  final String title;
  final String subtitle;
  final String time;
  const _Activity(this.title, this.subtitle, this.time);
}

// <-- NEW MODEL -->
class _OrderItem {
  final String name;
  final String description;
  final double price;
  final String imageUrl; // ใช้ asset path หรือ network url
  const _OrderItem({
    required this.name,
    required this.description,
    required this.price,
    required this.imageUrl,
  });
}

class _Order {
  final String id;
  final String title;
  final String date;
  final IconData icon;
  final List<_OrderItem> items; // <-- Added: for summary dialog
  final Map<String, double> summary; // <-- Added: for summary dialog

  const _Order({
    required this.id,
    required this.title,
    required this.date,
    required this.icon,
    required this.items,
    required this.summary,
  });
}

class _Subscription {
  final String id;
  final String planName;
  final String renewalDate;
  final IconData icon;
  final String status; // <-- Added: 'Active' or 'Expired'
  const _Subscription({
    required this.id,
    required this.planName,
    required this.renewalDate,
    required this.icon,
    required this.status,
  });
}

// ----------------- PAGE STATE -----------------

class _ProgressPageState extends State<ProgressPage> {
  // --- Mock stats (Updated to match Figma) ---
  final _kpis = const [
    _Kpi(label: 'Orders', value: '2', icon: Icons.shopping_basket_outlined),
    _Kpi(label: 'On-time', value: '95', icon: Icons.schedule_outlined),
    _Kpi(
      label: 'Avg Basket',
      value: '\$300',
      icon: Icons.account_balance_wallet_outlined,
    ),
  ];

  final _sparkData = const [7, 8, 6, 9, 12, 10, 13, 12, 14, 16, 15, 18];

  final _habits = const [
    _Habit(label: 'Veggie servings / day', value: 0.75),
    _Habit(label: 'Fruit servings / day', value: 0.60),
    _Habit(label: 'Cook-at-home ratio', value: 0.55),
  ];

  final _activities = const [
    _Activity('Order delivered', 'Balanced Weekly Plan \$15', 'Today, 10:20'),
    _Activity('Checkout completed', 'Fresh Starter Plan', 'Sunday, 10:20'),
    _Activity('Added to card', 'Avocado Set x2', 'Sat, 10:20'),
    _Activity('Order delivered', 'Veggies bundle', 'Fri, 10:20'),
  ];

  // <-- Changed: Mock data for Orders Tab (Added items and summary) -->
  final _orders = const [
    _Order(
      id: 'ord_123',
      title: 'Balanced Weekly Plan \$200',
      date: 'Today, 10:15',
      icon: Icons.delivery_dining_outlined,
      items: [
        _OrderItem(
          name: 'Bell Pepper Red',
          description: '1kg, Price',
          price: 4.99,
          imageUrl:
              'https://i.pinimg.com/736x/e9/b4/41/e9b441e5fcea11048e9930a46269d5c0.jpg', // Placeholder
        ),
        _OrderItem(
          name: 'Ginger',
          description: '250gm, Price',
          price: 2.99,
          imageUrl:
              'https://i.pinimg.com/736x/13/94/2a/13942ac5cfffc5eb5589c98f4a65fc33.jpg', // Placeholder
        ),
        _OrderItem(
          name: 'Organic Bananas',
          description: '12kg, Price',
          price: 3.00,
          imageUrl:
              'https://i.pinimg.com/736x/02/49/5f/02495fb1b8bd32a24fb8eb483a18a074.jpg', // Placeholder
        ),
      ],
      summary: {
        'Subtotal': 10.00,
        'Delivery fee': 10.00,
        'Total': 10.00, // Figma data is inconsistent, using as-is
      },
    ),
    _Order(
      id: 'ord_456',
      title: 'Veggies bundle',
      date: 'Friday, 10:15',
      icon: Icons.delivery_dining_outlined,
      items: [
        _OrderItem(
          name: 'Veggies Bundle',
          description: '1 set',
          price: 15.00,
          imageUrl:
              'https://i.pinimg.com/736x/4e/ff/d1/4effd19380002ba50643b3dd6be0be8f.jpg', // Placeholder
        ),
      ],
      summary: {'Subtotal': 15.00, 'Delivery fee': 5.00, 'Total': 20.00},
    ),
  ];

  // <-- Changed: Mock data for Subscriptions Tab (Added status) -->
  final _subscriptions = const [
    _Subscription(
      id: 'sub_fresh_77',
      planName: 'Fresh Starter Plan',
      renewalDate: '30 Oct 2025', // Changed to just date
      icon: Icons.ramen_dining_outlined,
      status: 'Expired',
    ),
    _Subscription(
      id: 'sub_balanced_88',
      planName: 'Balanced weekly Plan',
      renewalDate: '20 Dec 2025', // Changed to just date
      icon: Icons.ramen_dining_outlined,
      status: 'Active',
    ),
  ];

  final int _tabIndex = 0;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final primaryColor = Colors.green[700];

    return DefaultTabController(
      length: 3,
      initialIndex: _tabIndex,
      child: Scaffold(
        backgroundColor: const Color(0xFFF9F9F9), // Figma background color
        appBar: AppBar(
          title: const Text(
            'Progress',
            style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold),
          ),
          elevation: 0,
          backgroundColor: Colors.transparent,
          foregroundColor: Colors.black87,
          bottom: TabBar(
            isScrollable: true,
            indicatorColor: primaryColor,
            labelColor: primaryColor,
            labelStyle: const TextStyle(
              fontWeight: FontWeight.w600,
              fontSize: 16,
            ),
            unselectedLabelColor: Colors.grey[600],
            dividerColor: Colors.grey[300],
            tabs: const [
              Tab(text: 'Overview'),
              Tab(text: 'Orders'),
              Tab(text: 'Subscriptions'),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            _OverviewTab(
              kpis: _kpis,
              sparkData: _sparkData,
              habits: _habits,
              activities: _activities,
            ),
            _OrdersTab(orders: _orders),
            _SubscriptionsTab(subscriptions: _subscriptions),
          ],
        ),
      ),
    );
  }
}

// ----------------- OVERVIEW TAB -----------------
// (This widget is unchanged)
class _OverviewTab extends StatelessWidget {
  final List<_Kpi> kpis;
  final List<int> sparkData;
  final List<_Habit> habits;
  final List<_Activity> activities;

  const _OverviewTab({
    required this.kpis,
    required this.sparkData,
    required this.habits,
    required this.activities,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, c) {
        final isWide = c.maxWidth >= 1000;
        final isTablet = c.maxWidth >= 700 && c.maxWidth < 1000;

        return ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // --- KPIs ---
            _ResponsiveRow(
              gap: 12,
              children: kpis
                  .map((k) => _CardX(child: _KpiTile(kpi: k)))
                  .toList(),
            ),

            const SizedBox(height: 12),

            // --- Trend + Streak ---
            _ResponsiveRow(
              gap: 12,
              flexValues: [isWide ? 7 : 1, isWide ? 5 : 1],
              children: [
                _CardX(child: _TrendBlock(data: sparkData)),
                _CardX(child: const _StreakBlock(days: 8)),
              ],
            ),

            const SizedBox(height: 12),

            // --- Habits ---
            _CardX(
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  children: [
                    _SectionHeader(
                      title: 'Healthy Habits',
                      action: TextButton(
                        onPressed: () {
                          // TODO: Handle navigation
                          // context.go('/plans')
                        },
                        child: const Text('Improve with a plan'),
                      ),
                    ),
                    const SizedBox(height: 8),
                    if (isWide || isTablet)
                      Row(
                        children: habits
                            .map((h) => Expanded(child: _HabitTile(h: h)))
                            .toList(),
                      )
                    else
                      Column(
                        children: habits
                            .map(
                              (h) => Padding(
                                padding: const EdgeInsets.only(bottom: 12),
                                child: _HabitTile(h: h),
                              ),
                            )
                            .toList(),
                      ),
                  ],
                ),
              ),
            ),

            const SizedBox(height: 12),

            // --- Recent activity ---
            _CardX(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(12, 12, 12, 4),
                child: Column(
                  children: [
                    const _SectionHeader(title: 'Recent Activity'),
                    const SizedBox(height: 8),
                    for (final a in activities)
                      ListTile(
                        dense: true,
                        leading: CircleAvatar(
                          backgroundColor: Colors.green[50],
                          child: Icon(
                            Icons.ramen_dining_outlined, // Figma bowl icon
                            color: Colors.green[800],
                          ),
                        ),
                        title: Text(
                          a.title,
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        ),
                        subtitle: Text(a.subtitle),
                        trailing: Text(
                          a.time,
                          style: TextStyle(
                            color: Colors.grey[600],
                            fontSize: 12,
                          ),
                        ),
                      ),
                  ],
                ),
              ),
            ),
          ],
        );
      },
    );
  }
}

// ----------------- ORDERS TAB -----------------

class _OrdersTab extends StatelessWidget {
  final List<_Order> orders;
  const _OrdersTab({required this.orders});

  @override
  Widget build(BuildContext context) {
    if (orders.isEmpty) {
      return const Center(child: Text('No orders yet'));
    }
    return ListView.separated(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      itemCount: orders.length,
      separatorBuilder: (_, __) => const Divider(height: 1),
      itemBuilder: (_, i) {
        final order = orders[i];
        return ListTile(
          contentPadding: const EdgeInsets.symmetric(vertical: 8),
          leading: Icon(order.icon, color: Colors.green[700], size: 36),
          title: Text(
            order.title,
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
          subtitle: Text(order.date),
          trailing: Text(
            'Delivered',
            style: TextStyle(color: Colors.grey[700]),
          ),
          onTap: () {
            // <-- Changed
            // Show the order summary dialog
            _showModalPopup(context, _OrderSummaryDialog(order: order));
          },
          hoverColor: Colors.grey[100],
          splashColor: Colors.green[100],
        );
      },
    );
  }
}

// ----------------- SUBSCRIPTIONS TAB -----------------

class _SubscriptionsTab extends StatelessWidget {
  final List<_Subscription> subscriptions;
  const _SubscriptionsTab({required this.subscriptions});

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: subscriptions.length,
      separatorBuilder: (_, __) => const SizedBox(height: 12),
      itemBuilder: (_, i) {
        // <-- Changed: Pass full subscription object
        return _SubscriptionCard(sub: subscriptions[i]);
      },
    );
  }
}

class _SubscriptionCard extends StatelessWidget {
  final _Subscription sub;
  const _SubscriptionCard({required this.sub});

  @override
  Widget build(BuildContext context) {
    // <-- Changed: Use _ManageSubscriptionDialog
    final dialog = _ManageSubscriptionDialog(subscription: sub);

    return _CardX(
      child: InkWell(
        onTap: () {
          // <-- Changed
          _showModalPopup(context, dialog);
        },
        borderRadius: BorderRadius.circular(12),
        splashColor: Colors.green[100],
        hoverColor: Colors.grey[50],
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                decoration: BoxDecoration(
                  color: Colors.green[50],
                  borderRadius: BorderRadius.circular(12),
                ),
                padding: const EdgeInsets.all(10),
                child: Icon(sub.icon, color: Colors.green[700]),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      sub.planName,
                      style: const TextStyle(
                        fontWeight: FontWeight.w800,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Next renewal: ${sub.renewalDate}', // <-- Changed
                      style: TextStyle(color: Colors.grey[600]),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              FilledButton(
                onPressed: () {
                  // <-- Changed
                  _showModalPopup(context, dialog);
                },
                style: FilledButton.styleFrom(
                  backgroundColor: Colors.green[600],
                ),
                child: const Text('Manage'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/* ========================= NEW DIALOG WIDGETS ========================= */

// <-- NEW WIDGET: Helper function for modal pop-up -->
Future<void> _showModalPopup(BuildContext context, Widget child) {
  return showGeneralDialog(
    context: context,
    barrierDismissible: true,
    barrierLabel: MaterialLocalizations.of(context).modalBarrierDismissLabel,
    barrierColor: Colors.black.withOpacity(0.6), // Figma background dim
    transitionDuration: const Duration(milliseconds: 200),
    pageBuilder: (context, anim1, anim2) {
      return ScaleTransition(
        scale: CurvedAnimation(parent: anim1, curve: Curves.easeOutCubic),
        child: FadeTransition(
          opacity: anim1,
          child: Dialog(
            backgroundColor: Colors.transparent,
            elevation: 0,
            child: child,
          ),
        ),
      );
    },
  );
}

// <-- NEW WIDGET: Based on image_978f6a.png -->
class _OrderSummaryDialog extends StatelessWidget {
  final _Order order;
  const _OrderSummaryDialog({required this.order});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 400, // Max width for the dialog
      decoration: BoxDecoration(
        color: const Color(0xFFEFEFEF), // Figma card color
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // --- Header ---
          const Padding(
            padding: EdgeInsets.all(16.0),
            child: Text(
              'Order Summary',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
          ),
          // --- Item List ---
          ...order.items.map((item) {
            return Column(
              children: [
                ListTile(
                  // <-- Changed: Using Image.network instead of Image.asset -->
                  leading: Image.network(
                    item.imageUrl, // Use the URL from the _OrderItem
                    width: 50,
                    height: 50,
                    fit: BoxFit.cover, // Ensures the image fills the space
                    errorBuilder: (c, e, s) => Container(
                      width: 50,
                      height: 50,
                      color: Colors.grey[300],
                      child: const Icon(Icons.broken_image, color: Colors.grey),
                    ),
                  ),
                  title: Text(item.name),
                  subtitle: Text(item.description),
                  trailing: Text(
                    '\$${item.price.toStringAsFixed(2)}',
                    style: const TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 15,
                    ),
                  ),
                ),
                const Divider(height: 1, indent: 16, endIndent: 16),
              ],
            );
          }).toList(), // Add .toList() here
          // --- Summary ---
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                _buildSummaryRow('Subtotal', order.summary['Subtotal'] ?? 0),
                const SizedBox(height: 8),
                _buildSummaryRow(
                  'Delivery fee',
                  order.summary['Delivery fee'] ?? 0,
                ),
                const SizedBox(height: 8),
                _buildSummaryRow(
                  'Total',
                  order.summary['Total'] ?? 0,
                  isTotal: true,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSummaryRow(String label, double value, {bool isTotal = false}) {
    final style = TextStyle(
      fontWeight: isTotal ? FontWeight.bold : FontWeight.normal,
      fontSize: isTotal ? 16 : 14,
    );
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: style),
        Text('\$${value.toStringAsFixed(2)}', style: style),
      ],
    );
  }
}

// <-- NEW WIDGET: Based on image_97e121.png -->
class _ManageSubscriptionDialog extends StatelessWidget {
  final _Subscription subscription;
  const _ManageSubscriptionDialog({required this.subscription});

  @override
  Widget build(BuildContext context) {
    final bool isExpired = subscription.status == 'Expired';
    return Container(
      width: 400, // Max width
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
      ),
      padding: const EdgeInsets.all(20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(subscription.icon, color: Colors.green[700], size: 40),
          const SizedBox(height: 12),
          Text(
            subscription.planName,
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 20),
          _buildInfoRow('plan', subscription.planName),
          const Divider(),
          _buildInfoRow(
            'Status',
            subscription.status,
            valueColor: isExpired ? Colors.red : Colors.green,
          ),
          const Divider(),
          _buildInfoRow(
            isExpired ? 'Expires' : 'Renews',
            subscription.renewalDate,
          ),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: FilledButton(
              onPressed: () {
                // TODO: Handle action
                Navigator.of(context).pop(); // Close dialog
              },
              style: FilledButton.styleFrom(
                backgroundColor: isExpired
                    ? Colors.green[600]
                    : Colors.red[700],
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
              child: Text(isExpired ? 'Renew' : 'Cancel'),
            ),
          ),
        ],
      ),
    );
  }

  // --- vvv THIS FUNCTION IS FIXED vvv ---
  Widget _buildInfoRow(String label, String value, {Color? valueColor}) {
    // <-- Changed: Capitalize the string manually -->
    final String capitalizedLabel = label.isEmpty
        ? ''
        : '${label[0].toUpperCase()}${label.substring(1)}';

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            capitalizedLabel, // <-- Use the capitalized string
            style: TextStyle(color: Colors.grey[600]),
            // 'textTransform' property removed
          ),
          Text(
            value,
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: valueColor ?? Colors.black87,
            ),
          ),
        ],
      ),
    );
  }
}

/* ========================= UNCHANGED WIDGETS ========================= */

class _CardX extends StatelessWidget {
  final Widget child;
  const _CardX({required this.child});

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 0,
      color: Colors.white,
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: child,
    );
  }
}

class _ResponsiveRow extends StatelessWidget {
  final List<Widget> children;
  final List<int>? flexValues;
  final double gap;
  const _ResponsiveRow({
    required this.children,
    this.flexValues,
    this.gap = 12,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, c) {
        final isStack = c.maxWidth < 700; // stack on narrow widths
        if (isStack) {
          return Column(
            children: [
              for (int i = 0; i < children.length; i++) ...[
                SizedBox(width: double.infinity, child: children[i]),
                if (i != children.length - 1) SizedBox(height: gap),
              ],
            ],
          );
        }
        return Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            for (int i = 0; i < children.length; i++) ...[
              Expanded(
                flex: flexValues != null ? flexValues![i] : 1,
                child: children[i],
              ),
              if (i != children.length - 1) SizedBox(width: gap),
            ],
          ],
        );
      },
    );
  }
}

class _KpiTile extends StatelessWidget {
  final _Kpi kpi;
  const _KpiTile({required this.kpi});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(14),
      child: Row(
        children: [
          Container(
            decoration: BoxDecoration(
              color: Colors.orange[50],
              borderRadius: BorderRadius.circular(12),
            ),
            padding: const EdgeInsets.all(10),
            child: Icon(kpi.icon, color: Colors.orange[800]),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(kpi.label, style: TextStyle(color: Colors.grey[600])),
              const SizedBox(height: 4),
              Text(
                kpi.value,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w800,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _TrendBlock extends StatelessWidget {
  final List<int> data;
  const _TrendBlock({required this.data});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.fromLTRB(14, 14, 14, 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const _SectionHeader(title: 'Weekly Orders Trend'),
          const SizedBox(height: 12),
          SizedBox(
            height: 120,
            child: CustomPaint(
              painter: _SparklinePainter(data),
              child: Container(),
            ),
          ),
          const SizedBox(height: 8),
          Text('Last 12 weeks', style: TextStyle(color: Colors.grey[600])),
        ],
      ),
    );
  }
}

class _SparklinePainter extends CustomPainter {
  final List<int> data;
  _SparklinePainter(this.data);

  @override
  void paint(Canvas canvas, Size size) {
    if (data.isEmpty) return;
    final maxV = data.reduce(max).toDouble();
    final minV = data.reduce(min).toDouble();
    final span = (maxV - minV).clamp(1, double.infinity);

    final path = Path();
    for (int i = 0; i < data.length; i++) {
      final x = i * (size.width / (data.length - 1));
      final yNorm = (data[i] - minV) / span;
      final y = size.height - (yNorm * size.height);
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }

    final stroke = Paint()
      ..color = Colors.green[600]!
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2.5;

    final fill = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          Colors.green[600]!.withOpacity(0.4),
          Colors.green[600]!.withOpacity(0.05),
        ],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height))
      ..style = PaintingStyle.fill;

    // Fill under curve
    final fillPath = Path.from(path)
      ..lineTo(size.width, size.height)
      ..lineTo(0, size.height)
      ..close();

    canvas.drawPath(fillPath, fill);
    canvas.drawPath(path, stroke);
  }

  @override
  bool shouldRepaint(covariant _SparklinePainter oldDelegate) =>
      oldDelegate.data != data;
}

class _StreakBlock extends StatelessWidget {
  final int days;
  const _StreakBlock({required this.days});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(14),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const _SectionHeader(title: 'Healthy Eating Streak'),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: List.generate(14, (i) {
              final success = i < days;
              return Container(
                width: 22,
                height: 22,
                decoration: BoxDecoration(
                  color: success ? Colors.green[600] : Colors.grey.shade200,
                  shape: BoxShape.circle,
                ),
              );
            }),
          ),
          const SizedBox(height: 12),
          Text(
            '$days-day streak',
            style: const TextStyle(fontWeight: FontWeight.w700),
          ),
        ],
      ),
    );
  }
}

class _HabitTile extends StatelessWidget {
  final _Habit h;
  const _HabitTile({required this.h});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(h.label, style: const TextStyle(fontWeight: FontWeight.w700)),
          const SizedBox(height: 8),
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: LinearProgressIndicator(
              value: h.value,
              minHeight: 10,
              backgroundColor: Colors.grey.shade200,
              color: Colors.green[600],
            ),
          ),
          const SizedBox(height: 6),
          Text(
            '${(h.value * 100).toStringAsFixed(0)}%',
            style: TextStyle(color: Colors.grey[600]),
          ),
        ],
      ),
    );
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;
  final Widget? action;
  const _SectionHeader({required this.title, this.action});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(
          title,
          style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800),
        ),
        const Spacer(),
        if (action != null) action!,
      ],
    );
  }
}
