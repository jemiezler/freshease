// core/auth/auth_utils.dart
import 'package:shared_preferences/shared_preferences.dart';

class AuthUtils {
  /// Check if user is currently authenticated
  static Future<bool> isAuthenticated() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('access_token');
    return token != null && token.isNotEmpty;
  }

  /// Get the stored access token
  static Future<String?> getAccessToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('access_token');
  }

  /// Clear all authentication data
  static Future<void> clearAuth() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('access_token');
    await prefs.remove('refresh_token');
    await prefs.remove('id_token');
  }

  /// Store authentication token
  static Future<void> storeToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('access_token', token);
  }

  /// Check authentication status with detailed info
  static Future<AuthStatus> getAuthStatus() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('access_token');

    if (token == null || token.isEmpty) {
      return AuthStatus.notAuthenticated;
    }

    return AuthStatus.authenticated;
  }
}

enum AuthStatus { authenticated, notAuthenticated, expired }
