import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../../core/network/token_store.dart';
import '../../../auth/presentation/viewmodels/auth_viewmodel.dart';
import '../../domain/entities/thread.dart';
import '../viewmodels/thread_viewmodel.dart';

class ThreadChatScreen extends ConsumerStatefulWidget {
  const ThreadChatScreen({super.key, required this.id});

  final String id;

  @override
  ConsumerState<ThreadChatScreen> createState() => _ThreadChatScreenState();
}

class _ThreadChatScreenState extends ConsumerState<ThreadChatScreen> {
  final _draft = TextEditingController();
  final _scroll = ScrollController();
  WebSocketChannel? _channel;
  StreamSubscription<dynamic>? _sub;
  List<ThreadMessage> _messages = [];
  BookingThread? _thread;
  bool _loading = true;
  bool _connected = false;

  @override
  void initState() {
    super.initState();
    _bootstrap();
  }

  Future<void> _bootstrap() async {
    final repo = ref.read(threadRepositoryProvider);
    try {
      final d = await repo.get(widget.id);
      if (!mounted) return;
      setState(() {
        _thread = d.thread;
        _messages = d.messages;
        _loading = false;
      });
      await _connectSocket();
    } catch (e) {
      if (!mounted) return;
      setState(() => _loading = false);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Load failed: $e')),
      );
    }
  }

  Future<void> _connectSocket() async {
    final token = await TokenStore().readAccess();
    if (token == null) return;
    final url = ApiEndpoints.baseUrl
            .replaceFirst(RegExp(r'^http'), 'ws') +
        '/api/ws?token=' +
        Uri.encodeQueryComponent(token);
    try {
      _channel = WebSocketChannel.connect(Uri.parse(url));
      setState(() => _connected = true);
      _sub = _channel!.stream.listen(
        (raw) {
          try {
            final env = jsonDecode(raw as String) as Map<String, dynamic>;
            if (env['op'] != 'thread.message') return;
            final data = env['data'] as Map<String, dynamic>;
            if (data['threadId'] != widget.id) return;
            final m = ThreadMessage(
              id: data['id'] as String,
              threadId: data['threadId'] as String,
              senderId: data['senderId'] as String,
              text: data['text'] as String,
              createdAt: data['createdAt'] as String,
            );
            if (_messages.any((x) => x.id == m.id)) return;
            setState(() => _messages = [..._messages, m]);
          } catch (_) {/* ignore */}
        },
        onDone: () {
          setState(() => _connected = false);
          // Naive reconnect after 2s.
          Future.delayed(const Duration(seconds: 2), () {
            if (mounted) _connectSocket();
          });
        },
        onError: (_) => setState(() => _connected = false),
      );
    } catch (e) {
      setState(() => _connected = false);
    }
  }

  @override
  void dispose() {
    _sub?.cancel();
    _channel?.sink.close();
    _draft.dispose();
    _scroll.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    final text = _draft.text.trim();
    if (text.isEmpty || _thread == null) return;
    _draft.clear();
    final channel = _channel;
    if (channel != null && _connected) {
      channel.sink.add(jsonEncode({
        'op': 'message',
        'data': {'threadId': _thread!.id, 'text': text},
      }));
    } else {
      try {
        final m = await ref.read(threadRepositoryProvider).send(widget.id, text);
        setState(() => _messages = [..._messages, m]);
      } catch (e) {
        if (!mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Send failed: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final me = ref.watch(authViewModelProvider).user?.id;
    if (_loading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    return Scaffold(
      appBar: AppBar(
        title: const Text('Conversation'),
        actions: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Center(
              child: Text(
                _connected ? '● live' : '○ offline',
                style: TextStyle(
                  fontSize: 11,
                  color: _connected ? Colors.green : Colors.grey,
                ),
              ),
            ),
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              controller: _scroll,
              padding: const EdgeInsets.all(12),
              itemCount: _messages.length,
              itemBuilder: (_, i) => _bubble(_messages[i], mine: _messages[i].senderId == me),
            ),
          ),
          SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(8),
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _draft,
                      decoration: const InputDecoration(hintText: 'Message…'),
                      onSubmitted: (_) => _send(),
                    ),
                  ),
                  IconButton(icon: const Icon(Icons.send), onPressed: _send),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _bubble(ThreadMessage m, {required bool mine}) {
    return Align(
      alignment: mine ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.symmetric(vertical: 4),
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: mine ? Colors.indigo[100] : Colors.grey[200],
          borderRadius: BorderRadius.circular(12),
        ),
        child: Text(m.text),
      ),
    );
  }
}
