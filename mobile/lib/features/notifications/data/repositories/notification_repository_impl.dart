import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../domain/entities/notification.dart';

abstract class NotificationRepository {
  Future<List<AppNotification>> inbox();
}

class NotificationRepositoryImpl implements NotificationRepository {
  NotificationRepositoryImpl(this._dio);

  final Dio _dio;

  @override
  Future<List<AppNotification>> inbox() async {
    final res = await _dio.get(ApiEndpoints.notifications);
    return ((res.data['items'] as List?) ?? const [])
        .map((e) => e as Map<String, dynamic>)
        .map((j) => AppNotification(
              id: j['id'] as String,
              type: j['type'] as String,
              channel: j['channel'] as String,
              title: j['title'] as String,
              body: j['body'] as String,
              status: j['status'] as String,
              createdAt: j['createdAt'] as String,
            ))
        .toList();
  }
}
