import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:frontend/app/di.dart';
import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/core/widgets/global_appbar.dart';
import 'package:frontend/core/genai/widgets.dart';
import 'package:frontend/core/platform_helper.dart';

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
    // Meal plans are now auto-generated in HealthController.init()
    hc.authorize().then((_) => hc.fetchAll24hAndComputeKpis());
  }

  Future<void> _refresh() => hc.fetchAll24hAndComputeKpis();

  String _formatCacheAge(Duration duration) {
    if (duration.inMinutes < 1) {
      return 'just now';
    } else if (duration.inMinutes < 60) {
      return '${duration.inMinutes}m ago';
    } else if (duration.inHours < 24) {
      return '${duration.inHours}h ago';
    } else {
      return '${duration.inDays}d ago';
    }
  }

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
                    text: kIsWeb
                        ? 'Health data is not available on web platform.'
                        : (PlatformHelper.isAndroid
                              ? 'Authorization denied. Open Health Connect and grant permissions.'
                              : 'Authorization denied. Enable HealthKit permissions in Settings.'),
                  ),
                if (hc.state == HealthState.fetching)
                  const Padding(
                    padding: EdgeInsets.only(top: 12),
                    child: Center(child: CircularProgressIndicator()),
                  ),

                // GenAI Meal Plan Section
                if (hc.state == HealthState.ready)
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SizedBox(height: 24),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            'AI Meal Plan',
                            style: Theme.of(context).textTheme.titleLarge
                                ?.copyWith(fontWeight: FontWeight.bold),
                          ),
                          Row(
                            children: [
                              if (hc.cacheAge != null)
                                Text(
                                  'Updated ${_formatCacheAge(hc.cacheAge!)}',
                                  style: Theme.of(context).textTheme.bodySmall
                                      ?.copyWith(
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.onSurfaceVariant,
                                      ),
                                ),
                              const SizedBox(width: 8),
                              IconButton(
                                onPressed: () => hc.refreshAllPlans(),
                                icon: const Icon(Icons.refresh),
                                tooltip: 'Refresh meal plans',
                              ),
                            ],
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Based on your activity: ${hc.stepsToday} steps, ${hc.kcalToday.toStringAsFixed(0)} kcal burned',
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: Theme.of(context).colorScheme.onSurfaceVariant,
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Show loading state
                      if (hc.isGeneratingMealPlan)
                        const Center(
                          child: Padding(
                            padding: EdgeInsets.all(32),
                            child: Column(
                              children: [
                                CircularProgressIndicator(),
                                SizedBox(height: 16),
                                Text(
                                  'Generating your personalized meal plan...',
                                ),
                              ],
                            ),
                          ),
                        ),

                      // Show error if any
                      if (hc.mealPlanError != null)
                        Container(
                          padding: const EdgeInsets.all(12),
                          decoration: BoxDecoration(
                            color: Theme.of(context).colorScheme.errorContainer,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Row(
                            children: [
                              Icon(
                                Icons.error_outline,
                                color: Theme.of(
                                  context,
                                ).colorScheme.onErrorContainer,
                              ),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  hc.mealPlanError!,
                                  style: TextStyle(
                                    color: Theme.of(
                                      context,
                                    ).colorScheme.onErrorContainer,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),

                      // Display generated meal plan
                      if (hc.currentMealPlan != null)
                        Column(
                          children: [
                            ...hc.currentMealPlan!.plan.map(
                              (mealPlan) => MealPlanCard(mealPlan: mealPlan),
                            ),
                          ],
                        ),

                      // Show message if no plan available
                      if (hc.currentMealPlan == null &&
                          !hc.isGeneratingMealPlan &&
                          hc.mealPlanError == null)
                        Container(
                          padding: const EdgeInsets.all(24),
                          decoration: BoxDecoration(
                            color: Theme.of(context)
                                .colorScheme
                                .surfaceContainerHighest
                                .withValues(alpha: 0.3),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Column(
                            children: [
                              Icon(
                                Icons.restaurant_menu,
                                size: 48,
                                color: Theme.of(
                                  context,
                                ).colorScheme.onSurfaceVariant,
                              ),
                              const SizedBox(height: 16),
                              Text(
                                'Complete your profile to get personalized meal plans',
                                style: Theme.of(context).textTheme.bodyLarge
                                    ?.copyWith(
                                      color: Theme.of(
                                        context,
                                      ).colorScheme.onSurfaceVariant,
                                    ),
                                textAlign: TextAlign.center,
                              ),
                              const SizedBox(height: 8),
                              Text(
                                'Add your goal, height, and weight in your profile settings',
                                style: Theme.of(context).textTheme.bodyMedium
                                    ?.copyWith(
                                      color: Theme.of(
                                        context,
                                      ).colorScheme.onSurfaceVariant,
                                    ),
                                textAlign: TextAlign.center,
                              ),
                            ],
                          ),
                        ),
                    ],
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
              color: theme.colorScheme.shadow.withValues(alpha: .06),
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
