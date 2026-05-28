import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:image_picker/image_picker.dart';

import '../../../core/i18n/i18n.dart';
import '../../../core/network/api_endpoints.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/ui/ui.dart';
import '../data/vendor_self_repository.dart';

const _kCategories = ['Venue', 'Catering', 'Photo', 'Video', 'Decor', 'Music', 'Cakes'];
const _kUnits = ['fixed', 'hour', 'item', 'person', 'day'];

class VendorSelfScreen extends ConsumerStatefulWidget {
  const VendorSelfScreen({super.key});

  @override
  ConsumerState<VendorSelfScreen> createState() => _VendorSelfScreenState();
}

class _VendorSelfScreenState extends ConsumerState<VendorSelfScreen> {
  final _name = TextEditingController();
  final _city = TextEditingController(text: 'Almaty');
  final _price = TextEditingController();
  final _description = TextEditingController();
  String _category = 'Venue';
  bool _loading = true;
  bool _saving = false;
  String? _error;
  Map<String, dynamic>? _vendor;
  List<Map<String, dynamic>> _services = [];

  late VendorSelfRepository _repo;

  @override
  void initState() {
    super.initState();
    _repo = VendorSelfRepository(ref.read(dioProvider));
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      _vendor = await _repo.myVendor();
      if (_vendor != null) {
        _name.text = (_vendor!['name'] ?? '').toString();
        _city.text = (_vendor!['city'] ?? 'Almaty').toString();
        _price.text = (_vendor!['priceFrom'] ?? 0).toString();
        _description.text = (_vendor!['description'] ?? '').toString();
        _category = (_vendor!['category'] ?? 'Venue').toString();
        _services = await _repo.myServices();
      }
    } catch (e) {
      _error = e.toString();
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _save() async {
    setState(() => _saving = true);
    try {
      _vendor = await _repo.upsert(
        name: _name.text.trim(),
        category: _category,
        city: _city.text.trim(),
        priceFrom: int.tryParse(_price.text) ?? 0,
        description: _description.text.trim(),
      );
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(tr(ref, 'common_save'))),
      );
      await _load();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _uploadPhoto() async {
    final picker = ImagePicker();
    final picked = await picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1600,
      imageQuality: 85,
    );
    if (picked == null) return;
    final bytes = await picked.readAsBytes();
    await _repo.uploadPhoto(bytes: bytes, filename: picked.name);
    await _load();
  }

  Future<void> _deletePhoto(String id) async {
    await _repo.deletePhoto(id);
    await _load();
  }

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    if (_loading) {
      return Scaffold(
        appBar: AppBar(title: Text(tr(ref, 'vendor_profile_title'))),
        body: Center(child: CupertinoActivityIndicator(color: p.mutedFg)),
      );
    }

    final photoIds = (_vendor?['photoIds'] as List?)?.cast<String>() ?? const <String>[];
    final status = _vendor?['status'] as String? ?? 'pending';

