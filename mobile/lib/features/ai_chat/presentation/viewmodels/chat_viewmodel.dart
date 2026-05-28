import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/repositories/chat_repository_impl.dart';
import '../../domain/entities/chat_message.dart';

final chatRepositoryProvider = Provider<ChatRepository>((ref) {
  return ChatRepositoryImpl(ref.watch(dioProvider));
});

class ChatState {
  const ChatState({
    this.chatId,
    this.messages = const [],
    this.thinking = false,
    this.error,
  });
  final String? chatId;
  final List<ChatMessage> messages;
  final bool thinking;
  final String? error;

  ChatState copyWith({
    String? chatId,
    List<ChatMessage>? messages,
    bool? thinking,
    String? error,
    bool clearError = false,
    bool clearChat = false,
  }) {
    return ChatState(
      chatId: clearChat ? null : (chatId ?? this.chatId),
      messages: messages ?? this.messages,
      thinking: thinking ?? this.thinking,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class ChatViewModel extends StateNotifier<ChatState> {
  ChatViewModel(this._repo) : super(const ChatState());

  final ChatRepository _repo;

  void newChat() {
    state = const ChatState();
  }

  Future<void> loadChat(String id) async {
    state = const ChatState(thinking: true);
    try {
      final data = await _repo.getChat(id);
      final raw = (data['messages'] as List?) ?? const [];
      final msgs = raw.map((m) {
        final mm = Map<String, dynamic>.from(m as Map);
        final blocks = ((mm['blocks'] as List?) ?? const [])
            .map((e) => Map<String, dynamic>.from(e as Map))
            .toList();
        return ChatMessage(role: mm['role']?.toString() ?? 'ai', text: mm['text']?.toString() ?? '', blocks: blocks);
      }).toList();
      state = ChatState(chatId: id, messages: msgs);
    } catch (e) {
      state = ChatState(error: e.toString());
    }
  }

  Future<void> send(String text) async {
    final next = [...state.messages, ChatMessage(role: 'user', text: text)];
    state = state.copyWith(messages: next, thinking: true, clearError: true);
    try {
      final r = await _repo.send(text, chatId: state.chatId);
      state = state.copyWith(
        chatId: r.chatId.isNotEmpty ? r.chatId : state.chatId,
        messages: [...next, r.reply],
        thinking: false,
      );
    } catch (e) {
      state = state.copyWith(thinking: false, error: e.toString());
    }
  }
}

final chatProvider = StateNotifierProvider<ChatViewModel, ChatState>((ref) {
  return ChatViewModel(ref.watch(chatRepositoryProvider));
});
