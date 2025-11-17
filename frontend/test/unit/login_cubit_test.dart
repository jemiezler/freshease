import 'package:flutter_test/flutter_test.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:frontend/features/auth/presentation/state/login_cubit.dart';
import 'package:frontend/features/auth/domain/repositories/auth_repository.dart';
import 'package:frontend/features/auth/domain/entities/user.dart';

import 'login_cubit_test.mocks.dart';

@GenerateMocks([AuthRepository])
void main() {
  late MockAuthRepository mockRepository;
  late LoginCubit loginCubit;

  setUp(() {
    mockRepository = MockAuthRepository();
    loginCubit = LoginCubit(mockRepository);
  });

  tearDown(() {
    loginCubit.close();
  });

  group('LoginCubit', () {
    test('initial state is correct', () {
      expect(loginCubit.state, const LoginState());
      expect(loginCubit.state.loading, false);
      expect(loginCubit.state.user, null);
      expect(loginCubit.state.error, null);
    });

    blocTest<LoginCubit, LoginState>(
      'emits loading then success when googleLogin succeeds',
      build: () {
        when(mockRepository.signInWithGoogle()).thenAnswer(
          (_) async => User(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
          ),
        );
        return loginCubit;
      },
      act: (cubit) => cubit.googleLogin(),
      expect: () => [
        const LoginState(
          loading: true,
          error: null,
          isSuccess: false,
          isFailure: false,
        ),
        LoginState(
          loading: false,
          user: User(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
          ),
          isSuccess: true,
          isFailure: false,
        ),
      ],
      verify: (_) {
        verify(mockRepository.signInWithGoogle()).called(1);
      },
    );

    blocTest<LoginCubit, LoginState>(
      'emits loading then failure when googleLogin fails',
      build: () {
        when(mockRepository.signInWithGoogle()).thenThrow(
          Exception('Authentication failed'),
        );
        return loginCubit;
      },
      act: (cubit) => cubit.googleLogin(),
      expect: () => [
        const LoginState(
          loading: true,
          error: null,
          isSuccess: false,
          isFailure: false,
        ),
        LoginState(
          loading: false,
          error: 'Exception: Authentication failed',
          isSuccess: false,
          isFailure: true,
        ),
      ],
    );

    blocTest<LoginCubit, LoginState>(
      'emits loading then success when verifyExistingToken succeeds',
      build: () {
        when(mockRepository.verifyToken()).thenAnswer(
          (_) async => User(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
          ),
        );
        return loginCubit;
      },
      act: (cubit) => cubit.verifyExistingToken(),
      expect: () => [
        const LoginState(
          loading: true,
          error: null,
          isSuccess: false,
          isFailure: false,
        ),
        LoginState(
          loading: false,
          user: User(
            id: '1',
            email: 'test@example.com',
            name: 'Test User',
          ),
          isSuccess: true,
          isFailure: false,
        ),
      ],
      verify: (_) {
        verify(mockRepository.verifyToken()).called(1);
      },
    );

    blocTest<LoginCubit, LoginState>(
      'clearError emits state with null error',
      build: () {
        loginCubit.emit(const LoginState(error: 'Some error', isFailure: true));
        return loginCubit;
      },
      act: (cubit) => cubit.clearError(),
      expect: () => [
        const LoginState(error: null, isFailure: false),
      ],
    );
  });
}

