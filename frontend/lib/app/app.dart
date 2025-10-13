import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../features/auth/presentation/pages/login_page.dart';
import '../features/home/presentation/pages/home_page.dart';

final _router = GoRouter(
  routes: [
    // Public route (no bottom bar)
    GoRoute(path: '/login', builder: (_, _) => const LoginPage()),

    // App shell (shows bottom bar)
    StatefulShellRoute.indexedStack(
      builder: (context, state, navigationShell) {
        return Scaffold(
          body: navigationShell, // keeps each branch alive
          bottomNavigationBar: NavigationBar(
            height: 65,
            selectedIndex: navigationShell.currentIndex,
            onDestinationSelected: (index) {
              // Navigate to a branch; keep current location if reselecting
              navigationShell.goBranch(
                index,
                initialLocation: index != navigationShell.currentIndex,
              );
            },
            destinations: const [
              NavigationDestination(icon: Icon(Icons.home), label: 'Home'),
              NavigationDestination(
                icon: Icon(Icons.favorite),
                label: 'Favorites',
              ),
              NavigationDestination(icon: Icon(Icons.person), label: 'Profile'),
            ],
          ),
        );
      },
      branches: [
        // Tab 0: Home
        StatefulShellBranch(
          routes: [GoRoute(path: '/', builder: (_, _) => const HomePage())],
        ),
        // Tab 1: Favorites
        StatefulShellBranch(
          routes: [
            GoRoute(
              path: '/favorites',
              builder: (_, _) =>
                  const Scaffold(body: Center(child: Text('Favorites Page'))),
            ),
          ],
        ),
        // Tab 2: Profile
        StatefulShellBranch(
          routes: [
            GoRoute(
              path: '/profile',
              builder: (_, _) =>
                  const Scaffold(body: Center(child: Text('Profile Page'))),
            ),
          ],
        ),
      ],
    ),
  ],
);

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'FreshEase',
      theme: ThemeData(useMaterial3: true, colorSchemeSeed: Colors.teal),
      routerConfig: _router,
      debugShowCheckedModeBanner: false,
    );
  }
}
