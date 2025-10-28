// data/repositories/user_repository_impl.dart
import '../../domain/repositories/user_repository.dart';
import '../../domain/entities/user_profile.dart';
import '../models/user_profile_dto.dart';
import '../sources/user_api.dart';

class UserRepositoryImpl implements UserRepository {
  final UserApi api;
  UserRepositoryImpl(this.api);

  @override
  Future<UserProfile> getCurrentUser() async {
    final json = await api.getCurrentUser();
    return UserProfileDto.fromJson(json).toEntity();
  }

  @override
  Future<UserProfile> getUserById(String userId) async {
    final json = await api.getUserById(userId);
    return UserProfileDto.fromJson(json).toEntity();
  }

  @override
  Future<UserProfile> updateUser(
    String userId,
    Map<String, dynamic> userData,
  ) async {
    final json = await api.updateUser(userId, userData);
    return UserProfileDto.fromJson(json).toEntity();
  }

  @override
  Future<UserProfile> updateCurrentUser(Map<String, dynamic> userData) async {
    final json = await api.updateCurrentUser(userData);
    return UserProfileDto.fromJson(json).toEntity();
  }
}
