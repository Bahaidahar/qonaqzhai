import 'dart:async';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../viewmodels/booking_viewmodel.dart';

class BookingFormScreen extends ConsumerStatefulWidget {
  const BookingFormScreen({
    super.key,
    required this.vendorId,
    this.priceFrom = 0,
    this.serviceId,
    this.serviceUnit,
  });

  final String vendorId;
  final int priceFrom;
  final String? serviceId;
  final String? serviceUnit;

  @override
  ConsumerState<BookingFormScreen> createState() => _BookingFormScreenState();
}

class _BookingFormScreenState extends ConsumerState<BookingFormScreen> {
  final _date = TextEditingController();
  final _guests = TextEditingController(text: '50');
  final _note = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _guests.addListener(() => setState(() {}));
  }

  @override
  void dispose() {
    _date.dispose();
    _guests.dispose();
    _note.dispose();
    super.dispose();
  }

  int get _guestCount {
    final raw = int.tryParse(_guests.text);
    return (raw == null || raw < 1) ? 1 : raw;
  }

  int get _amount {
    // Per-person pricing scales with guest count, matching the web pricing
    // rule on `vendors/[id]` BookingPanel.
    if (widget.serviceUnit == 'person' && widget.priceFrom > 0) {
      return widget.priceFrom * _guestCount;
    }
    return widget.priceFrom;
  }

  Future<void> _pickDate() async {
    final picked = await showDatePicker(
      context: context,
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365 * 2)),
      initialDate: DateTime.now().add(const Duration(days: 30)),
    );
    if (picked != null) {
      _date.text =
          '${picked.year.toString().padLeft(4, '0')}-${picked.month.toString().padLeft(2, '0')}-${picked.day.toString().padLeft(2, '0')}';
      setState(() {});
    }
  }

  Future<void> _submit() async {
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final repo = ref.read(bookingRepositoryProvider);
      await repo.create(
        vendorId: widget.vendorId,
        eventDate: _date.text.trim(),
        guestCount: _guestCount,
        note: _note.text.trim().isEmpty ? null : _note.text.trim(),
        amount: _amount > 0 ? _amount : null,
        serviceId: widget.serviceId,
      );
      if (!mounted) return;
      unawaited(ref.read(bookingsProvider.notifier).refresh());
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Booking request sent')),
      );
      context.go('/bookings');
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    final showPriceCard = _amount > 0;
    return Scaffold(
      appBar: AppBar(
        title: Text('Book vendor',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
        children: [
          const AppPageHeader(
            title: 'Request booking',
            subtitle: 'The vendor reviews and accepts before any charge.',
          ),
          const SizedBox(height: 20),
          if (showPriceCard)
            AppCard(
              padding: const EdgeInsets.fromLTRB(16, 14, 16, 14),
              child: Row(
                children: [
                  Icon(CupertinoIcons.money_dollar_circle,
                      color: p.primary, size: 18),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.serviceUnit == 'person'
                              ? 'Estimated total (per person)'
                              : 'Estimated price',
                          style: GoogleFonts.manrope(
                              fontSize: 12.5, color: p.mutedFg),
                        ),
                        if (widget.serviceUnit == 'person' &&
                            widget.priceFrom > 0)
                          Text(
                            '${_kzt(widget.priceFrom)} ₸ × $_guestCount guests',
                            style: GoogleFonts.manrope(
                                fontSize: 11, color: p.mutedFg),
                          ),
                      ],
                    ),
                  ),
                  Text(
                    '${_kzt(_amount)} ₸',
                    style: GoogleFonts.manrope(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: p.fg,
                    ),
                  ),
                ],
              ),
            ),
          if (showPriceCard) const SizedBox(height: 16),
          _Lbl(
            label: 'Event date',
            child: TextField(
              controller: _date,
              readOnly: true,
              onTap: _pickDate,
              decoration: InputDecoration(
                hintText: 'YYYY-MM-DD',
                suffixIcon:
                    Icon(CupertinoIcons.calendar, color: p.mutedFg, size: 18),
              ),
            ),
          ),
          const SizedBox(height: 12),
          _Lbl(
            label: 'Guests',
            child: TextField(
              controller: _guests,
              keyboardType: TextInputType.number,
            ),
          ),
          const SizedBox(height: 12),
          _Lbl(
            label: 'Note (optional)',
            child: TextField(
              controller: _note,
              maxLines: 3,
            ),
          ),
          if (_error != null) ...[
            const SizedBox(height: 14),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: BoxDecoration(
                color: p.destructive.withValues(alpha: 0.08),
                borderRadius: BorderRadius.circular(10),
                border: Border.all(color: p.destructive.withValues(alpha: 0.3)),
              ),
              child: Text(_error!,
                  style: GoogleFonts.manrope(fontSize: 12, color: p.destructive)),
            ),
          ],
          const SizedBox(height: 22),
          FilledButton.icon(
            onPressed: _busy || _date.text.isEmpty ? null : _submit,
            icon: _busy
                ? const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(
                        strokeWidth: 2, color: Colors.white),
                  )
                : const Icon(CupertinoIcons.paperplane_fill, size: 16),
            label: const Text('Send request'),
          ),
        ],
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
