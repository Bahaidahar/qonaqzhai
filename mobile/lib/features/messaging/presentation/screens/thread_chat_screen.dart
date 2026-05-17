import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

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
  Timer? _poll;

  @override
  void initState() {
    super.initState();
    _poll = Timer.periodic(const Duration(seconds: 4), (_) {
      if (mounted) ref.invalidate(threadDetailProvider(widget.id));
    });
  }

  @override
  void dispose() {
    _poll?.cancel();
    _draft.dispose();
    _scroll.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    final text = _draft.text.trim();
    if (text.isEmpty) return;
    _draft.clear();
    try {
      await ref.read(threadRepositoryProvider).send(widget.id, text);
      ref.invalidate(threadDetailProvider(widget.id));
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Send failed: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final me = ref.watch(authViewModelProvider).user?.id;
    final async = ref.watch(threadDetailProvider(widget.id));

    return Scaffold(
      appBar: AppBar(title: const Text('Conversation')),
      body: async.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text(e.toString())),
        data: (d) => Column(
          children: [
            Expanded(
              child: ListView.builder(
                controller: _scroll,
                padding: const EdgeInsets.all(12),
                itemCount: d.messages.length,
                itemBuilder: (_, i) {
                  final m = d.messages[i];
                  return _bubble(m, mine: m.senderId == me);
                },
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
