import 'package:dio/dio.dart';

import '../../domain/entities/thread.dart';

abstract class ThreadRepository {
  Future<List<BookingThread>> list();
  Future<({BookingThread thread, List<ThreadMessage> messages})> get(String id);
  Future<ThreadMessage> send(String threadId, String text);
}

class ThreadRepositoryImpl implements ThreadRepository {
  ThreadRepositoryImpl(this._dio);

  final Dio _dio;

  BookingThread _thread(Map<String, dynamic> j) => BookingThread(
        id: j['id'] as String,
        bookingId: j['bookingId'] as String,
        customerId: j['customerId'] as String,
        vendorId: j['vendorId'] as String,
        createdAt: j['createdAt'] as String,
        updatedAt: j['updatedAt'] as String,
      );

  ThreadMessage _msg(Map<String, dynamic> j) => ThreadMessage(
        id: j['id'] as String,
        threadId: j['threadId'] as String,
        senderId: j['senderId'] as String,
        text: j['text'] as String,
        createdAt: j['createdAt'] as String,
      );

  @override
  Future<List<BookingThread>> list() async {
    final res = await _dio.get('/api/threads');
    return ((res.data['items'] as List?) ?? const [])
        .map((e) => _thread(e as Map<String, dynamic>))
        .toList();
  }

  @override
  Future<({BookingThread thread, List<ThreadMessage> messages})> get(String id) async {
    final res = await _dio.get('/api/threads/$id');
    return (
      thread: _thread(res.data['thread'] as Map<String, dynamic>),
      messages: ((res.data['messages'] as List?) ?? const [])
          .map((e) => _msg(e as Map<String, dynamic>))
          .toList(),
    );
  }

  @override
  Future<ThreadMessage> send(String threadId, String text) async {
    final res = await _dio.post('/api/threads/$threadId/messages', data: {'text': text});
    return _msg(res.data as Map<String, dynamic>);
  }
}
