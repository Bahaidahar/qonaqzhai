import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/repositories/review_repository_impl.dart';
import '../../domain/entities/review.dart';

final reviewRepositoryProvider = Provider<ReviewRepository>((ref) {
  return ReviewRepositoryImpl(ref.watch(dioProvider));
});

final vendorReviewsProvider =
    FutureProvider.family<List<Review>, String>((ref, vendorId) {
  return ref.watch(reviewRepositoryProvider).listByVendor(vendorId);
});
