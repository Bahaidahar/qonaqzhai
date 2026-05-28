import 'package:dio/dio.dart';

import '../../../core/network/api_endpoints.dart';

class PaymentRepository {
  PaymentRepository(this._dio);
  final Dio _dio;

  /// Mock pay-with-saved-card: marks booking as paid immediately.
  Future<Map<String, dynamic>> payMock(String bookingId, {String? cardId}) async {
    final res = await _dio.post(
      ApiEndpoints.mockPayment(bookingId),
      data: cardId != null ? {'cardId': cardId} : null,
    );
    return Map<String, dynamic>.from(res.data as Map);
  }

  /// PayBox-hosted checkout: returns a redirect URL the user opens externally.
  Future<String> startCheckout(String bookingId) async {
    final res = await _dio.post(ApiEndpoints.startPayment(bookingId));
    return (res.data as Map)['redirectUrl'] as String;
  }
}
