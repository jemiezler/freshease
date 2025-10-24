// data/sources/auth_api.dart
import 'package:dio/dio.dart';
import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../../core/network/dio_client.dart';

class AuthApi {
  final Dio _dio;
  AuthApi(DioClient client) : _dio = client.dio;

  static const _scheme = 'freshease';
  static const _redirectUri = '$_scheme://auth/callback';
  static const _provider = 'google';

  Future<Map<String, dynamic>> signInWithGoogle() async {
    // 1) Start OIDC via your backend
    final startUrl = '${_dio.options.baseUrl}/auth/$_provider/start';

    final callbackUrl = await FlutterWebAuth2.authenticate(
      url: startUrl,
      callbackUrlScheme: _scheme,
    );

    // 2) Parse the app callback
    final uri = Uri.parse(callbackUrl);
    final access = uri.queryParameters['access_token'];
    final refresh = uri.queryParameters['refresh_token'];
    final idToken = uri.queryParameters['id_token']; // optional
    final code = uri.queryParameters['code']; // fallback

    final prefs = await SharedPreferences.getInstance();

    if (access != null) {
      await prefs.setString('access_token', access);
      if (refresh != null) await prefs.setString('refresh_token', refresh);
      if (idToken != null) await prefs.setString('id_token', idToken);
    } else if (code != null) {
      // If your backend returns a code to the app, expose an exchange endpoint:
      // POST /auth/:provider/exchange { code, redirect_uri }
      final res = await _dio.post(
        '/auth/$_provider/exchange',
        data: {'code': code, 'redirect_uri': _redirectUri},
      );
      final data = res.data as Map<String, dynamic>;
      final at = data['accessToken'] ?? data['access_token'];
      final rt = data['refreshToken'] ?? data['refresh_token'];
      final it = data['idToken'] ?? data['id_token'];
      if (at == null) {
        throw DioException(
          requestOptions: res.requestOptions,
          error: 'No access token in exchange response',
        );
      }
      await prefs.setString('access_token', at);
      if (rt != null) await prefs.setString('refresh_token', rt);
      if (it != null) await prefs.setString('id_token', it);
    } else {
      throw Exception('OAuth callback missing token/code');
    }

    // 3) Fetch current user (adjust path to your backend, e.g. /auth/me)
    final me = await _dio.get('/auth/me');
    return me.data as Map<String, dynamic>;
  }
}
