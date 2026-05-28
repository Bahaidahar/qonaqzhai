import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../core/i18n/i18n.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/ui/ui.dart';
import '../../cards/presentation/cards_viewmodel.dart';
import '../data/payment_repository.dart';

class PaymentScreen extends ConsumerStatefulWidget {
  const PaymentScreen({super.key, required this.bookingId, required this.amount});
  final String bookingId;
  final int amount;

  @override
  ConsumerState<PaymentScreen> createState() => _PaymentScreenState();
}

class _PaymentScreenState extends ConsumerState<PaymentScreen> {
  String? _cardId;
  bool _busy = false;
  String? _error;

  @override
  Widget build(BuildContext context) {
    final cards = ref.watch(cardsViewModelProvider);
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(tr(ref, 'payment_title'),
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      body: Padding(
        padding: const EdgeInsets.fromLTRB(20, 12, 20, 20),
        child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
          const AppPageHeader(
            title: 'Pay booking',
            subtitle: 'Funds are held until your event is completed.',
          ),
          const SizedBox(height: 18),
          AppCard(
            padding: const EdgeInsets.fromLTRB(18, 16, 18, 16),
            child: Row(
              children: [
                Container(
                  width: 44,
                  height: 44,
                  decoration: BoxDecoration(
                    color: p.primary.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(color: p.primary.withValues(alpha: 0.25)),
                  ),
                  alignment: Alignment.center,
                  child:
                      Icon(CupertinoIcons.money_dollar, color: p.primary, size: 18),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        tr(ref, 'payment_amount'),
                        style: GoogleFonts.manrope(
                          fontSize: 11.5,
                          color: p.mutedFg,
                          letterSpacing: 0.2,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '${_kzt(widget.amount)} ₸',
                        style: GoogleFonts.manrope(
                          fontSize: 24,
                          fontWeight: FontWeight.w700,
                          color: p.fg,
                          letterSpacing: -0.6,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),
          const AppSectionHeader('Card'),
          const SizedBox(height: 10),
          if (cards.loading)
            LinearProgressIndicator(
              backgroundColor: p.muted,
              color: p.primary,
              minHeight: 2,
            ),
          if (cards.items.isEmpty && !cards.loading)
            AppEmptyState(
              message: tr(ref, 'cards_empty'),
              icon: CupertinoIcons.creditcard,
            ),
          Expanded(
            child: ListView.separated(
              padding: const EdgeInsets.only(top: 4),
              itemCount: cards.items.length,
              separatorBuilder: (_, __) => const SizedBox(height: 8),
              itemBuilder: (_, i) {
                final c = cards.items[i];
                final selected = _cardId == c.id;
                return Material(
                  color: Colors.transparent,
                  child: InkWell(
                    onTap: () => setState(() => _cardId = c.id),
                    borderRadius: BorderRadius.circular(12),
                    child: Container(
                      padding: const EdgeInsets.fromLTRB(14, 12, 14, 12),
                      decoration: BoxDecoration(
                        color: selected ? p.primary.withValues(alpha: 0.08) : p.card,
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: selected ? p.primary : p.border,
                          width: selected ? 1.4 : 1,
                        ),
                      ),
                      child: Row(
                        children: [
                          Icon(
                            selected
                                ? CupertinoIcons.checkmark_circle_fill
                                : CupertinoIcons.circle,
                            size: 18,
                            color: selected ? p.primary : p.mutedFg,
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  '${c.brand.toUpperCase()} •••• ${c.last4}',
                                  style: GoogleFonts.manrope(
                                    fontWeight: FontWeight.w700,
                                    fontSize: 14,
                                    color: p.fg,
                                  ),
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
                          if (c.isDefault)
                            const AppBadge(label: 'Default', tone: AppBadgeTone.info),
                        ],
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
          if (_error != null) ...[
            Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: Text(_error!,
                  style: GoogleFonts.manrope(fontSize: 12, color: p.destructive)),
            ),
          ],
          FilledButton.icon(
            onPressed: _busy ? null : _pay,
            icon: _busy
                ? const SizedBox(
                    height: 16,
                    width: 16,
                    child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                  )
                : const Icon(CupertinoIcons.lock_fill, size: 16),
            label: Text(tr(ref, 'payment_pay_mock')),
          ),
        ]),
      ),
    );
  }

  Future<void> _pay() async {
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final repo = PaymentRepository(ref.read(dioProvider));
      await repo.payMock(widget.bookingId, cardId: _cardId);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(tr(ref, 'payment_paid'))),
        );
        Navigator.of(context).pop(true);
      }
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _busy = false);
    }
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
