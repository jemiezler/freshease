import 'package:flutter_test/flutter_test.dart';
import 'package:dio/dio.dart';
import 'package:frontend/core/network/api_helper.dart';

void main() {
  group('ApiHelper', () {
    test('handleDioError returns authentication error for 401', () {
      final dioException = DioException(
        requestOptions: RequestOptions(path: '/api/test'),
        response: Response(
          requestOptions: RequestOptions(path: '/api/test'),
          statusCode: 401,
          data: {'message': 'Unauthorized'},
        ),
      );

      final result = ApiHelper.handleDioError(dioException);

      expect(result.toString(), contains('Authentication required'));
    });

    test('handleDioError returns forbidden error for 403', () {
      final dioException = DioException(
        requestOptions: RequestOptions(path: '/api/test'),
        response: Response(
          requestOptions: RequestOptions(path: '/api/test'),
          statusCode: 403,
          data: {'message': 'Forbidden'},
        ),
      );

      final result = ApiHelper.handleDioError(dioException);

      expect(result.toString(), contains('Access forbidden'));
    });

    test('handleDioError returns not found error for 404', () {
      final dioException = DioException(
        requestOptions: RequestOptions(path: '/api/test'),
        response: Response(
          requestOptions: RequestOptions(path: '/api/test'),
          statusCode: 404,
        ),
      );

      final result = ApiHelper.handleDioError(dioException);

      expect(result.toString(), contains('Resource not found'));
    });

    test('handleDioError returns bad request error for 400', () {
      final dioException = DioException(
        requestOptions: RequestOptions(path: '/api/test'),
        response: Response(
          requestOptions: RequestOptions(path: '/api/test'),
          statusCode: 400,
          data: {'message': 'Invalid input'},
        ),
      );

      final result = ApiHelper.handleDioError(dioException);

      expect(result.toString(), contains('Invalid request'));
    });

    test('handleDioError returns server error for 500', () {
      final dioException = DioException(
        requestOptions: RequestOptions(path: '/api/test'),
        response: Response(
          requestOptions: RequestOptions(path: '/api/test'),
          statusCode: 500,
          data: {'message': 'Internal server error'},
        ),
      );

      final result = ApiHelper.handleDioError(dioException);

      expect(result.toString(), contains('Server error'));
    });

    test('handleDioError returns network error when no response', () {
      final dioException = DioException(
        requestOptions: RequestOptions(path: '/api/test'),
        type: DioExceptionType.connectionTimeout,
        message: 'Connection timeout',
      );

      final result = ApiHelper.handleDioError(dioException);

      expect(result.toString(), contains('Network Error'));
    });

    test('extractData extracts data from wrapped response', () {
      final responseData = {
        'data': {
          'id': '1',
          'name': 'Test',
        },
      };

      final result = ApiHelper.extractData<Map<String, dynamic>>(
        responseData,
        (json) => json,
      );

      expect(result['id'], '1');
      expect(result['name'], 'Test');
    });

    test('extractData extracts data from direct response', () {
      final responseData = {
        'id': '1',
        'name': 'Test',
      };

      final result = ApiHelper.extractData<Map<String, dynamic>>(
        responseData,
        (json) => json,
      );

      expect(result['id'], '1');
      expect(result['name'], 'Test');
    });

    test('extractData throws exception for invalid format', () {
      expect(
        () => ApiHelper.extractData<Map<String, dynamic>>(
          'invalid',
          (json) => json,
        ),
        throwsException,
      );
    });

    test('extractList extracts list from wrapped response', () {
      final responseData = {
        'data': [
          {'id': '1', 'name': 'Item 1'},
          {'id': '2', 'name': 'Item 2'},
        ],
      };

      final result = ApiHelper.extractList<Map<String, dynamic>>(
        responseData,
        (json) => json,
      );

      expect(result.length, 2);
      expect(result[0]['id'], '1');
      expect(result[1]['id'], '2');
    });

    test('extractList extracts list from direct response', () {
      final responseData = [
        {'id': '1', 'name': 'Item 1'},
        {'id': '2', 'name': 'Item 2'},
      ];

      final result = ApiHelper.extractList<Map<String, dynamic>>(
        responseData,
        (json) => json,
      );

      expect(result.length, 2);
      expect(result[0]['id'], '1');
    });

    test('extractList throws exception for invalid format', () {
      expect(
        () => ApiHelper.extractList<Map<String, dynamic>>(
          'invalid',
          (json) => json,
        ),
        throwsException,
      );
    });

    test('extractList handles wrapped list in data field', () {
      final responseData = {
        'data': [
          {'id': '1', 'name': 'Item 1'},
        ],
      };

      final result = ApiHelper.extractList<Map<String, dynamic>>(
        responseData,
        (json) => json,
      );

      expect(result.length, 1);
      expect(result[0]['id'], '1');
    });
  });
}

