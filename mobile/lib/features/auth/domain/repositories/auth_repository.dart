import '../entities/user.dart';

abstract class AuthRepository {
  Future<User> signup({
    required String email,
    required String password,
    required UserRole role,
    String? name,
  });

  Future<User> login({required String email, required String password});

  Future<User?> me();

  Future<void> logout();

  Future<void> forgotPassword(String email);

  Future<void> resetPassword({required String token, required String newPassword});

  Future<void> registerPushToken(String token, {String platform});
}
