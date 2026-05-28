import 'dart:async';
import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../../core/network/token_store.dart';
import '../../../../core/theme/app_theme.dart';
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
    final wsBase = ApiEndpoints.baseUrl.replaceFirst(RegExp(r'^http'), 'ws');
    final url = '$wsBase/api/ws?token=${Uri.encodeQueryComponent(token)}';
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
    final p = AppPalette.of(context);
    if (_loading) {
      return Scaffold(body: Center(child: CupertinoActivityIndicator(color: p.mutedFg)));
    }
    return Scaffold(
      appBar: AppBar(
        title: Text('Conversation',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
        actions: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Center(
              child: Row(
                children: [
                  Container(
                    width: 8,
                    height: 8,
                    decoration: BoxDecoration(
                      color: _connected ? const Color(0xFF34C759) : p.mutedFg,
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 6),
                  Text(
                    _connected ? 'live' : 'offline',
                    style: GoogleFonts.manrope(
                      fontSize: 11,
                      color: p.mutedFg,
                      fontWeight: FontWeight.w600,
                      letterSpacing: 0.4,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: ListView.builder(
                controller: _scroll,
                padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                itemCount: _messages.length,
                itemBuilder: (_, i) =>
                    _bubble(_messages[i], mine: _messages[i].senderId == me),
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(12, 4, 12, 12),
              child: Container(
                padding: const EdgeInsets.fromLTRB(6, 6, 6, 6),
                decoration: BoxDecoration(
                  color: p.card,
                  borderRadius: BorderRadius.circular(18),
                  border: Border.all(color: p.border),
                ),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _draft,
                        style: GoogleFonts.manrope(fontSize: 14, color: p.fg),
                        decoration: InputDecoration(
                          filled: false,
                          hintText: 'Message…',
                          hintStyle:
                              GoogleFonts.manrope(fontSize: 14, color: p.mutedFg),
                          border: InputBorder.none,
                          enabledBorder: InputBorder.none,
                          focusedBorder: InputBorder.none,
                          contentPadding: const EdgeInsets.symmetric(
                              horizontal: 12, vertical: 12),
                        ),
                        onSubmitted: (_) => _send(),
                      ),
                    ),
                    SizedBox(
                      width: 38,
                      height: 38,
                      child: Material(
                        color: p.primary,
                        shape: const CircleBorder(),
                        child: InkWell(
                          customBorder: const CircleBorder(),
                          onTap: _send,
                          child: Icon(
                            CupertinoIcons.arrow_up,
                            size: 18,
                            color: p.onPrimary,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _bubble(ThreadMessage m, {required bool mine}) {
    final p = AppPalette.of(context);
    return Align(
      alignment: mine ? Alignment.centerRight : Alignment.centerLeft,
      child: ConstrainedBox(
        constraints: BoxConstraints(
            maxWidth: MediaQuery.of(context).size.width * 0.82),
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 4),
          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
          decoration: BoxDecoration(
            color: mine ? p.primary : p.card,
            borderRadius: BorderRadius.circular(16),
            border: mine ? null : Border.all(color: p.border),
          ),
          child: Text(
            m.text,
            style: GoogleFonts.manrope(
              fontSize: 14,
              color: mine ? p.onPrimary : p.fg,
              height: 1.45,
            ),
          ),
        ),
      ),
    );
  }
}
