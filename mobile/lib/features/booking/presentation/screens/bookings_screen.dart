import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../../../auth/domain/entities/user.dart';
import '../../../auth/presentation/viewmodels/auth_viewmodel.dart';
import '../viewmodels/booking_viewmodel.dart';

class BookingsScreen extends ConsumerWidget {
  const BookingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(bookingsProvider);
    final p = AppPalette.of(context);
    final role = ref.watch(authViewModelProvider).user?.role ?? UserRole.customer;
    final isVendor = role == UserRole.vendor;

    return Scaffold(
      appBar: AppBar(
        title: Text(isVendor ? 'Booking inbox' : 'Bookings',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
        actions: [
          if (!isVendor)
            IconButton(
              icon: const Icon(CupertinoIcons.bell),
              onPressed: () => context.push('/notifications'),
            ),
        ],
      ),
      body: RefreshIndicator(
        color: p.primary,
        onRefresh: () => ref.read(bookingsProvider.notifier).refresh(),
        child: () {
          if (state.loading && state.items.isEmpty) {
            return Center(child: CupertinoActivityIndicator(color: p.mutedFg));
          }
          if (state.error != null) {
            return Center(
              child: Text(state.error!,
                  style: GoogleFonts.manrope(color: p.destructive)),
            );
          }
          if (state.items.isEmpty) {
            return ListView(
              padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
              children: [
                AppPageHeader(
                  title: isVendor ? 'Booking inbox' : 'Bookings',
                  subtitle: isVendor
                      ? 'Incoming requests from customers.'
                      : 'Your requests and confirmed events.',
                ),
                const SizedBox(height: 24),
                const AppEmptyState(
                  message: 'No bookings yet.',
                  icon: CupertinoIcons.calendar,
                ),
              ],
            );
          }
          return ListView.separated(
            padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
            itemCount: state.items.length + 1,
            separatorBuilder: (_, __) => const SizedBox(height: 10),
            itemBuilder: (_, idx) {
              if (idx == 0) {
                return Padding(
                  padding: const EdgeInsets.only(bottom: 8),
                  child: AppPageHeader(
                    title: isVendor ? 'Booking inbox' : 'Bookings',
                    subtitle: isVendor
                        ? 'Incoming requests from customers.'
                        : 'Your requests and confirmed events.',
                  ),
                );
              }
              final b = state.items[idx - 1];
              return _BookingCard(
                eventDate: b.eventDate,
                guests: b.guestCount,
                status: b.status,
                amount: b.amount,
                role: role,
                onPay: !isVendor &&
                        (b.status == 'pending' ||
                            b.status == 'accepted' ||
                            b.status == 'completed') &&
                        b.amount > 0
                    ? () async {
                        final paid =
                            await context.push<bool>('/pay?booking=${b.id}&amount=${b.amount}');
                        if (paid == true) {
                          await ref.read(bookingsProvider.notifier).refresh();
                        }
                      }
                    : null,
                onReview: !isVendor && b.status == 'completed'
                    ? () => context.push('/reviews/new?booking=${b.id}&vendor=${b.vendorId}')
                    : null,
                onCancel: !isVendor && b.status == 'pending'
                    ? () => ref.read(bookingsProvider.notifier).cancel(b.id)
                    : null,
                onAccept: isVendor && b.status == 'pending'
                    ? () =>
                        ref.read(bookingsProvider.notifier).setStatus(b.id, 'accepted')
                    : null,
                onDecline: isVendor && b.status == 'pending'
                    ? () =>
                        ref.read(bookingsProvider.notifier).setStatus(b.id, 'declined')
                    : null,
              );
            },
          );
        }(),
      ),
    );
  }
}

class _BookingCard extends StatelessWidget {
  const _BookingCard({
    required this.eventDate,
    required this.guests,
    required this.status,
    required this.amount,
    required this.role,
    this.onPay,
    this.onReview,
    this.onCancel,
    this.onAccept,
    this.onDecline,
  });

