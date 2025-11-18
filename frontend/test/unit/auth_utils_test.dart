import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:frontend/core/auth/auth_utils.dart';

void main() {
  group('AuthUtils', () {
    setUp(() async {
      // Clear SharedPreferences before each test
      SharedPreferences.setMockInitialValues({});
      final prefs = await SharedPreferences.getInstance();
      await prefs.clear();
    });

    test('isAuthenticated returns false when no token exists', () async {
      final result = await AuthUtils.isAuthenticated();
      expect(result, false);
    });

    test('isAuthenticated returns true when token exists', () async {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', 'test-token');

      final result = await AuthUtils.isAuthenticated();
      expect(result, true);
    });

    test('isAuthenticated returns false when token is empty string', () async {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', '');

      final result = await AuthUtils.isAuthenticated();
      expect(result, false);
    });

    test('getAccessToken returns null when no token exists', () async {
      final result = await AuthUtils.getAccessToken();
      expect(result, isNull);
    });

    test('getAccessToken returns token when it exists', () async {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', 'test-token-123');

      final result = await AuthUtils.getAccessToken();
      expect(result, 'test-token-123');
    });

    test('storeToken saves token to SharedPreferences', () async {
      await AuthUtils.storeToken('new-token');

      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('access_token');
      expect(token, 'new-token');
    });

    test('clearAuth removes all auth tokens', () async {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', 'token');
      await prefs.setString('refresh_token', 'refresh');
      await prefs.setString('id_token', 'id');

      await AuthUtils.clearAuth();

      expect(prefs.getString('access_token'), isNull);
      expect(prefs.getString('refresh_token'), isNull);
      expect(prefs.getString('id_token'), isNull);
    });

    test('getAuthStatus returns notAuthenticated when no token', () async {
      final result = await AuthUtils.getAuthStatus();
      expect(result, AuthStatus.notAuthenticated);
    });

    test('getAuthStatus returns authenticated when token exists', () async {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', 'test-token');

      final result = await AuthUtils.getAuthStatus();
      expect(result, AuthStatus.authenticated);
    });

    test('getAuthStatus returns notAuthenticated when token is empty', () async {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', '');

      final result = await AuthUtils.getAuthStatus();
      expect(result, AuthStatus.notAuthenticated);
    });
  });
}

