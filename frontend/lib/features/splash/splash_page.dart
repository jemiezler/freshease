import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF5FBE7E),
      body: Column(
        children: [
          // Logo and branding - centered in one container
          Padding(
            padding: const EdgeInsets.only(top: 60),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                // Logo
                Image.asset('lib/assets/logo.png', width: 60, height: 60),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: const [
                    Text(
                      'FreashEase',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 32,
                        fontWeight: FontWeight.bold,
                        letterSpacing: 0.2,
                      ),
                    ),
                    Text(
                      'Farm-to-door grocery delivery',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),

          // Main character illustration
          Expanded(
            child: Stack(
              children: [
                // Character illustration
                Positioned.fill(
                  child: Align(
                    alignment: Alignment.bottomLeft,
                    child: Image.asset(
                      'lib/assets/delivery_character.png',
                      fit: BoxFit.fitHeight,
                      alignment: Alignment.bottomLeft,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  @override
  void initState() {
    super.initState();
    _navigate();
  }

  Future<void> _navigate() async {
    // DEV MODE: Set to true to freeze on splash screen
    const bool devMode = false;

    if (devMode) {
      // Stay on splash screen for development
      return;
    }

    // Show splash for 2 seconds
    await Future.delayed(const Duration(seconds: 2));

    if (!mounted) return;

    // Check if user has completed onboarding
    final prefs = await SharedPreferences.getInstance();
    final hasSeenOnboarding = prefs.getBool('hasSeenOnboarding') ?? false;

    if (!mounted) return;

    if (hasSeenOnboarding) {
      // Skip onboarding, go to login
      context.go('/login');
    } else {
      // First time, show onboarding
      context.go('/onboarding');
    }
  }
}
