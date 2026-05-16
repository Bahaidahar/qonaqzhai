import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../models/user_dto.dart';

class AuthRemoteDataSource {
  AuthRemoteDataSource(this._dio);

  final Dio _dio;

  Future<AuthResponseDto> signup({
    required String email,
    required String password,
    required String role,
    String? name,
  }) async {
    final res = await _dio.post(
      ApiEndpoints.signup,
      data: {'email': email, 'password': password, 'role': role, if (name != null) 'name': name},
      options: Options(extra: {'requiresAuth': false}),
    );
    return AuthResponseDto.fromJson(res.data as Map<String, dynamic>);
  }

  Future<AuthResponseDto> login(String email, String password) async {
    final res = await _dio.post(
      ApiEndpoints.login,
      data: {'email': email, 'password': password},
      options: Options(extra: {'requiresAuth': false}),
    );
    return AuthResponseDto.fromJson(res.data as Map<String, dynamic>);
  }

  Future<UserDto> me() async {
    final res = await _dio.get(ApiEndpoints.me);
    return UserDto.fromJson(res.data as Map<String, dynamic>);
  }

  Future<void> logout(String refreshToken) async {
    await _dio.post(
      ApiEndpoints.logout,
      data: {'refreshToken': refreshToken},
      options: Options(extra: {'requiresAuth': false}),
    );
  }

  Future<void> forgotPassword(String email) async {
    await _dio.post(
      ApiEndpoints.forgotPassword,
      data: {'email': email},
      options: Options(extra: {'requiresAuth': false}),
    );
  }

  Future<void> resetPassword(String token, String newPassword) async {
    await _dio.post(
      ApiEndpoints.resetPassword,
      data: {'token': token, 'newPassword': newPassword},
      options: Options(extra: {'requiresAuth': false}),
    );
  }

  Future<void> registerPushToken(String token, String platform) async {
    await _dio.post(
      ApiEndpoints.fcmTokens,
      data: {'token': token, 'platform': platform},
    );
  }
}
