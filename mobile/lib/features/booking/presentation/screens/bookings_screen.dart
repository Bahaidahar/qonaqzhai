import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';

import '../viewmodels/booking_viewmodel.dart';

class BookingsScreen extends ConsumerWidget {
  const BookingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(bookingsProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('My bookings')),
      body: RefreshIndicator(
        onRefresh: () => ref.read(bookingsProvider.notifier).refresh(),
        child: state.loading && state.items.isEmpty
            ? const Center(child: CircularProgressIndicator())
            : state.error != null
                ? Center(child: Text(state.error!))
                : state.items.isEmpty
                    ? ListView(
                        children: const [
                          SizedBox(height: 80),
                          Icon(Icons.event_busy, size: 48, color: Colors.grey),
                          SizedBox(height: 12),
                          Center(child: Text('No bookings yet')),
                        ],
                      )
                    : ListView.builder(
                        itemCount: state.items.length,
                        itemBuilder: (_, i) {
                          final b = state.items[i];
                          return Card(
                            margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                            child: ListTile(
                              title: Text(b.eventDate),
                              subtitle: Text('${b.guestCount} guests · ${b.status} · ${b.amount} ₸'),
                              trailing: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  if (b.status == 'pending' && b.amount > 0)
                                    IconButton(
                                      icon: const Icon(Icons.credit_card),
                                      tooltip: 'Pay',
                                      onPressed: () async {
                                        try {
                                          final url = await ref
                                              .read(bookingsProvider.notifier)
                                              .startPayment(b.id);
                                          if (!context.mounted) return;
                                          final uri = Uri.parse(url);
                                          if (await canLaunchUrl(uri)) {
                                            await launchUrl(uri, mode: LaunchMode.externalApplication);
                                          } else if (context.mounted) {
                                            ScaffoldMessenger.of(context).showSnackBar(
                                              SnackBar(content: Text('Cannot open: $url')),
                                            );
                                          }
                                        } catch (e) {
                                          if (!context.mounted) return;
                                          ScaffoldMessenger.of(context).showSnackBar(
                                            SnackBar(content: Text('Payment unavailable: $e')),
                                          );
                                        }
                                      },
                                    ),
                                  if (b.status == 'completed')
                                    IconButton(
                                      icon: const Icon(Icons.star_outline),
                                      tooltip: 'Review',
                                      onPressed: () => context.push(
                                        '/reviews/new?booking=${b.id}&vendor=${b.vendorId}',
                                      ),
                                    ),
                                  if (b.status == 'pending')
                                    IconButton(
                                      icon: const Icon(Icons.cancel),
                                      tooltip: 'Cancel',
                                      onPressed: () =>
                                          ref.read(bookingsProvider.notifier).cancel(b.id),
                                    ),
                                ],
                              ),
                            ),
                          );
                        },
                      ),
      ),
    );
  }
}
