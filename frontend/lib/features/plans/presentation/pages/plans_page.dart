import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

// --- Define App Colors from Figma ---
// ⚠️ หมายเหตุ: คุณเปลี่ยน _scaffoldBgColor เป็นสีขาว ผมจะคงตามนี้นะครับ
const Color _scaffoldBgColor = Color.fromARGB(255, 255, 255, 255);
const Color _appBarBgColor = Color(0xFFF7F3F0);
const Color _cardBgColor = Colors.white;
// ⚠️ _buttonColor และ _buttonTextColor ไม่ได้ใช้สำหรับปุ่ม plan card แล้ว
// const Color _buttonColor = Color(0xFFB0D9B1);
// const Color _buttonTextColor = Color(0xFF333333);
const Color _primaryTextColor = Color(0xFF333333);
const Color _secondaryTextColor = Color(0xFF616161); // Colors.grey.shade700
const Color _checkColor = Color(0xFF4CAF50); // Colors.green.shade600

class PlansPage extends StatefulWidget {
  const PlansPage({super.key});

  @override
  State<PlansPage> createState() => _PlansPageState();
}

class _PlansPageState extends State<PlansPage> {
  // Data model with colors that will be used for buttons
  final _plans = const [
    {
      "id": 1,
      "title": "Fresh Starter Plan",
      "subtitle": "7 days Clean Eating Challenge",
      "price": 499.0,
      "features": [
        "Daily meal plan & recipes",
        "Auto-generated grocery bundle",
        "Basic nutrition coaching",
      ],
      "color": Colors.green, // Now used for button
    },
    {
      "id": 2,
      "title": "Balanced weekly Plan",
      "subtitle": "Perfect for busy professionals",
      "price": 899.0,
      "features": [
        "5 Curated dinners per week",
        "One-click grocery ordering",
        "AI meal recommendations",
      ],
      "color": Colors.orange, // Now used for button
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
      "color": Colors.teal, // Now used for button
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: _scaffoldBgColor,
      appBar: AppBar(
        title: const Text(
          'FreshEasePlans',
          style: TextStyle(
            color: _primaryTextColor,
            fontWeight: FontWeight.bold,
          ),
        ),
        backgroundColor: _appBarBgColor,
        elevation: 0,
        // Assuming default icon theme is fine (dark icons on light bg)
      ),
      body: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
        itemCount: _plans.length,
        itemBuilder: (_, i) {
          final plan = _plans[i];
          return Padding(
            // Add spacing between cards
            padding: const EdgeInsets.only(bottom: 16.0),
            child: _PlanCard(
              title: plan["title"] as String,
              subtitle: plan["subtitle"] as String,
              price: plan["price"] as double,
              features: (plan["features"] as List).cast<String>(),
              // --- 🟢 FIX: Pass the color to the card ---
              buttonColor: plan["color"] as Color,
              // ------------------------------------------
              onSubscribe: () =>
                  context.go('/plans/${plan["id"]}', extra: plan),
            ),
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
  final VoidCallback onSubscribe;
  // --- 🟢 FIX: Add buttonColor property ---
  final Color buttonColor;
  // -------------------------------------

  // --- Removed 'const' because Image.network is not const ---
  const _PlanCard({
    super.key,
    required this.title,
    required this.subtitle,
    required this.price,
    required this.features,
    required this.onSubscribe,
    // --- 🟢 FIX: Require buttonColor in constructor ---
    required this.buttonColor,
    // -----------------------------------------------
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      color: _cardBgColor,
      child: Stack(
        children: [
          // Salad bowl icon in top right
          Positioned(
            top: 16,
            right: 16,
            child: Image.network(
              'https://i.pinimg.com/736x/80/e2/0c/80e20c8ed50251fc4274d067df7c8a11.jpg',
              width: 48,
              height: 48,
              // --- 🟢 FIX: Removed 'color' property so image is visible ---
              // color: buttonColor.withOpacity(0.2), // This tints the image
              // -----------------------------------------------------------
              errorBuilder: (context, error, stackTrace) => Icon(
                Icons.image_not_supported_outlined,
                color: Colors.grey.shade200,
                size: 40,
              ),
            ),
          ),
          // Main content
          Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Title
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 19,
                    fontWeight: FontWeight.w800,
                    color: _primaryTextColor,
                  ),
                ),
                const SizedBox(height: 4),
                // Subtitle
                Text(
                  subtitle,
                  style: const TextStyle(
                    color: _secondaryTextColor,
                    fontSize: 15,
                  ),
                ),
                const SizedBox(height: 16),
                // Price
                Text(
                  '฿ ${price.toStringAsFixed(0)} / Plan',
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w700,
                    color: _primaryTextColor,
                  ),
                ),
                const SizedBox(height: 16),
                // Features
                Column(
                  children: features
                      .map((feature) => _FeatureItem(text: feature))
                      .toList(),
                ),
                const SizedBox(height: 16),
                // Button
                FilledButton(
                  onPressed: onSubscribe,
                  style: FilledButton.styleFrom(
                    minimumSize: const Size.fromHeight(44),
                    // --- 🟢 FIX: Use dynamic color ---
                    backgroundColor: buttonColor,
                    foregroundColor:
                        Colors.white, // Use white text for contrast
                    // ----------------------------------
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: const Text(
                    'View Details',
                    style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// Helper widget for a single feature item
class _FeatureItem extends StatelessWidget {
  final String text;
  const _FeatureItem({super.key, required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(Icons.check_circle, size: 18, color: _checkColor),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              text,
              style: const TextStyle(
                color: _primaryTextColor,
                fontSize: 15,
                height: 1.3,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
