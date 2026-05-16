import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../domain/entities/review.dart';

abstract class ReviewRepository {
  Future<List<Review>> listByVendor(String vendorId);
  Future<Review> submit({
    required String bookingId,
    required int rating,
    String? text,
  });
}

class ReviewRepositoryImpl implements ReviewRepository {
  ReviewRepositoryImpl(this._dio);

  final Dio _dio;

  Review _fromJson(Map<String, dynamic> json) => Review(
        id: json['id'] as String,
        bookingId: json['bookingId'] as String,
        customerId: json['customerId'] as String,
        vendorId: json['vendorId'] as String,
        rating: (json['rating'] as num).toInt(),
        text: (json['text'] as String?) ?? '',
        createdAt: json['createdAt'] as String,
      );

  @override
  Future<List<Review>> listByVendor(String vendorId) async {
    final res = await _dio.get(
      ApiEndpoints.vendorReviews(vendorId),
      options: Options(extra: {'requiresAuth': false}),
    );
    return ((res.data['items'] as List?) ?? const [])
        .map((e) => _fromJson(e as Map<String, dynamic>))
        .toList();
  }

  @override
  Future<Review> submit({
    required String bookingId,
    required int rating,
    String? text,
  }) async {
    final res = await _dio.post(ApiEndpoints.reviews, data: {
      'bookingId': bookingId,
      'rating': rating,
      if (text != null) 'text': text,
    });
    return _fromJson(res.data as Map<String, dynamic>);
  }
}
