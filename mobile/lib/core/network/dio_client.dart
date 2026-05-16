import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';

import 'api_endpoints.dart';
import 'api_exception.dart';
import 'token_store.dart';

final tokenStoreProvider = Provider<TokenStore>((ref) => TokenStore());

final dioProvider = Provider<Dio>((ref) {
  final tokens = ref.read(tokenStoreProvider);
  final dio = Dio(BaseOptions(
    baseUrl: ApiEndpoints.baseUrl,
    connectTimeout: const Duration(seconds: 15),
    receiveTimeout: const Duration(seconds: 30),
    contentType: 'application/json',
  ));

  dio.interceptors.add(_AuthInterceptor(tokens, dio));
  dio.interceptors.add(_ErrorInterceptor());
  dio.interceptors.add(PrettyDioLogger(
    requestBody: true,
    responseBody: false,
    requestHeader: false,
  ));
  return dio;
});

class _AuthInterceptor extends Interceptor {
  _AuthInterceptor(this._tokens, this._dio);

  final TokenStore _tokens;
  final Dio _dio;
  bool _refreshing = false;

  @override
  Future<void> onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    if (options.extra['requiresAuth'] != false) {
      final token = await _tokens.readAccess();
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }
    handler.next(options);
  }

  @override
  Future<void> onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && !_refreshing) {
      final refresh = await _tokens.readRefresh();
      if (refresh != null) {
        _refreshing = true;
        try {
          final res = await _dio.post(
            ApiEndpoints.refresh,
            data: {'refreshToken': refresh},
            options: Options(extra: {'requiresAuth': false}),
          );
          final data = res.data as Map<String, dynamic>;
          await _tokens.saveAccess((data['accessToken'] ?? data['token']) as String);
          if (data['refreshToken'] != null) {
            await _tokens.saveRefresh(data['refreshToken'] as String);
          }
          // retry the original request
          final retried = await _dio.fetch(err.requestOptions);
          handler.resolve(retried);
          return;
        } catch (_) {
          await _tokens.clear();
        } finally {
          _refreshing = false;
        }
      }
    }
    handler.next(err);
  }
}

class _ErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final status = err.response?.statusCode ?? 0;
    final body = err.response?.data;
    final message = (body is Map && body['error'] is String)
        ? body['error'] as String
        : err.message ?? 'network error';
    handler.next(DioException(
      requestOptions: err.requestOptions,
      response: err.response,
      type: err.type,
      error: ApiException(status, message),
    ));
  }
}
