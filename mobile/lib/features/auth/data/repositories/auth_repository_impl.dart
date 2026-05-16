import '../../../../core/network/token_store.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_datasource.dart';

class AuthRepositoryImpl implements AuthRepository {
  AuthRepositoryImpl(this._remote, this._tokens);

  final AuthRemoteDataSource _remote;
  final TokenStore _tokens;

  Future<User> _persistAndReturn(dynamic res) async {
    final dto = res;
    await _tokens.saveAccess(dto.accessToken as String);
    if (dto.refreshToken is String) {
      await _tokens.saveRefresh(dto.refreshToken as String);
    }
    return (dto.user as dynamic).toDomain() as User;
  }

  @override
  Future<User> signup({
    required String email,
    required String password,
    required UserRole role,
    String? name,
  }) async {
    final dto = await _remote.signup(
      email: email,
      password: password,
      role: role.name,
      name: name,
    );
    return _persistAndReturn(dto);
  }

  @override
  Future<User> login({required String email, required String password}) async {
    final dto = await _remote.login(email, password);
    return _persistAndReturn(dto);
  }

  @override
  Future<User?> me() async {
    final token = await _tokens.readAccess();
    if (token == null) return null;
    try {
      final dto = await _remote.me();
      return dto.toDomain();
    } catch (_) {
      return null;
    }
  }

  @override
  Future<void> logout() async {
    final refresh = await _tokens.readRefresh();
    if (refresh != null) {
      try {
        await _remote.logout(refresh);
      } catch (_) {/* best-effort */}
    }
    await _tokens.clear();
  }

  @override
  Future<void> forgotPassword(String email) => _remote.forgotPassword(email);

  @override
  Future<void> resetPassword({required String token, required String newPassword}) =>
      _remote.resetPassword(token, newPassword);

  @override
  Future<void> registerPushToken(String token, {String platform = 'unknown'}) =>
      _remote.registerPushToken(token, platform);
}
