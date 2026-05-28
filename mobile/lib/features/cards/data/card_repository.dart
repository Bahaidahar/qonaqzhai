import 'package:dio/dio.dart';

import '../../../core/network/api_endpoints.dart';
import '../domain/entities/card.dart';

class CardRepository {
  CardRepository(this._dio);
  final Dio _dio;

  Future<List<PaymentCard>> list() async {
    final res = await _dio.get(ApiEndpoints.cards);
    final items = (res.data['items'] as List?) ?? const [];
    return items.map((e) => PaymentCard.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<PaymentCard> add({
    required String number,
    required int expMonth,
    required int expYear,
    required String holder,
    bool makeDefault = false,
  }) async {
    final res = await _dio.post(ApiEndpoints.cards, data: {
      'number': number,
      'expMonth': expMonth,
      'expYear': expYear,
      'holder': holder,
      'makeDefault': makeDefault,
    });
    return PaymentCard.fromJson(res.data as Map<String, dynamic>);
  }

  Future<void> remove(String id) => _dio.delete(ApiEndpoints.card(id));

  Future<void> makeDefault(String id) => _dio.post(ApiEndpoints.cardDefault(id));
}
