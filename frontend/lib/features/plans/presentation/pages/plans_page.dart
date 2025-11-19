import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class PlansPage extends StatefulWidget {
  const PlansPage({super.key});

  @override
  State<PlansPage> createState() => _PlansPageState();
}

class _PlansPageState extends State<PlansPage> {
  final _plans = const [
    {
      "id": 1,
      "title": "Fresh Starter Plan",
      "subtitle": "7-Day Clean Eating Challenge",
      "price": 499.0,
      "features": [
        "Daily meal plan & recipes",
        "Auto-generated grocery bundle",
        "Basic nutrition coaching",
      ],
      "color": Colors.green,
    },
    {
      "id": 2,
      "title": "Balanced Weekly Plan",
      "subtitle": "Perfect for busy professionals",
      "price": 899.0,
      "features": [
        "5 curated dinners per week",
        "One-click grocery ordering",
        "AI meal recommendations",
      ],
      "color": Colors.orange,
    },
    {
      "id": 3,
      "title": "Pro Wellness Plan",
      "subtitle": "Monthly personalized plan",
      "price": 2499.0,
      "features": [
        "Full AI personalization",
        "Private nutrition coach",
        "Progress tracking dashboard",
      ],
      "color": Colors.teal,
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'FreshEase Plans',
          style: TextStyle(color: Colors.white),
        ),
        backgroundColor: Theme.of(context).colorScheme.primary,
      ),
      body: LayoutBuilder(
        builder: (context, constraints) {
          final isWide = constraints.maxWidth >= 900;
          final crossAxisCount = isWide
              ? 3
              : (constraints.maxWidth > 600 ? 2 : 1);
          return GridView.builder(
            padding: const EdgeInsets.all(16),
            gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: crossAxisCount,
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
              childAspectRatio: isWide ? 1.1 : 0.95,
            ),
            itemCount: _plans.length,
            itemBuilder: (_, i) {
              final plan = _plans[i];
              return _PlanCard(
                title: plan["title"] as String,
                subtitle: plan["subtitle"] as String,
                price: plan["price"] as double,
                features: (plan["features"] as List).cast<String>(),
                color: plan["color"] as Color,
                onSubscribe: () =>
                    context.go('/plans/${plan["id"]}', extra: plan),
              );
            },
          );
        },
      ),
    );
  }
}

class _PlanCard extends StatelessWidget {
  final String title;
  final String subtitle;
  final double price;
  final List<String> features;
  final Color color;
  final VoidCallback onSubscribe;

  const _PlanCard({
    required this.title,
    required this.subtitle,
    required this.price,
    required this.features,
    required this.color,
    required this.onSubscribe,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              decoration: BoxDecoration(
                color: color.withValues(alpha: .15),
                borderRadius: BorderRadius.circular(12),
              ),
              padding: const EdgeInsets.all(8),
              child: Icon(Icons.eco, color: color, size: 28),
            ),
            const SizedBox(height: 12),
            Text(
              title,
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
            ),
            const SizedBox(height: 4),
            Text(subtitle, style: TextStyle(color: Colors.grey[600])),
            const SizedBox(height: 12),
            Text(
              '฿${price.toStringAsFixed(0)} / plan',
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w700,
                color: Colors.green,
              ),
            ),
            const SizedBox(height: 12),
            Expanded(
              child: ListView.builder(
                itemCount: features.length,
                physics: const NeverScrollableScrollPhysics(),
                itemBuilder: (_, i) => Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Icon(
                        Icons.check_circle,
                        size: 16,
                        color: Colors.green,
                      ),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Text(
                          features[i],
                          style: const TextStyle(height: 1.2),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(height: 8),
            FilledButton.icon(
              onPressed: onSubscribe,
              icon: const Icon(Icons.arrow_forward),
              label: const Text('View Details'),
              style: FilledButton.styleFrom(
                minimumSize: const Size.fromHeight(44),
                backgroundColor: color.withValues(alpha: 0.9),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
