import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../viewmodels/booking_viewmodel.dart';

class BookingFormScreen extends ConsumerStatefulWidget {
  const BookingFormScreen({
    super.key,
    required this.vendorId,
    this.priceFrom = 0,
  });

  final String vendorId;
  final int priceFrom;

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
  void dispose() {
    _date.dispose();
    _guests.dispose();
    _note.dispose();
    super.dispose();
  }

  Future<void> _pickDate() async {
    final picked = await showDatePicker(
      context: context,
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365 * 2)),
      initialDate: DateTime.now().add(const Duration(days: 30)),
    );
    if (picked != null) {
      _date.text = '${picked.year.toString().padLeft(4, '0')}-${picked.month.toString().padLeft(2, '0')}-${picked.day.toString().padLeft(2, '0')}';
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
        guestCount: int.tryParse(_guests.text) ?? 0,
        note: _note.text.trim().isEmpty ? null : _note.text.trim(),
        amount: widget.priceFrom > 0 ? widget.priceFrom : null,
      );
      if (!mounted) return;
      // refresh inbox + bookings list
      ref.read(bookingsProvider.notifier).refresh();
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
    return Scaffold(
      appBar: AppBar(title: const Text('Book vendor')),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: ListView(
          children: [
            if (widget.priceFrom > 0)
              Text(
                'Estimated price: ${widget.priceFrom} ₸',
                style: Theme.of(context).textTheme.titleMedium,
              ),
            const SizedBox(height: 16),
            TextField(
              controller: _date,
              readOnly: true,
              onTap: _pickDate,
              decoration: const InputDecoration(
                labelText: 'Event date',
                suffixIcon: Icon(Icons.calendar_today),
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _guests,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(labelText: 'Guests'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _note,
              maxLines: 3,
              decoration: const InputDecoration(labelText: 'Note (optional)'),
            ),
            const SizedBox(height: 20),
            if (_error != null) ...[
              Text(_error!, style: const TextStyle(color: Colors.red)),
              const SizedBox(height: 12),
            ],
            FilledButton(
              onPressed: _busy || _date.text.isEmpty ? null : _submit,
              child: _busy
                  ? const SizedBox(width: 16, height: 16, child: CircularProgressIndicator(strokeWidth: 2))
                  : const Text('Send request'),
            ),
          ],
        ),
      ),
    );
  }
}
