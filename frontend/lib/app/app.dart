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
              NavigationDestination(
                icon: Icon(Icons.storefront),
                label: 'Shop',
              ),
              NavigationDestination(
                icon: Icon(Icons.shopping_cart),
                label: 'Cart',
              ),
              NavigationDestination(icon: Icon(Icons.person), label: 'Profile'),
              NavigationDestination(icon: Icon(Icons.list), label: 'Plans'),
              NavigationDestination(
                icon: Icon(Icons.show_chart),
                label: 'Progress',
              ),
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
              path: '/carts',
              builder: (_, _) =>
                  const Scaffold(body: Center(child: Text('Carts Page'))),
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
        StatefulShellBranch(
          routes: [
            GoRoute(
              path: '/plans',
              builder: (_, _) =>
                  const Scaffold(body: Center(child: Text('Plans Page'))),
            ),
          ],
        ),
        StatefulShellBranch(
          routes: [
            GoRoute(
              path: '/progress',
              builder: (_, _) =>
                  const Scaffold(body: Center(child: Text('Progress Page'))),
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
