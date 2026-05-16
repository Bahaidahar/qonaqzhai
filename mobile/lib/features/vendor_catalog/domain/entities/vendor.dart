class Vendor {
  const Vendor({
    required this.id,
    required this.userId,
    required this.name,
    required this.category,
    required this.city,
    required this.description,
    required this.priceFrom,
    required this.status,
    required this.ratingAvg,
    required this.ratingCount,
    required this.photoIds,
  });

  final String id;
  final String userId;
  final String name;
  final String category;
  final String city;
  final String description;
  final int priceFrom;
  final String status;
  final double ratingAvg;
  final int ratingCount;
  final List<String> photoIds;
}

class VendorSearchResult {
  const VendorSearchResult({
    required this.items,
    required this.total,
    required this.page,
    required this.limit,
  });

  final List<Vendor> items;
  final int total;
  final int page;
  final int limit;
}

class VendorQuery {
  const VendorQuery({
    this.query,
    this.category,
    this.city,
    this.priceMin,
    this.priceMax,
    this.ratingMin,
    this.sort,
    this.page = 1,
    this.limit = 20,
  });

  final String? query;
  final String? category;
  final String? city;
  final int? priceMin;
  final int? priceMax;
  final double? ratingMin;
  final String? sort;
  final int page;
  final int limit;
}
