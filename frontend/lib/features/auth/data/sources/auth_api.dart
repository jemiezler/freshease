// data/sources/auth_api.dart
import 'package:dio/dio.dart';
import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../../core/network/dio_client.dart';

class AuthApi {
  final Dio _dio;
  AuthApi(DioClient client) : _dio = client.dio;

  static const _scheme = 'freshease';
  static const _provider = 'google';

  Future<Map<String, dynamic>> signInWithGoogle() async {
    // 1) Start OIDC via your backend
    final startUrl = '${_dio.options.baseUrl}/api/auth/$_provider/start';

    final callbackUrl = await FlutterWebAuth2.authenticate(
      url: startUrl,
      callbackUrlScheme: _scheme,
    );

    // 2) Parse the app callback - backend returns JSON with accessToken
    final uri = Uri.parse(callbackUrl);
    final code = uri.queryParameters['code'];
    final state = uri.queryParameters['state'];

    if (code == null || state == null) {
      throw Exception('OAuth callback missing code or state');
    }

    // 3) Exchange code for tokens using backend exchange endpoint
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
}
