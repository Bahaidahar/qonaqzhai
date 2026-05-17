import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../domain/entities/vendor.dart';
import '../viewmodels/vendor_catalog_viewmodel.dart';

const _categories = ['All', 'Venue', 'Catering', 'Photo', 'Video', 'Decor', 'Music', 'Cakes'];
// MVP locked to Almaty.
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
    return Scaffold(
      appBar: AppBar(title: const Text('Vendors')),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(12, 8, 12, 0),
            child: TextField(
              controller: _search,
              decoration: InputDecoration(
                prefixIcon: const Icon(Icons.search),
                hintText: 'Search…',
                suffixIcon: IconButton(
                  icon: const Icon(Icons.tune),
                  onPressed: _openFilters,
                ),
              ),
              onSubmitted: (_) => _apply(),
            ),
          ),
          if (state.loading) const LinearProgressIndicator(),
          if (state.total > 0)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              child: Align(
                alignment: Alignment.centerLeft,
                child: Text(
                  '${state.total} vendors',
                  style: const TextStyle(fontSize: 12, color: Colors.grey),
                ),
              ),
            ),
          Expanded(
            child: state.error != null
                ? Center(child: Text(state.error!))
                : state.items.isEmpty && !state.loading
                    ? const Center(child: Text('No vendors'))
                    : ListView.builder(
                        itemCount: state.items.length,
                        itemBuilder: (_, i) {
                          final v = state.items[i];
                          return ListTile(
                            title: Text(v.name),
                            subtitle: Row(
                              children: [
                                Text('${v.category} · ${v.city}'),
                                if (v.ratingCount > 0) ...[
                                  const SizedBox(width: 8),
                                  const Icon(Icons.star, size: 14, color: Colors.amber),
                                  Text(' ${v.ratingAvg.toStringAsFixed(1)}',
                                      style: const TextStyle(fontSize: 12)),
                                ],
                              ],
                            ),
                            trailing: Text('${v.priceFrom} ₸'),
                            onTap: () => context.push('/vendors/${v.id}'),
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
    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      builder: (_) => StatefulBuilder(
        builder: (ctx, setSt) => Padding(
          padding: EdgeInsets.only(
            left: 16,
            right: 16,
            top: 16,
            bottom: 16 + MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text('Filters', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),
              const SizedBox(height: 12),
              const Text('Category'),
              Wrap(
                spacing: 6,
                children: _categories
                    .map((c) => ChoiceChip(
                          label: Text(c),
                          selected: _category == c,
                          onSelected: (_) => setSt(() => _category = c),
                        ))
                    .toList(),
              ),
              const SizedBox(height: 12),
              const Text('Sort'),
              Wrap(
                spacing: 6,
                children: _sorts.entries
                    .map((e) => ChoiceChip(
                          label: Text(e.value),
                          selected: _sort == e.key,
                          onSelected: (_) => setSt(() => _sort = e.key),
                        ))
                    .toList(),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: priceCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'Max price (₸)'),
              ),
              const SizedBox(height: 16),
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
