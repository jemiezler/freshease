// data/sources/auth_api.dart
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../../core/network/dio_client.dart';

class AuthApi {
  final Dio _dio;
  AuthApi(DioClient client) : _dio = client.dio;

  static const _scheme = 'freshease';
  static const _provider = 'google';

  /// Verify if the existing token is valid by calling /api/whoami
  /// Returns user data if token is valid, throws exception if invalid
  Future<Map<String, dynamic>> verifyToken() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('access_token');

      if (token == null || token.isEmpty) {
        throw Exception('No token found');
      }

      // Verify token by calling whoami endpoint
      final response = await _dio.get('/api/whoami');

      if (response.statusCode == 200) {
        final data = response.data;
        if (data is Map<String, dynamic>) {
          return data;
        } else {
          throw Exception('Invalid response format from server');
        }
      } else {
        throw Exception('Token verification failed: ${response.statusCode}');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response?.statusCode;
        if (statusCode == 401 || statusCode == 403) {
          // Token is invalid, clear it
          final prefs = await SharedPreferences.getInstance();
          await prefs.remove('access_token');
          await prefs.remove('refresh_token');
          await prefs.remove('id_token');
          throw Exception('Token expired or invalid');
        }
        throw Exception(
          'Token verification failed: ${e.response?.data?['message'] ?? e.message}',
        );
      } else {
        throw Exception(
          'Network error: ${e.message ?? 'Unable to connect to server'}',
        );
      }
    } catch (e) {
      if (e is Exception) {
        rethrow;
      }
      throw Exception('Token verification failed: $e');
    }
  }

  Future<Map<String, dynamic>> signInWithGoogle() async {
    // 1) Start OIDC via your backend
    final isWeb = kIsWeb;

    // For web, get the current origin where the Flutter app is running
    String? webCallbackUrl;
    if (isWeb) {
      // Use Uri.base to get the current origin (where the Flutter web app is hosted)
      final currentOrigin = Uri.base.origin;
      webCallbackUrl = '$currentOrigin/auth/callback';
      final startUrl =
          '${_dio.options.baseUrl}/api/auth/$_provider/start?platform=web&callback_url=${Uri.encodeComponent(webCallbackUrl)}';

      try {
        // For web, callbackUrlScheme should be just the scheme (http or https)
        // FlutterWebAuth2 will match URLs that start with this scheme
        // The full callback URL is passed to the backend via callback_url parameter

        final uri = Uri.parse(currentOrigin);
        final scheme = uri.scheme; // Extract just 'http' or 'https'

        final callbackUrl = await FlutterWebAuth2.authenticate(
          url: startUrl,
          callbackUrlScheme:
              scheme, // Use just the scheme for web (e.g., 'http' or 'https')
        );

        // Parse callback and continue with exchange
        return await _handleOAuthCallback(callbackUrl);
      } catch (e) {
        _handleAuthError(e);
      }
    } else {
      // Mobile: use custom scheme
      final startUrl = '${_dio.options.baseUrl}/api/auth/$_provider/start';

      try {
        final callbackUrl = await FlutterWebAuth2.authenticate(
          url: startUrl,
          callbackUrlScheme: _scheme,
        );

        // Parse callback and continue with exchange
        return await _handleOAuthCallback(callbackUrl);
      } catch (e) {
        _handleAuthError(e);
      }
    }
  }

  Future<Map<String, dynamic>> _handleOAuthCallback(String callbackUrl) async {
    // Check if callback is null or empty
    if (callbackUrl.isEmpty) {
      throw Exception('Empty callback URL received');
    }

    // 2) Parse the app callback - backend returns JSON with accessToken
    final uri = Uri.parse(callbackUrl);
    final code = uri.queryParameters['code'];
    final state = uri.queryParameters['state'];

    if (code == null || state == null) {
      throw Exception('OAuth callback missing code or state');
    }
    final res = await _dio.post(
      '/api/auth/$_provider/exchange',
      data: {'code': code, 'state': state},
    );

    final data = res.data as Map<String, dynamic>;
    final responseData = data['data'] as Map<String, dynamic>;
    final accessToken = responseData['accessToken'] as String;

    if (accessToken.isEmpty) {
      throw DioException(
        requestOptions: res.requestOptions,
        error: 'No access token in exchange response',
      );
    }

    // 4) Store the access token
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('access_token', accessToken);

    // 5) Fetch current user info
    final me = await _dio.get('/api/whoami');

    return me.data as Map<String, dynamic>;
  }

  Never _handleAuthError(dynamic e) {
    // Handle specific error types
    if (e.toString().contains('CANCELED') ||
        e.toString().contains('User canceled')) {
      throw Exception(
        'Login was canceled or interrupted. Please try again and complete the full login process.',
      );
    } else if (e.toString().contains('NETWORK_ERROR')) {
      throw Exception('Network error. Please check your internet connection');
    } else if (e.toString().contains('401')) {
      throw Exception('Authentication failed. Please try again');
    } else {
      throw Exception('Google Sign-in failed: ${e.toString()}');
    }
  }
}