  final String eventDate;
  final int guests;
  final String status;
  final int amount;
  final UserRole role;
  final VoidCallback? onPay;
  final VoidCallback? onReview;
  final VoidCallback? onCancel;
  final VoidCallback? onAccept;
  final VoidCallback? onDecline;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    final tone = switch (status) {
      'accepted' || 'completed' => AppBadgeTone.success,
      'declined' || 'canceled' || 'cancelled' => AppBadgeTone.danger,
      'paid' => AppBadgeTone.info,
      _ => AppBadgeTone.warning,
    };
    return AppCard(
      padding: const EdgeInsets.fromLTRB(16, 14, 16, 14),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(CupertinoIcons.calendar,
                            size: 14, color: p.mutedFg),
                        const SizedBox(width: 6),
                        Text(
                          eventDate,
                          style: GoogleFonts.manrope(
                            fontSize: 14.5,
                            fontWeight: FontWeight.w700,
                            color: p.fg,
                            letterSpacing: -0.2,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 6),
                    Row(
                      children: [
                        Icon(CupertinoIcons.person_2,
                            size: 13, color: p.mutedFg),
                        const SizedBox(width: 6),
                        Text(
                          '$guests guests',
                          style: GoogleFonts.manrope(
                              fontSize: 12.5, color: p.mutedFg),
                        ),
                        if (amount > 0) ...[
                          const SizedBox(width: 12),
                          Icon(CupertinoIcons.money_dollar,
                              size: 13, color: p.mutedFg),
                          const SizedBox(width: 4),
                          Text(
                            '${_kzt(amount)} ₸',
                            style: GoogleFonts.manrope(
                                fontSize: 12.5, color: p.mutedFg),
                          ),
                        ],
                      ],
                    ),
                  ],
                ),
              ),
              AppBadge(label: status, tone: tone),
            ],
          ),
          if (onAccept != null ||
              onDecline != null ||
              onPay != null ||
              onReview != null ||
              onCancel != null) ...[
            const SizedBox(height: 12),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                if (onDecline != null)
                  OutlinedButton.icon(
                    onPressed: onDecline,
                    icon: Icon(CupertinoIcons.xmark, size: 14, color: p.destructive),
                    label: Text(
                      'Decline',
                      style: GoogleFonts.manrope(
                          color: p.destructive,
                          fontSize: 12.5,
                          fontWeight: FontWeight.w600),
                    ),
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 14),
                      minimumSize: const Size(0, 36),
                      side:
                          BorderSide(color: p.destructive.withValues(alpha: 0.3)),
                    ),
                  ),
                if (onAccept != null)
                  FilledButton.icon(
                    onPressed: onAccept,
                    icon: const Icon(CupertinoIcons.checkmark_alt, size: 16),
                    label: const Text('Accept'),
                    style: FilledButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 14),
                      minimumSize: const Size(0, 36),
                      textStyle:
                          GoogleFonts.manrope(fontSize: 12.5, fontWeight: FontWeight.w600),
                    ),
                  ),
                if (onPay != null)
                  FilledButton.icon(
                    onPressed: onPay,
                    icon: const Icon(CupertinoIcons.creditcard, size: 16),
                    label: const Text('Pay'),
                    style: FilledButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 14),
                      minimumSize: const Size(0, 36),
                      textStyle:
                          GoogleFonts.manrope(fontSize: 12.5, fontWeight: FontWeight.w600),
                    ),
                  ),
                if (onReview != null)
                  OutlinedButton.icon(
                    onPressed: onReview,
                    icon: const Icon(CupertinoIcons.star, size: 16),
                    label: const Text('Review'),
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 14),
                      minimumSize: const Size(0, 36),
                      textStyle:
                          GoogleFonts.manrope(fontSize: 12.5, fontWeight: FontWeight.w600),
                    ),
                  ),
                if (onCancel != null)
                  OutlinedButton.icon(
                    onPressed: onCancel,
                    icon: Icon(CupertinoIcons.xmark, size: 14, color: p.destructive),
                    label: Text(
                      'Cancel',
                      style: GoogleFonts.manrope(
                          color: p.destructive,
                          fontSize: 12.5,
                          fontWeight: FontWeight.w600),
                    ),
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 14),
                      minimumSize: const Size(0, 36),
                      side: BorderSide(color: p.destructive.withValues(alpha: 0.3)),
                    ),
                  ),
              ],
            ),
          ],
        ],
      ),
    );
  }

  static String _kzt(int v) {
    if (v == 0) return '0';
    final s = v.toString();
    final buf = StringBuffer();
    for (var i = 0; i < s.length; i++) {
      if (i > 0 && (s.length - i) % 3 == 0) buf.write(' ');
      buf.write(s[i]);
    }
    return buf.toString();
  }
}
