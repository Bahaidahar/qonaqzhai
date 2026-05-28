import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../domain/entities/chat_message.dart';

abstract class ChatRepository {
  Future<({String chatId, ChatMessage reply})> send(String message, {String? chatId});
  Future<List<Map<String, dynamic>>> listChats();
  Future<Map<String, dynamic>> getChat(String id);
  Future<void> deleteChat(String id);
  Future<void> renameChat(String id, String title);
}

class ChatRepositoryImpl implements ChatRepository {
  ChatRepositoryImpl(this._dio);

  final Dio _dio;

  @override
  Future<({String chatId, ChatMessage reply})> send(String message, {String? chatId}) async {
    final res = await _dio.post(ApiEndpoints.chat, data: {
      'message': message,
      if (chatId != null) 'chatId': chatId,
    });
    final data = res.data as Map<String, dynamic>;
    final blocks = ((data['blocks'] as List?) ?? const [])
        .map((e) => Map<String, dynamic>.from(e as Map))
        .toList();
    return (
      chatId: (data['chatId'] as String?) ?? chatId ?? '',
      reply: ChatMessage(role: 'ai', text: (data['reply'] as String?) ?? '', blocks: blocks),
    );
  }

  @override
  Future<List<Map<String, dynamic>>> listChats() async {
    final res = await _dio.get(ApiEndpoints.chats);
    final items = (res.data['items'] as List?) ?? const [];
    return items.map((e) => Map<String, dynamic>.from(e as Map)).toList();
  }

  @override
  Future<Map<String, dynamic>> getChat(String id) async {
    final res = await _dio.get(ApiEndpoints.chatById(id));
    return Map<String, dynamic>.from(res.data as Map);
  }

  @override
  Future<void> deleteChat(String id) => _dio.delete(ApiEndpoints.chatById(id));

  @override
  Future<void> renameChat(String id, String title) =>
      _dio.patch(ApiEndpoints.chatById(id), data: {'title': title});
}
