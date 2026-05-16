import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../domain/entities/booking.dart';
import '../../domain/repositories/booking_repository.dart';

class BookingRepositoryImpl implements BookingRepository {
  BookingRepositoryImpl(this._dio);

  final Dio _dio;

  Booking _fromJson(Map<String, dynamic> json) => Booking(
        id: json['id'] as String,
        customerId: json['customerId'] as String,
        vendorId: json['vendorId'] as String,
        eventDate: json['eventDate'] as String,
        guestCount: (json['guestCount'] as num?)?.toInt() ?? 0,
        note: (json['note'] as String?) ?? '',
        status: json['status'] as String,
        amount: (json['amount'] as num?)?.toInt() ?? 0,
        paymentId: (json['paymentId'] as String?) ?? '',
      );

  @override
  Future<List<Booking>> myBookings() async {
    final res = await _dio.get(ApiEndpoints.bookings);
    return ((res.data['items'] as List?) ?? const [])
        .map((e) => _fromJson(e as Map<String, dynamic>))
        .toList();
  }

  @override
  Future<Booking> create({
    required String vendorId,
    required String eventDate,
    required int guestCount,
    String? note,
    int? amount,
  }) async {
    final res = await _dio.post(ApiEndpoints.bookings, data: {
      'vendorId': vendorId,
      'eventDate': eventDate,
      'guestCount': guestCount,
      if (note != null) 'note': note,
      if (amount != null) 'amount': amount,
    });
    return _fromJson(res.data as Map<String, dynamic>);
  }

  @override
  Future<Booking> updateStatus(String id, String status) async {
    final res = await _dio.patch(ApiEndpoints.booking(id), data: {'status': status});
    return _fromJson(res.data as Map<String, dynamic>);
  }

  @override
  Future<String> startPayment(String bookingId) async {
    final res = await _dio.post(ApiEndpoints.startPayment(bookingId));
    return (res.data['redirectUrl'] as String?) ?? '';
  }
}
