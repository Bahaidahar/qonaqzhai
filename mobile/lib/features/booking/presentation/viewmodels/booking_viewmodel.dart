import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/repositories/booking_repository_impl.dart';
import '../../domain/entities/booking.dart';
import '../../domain/repositories/booking_repository.dart';

final bookingRepositoryProvider = Provider<BookingRepository>((ref) {
  return BookingRepositoryImpl(ref.watch(dioProvider));
});

class BookingsState {
  const BookingsState({this.items = const [], this.loading = false, this.error});
  final List<Booking> items;
  final bool loading;
  final String? error;

  BookingsState copyWith({List<Booking>? items, bool? loading, String? error, bool clearError = false}) {
    return BookingsState(
      items: items ?? this.items,
      loading: loading ?? this.loading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class BookingsViewModel extends StateNotifier<BookingsState> {
  BookingsViewModel(this._repo) : super(const BookingsState(loading: true)) {
    refresh();
  }

  final BookingRepository _repo;

  Future<void> refresh() async {
    state = state.copyWith(loading: true, clearError: true);
    try {
      final list = await _repo.myBookings();
      state = state.copyWith(items: list, loading: false);
    } catch (e) {
      state = state.copyWith(loading: false, error: e.toString());
    }
  }

  Future<void> cancel(String id) async {
    await _repo.updateStatus(id, 'cancelled');
    await refresh();
  }

  Future<String> startPayment(String id) => _repo.startPayment(id);
}

final bookingsProvider =
    StateNotifierProvider<BookingsViewModel, BookingsState>((ref) {
  return BookingsViewModel(ref.watch(bookingRepositoryProvider));
});
