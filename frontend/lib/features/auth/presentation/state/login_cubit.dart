// presentation/login/login_cubit.dart
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../domain/repositories/auth_repository.dart';
import '../../domain/entities/user.dart';

class LoginState extends Equatable {
  final bool loading;
  final String? error;
  final User? user;
  final bool isSuccess;
  final bool isFailure;

  const LoginState({
    this.loading = false,
    this.error,
    this.user,
    this.isSuccess = false,
    this.isFailure = false,
  });

  LoginState copyWith({
    bool? loading,
    String? error,
    User? user,
    bool? isSuccess,
    bool? isFailure,
  }) => LoginState(
    loading: loading ?? this.loading,
    error: error,
    user: user,
    isSuccess: isSuccess ?? this.isSuccess,
    isFailure: isFailure ?? this.isFailure,
  );

  @override
  List<Object?> get props => [loading, error, user, isSuccess, isFailure];
}

class LoginCubit extends Cubit<LoginState> {
  final AuthRepository repo;
  LoginCubit(this.repo) : super(const LoginState());

  Future<void> googleLogin() async {
    if (isClosed) return;
    emit(
      state.copyWith(
        loading: true,
        error: null,
        isSuccess: false,
        isFailure: false,
      ),
    );
    try {
      final u = await repo.signInWithGoogle();
      if (!isClosed) {
        emit(LoginState(user: u, isSuccess: true, isFailure: false));
      }
    } catch (e) {
      if (!isClosed) {
        emit(
          LoginState(error: e.toString(), isSuccess: false, isFailure: true),
        );
      }
    }
  }

  /// Clear error state
  void clearError() {
    if (!isClosed) {
      emit(state.copyWith(error: null, isFailure: false));
    }
  }
}
