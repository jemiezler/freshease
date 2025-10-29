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
  final DateTime? dateOfBirth;
  final String? sex;
  final String? goal;
  final double? heightCm;
  final double? weightKg;
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
    this.dateOfBirth,
    this.sex,
    this.goal,
    this.heightCm,
    this.weightKg,
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
      dateOfBirth: json['date_of_birth'] != null
          ? DateTime.parse(json['date_of_birth'] as String)
          : null,
      sex: json['sex'] as String?,
      goal: json['goal'] as String?,
      heightCm: json['height_cm'] != null
          ? (json['height_cm'] as num).toDouble()
          : null,
      weightKg: json['weight_kg'] != null
          ? (json['weight_kg'] as num).toDouble()
          : null,
      status: json['status'] as String? ?? 'active',
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
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
      'date_of_birth': dateOfBirth?.toIso8601String(),
      'sex': sex,
      'goal': goal,
      'height_cm': heightCm,
      'weight_kg': weightKg,
      'status': status,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
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
      'date_of_birth': dateOfBirth?.toIso8601String(),
      'sex': sex,
      'goal': goal,
      'height_cm': heightCm,
      'weight_kg': weightKg,
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
      dateOfBirth: dateOfBirth,
      sex: sex,
      goal: goal,
      heightCm: heightCm,
      weightKg: weightKg,
      status: status,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }
}
