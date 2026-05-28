import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../../domain/entities/vendor.dart';
import '../viewmodels/vendor_catalog_viewmodel.dart';

const _categories = ['All', 'Venue', 'Catering', 'Photo', 'Video', 'Decor', 'Music', 'Cakes'];
const _fixedCity = 'Almaty';
const _sorts = {
  'newest': 'Newest',
  'price_asc': 'Price ↑',
  'price_desc': 'Price ↓',
  'rating_desc': 'Top rated',
};

class VendorCatalogScreen extends ConsumerStatefulWidget {
  const VendorCatalogScreen({super.key});

  @override
  ConsumerState<VendorCatalogScreen> createState() => _VendorCatalogScreenState();
}

class _VendorCatalogScreenState extends ConsumerState<VendorCatalogScreen> {
  final _search = TextEditingController();
  String _category = 'All';
  String _sort = 'newest';
  int? _priceMax;

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  void _apply() {
    ref.read(vendorCatalogProvider.notifier).load(VendorQuery(
          query: _search.text.trim().isEmpty ? null : _search.text.trim(),
          category: _category == 'All' ? null : _category,
          city: _fixedCity,
          priceMax: _priceMax,
          sort: _sort,
          limit: 30,
        ));
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(vendorCatalogProvider);
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text('Vendors',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
        actions: [
          IconButton(
            icon: const Icon(CupertinoIcons.bell),
            onPressed: () => context.push('/notifications'),
          ),
          IconButton(
            icon: const Icon(CupertinoIcons.slider_horizontal_3),
            onPressed: _openFilters,
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 8, 20, 12),
            child: TextField(
              controller: _search,
              style: GoogleFonts.manrope(fontSize: 14, color: p.fg),
              decoration: InputDecoration(
                prefixIcon: Icon(CupertinoIcons.search, color: p.mutedFg, size: 18),
                hintText: 'Search venues, photo, music…',
              ),
              onSubmitted: (_) => _apply(),
            ),
          ),
          SizedBox(
            height: 36,
            child: ListView.separated(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              scrollDirection: Axis.horizontal,
              itemCount: _categories.length,
              separatorBuilder: (_, __) => const SizedBox(width: 8),
              itemBuilder: (_, i) {
                final c = _categories[i];
                final selected = c == _category;
                return GestureDetector(
                  onTap: () {
                    setState(() => _category = c);
                    _apply();
                  },
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 14),
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: selected ? p.primary : p.muted,
                      borderRadius: BorderRadius.circular(999),
                      border: Border.all(
                        color: selected ? p.primary : p.border,
                      ),
                    ),
                    child: Text(
                      c,
                      style: GoogleFonts.manrope(
                        fontSize: 12.5,
                        fontWeight: FontWeight.w600,
                        color: selected ? p.onPrimary : p.fg,
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
          const SizedBox(height: 12),
          if (state.loading)
            LinearProgressIndicator(
              backgroundColor: p.muted,
              color: p.primary,
              minHeight: 2,
            ),
          if (state.total > 0)
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 4, 20, 0),
              child: Align(
                alignment: Alignment.centerLeft,
                child: Text(
                  '${state.total} vendors',
                  style: GoogleFonts.manrope(fontSize: 12, color: p.mutedFg),
                ),
              ),
            ),
          Expanded(
            child: state.error != null
                ? Center(
                    child: Text(
                      state.error!,
                      style: GoogleFonts.manrope(color: p.destructive),
                    ),
                  )
                : state.items.isEmpty && !state.loading
                    ? const Padding(
                        padding: EdgeInsets.all(20),
                        child: AppEmptyState(
                          message: 'No vendors match your filters yet.',
                          icon: CupertinoIcons.square_grid_2x2,
                        ),
                      )
                    : ListView.separated(
                        padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
                        itemCount: state.items.length,
                        separatorBuilder: (_, __) => const SizedBox(height: 10),
                        itemBuilder: (_, i) {
                          return _VendorRow(
                            vendor: state.items[i],
                            onTap: () => context.push('/vendors/${state.items[i].id}'),
                          );
                        },
                      ),
          ),
        ],
      ),
    );
  }

