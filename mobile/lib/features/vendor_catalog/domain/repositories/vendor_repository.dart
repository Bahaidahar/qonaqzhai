import '../entities/vendor.dart';

abstract class VendorRepository {
  Future<VendorSearchResult> search(VendorQuery query);
  Future<Vendor> byId(String id);
}
