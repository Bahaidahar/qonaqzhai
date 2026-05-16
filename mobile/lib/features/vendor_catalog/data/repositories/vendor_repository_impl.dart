import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../domain/entities/vendor.dart';
import '../../domain/repositories/vendor_repository.dart';

class VendorRepositoryImpl implements VendorRepository {
  VendorRepositoryImpl(this._dio);

  final Dio _dio;

  Vendor _fromJson(Map<String, dynamic> json) => Vendor(
        id: json['id'] as String,
        userId: json['userId'] as String,
        name: json['name'] as String,
        category: json['category'] as String,
        city: json['city'] as String,
        description: (json['description'] as String?) ?? '',
        priceFrom: (json['priceFrom'] as num?)?.toInt() ?? 0,
        status: json['status'] as String,
        ratingAvg: ((json['ratingAvg'] as num?) ?? 0).toDouble(),
        ratingCount: (json['ratingCount'] as num?)?.toInt() ?? 0,
        photoIds: ((json['photoIds'] as List?) ?? const [])
            .whereType<String>()
            .toList(),
      );

  @override
  Future<VendorSearchResult> search(VendorQuery q) async {
    final params = <String, dynamic>{
      if (q.query != null && q.query!.isNotEmpty) 'q': q.query,
      if (q.category != null && q.category!.isNotEmpty) 'category': q.category,
      if (q.city != null && q.city!.isNotEmpty) 'city': q.city,
      if (q.priceMin != null) 'price_min': q.priceMin,
      if (q.priceMax != null) 'price_max': q.priceMax,
      if (q.ratingMin != null) 'rating_min': q.ratingMin,
      if (q.sort != null) 'sort': q.sort,
      'page': q.page,
      'limit': q.limit,
    };
    final res = await _dio.get(
      ApiEndpoints.vendors,
      queryParameters: params,
      options: Options(extra: {'requiresAuth': false}),
    );
    final data = res.data as Map<String, dynamic>;
    final items = ((data['items'] as List?) ?? const [])
        .map((e) => _fromJson(e as Map<String, dynamic>))
        .toList();
    return VendorSearchResult(
      items: items,
      total: (data['total'] as num?)?.toInt() ?? items.length,
      page: (data['page'] as num?)?.toInt() ?? q.page,
      limit: (data['limit'] as num?)?.toInt() ?? q.limit,
    );
  }

  @override
  Future<Vendor> byId(String id) async {
    final res = await _dio.get(
      ApiEndpoints.vendor(id),
      options: Options(extra: {'requiresAuth': false}),
    );
    return _fromJson(res.data as Map<String, dynamic>);
  }
}
