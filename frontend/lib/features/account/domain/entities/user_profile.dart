// domain/entities/user_profile.dart
import 'package:equatable/equatable.dart';

class UserProfile extends Equatable {
  final String id;
  final String email;
  final String name;
  final String? phone;
  final String? bio;
  final String? avatar;
  final String? cover;
  final String status;
  final DateTime createdAt;
  final DateTime updatedAt;

  const UserProfile({
    required this.id,
    required this.email,
    required this.name,
    this.phone,
    this.bio,
    this.avatar,
    this.cover,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  @override
  List<Object?> get props => [
    id,
    email,
    name,
    phone,
    bio,
    avatar,
    cover,
    status,
    createdAt,
    updatedAt,
  ];

  /// Get user initials for avatar
  String get initials {
    final nameParts = name.trim().split(' ');
    if (nameParts.isEmpty) return 'U';
    if (nameParts.length == 1) return nameParts[0][0].toUpperCase();
    return '${nameParts[0][0]}${nameParts[1][0]}'.toUpperCase();
  }

  /// Get display name (fallback to email if name is empty)
  String get displayName {
    return name.trim().isNotEmpty ? name : email;
  }

  /// Check if user has complete profile
  bool get isCompleteProfile {
    return name.trim().isNotEmpty &&
        email.trim().isNotEmpty &&
        phone?.trim().isNotEmpty == true;
  }

  UserProfile copyWith({
    String? id,
    String? email,
    String? name,
    String? phone,
    String? bio,
    String? avatar,
    String? cover,
    String? status,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return UserProfile(
      id: id ?? this.id,
      email: email ?? this.email,
      name: name ?? this.name,
      phone: phone ?? this.phone,
      bio: bio ?? this.bio,
      avatar: avatar ?? this.avatar,
      cover: cover ?? this.cover,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
