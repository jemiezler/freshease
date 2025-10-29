// test/unit/user_cubit_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';
import 'package:frontend/features/account/presentation/state/user_cubit.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';
import '../mocks/mock_repositories.dart';
import '../mocks/mock_repositories.mocks.dart';

void main() {
  group('UserCubit', () {
    late UserCubit userCubit;
    late MockUserRepository mockUserRepository;

    setUp(() {
      mockUserRepository = MockUserRepository();
      userCubit = UserCubit(mockUserRepository);
    });

    tearDown(() {
      userCubit.close();
    });

    group('Initial State', () {
      test('should have initial state as not loading', () {
        expect(userCubit.state.loading, false);
        expect(userCubit.state.user, null);
        expect(userCubit.state.error, null);
      });
    });

    group('Load User Profile', () {
      blocTest<UserCubit, UserState>(
        'should emit loaded state when user profile is loaded successfully',
        build: () {
          when(
            mockUserRepository.getCurrentUser(),
          ).thenAnswer((_) async => MockData.mockUserProfile);
          return userCubit;
        },
        act: (cubit) => cubit.loadCurrentUser(),
        expect: () => [
          const UserState(loading: true, error: null),
          UserState(user: MockData.mockUserProfile, loading: false),
        ],
      );

      blocTest<UserCubit, UserState>(
        'should emit error state when user profile loading fails',
        build: () {
          when(
            mockUserRepository.getCurrentUser(),
          ).thenThrow(Exception('Failed to load user'));
          return userCubit;
        },
        act: (cubit) => cubit.loadCurrentUser(),
        expect: () => [
          const UserState(loading: true, error: null),
          const UserState(
            loading: false,
            error:
                'Failed to load user profile: Exception: Failed to load user',
          ),
        ],
      );
    });

    group('Update User Profile', () {
      blocTest<UserCubit, UserState>(
        'should emit updated state when user profile is updated successfully',
        build: () {
          when(
            mockUserRepository.getCurrentUser(),
          ).thenAnswer((_) async => MockData.mockUserProfile);
          when(
            mockUserRepository.updateCurrentUser(any),
          ).thenAnswer((_) async => MockData.mockUserProfile);
          return userCubit;
        },
        act: (cubit) async {
          // First load the user
          await cubit.loadCurrentUser();
          // Then update profile
          await cubit.updateProfile({
            'name': 'Updated Name',
            'bio': 'Updated bio',
          });
        },
        expect: () => [
          const UserState(loading: true, error: null),
          UserState(user: MockData.mockUserProfile, loading: false),
          UserState(user: MockData.mockUserProfile, isUpdating: true),
          UserState(user: MockData.mockUserProfile, isUpdating: false),
        ],
      );

      blocTest<UserCubit, UserState>(
        'should emit error state when user profile update fails',
        build: () {
          when(
            mockUserRepository.getCurrentUser(),
          ).thenAnswer((_) async => MockData.mockUserProfile);
          when(
            mockUserRepository.updateCurrentUser(any),
          ).thenThrow(Exception('Failed to update user'));
          return userCubit;
        },
        act: (cubit) async {
          // First load the user
          await cubit.loadCurrentUser();
          // Then update profile
          await cubit.updateProfile({'name': 'Updated Name'});
        },
        expect: () => [
          const UserState(loading: true, error: null),
          UserState(user: MockData.mockUserProfile, loading: false),
          UserState(user: MockData.mockUserProfile, isUpdating: true),
          UserState(
            user: MockData.mockUserProfile,
            isUpdating: false,
            error: 'Failed to update profile: Exception: Failed to update user',
          ),
        ],
      );
    });

    group('Specific Field Updates', () {
      blocTest<UserCubit, UserState>(
        'should update name correctly',
        build: () {
          when(
            mockUserRepository.getCurrentUser(),
          ).thenAnswer((_) async => MockData.mockUserProfile);
          when(
            mockUserRepository.updateCurrentUser(any),
          ).thenAnswer((_) async => MockData.mockUserProfile);
          return userCubit;
        },
        act: (cubit) async {
          // First load the user
          await cubit.loadCurrentUser();
          // Then update name
          await cubit.updateName('New Name');
        },
        expect: () => [
          const UserState(loading: true, error: null),
          UserState(user: MockData.mockUserProfile, loading: false),
          UserState(user: MockData.mockUserProfile, isUpdating: true),
          UserState(user: MockData.mockUserProfile, isUpdating: false),
        ],
      );
    });

    group('Error Handling', () {
      test('should clear error correctly', () {
        // Arrange
        userCubit.emit(const UserState(error: 'Test error'));

        // Act
        userCubit.clearError();

        // Assert
        expect(userCubit.state.error, null);
      });
    });
  });
}
