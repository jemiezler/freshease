import '../../domain/repositories/auth_repository.dart';
import '../../domain/entities/user.dart';
import '../models/user_dto.dart';
import '../sources/auth_api.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthApi api;
  AuthRepositoryImpl(this.api);

  @override
  Future<User> login({required String email, required String password}) async {
    final json = await api.login(email, password);
    return UserDto.fromJson(json).toEntity();
  }
}
