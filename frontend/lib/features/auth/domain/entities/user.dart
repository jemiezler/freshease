// domain/entities/user.dart
import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String email;
  final String? name;
  final String? avatar;

  const User({required this.id, required this.email, this.name, this.avatar});

  @override
  List<Object?> get props => [id, email, name, avatar];
}
