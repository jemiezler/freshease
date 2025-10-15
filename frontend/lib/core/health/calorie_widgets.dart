import 'package:flutter/material.dart';
import 'mock_health_repo.dart';

class CaloriesBlock extends StatelessWidget {
  final CalorieSnapshot? snapshot;
  final bool syncing;
  final VoidCallback onSync;

  const CaloriesBlock({
    super.key,
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
                  label: Text(syncing ? 'Syncing…' : 'Sync (Mock)'),
                ),
              ],
            ),
            const SizedBox(height: 12),

            if (snap == null) ...[
              Text(
                'Connect later — using mock numbers.\nTap “Sync (Mock)” to refresh.',
                style: TextStyle(color: Colors.grey[700]),
              ),
              const SizedBox(height: 8),
            ] else ...[
              _RowMetric(label: 'Intake (kcal)', value: snap.intakeKcal),
              const SizedBox(height: 6),
              _RowMetric(
                label: 'Active burned (kcal)',
                value: snap.activeBurnKcal,
              ),
              const Divider(height: 18),
              _NetRow(net: snap.netKcal),
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
    final h = dt.hour.toString().padLeft(2, '0');
    final m = dt.minute.toString().padLeft(2, '0');
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')} $h:$m';
  }
}

class _RowMetric extends StatelessWidget {
  final String label;
  final double value;
  const _RowMetric({required this.label, required this.value});

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

class _NetRow extends StatelessWidget {
  final double net;
  const _NetRow({required this.net});

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
