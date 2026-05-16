import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/repositories/chat_repository_impl.dart';
import '../../domain/entities/chat_message.dart';

final chatRepositoryProvider = Provider<ChatRepository>((ref) {
  return ChatRepositoryImpl(ref.watch(dioProvider));
});

class ChatState {
  const ChatState({this.messages = const [], this.thinking = false, this.error});
  final List<ChatMessage> messages;
  final bool thinking;
  final String? error;

  ChatState copyWith({List<ChatMessage>? messages, bool? thinking, String? error, bool clearError = false}) {
    return ChatState(
      messages: messages ?? this.messages,
      thinking: thinking ?? this.thinking,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class ChatViewModel extends StateNotifier<ChatState> {
  ChatViewModel(this._repo) : super(const ChatState());

  final ChatRepository _repo;

  Future<void> send(String text) async {
    final next = [...state.messages, ChatMessage(role: 'user', text: text)];
    state = state.copyWith(messages: next, thinking: true, clearError: true);
    try {
      final reply = await _repo.send(text);
      state = state.copyWith(messages: [...next, reply], thinking: false);
    } catch (e) {
      state = state.copyWith(thinking: false, error: e.toString());
    }
  }
}

final chatProvider = StateNotifierProvider<ChatViewModel, ChatState>((ref) {
  return ChatViewModel(ref.watch(chatRepositoryProvider));
});
