import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:frontend/features/auth/data/sources/auth_api.dart';

import 'auth_repository_impl_test.mocks.dart';

@GenerateMocks([AuthApi])
void main() {
  group('AuthRepositoryImpl', () {
    late AuthRepositoryImpl repository;
    late MockAuthApi mockApi;

    setUp(() {
      mockApi = MockAuthApi();
      repository = AuthRepositoryImpl(mockApi);
    });

    test('signInWithGoogle returns User from API', () async {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'avatar': 'avatar.jpg',
      };

      when(mockApi.signInWithGoogle()).thenAnswer((_) async => json);

      final result = await repository.signInWithGoogle();

      expect(result.id, 'user-1');
      expect(result.email, 'test@example.com');
      expect(result.name, 'Test User');
      verify(mockApi.signInWithGoogle()).called(1);
    });

    test('signInWithGoogle handles missing optional fields', () async {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
      };

      when(mockApi.signInWithGoogle()).thenAnswer((_) async => json);

      final result = await repository.signInWithGoogle();

      expect(result.id, 'user-1');
      expect(result.email, 'test@example.com');
      expect(result.name, isNull);
      expect(result.avatar, isNull);
    });

    test('verifyToken returns User from API', () async {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
      };

      when(mockApi.verifyToken()).thenAnswer((_) async => json);

      final result = await repository.verifyToken();

      expect(result.id, 'user-1');
      expect(result.email, 'test@example.com');
      verify(mockApi.verifyToken()).called(1);
    });
  });
}

