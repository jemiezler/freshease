import 'package:flutter/material.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:go_router/go_router.dart';

class PlanDetailPage extends StatelessWidget {
  final Map<String, dynamic> plan;
  const PlanDetailPage({super.key, required this.plan});

  @override
  Widget build(BuildContext context) {
    final color = plan["color"] as Color;
    final features = (plan["features"] as List).cast<String>();

    return Scaffold(
      appBar: AppBar(title: Text(plan["title"])),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 900),
          child: ListView(
            padding: const EdgeInsets.all(24),
            children: [
              Text(
                plan["title"],
                style: const TextStyle(
                  fontSize: 26,
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                plan["subtitle"],
                style: TextStyle(fontSize: 16, color: Colors.grey[600]),
              ),
              const SizedBox(height: 16),
              Container(
                decoration: BoxDecoration(
                  color: color.withValues(alpha: .1),
                  borderRadius: BorderRadius.circular(12),
                ),
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    const Text(
                      'What’s Included:',
                      style: TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 18,
                      ),
                    ),
                    const SizedBox(height: 12),
                    for (final f in features)
                      Padding(
                        padding: const EdgeInsets.only(bottom: 6),
                        child: Row(
                          children: [
                            const Icon(Icons.check, color: Colors.green),
                            const SizedBox(width: 8),
                            Expanded(child: Text(f)),
                          ],
                        ),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              Text(
                'Price: ฿${plan["price"]}',
                style: const TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.w700,
                  color: Colors.green,
                ),
              ),
              const SizedBox(height: 24),
              FilledButton.icon(
                onPressed: () {
                  final co = CheckoutScope.of(context);
                  co.setPlanCheckout(
                    PlanOrder(
                      id: plan["id"] as int,
                      title: plan["title"] as String,
                      price: plan["price"] as double,
                      subtitle: plan["subtitle"] as String,
                    ),
                  );

                  // Navigate directly into checkout flow (starting at address)
                  context.go('/cart/checkout/address');
                },
                icon: const Icon(Icons.shopping_bag_outlined),
                label: const Text('Subscribe Now'),
                style: FilledButton.styleFrom(
                  minimumSize: const Size.fromHeight(48),
                  backgroundColor: color.withValues(alpha: .9),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
