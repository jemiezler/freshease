// data/models/user_profile_dto.dart
import '../../domain/entities/user_profile.dart';

class UserProfileDto {
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

  UserProfileDto({
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

  factory UserProfileDto.fromJson(Map<String, dynamic> json) {
    return UserProfileDto(
      id: json['id'] as String,
      email: json['email'] as String,
      name: json['name'] as String,
      phone: json['phone'] as String?,
      bio: json['bio'] as String?,
      avatar: json['avatar'] as String?,
      cover: json['cover'] as String?,
      status: json['status'] as String? ?? 'active',
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'phone': phone,
      'bio': bio,
      'avatar': avatar,
      'cover': cover,
      'status': status,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  Map<String, dynamic> toUpdateJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'phone': phone,
      'bio': bio,
      'avatar': avatar,
      'cover': cover,
      'status': status,
    };
  }

  UserProfile toEntity() {
    return UserProfile(
      id: id,
      email: email,
      name: name,
      phone: phone,
      bio: bio,
      avatar: avatar,
      cover: cover,
      status: status,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }
}
