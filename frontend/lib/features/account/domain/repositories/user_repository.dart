// domain/repositories/user_repository.dart
import '../entities/user_profile.dart';

abstract class UserRepository {
  Future<UserProfile> getCurrentUser();
  Future<UserProfile> getUserById(String userId);
  Future<UserProfile> updateUser(String userId, Map<String, dynamic> userData);
  Future<UserProfile> updateCurrentUser(Map<String, dynamic> userData);
}
