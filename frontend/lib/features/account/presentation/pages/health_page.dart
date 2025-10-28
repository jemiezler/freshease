import 'dart:io';
import 'package:flutter/material.dart';
import 'package:frontend/app/di.dart';
import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/core/widgets/global_appbar.dart';

class HealthPage extends StatefulWidget {
  const HealthPage({super.key});
  @override
  State<HealthPage> createState() => _HealthPageState();
}

class _HealthPageState extends State<HealthPage> {
  late final HealthController hc;

  @override
  void initState() {
    super.initState();
    hc = getIt<HealthController>();
    // Ask for permissions, then fetch & compute KPIs
    hc.authorize().then((_) => hc.fetchAll24hAndComputeKpis());
  }

  Future<void> _refresh() => hc.fetchAll24hAndComputeKpis();

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: hc,
      builder: (_, _) {
        return Scaffold(
          appBar: GlobalAppBar(title: 'Health'),
          // appBar: AppBar(
          //   title: const Text('Health'),
          //   actions: [
          //     if (Platform.isAndroid && hc.hcStatus != null)
          //       Padding(
          //         padding: const EdgeInsets.symmetric(horizontal: 12),
          //         child: Center(
          //           child: Text('', style: theme.textTheme.labelSmall),
          //         ),
          //       ),
          //   ],
          // ),
          body: RefreshIndicator(
            onRefresh: _refresh,
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // KPIs only
                Row(
                  children: [
                    _KpiCard(
                      title: 'Steps Today',
                      value: '${hc.stepsToday}',
                      icon: Icons.directions_walk,
                    ),
                    const SizedBox(width: 12),
                    _KpiCard(
                      title: 'Total Calories',
                      value: _fmtKcal(hc.kcalToday),
                      icon: Icons.local_fire_department,
                    ),
                  ],
                ),

                const SizedBox(height: 16),

                // Simple state hints
                if (hc.state == HealthState.noData)
                  _Hint(
                    icon: Icons.info_outline,
                    text: 'No data in the last 24 hours.',
                  ),
                if (hc.state == HealthState.authDenied)
                  _Hint(
                    icon: Icons.lock_outline,
                    text: Platform.isAndroid
                        ? 'Authorization denied. Open Health Connect and grant permissions.'
                        : 'Authorization denied. Enable HealthKit permissions in Settings.',
                  ),
                if (hc.state == HealthState.fetching)
                  const Padding(
                    padding: EdgeInsets.only(top: 12),
                    child: Center(child: CircularProgressIndicator()),
                  ),
              ],
            ),
          ),
          floatingActionButton: FloatingActionButton.extended(
            onPressed: _refresh,
            icon: const Icon(Icons.refresh),
            label: const Text('Refresh'),
          ),
        );
      },
    );
  }

  String _fmtKcal(double v) {
    if (v.isNaN || v.isInfinite) return '0 kcal';
    if (v >= 1000) {
      // display in K for readability, e.g., 1.2K kcal
      return '${(v / 1000).toStringAsFixed(1)}K kcal';
    }
    // no trailing .0 for small values
    final str = v.toStringAsFixed(v.truncateToDouble() == v ? 0 : 1);
    return '$str kcal';
  }
}

class _KpiCard extends StatelessWidget {
  const _KpiCard({
    required this.title,
    required this.value,
    required this.icon,
  });

  final String title;
  final String value;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Expanded(
      child: Container(
        height: 96,
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              blurRadius: 12,
              offset: const Offset(0, 6),
              color: theme.colorScheme.shadow.withOpacity(.06),
            ),
          ],
        ),
        child: Row(
          children: [
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: theme.colorScheme.onPrimaryContainer),
            ),
            const SizedBox(width: 12),
            Flexible(
              child: DefaultTextStyle(
                style: theme.textTheme.bodyMedium!,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Opacity(
                      opacity: .7,
                      child: Text(
                        title,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: theme.textTheme.labelMedium,
                      ),
                    ),
                    Text(
                      value,
                      maxLines: 1,
                      overflow: TextOverflow.fade,
                      style: theme.textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Hint extends StatelessWidget {
  const _Hint({required this.icon, required this.text});
  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      margin: const EdgeInsets.only(top: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Icon(icon, size: 20, color: theme.colorScheme.primary),
          const SizedBox(width: 8),
          Expanded(child: Text(text)),
        ],
      ),
    );
  }
}
