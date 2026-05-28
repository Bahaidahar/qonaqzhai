import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../viewmodels/review_viewmodel.dart';

class ReviewSubmitScreen extends ConsumerStatefulWidget {
  const ReviewSubmitScreen({super.key, required this.bookingId, this.vendorId});

  final String bookingId;
  final String? vendorId;

  @override
  ConsumerState<ReviewSubmitScreen> createState() => _ReviewSubmitScreenState();
}

class _ReviewSubmitScreenState extends ConsumerState<ReviewSubmitScreen> {
  int _rating = 5;
  final _text = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref.read(reviewRepositoryProvider).submit(
            bookingId: widget.bookingId,
            rating: _rating,
            text: _text.text.trim().isEmpty ? null : _text.text.trim(),
          );
      if (!mounted) return;
      if (widget.vendorId != null) {
        ref.invalidate(vendorReviewsProvider(widget.vendorId!));
      }
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Thanks for the review!')),
      );
      context.pop();
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Text('Leave review',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
        children: [
          const AppPageHeader(
            title: 'How did it go?',
            subtitle: 'Your review helps other customers and the vendor.',
          ),
          const SizedBox(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: List.generate(
              5,
              (i) => GestureDetector(
                onTap: () => setState(() => _rating = i + 1),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 6),
                  child: Icon(
                    i < _rating ? CupertinoIcons.star_fill : CupertinoIcons.star,
                    size: 38,
                    color: i < _rating ? const Color(0xFFF59E0B) : p.border,
                  ),
                ),
              ),
            ),
          ),
          const SizedBox(height: 22),
          Text(
            'Comment (optional)',
            style: GoogleFonts.manrope(
              fontSize: 12.5,
              fontWeight: FontWeight.w600,
              color: p.fg,
            ),
          ),
          const SizedBox(height: 6),
          TextField(controller: _text, maxLines: 5),
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
          FilledButton(
            onPressed: _busy ? null : _submit,
            child: _busy
                ? const SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(
                        strokeWidth: 2, color: Colors.white),
                  )
                : const Text('Submit review'),
          ),
        ],
      ),
    );
  }
}
