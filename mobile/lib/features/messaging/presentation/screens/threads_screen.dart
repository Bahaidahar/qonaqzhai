import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../viewmodels/thread_viewmodel.dart';

class ThreadsScreen extends ConsumerWidget {
  const ThreadsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(threadsProvider);
    final p = AppPalette.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Text('Messages',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      body: RefreshIndicator(
        color: p.primary,
        onRefresh: () async => ref.invalidate(threadsProvider),
        child: async.when(
          loading: () => Center(child: CupertinoActivityIndicator(color: p.mutedFg)),
          error: (e, _) => Center(
            child: Text(e.toString(),
                style: GoogleFonts.manrope(color: p.destructive)),
          ),
          data: (items) {
            if (items.isEmpty) {
              return ListView(
                padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
                children: const [
                  AppPageHeader(
                    title: 'Messages',
                    subtitle: 'Chats open when a vendor accepts a booking.',
                  ),
                  SizedBox(height: 24),
                  AppEmptyState(
                    message:
                        'No conversations yet. Once a vendor accepts your booking, a chat opens here.',
                    icon: CupertinoIcons.chat_bubble_2,
                  ),
                ],
              );
            }
            return ListView.separated(
              padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
              itemCount: items.length + 1,
              separatorBuilder: (_, __) => const SizedBox(height: 8),
              itemBuilder: (_, idx) {
                if (idx == 0) {
                  return const Padding(
                    padding: EdgeInsets.only(bottom: 8),
                    child: AppPageHeader(
                      title: 'Messages',
                      subtitle: 'Direct chats with your vendors.',
                    ),
                  );
                }
                final t = items[idx - 1];
                return AppCard(
                  onTap: () => context.push('/threads/${t.id}'),
                  padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
                  child: Row(
                    children: [
                      Container(
                        width: 40,
                        height: 40,
                        decoration: BoxDecoration(
                          color: p.primary.withValues(alpha: 0.1),
                          shape: BoxShape.circle,
                        ),
                        alignment: Alignment.center,
                        child: Icon(CupertinoIcons.chat_bubble_2,
                            size: 18, color: p.primary),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Booking ${t.bookingId.substring(0, 8)}',
                              style: GoogleFonts.manrope(
                                fontSize: 14,
                                fontWeight: FontWeight.w700,
                                color: p.fg,
                              ),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              'Updated ${t.updatedAt}',
                              style: GoogleFonts.manrope(
                                  fontSize: 11.5, color: p.mutedFg),
                            ),
                          ],
                        ),
                      ),
                      Icon(CupertinoIcons.chevron_forward,
                          size: 16, color: p.mutedFg),
                    ],
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }
}
