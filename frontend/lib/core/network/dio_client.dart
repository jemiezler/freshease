import 'package:dio/dio.dart';

class DioClient {
  final Dio dio = Dio(
    BaseOptions(
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 20),
      // baseUrl: const String.fromEnvironment('BASE_URL', defaultValue: 'https://api.dev'),
    ),
  )..interceptors.add(LogInterceptor(requestBody: true, responseBody: true));
}
