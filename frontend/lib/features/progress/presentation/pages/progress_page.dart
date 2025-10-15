import 'dart:math';
import 'package:flutter/material.dart';
import 'package:frontend/core/health/mock_health_repo.dart';
import 'package:go_router/go_router.dart';

class ProgressPage extends StatefulWidget {
  const ProgressPage({super.key});

  @override
  State<ProgressPage> createState() => _ProgressPageState();
}

class _ProgressPageState extends State<ProgressPage> {
  // --- Mock stats (replace from API later) ---
  final _kpis = const [
    _Kpi(label: 'Orders', value: '12', icon: Icons.shopping_bag_outlined),
    _Kpi(label: 'On-time', value: '95%', icon: Icons.timelapse_outlined),
    _Kpi(label: 'Avg Basket', value: '฿368', icon: Icons.payments_outlined),
  ];

  final _sparkData = const [7, 8, 6, 9, 12, 10, 13, 12, 14, 16, 15, 18];

  final _habits = const [
    _Habit(label: 'Veggie servings / day', value: 0.75),
    _Habit(label: 'Fruit servings / day', value: 0.60),
    _Habit(label: 'Cook-at-home ratio', value: 0.55),
  ];

  final _activities = const [
    _Activity('Order delivered', 'Balanced Weekly Plan · ฿899', 'Today, 10:15'),
    _Activity('Checkout completed', 'Fresh Starter Plan', 'Sun, 16:40'),
    _Activity('Added to cart', 'Avocado Set ×2', 'Sat, 18:22'),
    _Activity('Order delivered', 'Veggies bundle', 'Fri, 12:05'),
  ];

  final int _tabIndex = 0;

  // ---- Calories mock integration ----
  final IHealthRepo _healthRepo = MockHealthRepo();
  CalorieSnapshot? _calSnapshot;
  bool _syncingCal = false;

  Future<void> _syncCalories() async {
    setState(() => _syncingCal = true);
    try {
      final snap = await _healthRepo.readToday();
      setState(() => _calSnapshot = snap);
    } finally {
      setState(() => _syncingCal = false);
    }
  }

  @override
  void initState() {
    super.initState();
    _syncCalories(); // initial mock load
  }

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 3,
      initialIndex: _tabIndex,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Progress'),
          bottom: const TabBar(
            isScrollable: true,
            tabs: [
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
              // NEW props
              calSnapshot: _calSnapshot,
              syncingCal: _syncingCal,
              onSyncCalories: _syncCalories,
            ),
            _OrdersTab(
              activities: _activities
                  .where((a) => a.title.contains('Order'))
                  .toList(),
            ),
            const _SubscriptionsTab(),
          ],
        ),
      ),
    );
  }
}

class _OverviewTab extends StatelessWidget {
  final List<_Kpi> kpis;
  final List<int> sparkData;
  final List<_Habit> habits;
  final List<_Activity> activities;

  // NEW: calories props
  final CalorieSnapshot? calSnapshot;
  final bool syncingCal;
  final VoidCallback onSyncCalories;

