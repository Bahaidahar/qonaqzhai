/// A bookable line-item attached to a vendor. The web app uses these on the
/// vendor detail page so customers can pick the exact package they want
/// (e.g. catering pricing per-person vs a fixed venue rental).
class VendorService {
  const VendorService({
    required this.id,
    required this.vendorId,
    required this.name,
    required this.description,
    required this.price,
    required this.unit,
  });

  final String id;
  final String vendorId;
  final String name;
  final String description;
  final int price;
  final String unit; // 'fixed' | 'hour' | 'item' | 'person' | 'day'
}