    return Scaffold(
      appBar: AppBar(
        title: Text(tr(ref, 'vendor_profile_title'),
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      body: RefreshIndicator(
        color: p.primary,
        onRefresh: _load,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(20, 12, 20, 32),
          children: [
            const AppPageHeader(
              title: 'My profile',
              subtitle: 'Your public listing in the qonaqzhai catalog.',
            ),
            const SizedBox(height: 18),
            if (_vendor != null) ...[
              _StatusBadge(status: status),
              const SizedBox(height: 18),
            ],
            if (_error != null) ...[
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
              const SizedBox(height: 14),
            ],
            AppCard(
              padding: const EdgeInsets.fromLTRB(16, 18, 16, 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const AppSectionHeader('Basics'),
                  const SizedBox(height: 14),
                  _Field(
                    label: tr(ref, 'vendor_name'),
                    child: TextField(
                      controller: _name,
                      decoration: const InputDecoration(hintText: 'Rixos Almaty Ballroom'),
                    ),
                  ),
                  const SizedBox(height: 14),
                  _Field(
                    label: tr(ref, 'vendor_category'),
                    child: DropdownButtonFormField<String>(
                      initialValue: _category,
                      isExpanded: true,
                      items: _kCategories
                          .map((c) => DropdownMenuItem(value: c, child: Text(c)))
                          .toList(),
                      onChanged: (v) => setState(() => _category = v ?? 'Venue'),
                    ),
                  ),
                  const SizedBox(height: 14),
                  _Field(
                    label: tr(ref, 'vendor_city'),
                    child: TextField(controller: _city),
                  ),
                  const SizedBox(height: 14),
                  _Field(
                    label: tr(ref, 'vendor_price_from'),
                    child: TextField(
                      controller: _price,
                      keyboardType: TextInputType.number,
                      decoration: const InputDecoration(hintText: '500000'),
                    ),
                  ),
                  const SizedBox(height: 14),
                  _Field(
                    label: tr(ref, 'vendor_description'),
                    child: TextField(
                      controller: _description,
                      maxLines: 4,
                      decoration: const InputDecoration(
                          hintText: 'Premier venue in the heart of Almaty…'),
                    ),
                  ),
                  const SizedBox(height: 18),
                  Align(
                    alignment: Alignment.centerRight,
                    child: FilledButton.icon(
                      onPressed: _saving ? null : _save,
                      icon: const Icon(CupertinoIcons.checkmark_alt, size: 18),
                      label: Text(_saving ? '…' : tr(ref, 'vendor_save_profile')),
                    ),
                  ),
                ],
              ),
            ),
            if (_vendor != null) ...[
              const SizedBox(height: 20),
              AppCard(
                padding: const EdgeInsets.fromLTRB(16, 18, 16, 18),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Expanded(child: AppSectionHeader('Photos')),
                        OutlinedButton.icon(
                          onPressed: _uploadPhoto,
                          icon: const Icon(CupertinoIcons.cloud_upload, size: 16),
                          label: Text(tr(ref, 'vendor_upload_photo')),
                          style: OutlinedButton.styleFrom(
                            padding: const EdgeInsets.symmetric(horizontal: 12),
                            minimumSize: const Size(0, 36),
                            textStyle:
                                GoogleFonts.manrope(fontSize: 12, fontWeight: FontWeight.w600),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    if (photoIds.isEmpty)
                      const AppEmptyState(
                        message: 'No photos yet.',
                        icon: CupertinoIcons.photo,
                      )
                    else
                      Wrap(
                        spacing: 10,
                        runSpacing: 10,
                        children: [
                          for (final pid in photoIds)
                            Stack(
                              children: [
                                ClipRRect(
                                  borderRadius: BorderRadius.circular(10),
                                  child: CachedNetworkImage(
                                    imageUrl:
                                        ApiEndpoints.baseUrl + ApiEndpoints.photo(pid),
                                    width: 96,
                                    height: 96,
                                    fit: BoxFit.cover,
                                    errorWidget: (_, __, ___) => Container(
                                      width: 96,
                                      height: 96,
                                      color: p.muted,
                                    ),
                                  ),
                                ),
                                Positioned(
                                  right: 4,
                                  top: 4,
                                  child: GestureDetector(
                                    onTap: () => _deletePhoto(pid),
                                    child: Container(
                                      padding: const EdgeInsets.all(5),
                                      decoration: BoxDecoration(
                                        color: p.card.withValues(alpha: 0.9),
                                        borderRadius: BorderRadius.circular(8),
                                        border: Border.all(color: p.border),
                                      ),
                                      child: Icon(CupertinoIcons.trash,
                                          size: 14, color: p.destructive),
                                    ),
                                  ),
                                ),
                              ],
                            ),
                        ],
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 20),
              AppCard(
                padding: const EdgeInsets.fromLTRB(16, 18, 16, 12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Expanded(child: AppSectionHeader('Services')),
                        OutlinedButton.icon(
                          onPressed: () => _showServiceSheet(null),
                          icon: const Icon(CupertinoIcons.add, size: 16),
                          label: Text(tr(ref, 'vendor_service_add')),
                          style: OutlinedButton.styleFrom(
                            padding: const EdgeInsets.symmetric(horizontal: 12),
                            minimumSize: const Size(0, 36),
                            textStyle:
                                GoogleFonts.manrope(fontSize: 12, fontWeight: FontWeight.w600),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    if (_services.isEmpty)
                      const AppEmptyState(
                        message: 'No services yet.',
                        icon: CupertinoIcons.square_list,
                      )
                    else
                      for (final s in _services) ...[
                        _ServiceRow(
                          name: s['name']?.toString() ?? '',
                          description: s['description']?.toString() ?? '',
                          price: '${s['price']} ₸ · ${s['unit']}',
                          onEdit: () => _showServiceSheet(s),
                          onDelete: () async {
                            await _repo.deleteService(s['id'].toString());
                            await _load();
                          },
                        ),
                        const SizedBox(height: 8),
                      ],
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  void _showServiceSheet(Map<String, dynamic>? existing) {
    final isEdit = existing != null;
    final name = TextEditingController(text: (existing?['name'] ?? '').toString());
    final desc =
        TextEditingController(text: (existing?['description'] ?? '').toString());
    final price = TextEditingController(
        text: ((existing?['price'] as num?)?.toInt() ?? 0).toString());
    String unit = (existing?['unit'] as String?) ?? 'fixed';
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
                isEdit
                    ? (tr(ref, 'vendor_service_edit') == 'vendor_service_edit'
                        ? 'Edit service'
                        : tr(ref, 'vendor_service_edit'))
                    : tr(ref, 'vendor_service_add'),
                style: GoogleFonts.manrope(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  letterSpacing: -0.4,
                  color: p.fg,
                ),
              ),
              const SizedBox(height: 16),
              _Field(label: tr(ref, 'service_name'), child: TextField(controller: name)),
              const SizedBox(height: 12),
              _Field(
                label: tr(ref, 'service_description'),
                child: TextField(controller: desc, maxLines: 2),
              ),
              const SizedBox(height: 12),
              _Field(
                label: tr(ref, 'service_price'),
                child: TextField(controller: price, keyboardType: TextInputType.number),
              ),
              const SizedBox(height: 12),
              _Field(
                label: tr(ref, 'service_unit'),
                child: DropdownButtonFormField<String>(
                  initialValue: unit,
                  isExpanded: true,
                  items:
                      _kUnits.map((u) => DropdownMenuItem(value: u, child: Text(u))).toList(),
                  onChanged: (v) => setState(() => unit = v ?? 'fixed'),
                ),
              ),
              const SizedBox(height: 20),
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
                      if (isEdit) {
                        await _repo.updateService(
                          id: existing['id'].toString(),
                          name: name.text.trim(),
                          description: desc.text.trim(),
                          price: int.tryParse(price.text) ?? 0,
                          unit: unit,
                        );
                      } else {
                        await _repo.addService(
                          name: name.text.trim(),
                          description: desc.text.trim(),
                          price: int.tryParse(price.text) ?? 0,
                          unit: unit,
                        );
                      }
                      if (ctx.mounted) Navigator.pop(ctx);
                      await _load();
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

class _Field extends StatelessWidget {
  const _Field({required this.label, required this.child});
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

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});
  final String status;

  @override
  Widget build(BuildContext context) {
    final (tone, label, icon) = switch (status) {
      'approved' => (
          AppBadgeTone.success,
          'Approved — your profile is live in the catalog.',
          CupertinoIcons.checkmark_seal_fill,
        ),
      'rejected' => (
          AppBadgeTone.danger,
          'Rejected — contact support to resolve.',
          CupertinoIcons.exclamationmark_triangle_fill,
        ),
      _ => (
          AppBadgeTone.warning,
          'Pending admin approval — your profile is hidden from customers until reviewed.',
          CupertinoIcons.clock_fill,
        ),
    };
    return _StatusPanel(tone: tone, message: label, icon: icon);
  }
}

class _StatusPanel extends StatelessWidget {
  const _StatusPanel({required this.tone, required this.message, required this.icon});
  final AppBadgeTone tone;
  final String message;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    final (Color border, Color bg, Color fg) = switch (tone) {
      AppBadgeTone.success => (
          const Color(0x4D34C759),
          const Color(0x1A34C759),
          const Color(0xFF1B7A3A),
        ),
      AppBadgeTone.warning => (
          const Color(0x4DF59E0B),
          const Color(0x1AF59E0B),
          const Color(0xFFB45309),
        ),
      AppBadgeTone.danger => (
          const Color(0x4DDC2626),
          const Color(0x1ADC2626),
          const Color(0xFFB91C1C),
        ),
      _ => (p.border, p.muted, p.mutedFg),
    };
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: border),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 14, color: fg),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              message,
              style: GoogleFonts.manrope(
                fontSize: 12.5,
                color: fg,
                height: 1.4,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _ServiceRow extends StatelessWidget {
  const _ServiceRow({
    required this.name,
    required this.description,
    required this.price,
    required this.onEdit,
    required this.onDelete,
  });

  final String name;
  final String description;
  final String price;
  final VoidCallback onEdit;
  final Future<void> Function() onDelete;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Container(
      padding: const EdgeInsets.fromLTRB(12, 10, 8, 10),
      decoration: BoxDecoration(
        color: p.muted,
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: p.border),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(name,
                    style: GoogleFonts.manrope(
                        fontWeight: FontWeight.w700, fontSize: 13.5, color: p.fg)),
                if (description.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 2),
                    child: Text(
                      description,
                      style: GoogleFonts.manrope(fontSize: 12, color: p.mutedFg),
                    ),
                  ),
                const SizedBox(height: 4),
                Text(price,
                    style: GoogleFonts.manrope(
                        fontSize: 12, fontWeight: FontWeight.w600, color: p.fg)),
              ],
            ),
          ),
          IconButton(
            icon: Icon(CupertinoIcons.pencil, color: p.mutedFg, size: 18),
            onPressed: onEdit,
          ),
          IconButton(
            icon: Icon(CupertinoIcons.trash, color: p.mutedFg, size: 18),
            onPressed: () => onDelete(),
          ),
        ],
      ),
    );
  }
}
