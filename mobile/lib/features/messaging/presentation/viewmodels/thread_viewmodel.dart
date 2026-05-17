import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/repositories/thread_repository_impl.dart';
import '../../domain/entities/thread.dart';

final threadRepositoryProvider = Provider<ThreadRepository>((ref) {
  return ThreadRepositoryImpl(ref.watch(dioProvider));
});

final threadsProvider = FutureProvider<List<BookingThread>>((ref) {
  return ref.watch(threadRepositoryProvider).list();
});

final threadDetailProvider =
    FutureProvider.family<({BookingThread thread, List<ThreadMessage> messages}), String>(
  (ref, id) => ref.watch(threadRepositoryProvider).get(id),
);
