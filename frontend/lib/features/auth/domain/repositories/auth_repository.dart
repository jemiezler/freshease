// domain/repositories/auth_repository.dart
import '../entities/user.dart';

abstract class AuthRepository {
  Future<User> signInWithGoogle();
  Future<User> verifyToken();
}
