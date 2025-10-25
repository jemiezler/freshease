import 'package:flutter/material.dart';
import 'package:frontend/core/state/checkout_controller.dart';
import 'package:go_router/go_router.dart';

// --- Define App Colors from Figma ---
const Color _scaffoldBgColor = Color.fromARGB(255, 252, 252, 252);
const Color _appBarBgColor = Color(0xFFF7F3F0);
const Color _primaryTextColor = Color(0xFF333333);
const Color _secondaryTextColor = Color(0xFF616161);
const Color _checkColor = Color(0xFF4CAF50); // Colors.green.shade600

class PlanDetailPage extends StatelessWidget {
  final Map<String, dynamic> plan;
  const PlanDetailPage({super.key, required this.plan});

  @override
  Widget build(BuildContext context) {
    // Extract all data first
    final color = plan["color"] as Color;
    final features = (plan["features"] as List).cast<String>();
    final title = plan["title"] as String;
    final subtitle = plan["subtitle"] as String;
    final price = plan["price"] as double;

    return Scaffold(
      // 1. Set background color
      backgroundColor: _scaffoldBgColor,
      appBar: AppBar(
        // 2. Fix AppBar title and style
        title: const Text(
          'FreshEasePlan', // Fixed title from Figma
          style: TextStyle(
            color: _primaryTextColor,
            fontWeight: FontWeight.bold,
          ),
        ),
        backgroundColor: _appBarBgColor,
        elevation: 0,
        iconTheme: const IconThemeData(
          color: _primaryTextColor, // Dark back arrow
        ),
      ),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 900),
          child: ListView(
            padding: const EdgeInsets.all(24),
            children: [
              // 3. Content Order: Title
              Text(
                title,
                style: const TextStyle(
                  fontSize: 26,
                  fontWeight: FontWeight.w800,
                  color: _primaryTextColor, // Ensure dark text
                ),
              ),
              const SizedBox(height: 8),
              // 4. Content Order: Subtitle
              Text(
                subtitle,
                style: const TextStyle(
                  fontSize: 16,
                  color: _secondaryTextColor, // Use theme color
                ),
              ),
              const SizedBox(height: 16),
              // 5. Content Order: Price (Moved & Restyled)
              Text(
                '฿ ${price.toStringAsFixed(0)} / Plan', // New format
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  color: _primaryTextColor, // Dark text, not green
                ),
              ),
              const SizedBox(height: 24),
              // 6. Content Order: "What's Include" Box
              Container(
                decoration: BoxDecoration(
                  // Use the plan color with more opacity
                  color: color.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(16),
                ),
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      "What's Include:", // Match Figma (no 'd')
                      style: TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 18,
                        color: _primaryTextColor,
                      ),
                    ),
                    const SizedBox(height: 12),
                    for (final f in features)
                      Padding(
                        padding: const EdgeInsets.only(bottom: 8),
                        child: Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Icon(
                              Icons
                                  .check_circle, // Match previous screen's icon
                              color: _checkColor,
                              size: 18,
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                f,
                                style: const TextStyle(
                                  color: _primaryTextColor,
                                  fontSize: 15,
                                  height: 1.3,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              // 7. Button (No Icon, solid color)
              FilledButton(
                onPressed: () {
                  final co = CheckoutScope.of(context);
                  co.setPlanCheckout(
                    PlanOrder(
                      id: plan["id"] as int,
                      title: title,
                      price: price,
                      subtitle: subtitle,
                    ),
                  );

                  // Navigate directly into checkout flow (starting at address)
                  context.go('/cart/checkout/address');
                },
                style: FilledButton.styleFrom(
                  minimumSize: const Size.fromHeight(48),
                  // Use solid plan color
                  backgroundColor: color,
                  // Use white text for contrast
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                // Remove icon, use child
                child: const Text(
                  'Subscribe Now',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
