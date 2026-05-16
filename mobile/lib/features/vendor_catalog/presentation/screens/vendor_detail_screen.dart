import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../reviews/presentation/viewmodels/review_viewmodel.dart';
import '../../data/repositories/vendor_repository_impl.dart';
import '../../domain/entities/vendor.dart';
import '../viewmodels/vendor_catalog_viewmodel.dart';

final vendorByIdProvider = FutureProvider.family<Vendor, String>((ref, id) {
  return ref.watch(vendorRepositoryProvider).byId(id);
});

class VendorDetailScreen extends ConsumerWidget {
  const VendorDetailScreen({super.key, required this.id});

  final String id;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final vendorAsync = ref.watch(vendorByIdProvider(id));
    final reviewsAsync = ref.watch(vendorReviewsProvider(id));

    return Scaffold(
      appBar: AppBar(title: const Text('Vendor')),
      floatingActionButton: vendorAsync.maybeWhen(
        data: (v) => FloatingActionButton.extended(
          onPressed: () => context.push('/bookings/new?vendor=${v.id}&price=${v.priceFrom}'),
          icon: const Icon(Icons.event_available),
          label: const Text('Book now'),
        ),
        orElse: () => null,
      ),
      body: vendorAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text(e.toString())),
        data: (v) => ListView(
          padding: const EdgeInsets.all(16),
          children: [
            if (v.photoIds.isNotEmpty)
              SizedBox(
                height: 200,
                child: ListView.separated(
                  scrollDirection: Axis.horizontal,
                  itemCount: v.photoIds.length,
                  separatorBuilder: (_, __) => const SizedBox(width: 8),
                  itemBuilder: (_, i) => ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: CachedNetworkImage(
                      imageUrl: '${ApiEndpoints.baseUrl}${ApiEndpoints.photo(v.photoIds[i])}',
                      width: 280,
                      fit: BoxFit.cover,
                    ),
                  ),
                ),
              ),
            const SizedBox(height: 16),
            Text(v.name, style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: [
                Chip(label: Text(v.category)),
                Chip(label: Text(v.city)),
                if (v.ratingCount > 0)
                  Chip(
                    label: Text('⭐ ${v.ratingAvg.toStringAsFixed(1)} (${v.ratingCount})'),
                  ),
              ],
            ),
            if (v.priceFrom > 0)
              Padding(
                padding: const EdgeInsets.only(top: 8),
                child: Text(
                  'From ${v.priceFrom} ₸',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
              ),
            const SizedBox(height: 12),
            Text(v.description),
            const Divider(height: 32),
            Text('Reviews', style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            reviewsAsync.when(
              loading: () => const Padding(
                padding: EdgeInsets.all(16),
                child: Center(child: CircularProgressIndicator()),
              ),
              error: (e, _) => Text(e.toString()),
              data: (list) => list.isEmpty
                  ? const Text('No reviews yet')
                  : Column(
                      children: list
                          .map((r) => Card(
                                child: ListTile(
                                  title: Row(
                                    children: List.generate(
                                      5,
                                      (i) => Icon(
                                        Icons.star,
                                        size: 16,
                                        color: i < r.rating ? Colors.amber : Colors.grey,
                                      ),
                                    ),
                                  ),
                                  subtitle: Text(r.text),
                                ),
                              ))
                          .toList(),
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
