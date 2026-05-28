import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../viewmodels/notifications_viewmodel.dart';

class NotificationsScreen extends ConsumerWidget {
  const NotificationsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final inboxAsync = ref.watch(inboxProvider);
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text('Notifications',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      body: RefreshIndicator(
        color: p.primary,
        onRefresh: () async => ref.invalidate(inboxProvider),
        child: inboxAsync.when(
          loading: () => Center(child: CupertinoActivityIndicator(color: p.mutedFg)),
          error: (e, _) => ListView(
            padding: const EdgeInsets.all(20),
            children: [
              Text(e.toString(), style: GoogleFonts.manrope(color: p.destructive)),
            ],
          ),
          data: (items) {
            if (items.isEmpty) {
              return ListView(
                padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
                children: const [
                  AppPageHeader(
                    title: 'Notifications',
                    subtitle: 'Booking updates, payments, reviews — in one feed.',
                  ),
                  SizedBox(height: 24),
                  AppEmptyState(
                    message: 'No notifications yet.',
                    icon: CupertinoIcons.bell,
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
                      title: 'Notifications',
                      subtitle: 'Booking updates, payments, reviews — in one feed.',
                    ),
                  );
                }
                final n = items[idx - 1];
                final body = n.body.replaceAll(RegExp(r'<[^>]+>'), '').trim();
                return AppCard(
                  padding: const EdgeInsets.fromLTRB(14, 12, 14, 12),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: Text(
                              n.title,
                              style: GoogleFonts.manrope(
                                fontSize: 14,
                                fontWeight: FontWeight.w700,
                                color: p.fg,
                              ),
                            ),
                          ),
                          AppBadge(label: n.type, tone: AppBadgeTone.neutral),
                        ],
                      ),
                      if (body.isNotEmpty) ...[
                        const SizedBox(height: 4),
                        Text(
                          body,
                          style: GoogleFonts.manrope(
                            fontSize: 12.5,
                            color: p.mutedFg,
                            height: 1.45,
                          ),
                        ),
                      ],
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
