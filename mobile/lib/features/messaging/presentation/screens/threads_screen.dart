import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../viewmodels/thread_viewmodel.dart';

class ThreadsScreen extends ConsumerWidget {
  const ThreadsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(threadsProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('Messages')),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(threadsProvider),
        child: async.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (e, _) => Center(child: Text(e.toString())),
          data: (items) => items.isEmpty
              ? ListView(children: const [
                  SizedBox(height: 80),
                  Icon(Icons.chat_bubble_outline, size: 48, color: Colors.grey),
                  SizedBox(height: 12),
                  Center(child: Text('Threads open when the vendor accepts a booking.')),
                ])
              : ListView.separated(
                  itemCount: items.length,
                  separatorBuilder: (_, __) => const Divider(height: 1),
                  itemBuilder: (_, i) {
                    final t = items[i];
                    return ListTile(
                      title: Text('Booking ${t.bookingId.substring(0, 8)}'),
                      subtitle: Text('Updated ${t.updatedAt}'),
                      onTap: () => context.push('/threads/${t.id}'),
                    );
                  },
                ),
        ),
      ),
    );
  }
}
