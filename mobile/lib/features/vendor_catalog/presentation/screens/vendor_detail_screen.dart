import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/ui/ui.dart';
import '../../../reviews/presentation/viewmodels/review_viewmodel.dart';
import '../../domain/entities/service.dart';
import '../../domain/entities/vendor.dart';
import '../viewmodels/vendor_catalog_viewmodel.dart';

final vendorByIdProvider = FutureProvider.family<Vendor, String>((ref, id) {
  return ref.watch(vendorRepositoryProvider).byId(id);
});

class VendorDetailScreen extends ConsumerStatefulWidget {
  const VendorDetailScreen({super.key, required this.id});

  final String id;

  @override
  ConsumerState<VendorDetailScreen> createState() => _VendorDetailScreenState();
}

class _VendorDetailScreenState extends ConsumerState<VendorDetailScreen> {
  VendorService? _selectedService;

  @override
  Widget build(BuildContext context) {
    final vendorAsync = ref.watch(vendorByIdProvider(widget.id));
    final reviewsAsync = ref.watch(vendorReviewsProvider(widget.id));
    final servicesAsync = ref.watch(vendorServicesProvider(widget.id));
    final p = AppPalette.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text('Vendor',
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
      ),
      bottomNavigationBar: vendorAsync.maybeWhen(
        data: (v) => SafeArea(
          minimum: const EdgeInsets.fromLTRB(20, 8, 20, 16),
          child: SizedBox(
            width: double.infinity,
            child: FilledButton.icon(
              onPressed: () {
                final price = _selectedService?.price ?? v.priceFrom;
                final params = <String, String>{
                  'vendor': v.id,
                  'price': price.toString(),
                };
                if (_selectedService != null) {
                  params['service'] = _selectedService!.id;
                  params['unit'] = _selectedService!.unit;
                }
                final qs = Uri(queryParameters: params).query;
                context.push('/bookings/new?$qs');
              },
              icon: const Icon(CupertinoIcons.calendar_badge_plus, size: 18),
              label: Text(
                _selectedService == null
                    ? 'Book now'
                    : 'Book — ${_kzt(_selectedService!.price)} ₸',
                style: GoogleFonts.manrope(fontWeight: FontWeight.w700),
              ),
            ),
          ),
        ),
        orElse: () => const SizedBox.shrink(),
      ),
      body: vendorAsync.when(
        loading: () => Center(child: CupertinoActivityIndicator(color: p.mutedFg)),
        error: (e, _) => Center(
          child: Text(e.toString(),
              style: GoogleFonts.manrope(color: p.destructive)),
        ),
        data: (v) => ListView(
          padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
          children: [
            if (v.photoIds.isNotEmpty) ...[
              SizedBox(
                height: 220,
                child: ListView.separated(
                  scrollDirection: Axis.horizontal,
                  itemCount: v.photoIds.length,
                  separatorBuilder: (_, __) => const SizedBox(width: 8),
                  itemBuilder: (_, i) => ClipRRect(
                    borderRadius: BorderRadius.circular(14),
                    child: CachedNetworkImage(
                      imageUrl:
                          '${ApiEndpoints.baseUrl}${ApiEndpoints.photo(v.photoIds[i])}',
                      width: 300,
                      fit: BoxFit.cover,
                      placeholder: (_, __) => Container(color: p.muted),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 20),
            ],
            Text(
              v.name,
              style: GoogleFonts.manrope(
                fontSize: 28,
                fontWeight: FontWeight.w700,
                color: p.fg,
                letterSpacing: -0.8,
                height: 1.1,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${v.category} · ${v.city}',
              style: GoogleFonts.manrope(fontSize: 13, color: p.mutedFg),
            ),
            const SizedBox(height: 14),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                if (v.priceFrom > 0)
                  AppBadge(
                      label: 'from ${_kzt(v.priceFrom)} ₸',
                      tone: AppBadgeTone.info),
                if (v.ratingCount > 0)
                  AppBadge(
                    label:
                        '${v.ratingAvg.toStringAsFixed(1)} · ${v.ratingCount} reviews',
                    tone: AppBadgeTone.success,
                    icon: CupertinoIcons.star_fill,
                  ),
              ],
            ),
            if (v.description.isNotEmpty) ...[
              const SizedBox(height: 20),
              const AppSectionHeader('About'),
              const SizedBox(height: 10),
              Text(
                v.description,
                style: GoogleFonts.manrope(
                  fontSize: 14,
                  color: p.fg,
                  height: 1.55,
                ),
              ),
            ],

            const SizedBox(height: 24),
            const AppSectionHeader('Services'),
            const SizedBox(height: 10),
            servicesAsync.when(
              loading: () => Padding(
                padding: const EdgeInsets.all(12),
                child: Center(
                    child: CupertinoActivityIndicator(color: p.mutedFg)),
              ),
              error: (e, _) => const AppEmptyState(
                message: 'Could not load services.',
                icon: CupertinoIcons.square_list,
              ),
              data: (services) => services.isEmpty
                  ? const AppEmptyState(
                      message:
                          'This vendor has not listed any service packages yet.',
                      icon: CupertinoIcons.square_list,
                    )
                  : Column(
                      children: [
                        for (final s in services) ...[
                          _ServiceTile(
                            service: s,
                            selected: _selectedService?.id == s.id,
                            onTap: () => setState(() => _selectedService =
                                _selectedService?.id == s.id ? null : s),
                          ),
                          const SizedBox(height: 8),
                        ],
                      ],
                    ),
            ),

            const SizedBox(height: 24),
            const AppSectionHeader('Reviews'),
            const SizedBox(height: 12),
            reviewsAsync.when(
              loading: () => Center(
                  child: Padding(
                padding: const EdgeInsets.all(16),
                child: CupertinoActivityIndicator(color: p.mutedFg),
              )),
              error: (e, _) => Text(e.toString(),
                  style: GoogleFonts.manrope(color: p.destructive)),
              data: (list) => list.isEmpty
                  ? const AppEmptyState(
                      message: 'No reviews yet.',
                      icon: CupertinoIcons.chat_bubble_text,
                    )
                  : Column(
                      children: [
                        for (final r in list)
                          Padding(
                            padding: const EdgeInsets.only(bottom: 10),
                            child: AppCard(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Row(
                                    children: List.generate(5, (i) {
                                      return Padding(
                                        padding: const EdgeInsets.only(right: 3),
                                        child: Icon(
                                          i < r.rating
                                              ? CupertinoIcons.star_fill
                                              : CupertinoIcons.star,
                                          size: 14,
                                          color: i < r.rating
                                              ? const Color(0xFFF59E0B)
                                              : p.border,
                                        ),
                                      );
                                    }),
                                  ),
                                  if (r.text.isNotEmpty) ...[
                                    const SizedBox(height: 8),
                                    Text(
                                      r.text,
                                      style: GoogleFonts.manrope(
                                        fontSize: 13.5,
                                        height: 1.5,
                                        color: p.fg,
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                            ),
                          ),
                      ],
                    ),
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

class _ServiceTile extends StatelessWidget {
  const _ServiceTile({
    required this.service,
    required this.selected,
    required this.onTap,
  });

  final VendorService service;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
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
            crossAxisAlignment: CrossAxisAlignment.start,
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
                      service.name,
                      style: GoogleFonts.manrope(
                        fontSize: 14,
                        fontWeight: FontWeight.w700,
                        color: p.fg,
                      ),
                    ),
                    if (service.description.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(
                        service.description,
                        style: GoogleFonts.manrope(
                          fontSize: 12,
                          color: p.mutedFg,
                          height: 1.4,
                        ),
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(width: 10),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${_kzt(service.price)} ₸',
                    style: GoogleFonts.manrope(
                      fontSize: 13,
                      fontWeight: FontWeight.w700,
                      color: p.fg,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '/ ${service.unit}',
                    style: GoogleFonts.manrope(
                      fontSize: 11,
                      color: p.mutedFg,
                    ),
                  ),
                ],
              ),
            ],
          ),
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
