import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:dio/dio.dart';
import 'package:frontend/features/account/data/sources/user_api.dart';
import 'package:frontend/core/network/dio_client.dart';

import 'user_api_test.mocks.dart';

@GenerateMocks([DioClient, Dio])
void main() {
  group('UserApi', () {
    late UserApi userApi;
    late MockDioClient mockDioClient;
    late MockDio mockDio;

    setUp(() {
      mockDioClient = MockDioClient();
      mockDio = MockDio();
      when(mockDioClient.dio).thenReturn(mockDio);
      userApi = UserApi(mockDioClient);
    });


    test('getCurrentUser returns user data on success', () async {
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

      final result = await userApi.getCurrentUser();

      expect(result['id'], 'user-1');
      expect(result['email'], 'test@example.com');
      expect(result['name'], 'Test User');
      verify(mockDio.get('/api/whoami')).called(1);
    });

    test('getCurrentUser throws exception on 401', () async {
      when(mockDio.get('/api/whoami')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 401,
          data: {'message': 'Unauthorized'},
        ),
      ));

      expect(
        () => userApi.getCurrentUser(),
        throwsException,
      );
    });

    test('getCurrentUser throws exception on 404', () async {
      when(mockDio.get('/api/whoami')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 404,
        ),
      ));

      expect(
        () => userApi.getCurrentUser(),
        throwsException,
      );
    });

    test('getUserById returns user data on success', () async {
      final responseData = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
      };

      when(mockDio.get('/api/users/user-1')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: responseData,
          ));

      final result = await userApi.getUserById('user-1');

      expect(result['id'], 'user-1');
      verify(mockDio.get('/api/users/user-1')).called(1);
    });

    test('getUserById throws exception on 404', () async {
      when(mockDio.get('/api/users/user-1')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 404,
        ),
      ));

      expect(
        () => userApi.getUserById('user-1'),
        throwsException,
      );
    });

    test('updateUser returns updated user data on success', () async {
      final updateData = {'name': 'Updated Name'};
      final responseData = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Updated Name',
      };

      when(mockDio.put(
        '/api/users/user-1',
        data: anyNamed('data'),
        options: anyNamed('options'),
      )).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: responseData,
          ));

      final result = await userApi.updateUser('user-1', updateData);

      expect(result['name'], 'Updated Name');
      verify(mockDio.put(
        '/api/users/user-1',
        data: anyNamed('data'),
        options: anyNamed('options'),
      )).called(1);
    });

    test('updateUser removes null and empty values', () async {
      final updateData = {
        'name': 'Updated Name',
        'phone': null,
        'bio': '',
      };

      when(mockDio.put(
        '/api/users/user-1',
        data: anyNamed('data'),
        options: anyNamed('options'),
      )).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: {'id': 'user-1', 'name': 'Updated Name'},
          ));

      await userApi.updateUser('user-1', updateData);

      verify(mockDio.put(
        '/api/users/user-1',
        data: argThat(
          predicate<Map<String, dynamic>>((data) =>
              data.containsKey('name') &&
              !data.containsKey('phone') &&
              !data.containsKey('bio')),
        ),
        options: anyNamed('options'),
      )).called(1);
    });

    test('updateUser throws exception on 400', () async {
      when(mockDio.put(
        '/api/users/user-1',
        data: anyNamed('data'),
        options: anyNamed('options'),
      )).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        response: Response(
          requestOptions: RequestOptions(path: '/'),
          statusCode: 400,
          data: {'message': 'Invalid data'},
        ),
      ));

      expect(
        () => userApi.updateUser('user-1', {'name': 'Test'}),
        throwsException,
      );
    });

    test('updateCurrentUser gets current user first then updates', () async {
      final currentUserData = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
      };
      final updatedData = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Updated Name',
      };

      // Set up mocks first, then create instance
      when(mockDio.get('/api/whoami')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: currentUserData,
          ));

      when(mockDio.put(
        '/api/users/user-1',
        data: anyNamed('data'),
        options: anyNamed('options'),
      )).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: updatedData,
          ));

      final result = await userApi.updateCurrentUser({'name': 'Updated Name'});

      expect(result['name'], 'Updated Name');
      verify(mockDio.get('/api/whoami')).called(1);
      verify(mockDio.put(
        '/api/users/user-1',
        data: anyNamed('data'),
        options: anyNamed('options'),
      )).called(1);
    });

    test('updateCurrentUser throws exception when user ID not found', () async {
      when(mockDio.get('/api/whoami')).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/'),
            statusCode: 200,
            data: {'email': 'test@example.com'}, // Missing 'id'
          ));

      expect(
        () => userApi.updateCurrentUser({'name': 'Test'}),
        throwsException,
      );
    });

    test('handles network errors gracefully', () async {
      when(mockDio.get('/api/whoami')).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        type: DioExceptionType.connectionError,
        error: 'Network error',
      ));

      expect(
        () => userApi.getCurrentUser(),
        throwsException,
      );
    });
  });
}