  Future<void> _openFilters() async {
    final priceCtrl = TextEditingController(text: _priceMax?.toString() ?? '');
    final p = AppPalette.of(context);
    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      backgroundColor: p.card,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (_) => StatefulBuilder(
        builder: (ctx, setSt) => Padding(
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
              const SizedBox(height: 16),
              Text(
                'Filters',
                style: GoogleFonts.manrope(
                  fontSize: 22,
                  fontWeight: FontWeight.w700,
                  letterSpacing: -0.5,
                  color: p.fg,
                ),
              ),
              const SizedBox(height: 20),
              const AppSectionHeader('Sort'),
              const SizedBox(height: 10),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: _sorts.entries
                    .map((e) => _PillChoice(
                          label: e.value,
                          selected: _sort == e.key,
                          onTap: () => setSt(() => _sort = e.key),
                        ))
                    .toList(),
              ),
              const SizedBox(height: 20),
              const AppSectionHeader('Max price'),
              const SizedBox(height: 10),
              TextField(
                controller: priceCtrl,
                keyboardType: TextInputType.number,
                style: GoogleFonts.manrope(fontSize: 14, color: p.fg),
                decoration: const InputDecoration(hintText: '500000'),
              ),
              const SizedBox(height: 24),
              FilledButton(
                onPressed: () {
                  _priceMax = int.tryParse(priceCtrl.text);
                  Navigator.pop(ctx);
                  setState(() {});
                  _apply();
                },
                child: const Text('Apply'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _PillChoice extends StatelessWidget {
  const _PillChoice({required this.label, required this.selected, required this.onTap});
  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
        decoration: BoxDecoration(
          color: selected ? p.primary : p.muted,
          borderRadius: BorderRadius.circular(999),
          border: Border.all(color: selected ? p.primary : p.border),
        ),
        child: Text(
          label,
          style: GoogleFonts.manrope(
            fontSize: 12.5,
            fontWeight: FontWeight.w600,
            color: selected ? p.onPrimary : p.fg,
          ),
        ),
      ),
    );
  }
}

class _VendorRow extends StatelessWidget {
  const _VendorRow({required this.vendor, required this.onTap});
  final Vendor vendor;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    final photoId = vendor.photoIds.isNotEmpty ? vendor.photoIds.first : null;
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(14),
        child: Container(
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: p.card,
            borderRadius: BorderRadius.circular(14),
            border: Border.all(color: p.border),
          ),
          child: Row(
            children: [
              ClipRRect(
                borderRadius: BorderRadius.circular(10),
                child: SizedBox(
                  width: 76,
                  height: 76,
                  child: photoId == null
                      ? Container(
                          color: p.muted,
                          alignment: Alignment.center,
                          child: Icon(CupertinoIcons.photo, color: p.mutedFg, size: 20),
                        )
                      : CachedNetworkImage(
                          imageUrl: '${ApiEndpoints.baseUrl}${ApiEndpoints.photo(photoId)}',
                          fit: BoxFit.cover,
                          placeholder: (_, __) => Container(color: p.muted),
                          errorWidget: (_, __, ___) => Container(
                            color: p.muted,
                            alignment: Alignment.center,
                            child: Icon(CupertinoIcons.photo, color: p.mutedFg, size: 18),
                          ),
                        ),
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      vendor.name,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: GoogleFonts.manrope(
                        fontSize: 15,
                        fontWeight: FontWeight.w700,
                        color: p.fg,
                        letterSpacing: -0.2,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${vendor.category} · ${vendor.city}',
                      style: GoogleFonts.manrope(fontSize: 12, color: p.mutedFg),
                    ),
                    const SizedBox(height: 6),
                    Row(
                      children: [
                        Text(
                          'from ${_formatKzt(vendor.priceFrom)} ₸',
                          style: GoogleFonts.manrope(
                            fontSize: 12.5,
                            fontWeight: FontWeight.w600,
                            color: p.fg,
                          ),
                        ),
                        if (vendor.ratingCount > 0) ...[
                          const SizedBox(width: 10),
                          const Icon(CupertinoIcons.star_fill,
                              size: 11, color: Color(0xFFF59E0B)),
                          const SizedBox(width: 2),
                          Text(
                            vendor.ratingAvg.toStringAsFixed(1),
                            style: GoogleFonts.manrope(
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                              color: p.fg,
                            ),
                          ),
                        ],
                      ],
                    ),
                  ],
                ),
              ),
              Icon(CupertinoIcons.chevron_forward, size: 16, color: p.mutedFg),
            ],
          ),
        ),
      ),
    );
  }

  static String _formatKzt(int v) {
    final s = v.toString();
    final buf = StringBuffer();
    for (var i = 0; i < s.length; i++) {
      if (i > 0 && (s.length - i) % 3 == 0) buf.write(' ');
      buf.write(s[i]);
    }
    return buf.toString();
  }
}
