import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/repositories/vendor_repository_impl.dart';
import '../../domain/entities/service.dart';
import '../../domain/entities/vendor.dart';
import '../../domain/repositories/vendor_repository.dart';

final vendorRepositoryProvider = Provider<VendorRepository>((ref) {
  return VendorRepositoryImpl(ref.watch(dioProvider));
});

/// Public vendor services list — used by the customer-facing detail page.
final vendorServicesProvider =
    FutureProvider.family<List<VendorService>, String>((ref, id) {
  return ref.watch(vendorRepositoryProvider).services(id);
});

class VendorCatalogState {
  const VendorCatalogState({
    this.query = const VendorQuery(),
    this.items = const [],
    this.total = 0,
    this.loading = false,
    this.error,
  });

  final VendorQuery query;
  final List<Vendor> items;
  final int total;
  final bool loading;
  final String? error;

  VendorCatalogState copyWith({
    VendorQuery? query,
    List<Vendor>? items,
    int? total,
    bool? loading,
    String? error,
    bool clearError = false,
  }) {
    return VendorCatalogState(
      query: query ?? this.query,
      items: items ?? this.items,
      total: total ?? this.total,
      loading: loading ?? this.loading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class VendorCatalogViewModel extends StateNotifier<VendorCatalogState> {
  VendorCatalogViewModel(this._repo) : super(const VendorCatalogState()) {
    load();
  }

  final VendorRepository _repo;

  Future<void> load([VendorQuery? q]) async {
    final query = q ?? state.query;
    state = state.copyWith(loading: true, clearError: true, query: query);
    try {
      final res = await _repo.search(query);
      state = state.copyWith(items: res.items, total: res.total, loading: false);
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
    }
  }

  void setSearch(String text) => load(VendorQuery(
        query: text,
        category: state.query.category,
        city: state.query.city,
        priceMin: state.query.priceMin,
        priceMax: state.query.priceMax,
        sort: state.query.sort,
      ));

  void setCategory(String? cat) => load(VendorQuery(
        query: state.query.query,
        category: cat,
        city: state.query.city,
        sort: state.query.sort,
      ));
}

final vendorCatalogProvider =
    StateNotifierProvider<VendorCatalogViewModel, VendorCatalogState>((ref) {
  return VendorCatalogViewModel(ref.watch(vendorRepositoryProvider));
});
