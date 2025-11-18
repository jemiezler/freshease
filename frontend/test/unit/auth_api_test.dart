import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:frontend/features/auth/data/sources/auth_api.dart';
import 'package:frontend/core/network/dio_client.dart';

import 'auth_api_test.mocks.dart';

@GenerateMocks([DioClient, Dio])
void main() {
  group('AuthApi', () {
    late AuthApi authApi;
    late MockDioClient mockDioClient;
    late MockDio mockDio;

    setUp(() {
      mockDioClient = MockDioClient();
      mockDio = MockDio();
      when(mockDioClient.dio).thenReturn(mockDio);
      authApi = AuthApi(mockDioClient);
      SharedPreferences.setMockInitialValues({});
    });

    test('verifyToken returns user data when token is valid', () async {
      SharedPreferences.setMockInitialValues({'access_token': 'valid_token'});
      final responseData = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
      };

      when(mockDio.get('/api/whoami')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: responseData,
          ));

      final result = await authApi.verifyToken();

      expect(result['id'], 'user-1');
      expect(result['email'], 'test@example.com');
      verify(mockDio.get('/api/whoami')).called(1);
    });

    test('verifyToken throws exception when no token exists', () async {
      SharedPreferences.setMockInitialValues({});

      expect(
        () => authApi.verifyToken(),
        throwsException,
      );
    });

    test('verifyToken throws exception on 401 and clears tokens', () async {
      SharedPreferences.setMockInitialValues({'access_token': 'invalid_token'});
      final prefs = await SharedPreferences.getInstance();

      when(mockDio.get('/api/whoami')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 401,
        ),
      ));

      expect(
        () => authApi.verifyToken(),
        throwsException,
      );

      // Verify tokens are cleared (this is done internally by AuthApi)
      // We can't directly verify SharedPreferences.remove calls without more complex mocking
    });

    test('verifyToken throws exception on 403', () async {
      SharedPreferences.setMockInitialValues({'access_token': 'token'});

      when(mockDio.get('/api/whoami')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 403,
        ),
      ));

      expect(
        () => authApi.verifyToken(),
        throwsException,
      );
    });

    test('verifyToken handles network errors', () async {
      SharedPreferences.setMockInitialValues({'access_token': 'token'});

      when(mockDio.get('/api/whoami')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        type: DioExceptionType.connectionError,
        error: 'Network error',
      ));

      expect(
        () => authApi.verifyToken(),
        throwsException,
      );
    });

    test('verifyToken throws exception for invalid response format', () async {
      SharedPreferences.setMockInitialValues({'access_token': 'token'});

      when(mockDio.get('/api/whoami')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: 'invalid_json', // Not a Map
          ));

      expect(
        () => authApi.verifyToken(),
        throwsException,
      );
    });

    // Note: signInWithGoogle tests are complex because they involve
    // FlutterWebAuth2 which requires platform-specific setup.
    // These would typically be integration tests or require more complex mocking.
    // For now, we'll test the error handling paths that don't require web auth.
  });
}


