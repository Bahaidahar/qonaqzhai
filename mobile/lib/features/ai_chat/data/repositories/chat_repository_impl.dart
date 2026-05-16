import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../domain/entities/chat_message.dart';

abstract class ChatRepository {
  Future<ChatMessage> send(String message);
}

class ChatRepositoryImpl implements ChatRepository {
  ChatRepositoryImpl(this._dio);

  final Dio _dio;

  @override
  Future<ChatMessage> send(String message) async {
    final res = await _dio.post(ApiEndpoints.chat, data: {'message': message});
    final data = res.data as Map<String, dynamic>;
    final blocks = ((data['blocks'] as List?) ?? const [])
        .map((e) => Map<String, dynamic>.from(e as Map))
        .toList();
    return ChatMessage(
      role: 'ai',
      text: (data['reply'] as String?) ?? '',
      blocks: blocks,
    );
  }
}
