// presentation/state/user_cubit.dart
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../domain/repositories/user_repository.dart';
import '../../domain/entities/user_profile.dart';

class UserState extends Equatable {
  final UserProfile? user;
  final bool loading;
  final String? error;
  final bool isUpdating;

  const UserState({
    this.user,
    this.loading = false,
    this.error,
    this.isUpdating = false,
  });

  UserState copyWith({
    UserProfile? user,
    bool? loading,
    String? error,
    bool? isUpdating,
  }) {
    return UserState(
      user: user ?? this.user,
      loading: loading ?? this.loading,
      error: error,
      isUpdating: isUpdating ?? this.isUpdating,
    );
  }

  @override
  List<Object?> get props => [user, loading, error, isUpdating];
}

class UserCubit extends Cubit<UserState> {
  final UserRepository repository;

  UserCubit(this.repository) : super(const UserState());

  /// Load current user profile
  Future<void> loadCurrentUser() async {
    if (isClosed) return;
    emit(state.copyWith(loading: true, error: null));
    try {
      final user = await repository.getCurrentUser();
      if (!isClosed) {
        emit(state.copyWith(user: user, loading: false));
      }
    } catch (e) {
      if (!isClosed) {
        // Check if it's an authentication error
        if (e.toString().contains('401') ||
            e.toString().contains('missing bearer token')) {
          // Clear stored tokens and emit auth error
          final prefs = await SharedPreferences.getInstance();
          await prefs.remove('access_token');
          await prefs.remove('refresh_token');
          await prefs.remove('id_token');

          emit(
            state.copyWith(
              loading: false,
              error: 'Authentication expired. Please login again.',
            ),
          );
        } else {
          emit(
            state.copyWith(
              loading: false,
              error: 'Failed to load user profile: ${e.toString()}',
            ),
          );
        }
      }
    }
  }

  /// Update user profile
  Future<void> updateProfile(Map<String, dynamic> userData) async {
    if (state.user == null || isClosed) return;

    emit(state.copyWith(isUpdating: true, error: null));
    try {
      final updatedUser = await repository.updateCurrentUser(userData);
      if (!isClosed) {
        emit(state.copyWith(user: updatedUser, isUpdating: false));
      }
    } catch (e) {
      if (!isClosed) {
        // Check if it's an authentication error
        if (e.toString().contains('401') ||
            e.toString().contains('missing bearer token') ||
            e.toString().contains('Authentication expired')) {
          // Clear stored tokens and emit auth error
          final prefs = await SharedPreferences.getInstance();
          await prefs.remove('access_token');
          await prefs.remove('refresh_token');
          await prefs.remove('id_token');

          emit(
            state.copyWith(
              isUpdating: false,
              error: 'Authentication expired. Please login again.',
            ),
          );
        } else {
          emit(
            state.copyWith(
              isUpdating: false,
              error: 'Failed to update profile: ${e.toString()}',
            ),
          );
        }
      }
    }
  }

  /// Update specific fields
  Future<void> updateName(String name) async {
    await updateProfile({'name': name});
  }

  Future<void> updatePhone(String phone) async {
    await updateProfile({'phone': phone});
  }

  Future<void> updateBio(String bio) async {
    await updateProfile({'bio': bio});
  }

  Future<void> updateAvatar(String avatar) async {
    await updateProfile({'avatar': avatar});
  }

  /// Clear error
  void clearError() {
    emit(state.copyWith(error: null));
  }

  /// Clear authentication and redirect to login
  Future<void> clearAuthAndRedirect() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('access_token');
    await prefs.remove('refresh_token');
    await prefs.remove('id_token');

    emit(state.copyWith(user: null, loading: false, error: null));
  }
}
