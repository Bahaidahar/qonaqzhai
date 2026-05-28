import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../core/i18n/i18n.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/ui/ui.dart';
import 'cards_viewmodel.dart';

class CardsScreen extends ConsumerWidget {
  const CardsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(cardsViewModelProvider);
    final vm = ref.read(cardsViewModelProvider.notifier);
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(tr(ref, 'cards_title'),
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
        actions: [
          IconButton(
            icon: const Icon(CupertinoIcons.add),
            tooltip: tr(ref, 'cards_add'),
            onPressed: () => _showAddSheet(context, ref),
          ),
        ],
      ),
      body: RefreshIndicator(
        color: p.primary,
        onRefresh: vm.refresh,
        child: state.loading
            ? Center(child: CupertinoActivityIndicator(color: p.mutedFg))
            : state.items.isEmpty
                ? ListView(
                    padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
                    children: [
                      const AppPageHeader(
                        title: 'Saved cards',
                        subtitle: 'Pay faster by saving your card. Secure — only the last four digits are stored.',
                      ),
                      const SizedBox(height: 24),
                      AppEmptyState(
                        message: tr(ref, 'cards_empty'),
                        icon: CupertinoIcons.creditcard,
                      ),
                    ],
                  )
                : ListView.separated(
                    padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
                    itemCount: state.items.length + 1,
                    separatorBuilder: (_, __) => const SizedBox(height: 10),
                    itemBuilder: (_, idx) {
                      if (idx == 0) {
                        return const Padding(
                          padding: EdgeInsets.only(bottom: 8),
                          child: AppPageHeader(
                            title: 'Saved cards',
                            subtitle: 'Manage payment methods.',
                          ),
                        );
                      }
                      final c = state.items[idx - 1];
                      return AppCard(
                        padding: const EdgeInsets.fromLTRB(14, 14, 6, 14),
                        child: Row(
                          children: [
                            Container(
                              width: 44,
                              height: 44,
                              decoration: BoxDecoration(
                                color: p.primary.withValues(alpha: 0.1),
                                borderRadius: BorderRadius.circular(10),
                                border: Border.all(
                                  color: p.primary.withValues(alpha: 0.25),
                                ),
                              ),
                              alignment: Alignment.center,
                              child: Icon(CupertinoIcons.creditcard,
                                  color: p.primary, size: 18),
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Row(
                                    children: [
                                      Text(
                                        '${c.brand.toUpperCase()} •••• ${c.last4}',
                                        style: GoogleFonts.manrope(
                                          fontWeight: FontWeight.w700,
                                          color: p.fg,
                                          fontSize: 14,
                                        ),
                                      ),
                                      if (c.isDefault) ...[
                                        const SizedBox(width: 8),
                                        const AppBadge(
                                          label: 'Default',
                                          tone: AppBadgeTone.info,
                                        ),
                                      ],
                                    ],
                                  ),
                                  const SizedBox(height: 2),
                                  Text(
                                    '${c.holder} · ${c.expMonth.toString().padLeft(2, '0')}/${c.expYear}',
                                    style: GoogleFonts.manrope(
                                      fontSize: 12,
                                      color: p.mutedFg,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            PopupMenuButton<String>(
                              icon: Icon(CupertinoIcons.ellipsis, color: p.mutedFg),
                              onSelected: (v) async {
                                if (v == 'default') await vm.makeDefault(c.id);
                                if (v == 'delete') await vm.remove(c.id);
                              },
                              itemBuilder: (_) => [
                                if (!c.isDefault)
                                  PopupMenuItem(
                                    value: 'default',
                                    child: Text(tr(ref, 'cards_make_default')),
                                  ),
                                PopupMenuItem(
                                  value: 'delete',
                                  child: Text(tr(ref, 'common_delete')),
                                ),
                              ],
                            ),
                          ],
                        ),
                      );
                    },
                  ),
      ),
    );
  }

  void _showAddSheet(BuildContext context, WidgetRef ref) {
    final number = TextEditingController();
    final month = TextEditingController();
    final year = TextEditingController();
    final holder = TextEditingController();
    bool def = false;
    final p = AppPalette.of(context);
    showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      backgroundColor: p.card,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => StatefulBuilder(builder: (ctx, setState) {
        return Padding(
          padding: EdgeInsets.fromLTRB(
            20,
            12,
            20,
            20 + MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Center(
                child: Container(
                  width: 36,
                  height: 4,
                  decoration: BoxDecoration(
                    color: p.border,
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              ),
              const SizedBox(height: 14),
              Text(
                tr(ref, 'cards_add'),
                style: GoogleFonts.manrope(
                  fontSize: 22,
                  fontWeight: FontWeight.w700,
                  color: p.fg,
                  letterSpacing: -0.5,
                ),
              ),
              const SizedBox(height: 18),
              _Lbl(label: tr(ref, 'cards_number'), child: TextField(
                controller: number,
                keyboardType: TextInputType.number,
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                  LengthLimitingTextInputFormatter(19)
                ],
                decoration: const InputDecoration(hintText: '4111 1111 1111 1111'),
              )),
              const SizedBox(height: 12),
              Row(children: [
                Expanded(
                  child: _Lbl(
                    label: tr(ref, 'cards_exp_month'),
                    child: TextField(
                      controller: month,
                      keyboardType: TextInputType.number,
                      inputFormatters: [
                        FilteringTextInputFormatter.digitsOnly,
                        LengthLimitingTextInputFormatter(2)
                      ],
                      decoration: const InputDecoration(hintText: 'MM'),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _Lbl(
                    label: tr(ref, 'cards_exp_year'),
                    child: TextField(
                      controller: year,
                      keyboardType: TextInputType.number,
                      inputFormatters: [
                        FilteringTextInputFormatter.digitsOnly,
                        LengthLimitingTextInputFormatter(4)
                      ],
                      decoration: const InputDecoration(hintText: 'YYYY'),
                    ),
                  ),
                ),
              ]),
              const SizedBox(height: 12),
              _Lbl(
                label: tr(ref, 'cards_holder'),
                child: TextField(
                  controller: holder,
                  textCapitalization: TextCapitalization.characters,
                ),
              ),
              const SizedBox(height: 6),
              SwitchListTile.adaptive(
                contentPadding: EdgeInsets.zero,
                title: Text(
                  tr(ref, 'cards_make_default'),
                  style: GoogleFonts.manrope(fontSize: 13, color: p.fg),
                ),
                value: def,
                activeThumbColor: p.primary,
                onChanged: (v) => setState(() => def = v),
              ),
              const SizedBox(height: 14),
              Row(children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () => Navigator.pop(ctx),
                    child: Text(tr(ref, 'common_cancel')),
                  ),
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: FilledButton(
                    onPressed: () async {
                      final ok = await ref.read(cardsViewModelProvider.notifier).add(
                            number: number.text,
                            expMonth: int.tryParse(month.text) ?? 0,
                            expYear: int.tryParse(year.text) ?? 0,
                            holder: holder.text,
                            makeDefault: def,
                          );
                      if (ok && ctx.mounted) Navigator.pop(ctx);
                    },
                    child: Text(tr(ref, 'common_save')),
                  ),
                ),
              ]),
            ],
          ),
        );
      }),
    );
  }
}

class _Lbl extends StatelessWidget {
  const _Lbl({required this.label, required this.child});
  final String label;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: GoogleFonts.manrope(
            fontSize: 12.5,
            fontWeight: FontWeight.w600,
            color: p.fg,
          ),
        ),
        const SizedBox(height: 6),
        child,
      ],
    );
  }
}
