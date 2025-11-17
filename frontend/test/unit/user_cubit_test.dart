import 'package:flutter_test/flutter_test.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/account/presentation/state/user_cubit.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';

import 'user_cubit_test.mocks.dart';

@GenerateMocks([UserRepository])
void main() {
  late MockUserRepository mockRepository;
  late UserCubit userCubit;

  setUp(() {
    mockRepository = MockUserRepository();
    userCubit = UserCubit(mockRepository);
  });

  tearDown(() {
    userCubit.close();
  });

  group('UserCubit', () {
    test('initial state is correct', () {
      expect(userCubit.state, const UserState());
      expect(userCubit.state.loading, false);
      expect(userCubit.state.user, null);
      expect(userCubit.state.error, null);
    });

    blocTest<UserCubit, UserState>(
      'emits loading then success when loadCurrentUser succeeds',
      build: () {
        final fixedTime = DateTime(2025, 1, 1, 12, 0, 0);
        when(mockRepository.getCurrentUser()).thenAnswer(
          (_) async => UserProfile(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
            status: 'active',
            createdAt: fixedTime,
            updatedAt: fixedTime,
          ),
        );
        return userCubit;
      },
      act: (cubit) => cubit.loadCurrentUser(),
      expect: () {
        final fixedTime = DateTime(2025, 1, 1, 12, 0, 0);
        return [
          const UserState(loading: true, error: null),
          UserState(
            loading: false,
            user: UserProfile(
              id: '1',
              email: 'test@example.com',
              name: 'Test User',
              status: 'active',
              createdAt: fixedTime,
              updatedAt: fixedTime,
            ),
          ),
        ];
      },
      verify: (_) {
        verify(mockRepository.getCurrentUser()).called(1);
      },
    );

    blocTest<UserCubit, UserState>(
      'emits loading then error when loadCurrentUser fails',
      build: () {
        when(mockRepository.getCurrentUser()).thenThrow(
          Exception('Network error'),
        );
        return userCubit;
      },
      act: (cubit) => cubit.loadCurrentUser(),
      expect: () => [
        const UserState(loading: true, error: null),
        UserState(
          loading: false,
          error: 'Failed to load user profile: Exception: Network error',
        ),
      ],
    );

    blocTest<UserCubit, UserState>(
      'emits updating then success when updateProfile succeeds',
      build: () {
        final fixedTime = DateTime(2025, 1, 1, 12, 0, 0);
        final updatedTime = DateTime(2025, 1, 1, 12, 1, 0);
        when(mockRepository.getCurrentUser()).thenAnswer(
          (_) async => UserProfile(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
            status: 'active',
            createdAt: fixedTime,
            updatedAt: fixedTime,
          ),
        );
        when(mockRepository.updateCurrentUser(any)).thenAnswer(
          (_) async => UserProfile(
            id: '1',
            email: 'test@example.com',
            name: 'Updated User',
            status: 'active',
            createdAt: fixedTime,
            updatedAt: updatedTime,
          ),
        );
        return userCubit;
      },
      setUp: () {
        final fixedTime = DateTime(2025, 1, 1, 12, 0, 0);
        userCubit.emit(UserState(
          user: UserProfile(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
            status: 'active',
            createdAt: fixedTime,
            updatedAt: fixedTime,
          ),
        ));
      },
      act: (cubit) => cubit.updateProfile({'name': 'Updated User'}),
      expect: () {
        final fixedTime = DateTime(2025, 1, 1, 12, 0, 0);
        final updatedTime = DateTime(2025, 1, 1, 12, 1, 0);
        return [
          UserState(
            user: UserProfile(
              id: '1',
              email: 'test@example.com',
              name: 'Test User',
              status: 'active',
              createdAt: fixedTime,
              updatedAt: fixedTime,
            ),
            isUpdating: true,
            error: null,
          ),
          UserState(
            user: UserProfile(
              id: '1',
              email: 'test@example.com',
              name: 'Updated User',
              status: 'active',
              createdAt: fixedTime,
              updatedAt: updatedTime,
            ),
            isUpdating: false,
          ),
        ];
      },
      verify: (_) {
        verify(mockRepository.updateCurrentUser({'name': 'Updated User'})).called(1);
      },
    );

    blocTest<UserCubit, UserState>(
      'clearError emits state with null error',
      build: () {
        userCubit.emit(const UserState(error: 'Some error'));
        return userCubit;
      },
      act: (cubit) => cubit.clearError(),
      expect: () => [
        const UserState(error: null),
      ],
    );
  });
}

