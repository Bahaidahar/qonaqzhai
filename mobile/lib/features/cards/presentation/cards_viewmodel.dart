import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/dio_client.dart';
import '../data/card_repository.dart';
import '../domain/entities/card.dart';

final cardRepositoryProvider = Provider<CardRepository>((ref) {
  return CardRepository(ref.watch(dioProvider));
});

class CardsState {
  const CardsState({this.items = const [], this.loading = false, this.error});
  final List<PaymentCard> items;
  final bool loading;
  final String? error;

  CardsState copyWith({List<PaymentCard>? items, bool? loading, String? error, bool clearError = false}) {
    return CardsState(
      items: items ?? this.items,
      loading: loading ?? this.loading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class CardsViewModel extends StateNotifier<CardsState> {
  CardsViewModel(this._repo) : super(const CardsState(loading: true)) {
    refresh();
  }
  final CardRepository _repo;

  Future<void> refresh() async {
    state = state.copyWith(loading: true, clearError: true);
    try {
      state = state.copyWith(items: await _repo.list(), loading: false);
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
    }
  }

  Future<bool> add({
    required String number,
    required int expMonth,
    required int expYear,
    required String holder,
    bool makeDefault = false,
  }) async {
    try {
      await _repo.add(number: number, expMonth: expMonth, expYear: expYear, holder: holder, makeDefault: makeDefault);
      await refresh();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  Future<void> remove(String id) async {
    await _repo.remove(id);
    await refresh();
  }

  Future<void> makeDefault(String id) async {
    await _repo.makeDefault(id);
    await refresh();
  }
}

final cardsViewModelProvider = StateNotifierProvider<CardsViewModel, CardsState>((ref) {
  return CardsViewModel(ref.watch(cardRepositoryProvider));
});
