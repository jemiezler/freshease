// data/repositories/auth_repository_impl.dart
import '../../domain/repositories/auth_repository.dart';
import '../../domain/entities/user.dart';
import '../models/user_dto.dart';
import '../sources/auth_api.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthApi api;
  AuthRepositoryImpl(this.api);

  @override
  Future<User> signInWithGoogle() async {
    final json = await api.signInWithGoogle();
    return UserDto.fromJson(json).toEntity();
  }
}
