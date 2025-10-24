// presentation/login/login_cubit.dart
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../domain/repositories/auth_repository.dart';
import '../../domain/entities/user.dart';

class LoginState extends Equatable {
  final bool loading;
  final String? error;
  final User? user;
  const LoginState({this.loading = false, this.error, this.user});

  LoginState copyWith({bool? loading, String? error, User? user}) =>
      LoginState(loading: loading ?? this.loading, error: error, user: user);

  @override
  List<Object?> get props => [loading, error, user];
}

class LoginCubit extends Cubit<LoginState> {
  final AuthRepository repo;
  LoginCubit(this.repo) : super(const LoginState());

  Future<void> googleLogin() async {
    emit(state.copyWith(loading: true, error: null));
    try {
      final u = await repo.signInWithGoogle();
      emit(LoginState(user: u));
    } catch (e) {
      emit(LoginState(error: e.toString()));
    }
  }
}
