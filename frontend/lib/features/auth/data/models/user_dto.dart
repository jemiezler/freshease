// data/models/user_dto.dart
import '../../domain/entities/user.dart';

class UserDto {
  final String id;
  final String email;
  final String? name;
  final String? avatar;

  UserDto({required this.id, required this.email, this.name, this.avatar});

  factory UserDto.fromJson(Map<String, dynamic> json) => UserDto(
    id: json['id'] as String,
    email: json['email'] as String,
    name: json['name'] as String?,
    avatar: json['avatar'] as String?,
  );

  User toEntity() => User(id: id, email: email, name: name, avatar: avatar);
}
