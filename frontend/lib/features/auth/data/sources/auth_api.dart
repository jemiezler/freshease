import 'package:dio/dio.dart';
import '../../../../core/network/dio_client.dart';

class AuthApi {
  final Dio _dio;
  AuthApi(DioClient client) : _dio = client.dio;

  Future<Map<String, dynamic>> login(String email, String password) async {
    // Fake: replace with real endpoint
    // final res = await _dio.post('/auth/login', data: {'email': email, 'password': password});
    // return res.data as Map<String, dynamic>;
    await Future<void>.delayed(const Duration(milliseconds: 600));
    return {'id': 'u_123', 'email': email};
  }
}
