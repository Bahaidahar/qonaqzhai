import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/i18n/i18n.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';

class AiBlock extends ConsumerWidget {
  const AiBlock({super.key, required this.block});
  final Map<String, dynamic> block;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final type = block['type']?.toString() ?? '';
    final data = Map<String, dynamic>.from((block['data'] as Map?) ?? const {});
    switch (type) {
      case 'plan':
        return _PlanCard(data: data);
      case 'budget':
        return _BudgetCard(data: data);
      case 'vendors':
        return _VendorsCard(data: data);
      default:
        return const SizedBox.shrink();
    }
  }
}

class _PlanCard extends ConsumerWidget {
  const _PlanCard({required this.data});
  final Map<String, dynamic> data;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final p = AppPalette.of(context);
    return Padding(
      padding: const EdgeInsets.only(top: 8),
      child: AppCard(
        padding: const EdgeInsets.fromLTRB(14, 12, 14, 14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            AppSectionHeader(tr(ref, 'chat_block_plan')),
            const SizedBox(height: 8),
            Text(
              data['title']?.toString() ?? '',
              style: GoogleFonts.manrope(
                fontSize: 15,
                fontWeight: FontWeight.w700,
                color: p.fg,
              ),
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 6,
              runSpacing: 6,
              children: [
                if ((data['eventType'] ?? '').toString().isNotEmpty)
                  _Chip(icon: CupertinoIcons.gift, text: data['eventType'].toString()),
                if ((data['date'] ?? '').toString().isNotEmpty)
                  _Chip(icon: CupertinoIcons.calendar, text: data['date'].toString()),
                if ((data['city'] ?? '').toString().isNotEmpty)
                  _Chip(icon: CupertinoIcons.location, text: data['city'].toString()),
                if ((data['guests'] ?? 0) is num && (data['guests'] as num) > 0)
                  _Chip(icon: CupertinoIcons.person_2, text: '${data['guests']}'),
                if ((data['budget'] ?? 0) is num && (data['budget'] as num) > 0)
                  _Chip(
                      icon: CupertinoIcons.money_dollar,
                      text: '${data['budget']} ₸'),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  const _Chip({required this.icon, required this.text});
  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 5),
      decoration: BoxDecoration(
        color: p.muted,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: p.border),
      ),
      child: Row(mainAxisSize: MainAxisSize.min, children: [
        Icon(icon, size: 12, color: p.mutedFg),
        const SizedBox(width: 5),
        Text(text,
            style: GoogleFonts.manrope(
                fontSize: 11.5, fontWeight: FontWeight.w600, color: p.fg)),
      ]),
    );
  }
}

class _BudgetCard extends ConsumerWidget {
  const _BudgetCard({required this.data});
  final Map<String, dynamic> data;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final p = AppPalette.of(context);
    final total = (data['total'] as num? ?? 0).toInt();
    final cats = (data['categories'] as List?) ?? const [];
    return Padding(
      padding: const EdgeInsets.only(top: 8),
      child: AppCard(
        padding: const EdgeInsets.fromLTRB(14, 12, 14, 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            AppSectionHeader(tr(ref, 'chat_block_budget')),
            const SizedBox(height: 6),
            Text(
              '${_kzt(total)} ₸',
              style: GoogleFonts.manrope(
                fontSize: 18,
                fontWeight: FontWeight.w700,
                color: p.fg,
                letterSpacing: -0.4,
              ),
            ),
            const SizedBox(height: 10),
            for (final c in cats)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 3),
                child: Row(children: [
                  SizedBox(
                    width: 100,
                    child: Text(
                      c['name']?.toString() ?? '',
                      overflow: TextOverflow.ellipsis,
                      style: GoogleFonts.manrope(fontSize: 12, color: p.fg),
                    ),
                  ),
                  Expanded(
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(4),
                      child: LinearProgressIndicator(
                        value: ((c['pct'] as num?) ?? 0) / 100,
                        minHeight: 6,
                        backgroundColor: p.muted,
                        valueColor: AlwaysStoppedAnimation(p.primary),
                      ),
                    ),
                  ),
                  SizedBox(
                    width: 78,
                    child: Text(
                      '${_kzt((c['amount'] as num?)?.toInt() ?? 0)} ₸',
                      textAlign: TextAlign.end,
                      style: GoogleFonts.manrope(
                          fontSize: 12, fontWeight: FontWeight.w600, color: p.fg),
                    ),
                  ),
                ]),
              ),
          ],
        ),
      ),
    );
  }

  static String _kzt(int v) {
    final s = v.toString();
    final buf = StringBuffer();
    for (var i = 0; i < s.length; i++) {
      if (i > 0 && (s.length - i) % 3 == 0) buf.write(' ');
      buf.write(s[i]);
    }
    return buf.toString();
  }
}

class _VendorsCard extends ConsumerWidget {
  const _VendorsCard({required this.data});
  final Map<String, dynamic> data;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final p = AppPalette.of(context);
    final items = (data['items'] as List?) ?? const [];
    return Padding(
      padding: const EdgeInsets.only(top: 8),
      child: AppCard(
        padding: const EdgeInsets.fromLTRB(14, 12, 14, 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            AppSectionHeader(tr(ref, 'chat_block_vendors')),
            const SizedBox(height: 8),
            for (final raw in items)
              Builder(builder: (_) {
                final v = Map<String, dynamic>.from(raw as Map);
                return InkWell(
                  onTap: () {
                    final id = v['id']?.toString();
                    if (id != null) context.push('/vendors/$id');
                  },
                  borderRadius: BorderRadius.circular(8),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    child: Row(children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              v['name']?.toString() ?? '',
                              style: GoogleFonts.manrope(
                                  fontSize: 13.5,
                                  fontWeight: FontWeight.w700,
                                  color: p.fg),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              '${v['category'] ?? ''} · ${v['priceFrom'] ?? 0} ₸',
                              style: GoogleFonts.manrope(
                                  fontSize: 11.5, color: p.mutedFg),
                            ),
                          ],
                        ),
                      ),
                      Icon(CupertinoIcons.chevron_forward,
                          size: 14, color: p.mutedFg),
                    ]),
                  ),
                );
              }),
          ],
        ),
      ),
    );
  }
}
