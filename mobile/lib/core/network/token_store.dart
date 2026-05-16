import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Persists access + refresh tokens using the platform secure store.
class TokenStore {
  TokenStore([this._storage = const FlutterSecureStorage()]);

  final FlutterSecureStorage _storage;

  static const _kAccess = 'qz_access';
  static const _kRefresh = 'qz_refresh';

  Future<void> saveAccess(String token) => _storage.write(key: _kAccess, value: token);
  Future<void> saveRefresh(String token) => _storage.write(key: _kRefresh, value: token);
  Future<String?> readAccess() => _storage.read(key: _kAccess);
  Future<String?> readRefresh() => _storage.read(key: _kRefresh);

  Future<void> clear() async {
    await _storage.delete(key: _kAccess);
    await _storage.delete(key: _kRefresh);
  }
}
