import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/account/data/repositories/user_repository_impl.dart';
import 'package:frontend/features/account/data/sources/user_api.dart';
import 'package:frontend/features/account/data/models/user_profile_dto.dart';

import 'user_repository_impl_test.mocks.dart';

@GenerateMocks([UserApi])
void main() {
  group('UserRepositoryImpl', () {
    late UserRepositoryImpl repository;
    late MockUserApi mockApi;

    setUp(() {
      mockApi = MockUserApi();
      repository = UserRepositoryImpl(mockApi);
    });

    test('getCurrentUser returns UserProfile from API', () async {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'phone': '+1234567890',
        'bio': 'Test bio',
        'avatar': 'avatar.jpg',
        'cover': 'cover.jpg',
        'date_of_birth': '1990-01-01T00:00:00Z',
        'sex': 'male',
        'goal': 'weight_loss',
        'height_cm': 175.0,
        'weight_kg': 70.0,
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      when(mockApi.getCurrentUser()).thenAnswer((_) async => json);

      final result = await repository.getCurrentUser();

      expect(result.id, 'user-1');
      expect(result.email, 'test@example.com');
      expect(result.name, 'Test User');
      expect(result.phone, '+1234567890');
      verify(mockApi.getCurrentUser()).called(1);
    });

    test('getCurrentUser handles missing optional fields', () async {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Test User',
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      when(mockApi.getCurrentUser()).thenAnswer((_) async => json);

      final result = await repository.getCurrentUser();

      expect(result.id, 'user-1');
      expect(result.email, 'test@example.com');
      expect(result.phone, isNull);
      expect(result.bio, isNull);
    });

    test('getUserById returns UserProfile from API', () async {
      final json = {
        'id': 'user-2',
        'email': 'user2@example.com',
        'name': 'User 2',
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      when(mockApi.getUserById('user-2')).thenAnswer((_) async => json);

      final result = await repository.getUserById('user-2');

      expect(result.id, 'user-2');
      expect(result.email, 'user2@example.com');
      verify(mockApi.getUserById('user-2')).called(1);
    });

    test('updateUser returns updated UserProfile', () async {
      final updateData = {'name': 'Updated Name', 'phone': '+9876543210'};
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Updated Name',
        'phone': '+9876543210',
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-02T00:00:00Z',
      };

      when(mockApi.updateUser('user-1', updateData))
          .thenAnswer((_) async => json);

      final result = await repository.updateUser('user-1', updateData);

      expect(result.name, 'Updated Name');
      expect(result.phone, '+9876543210');
      verify(mockApi.updateUser('user-1', updateData)).called(1);
    });

    test('updateCurrentUser returns updated UserProfile', () async {
      final updateData = {'name': 'Updated Name'};
      final updatedJson = {
        'id': 'user-1',
        'email': 'test@example.com',
        'name': 'Updated Name',
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-02T00:00:00Z',
      };

      when(mockApi.updateCurrentUser(updateData))
          .thenAnswer((_) async => updatedJson);

      final result = await repository.updateCurrentUser(updateData);

      expect(result.name, 'Updated Name');
      verify(mockApi.updateCurrentUser(updateData)).called(1);
    });
  });
}

