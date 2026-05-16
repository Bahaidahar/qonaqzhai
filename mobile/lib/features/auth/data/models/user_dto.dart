import '../../domain/entities/user.dart';

class UserDto {
  UserDto({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
    required this.status,
  });

  final String id;
  final String email;
  final String name;
  final String role;
  final String status;

  factory UserDto.fromJson(Map<String, dynamic> json) => UserDto(
        id: json['id'] as String,
        email: json['email'] as String,
        name: (json['name'] as String?) ?? '',
        role: json['role'] as String,
        status: (json['status'] as String?) ?? 'active',
      );

  User toDomain() => User(
        id: id,
        email: email,
        name: name,
        role: UserRoleX.parse(role),
        status: status,
      );
}

class AuthResponseDto {
  AuthResponseDto({
    required this.accessToken,
    required this.refreshToken,
    required this.user,
  });

  final String accessToken;
  final String? refreshToken;
  final UserDto user;

  factory AuthResponseDto.fromJson(Map<String, dynamic> json) => AuthResponseDto(
        accessToken: (json['accessToken'] ?? json['token']) as String,
        refreshToken: json['refreshToken'] as String?,
        user: UserDto.fromJson(json['user'] as Map<String, dynamic>),
      );
}
