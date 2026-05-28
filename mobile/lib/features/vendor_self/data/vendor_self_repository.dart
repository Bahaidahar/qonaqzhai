import 'package:dio/dio.dart';

import '../../../core/network/api_endpoints.dart';

class VendorSelfRepository {
  VendorSelfRepository(this._dio);
  final Dio _dio;

  Future<Map<String, dynamic>?> myVendor() async {
    try {
      final res = await _dio.get(ApiEndpoints.myVendor);
      return Map<String, dynamic>.from(res.data as Map);
    } on DioException catch (e) {
      if (e.response?.statusCode == 404) return null;
      rethrow;
    }
  }

  Future<Map<String, dynamic>> upsert({
    required String name,
    required String category,
    required String city,
    required int priceFrom,
    required String description,
  }) async {
    final res = await _dio.post(ApiEndpoints.myVendor, data: {
      'name': name,
      'category': category,
      'city': city,
      'priceFrom': priceFrom,
      'description': description,
    });
    return Map<String, dynamic>.from(res.data as Map);
  }

  Future<List<Map<String, dynamic>>> myServices() async {
    final res = await _dio.get(ApiEndpoints.vendorServices);
    final items = (res.data['items'] as List?) ?? const [];
    return items.map((e) => Map<String, dynamic>.from(e as Map)).toList();
  }

  Future<void> addService({
    required String name,
    required String description,
    required int price,
    required String unit,
  }) async {
    await _dio.post(ApiEndpoints.vendorServices, data: {
      'name': name,
      'description': description,
      'price': price,
      'unit': unit,
    });
  }

  Future<void> updateService({
    required String id,
    required String name,
    required String description,
    required int price,
    required String unit,
  }) async {
    await _dio.patch(ApiEndpoints.vendorService(id), data: {
      'name': name,
      'description': description,
      'price': price,
      'unit': unit,
    });
  }

  Future<void> deleteService(String id) => _dio.delete(ApiEndpoints.vendorService(id));

  Future<void> uploadPhoto({required List<int> bytes, required String filename}) async {
    final form = FormData.fromMap({
      'photo': MultipartFile.fromBytes(bytes, filename: filename),
    });
    await _dio.post(ApiEndpoints.vendorPhotos, data: form);
  }

  Future<void> deletePhoto(String id) => _dio.delete(ApiEndpoints.vendorPhoto(id));
}
