import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/auth_remote_datasource.dart';
import '../../data/repositories/auth_repository_impl.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final dio = ref.watch(dioProvider);
  final tokens = ref.watch(tokenStoreProvider);
  return AuthRepositoryImpl(AuthRemoteDataSource(dio), tokens);
});

class AuthState {
  const AuthState({this.user, this.loading = false, this.error});

  final User? user;
  final bool loading;
  final String? error;

  AuthState copyWith({User? user, bool? loading, String? error, bool clearUser = false, bool clearError = false}) {
    return AuthState(
      user: clearUser ? null : (user ?? this.user),
      loading: loading ?? this.loading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class AuthViewModel extends StateNotifier<AuthState> {
  AuthViewModel(this._repo) : super(const AuthState(loading: true)) {
    _bootstrap();
  }

  final AuthRepository _repo;

  Future<void> _bootstrap() async {
    final u = await _repo.me();
    state = AuthState(user: u, loading: false);
  }

  Future<void> signup({
    required String email,
    required String password,
    required UserRole role,
    String? name,
  }) async {
    state = state.copyWith(loading: true, clearError: true);
    try {
      final u = await _repo.signup(email: email, password: password, role: role, name: name);
      state = AuthState(user: u, loading: false);
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
    }
  }

  Future<void> login(String email, String password) async {
    state = state.copyWith(loading: true, clearError: true);
    try {
      final u = await _repo.login(email: email, password: password);
      state = AuthState(user: u, loading: false);
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
    }
  }

  Future<void> logout() async {
    await _repo.logout();
    state = const AuthState();
  }

  Future<bool> forgotPassword(String email) async {
    try {
      await _repo.forgotPassword(email);
      return true;
    } catch (_) {
      return false;
    }
  }

  Future<bool> resetPassword({required String token, required String newPassword}) async {
    try {
      await _repo.resetPassword(token: token, newPassword: newPassword);
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  Future<void> registerPushToken(String token, {String platform = 'unknown'}) =>
      _repo.registerPushToken(token, platform: platform);
}

final authViewModelProvider = StateNotifierProvider<AuthViewModel, AuthState>((ref) {
  final repo = ref.watch(authRepositoryProvider);
  return AuthViewModel(repo);
});