  const _OverviewTab({
    required this.kpis,
    required this.sparkData,
    required this.habits,
    required this.activities,
    required this.calSnapshot,
    required this.syncingCal,
    required this.onSyncCalories,
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
            // --- Calories (Mock API) ---
            _CaloriesBlock(
              snapshot: calSnapshot,
              syncing: syncingCal,
              onSync: onSyncCalories,
            ),

            const SizedBox(height: 12),

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
              children: [
                _CardX(
                  flex: isWide ? 7 : 1,
                  child: _TrendBlock(data: sparkData),
                ),
                _CardX(
                  flex: isWide ? 5 : 1,
                  child: const _StreakBlock(days: 6),
                ),
              ],
            ),

            const SizedBox(height: 12),

            // --- Habits ---
            _CardX(
              child: Padding(
                padding: const EdgeInsets.all(8),
                child: Column(
                  children: [
                    _SectionHeader(
                      title: 'Healthy Habits',
                      action: TextButton(
                        onPressed: () => context.go('/plans'),
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
                padding: const EdgeInsets.all(8),
                child: Column(
                  children: [
                    const _SectionHeader(title: 'Recent Activity'),
                    const SizedBox(height: 8),
                    for (final a in activities)
                      ListTile(
                        dense: true,
                        leading: const CircleAvatar(
                          child: Icon(Icons.event_note),
                        ),
                        title: Text(
                          a.title,
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        ),
                        subtitle: Text(a.subtitle),
                        trailing: Text(
                          a.time,
                          style: TextStyle(color: Colors.grey[600]),
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

class _OrdersTab extends StatelessWidget {
  final List<_Activity> activities;
  const _OrdersTab({required this.activities});

  @override
  Widget build(BuildContext context) {
    if (activities.isEmpty) {
      return const Center(child: Text('No orders yet'));
    }
    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: activities.length,
      separatorBuilder: (_, __) => const Divider(),
      itemBuilder: (_, i) {
        final a = activities[i];
        return ListTile(
          leading: const Icon(Icons.local_shipping_outlined),
          title: Text(a.subtitle),
          subtitle: Text(a.time),
          trailing: const Text('Delivered'),
        );
      },
    );
  }
}

class _SubscriptionsTab extends StatelessWidget {
  const _SubscriptionsTab();

  @override
  Widget build(BuildContext context) {
    // Demo card – wire with real subscription data later
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _CardX(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Container(
                  decoration: BoxDecoration(
                    color: Colors.teal.withValues(alpha: .1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  padding: const EdgeInsets.all(10),
                  child: const Icon(Icons.eco, color: Colors.teal),
                ),
                const SizedBox(width: 12),
                const Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Balanced Weekly Plan',
                        style: TextStyle(
                          fontWeight: FontWeight.w800,
                          fontSize: 16,
                        ),
                      ),
                      SizedBox(height: 4),
                      Text('Next renewal: 30 Oct 2025'),
                    ],
                  ),
                ),
                FilledButton(
                  onPressed: null, // TODO: manage subscription
                  child: Text('Manage'),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}

/* ========================= WIDGETS ========================= */

class _CardX extends StatelessWidget {
  final Widget child;
  final int flex;
  const _CardX({required this.child, this.flex = 1});

  @override
  Widget build(BuildContext context) {
    return Expanded(
      flex: flex,
      child: Card(
        clipBehavior: Clip.antiAlias,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        child: child,
      ),
    );
  }
}

class _ResponsiveRow extends StatelessWidget {
  final List<Widget> children;
  final double gap;
  const _ResponsiveRow({required this.children, this.gap = 12});

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
              Expanded(child: children[i]),
              if (i != children.length - 1) SizedBox(width: gap),
            ],
          ],
        );
      },
    );
  }
}

class _Kpi {
  final String label;
  final String value;
  final IconData icon;
  const _Kpi({required this.label, required this.value, required this.icon});
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
              color: Theme.of(context).colorScheme.primaryContainer,
              borderRadius: BorderRadius.circular(12),
            ),
            padding: const EdgeInsets.all(10),
            child: Icon(kpi.icon, color: Theme.of(context).colorScheme.primary),
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
      ..color = Colors.teal
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2.5;

    final fill = Paint()
      ..shader = const LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [Color(0x6622AA99), Color(0x1122AA99)],
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
                  color: success ? Colors.green : Colors.grey.shade300,
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

class _Habit {
  final String label;
  final double value; // 0..1
  const _Habit({required this.label, required this.value});
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
              color: Colors.green,
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

class _Activity {
  final String title;
  final String subtitle;
  final String time;
  const _Activity(this.title, this.subtitle, this.time);
}

/* ---------- Calories Block (inline) ---------- */

class _CaloriesBlock extends StatelessWidget {
  final CalorieSnapshot? snapshot;
  final bool syncing;
  final VoidCallback onSync;

  const _CaloriesBlock({
    required this.snapshot,
    required this.syncing,
    required this.onSync,
  });

  @override
  Widget build(BuildContext context) {
    final snap = snapshot;

    return Card(
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(16, 14, 16, 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header + action
            Row(
              children: [
                const Text(
                  'Calories Today',
                  style: TextStyle(fontSize: 16, fontWeight: FontWeight.w800),
                ),
                const Spacer(),
                FilledButton.icon(
                  onPressed: syncing ? null : onSync,
                  icon: syncing
                      ? const SizedBox(
                          width: 14,
                          height: 14,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Icon(Icons.sync),
                  label: Text(syncing ? 'Syncing…' : 'Sync'),
                ),
              ],
            ),
            const SizedBox(height: 12),

            if (snap == null) ...[
              Text(
                'Connect later — using mock numbers.\nTap “Sync” to refresh.',
                style: TextStyle(color: Colors.grey[700]),
              ),
              const SizedBox(height: 8),
            ] else ...[
              _CalRow(label: 'Intake (kcal)', value: snap.intakeKcal),
              const SizedBox(height: 6),
              _CalRow(
                label: 'Active burned (kcal)',
                value: snap.activeBurnKcal,
              ),
              const Divider(height: 18),
              _NetChip(net: snap.netKcal),
              const SizedBox(height: 6),
              Text(
                'Last sync: ${_fmtTime(snap.syncedAt)}',
                style: TextStyle(color: Colors.grey[600], fontSize: 12),
              ),
            ],
          ],
        ),
      ),
    );
  }

  static String _fmtTime(DateTime dt) {
    String two(int n) => n.toString().padLeft(2, '0');
    return '${dt.year}-${two(dt.month)}-${two(dt.day)} ${two(dt.hour)}:${two(dt.minute)}';
  }
}

class _CalRow extends StatelessWidget {
  final String label;
  final double value;
  const _CalRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(label, style: TextStyle(color: Colors.grey[700])),
        const Spacer(),
        Text(
          value.toStringAsFixed(0),
          style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 18),
        ),
      ],
    );
  }
}

class _NetChip extends StatelessWidget {
  final double net;
  const _NetChip({required this.net});

  @override
  Widget build(BuildContext context) {
    final positive = net >= 0;
    final color = positive ? Colors.orange : Colors.teal;
    final sign = positive ? '+' : '−';
    final abs = net.abs();

    return Row(
      children: [
        Text(
          'Net (kcal)',
          style: TextStyle(
            color: Colors.grey[800],
            fontWeight: FontWeight.w700,
          ),
        ),
        const Spacer(),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
          decoration: BoxDecoration(
            color: color.withValues(alpha: .12),
            borderRadius: BorderRadius.circular(999),
          ),
          child: Row(
            children: [
              Icon(
                positive ? Icons.trending_up : Icons.trending_down,
                size: 16,
                color: color,
              ),
              const SizedBox(width: 6),
              Text(
                '$sign${abs.toStringAsFixed(0)}',
                style: TextStyle(color: color, fontWeight: FontWeight.w800),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
